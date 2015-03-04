package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/jtakamine/dodecahedronci/logutil"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const appDataDir = "/var/lib/dodec/"

type figFile struct {
	File   string
	Config map[interface{}]interface{}
}

type dockerFile struct {
	FigService string
	File       string
}

func generateRandID(length int) string {
	id := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		panic(err)
	}
	return hex.EncodeToString(id)
}

func build(repoDir string, uuid string, appName string, version string, writer *logutil.Writer) (err error) {
	w := writer.WriteType
	wIn := writer.Indent
	wOut := writer.Outdent

	w("Searching for fig files in "+repoDir+"...", logutil.Info)
	files, err := findFigFiles(repoDir)
	if err != nil {
		w("Error encountered while searching for fig files: "+err.Error(), logutil.Error)
		return err
	}
	w("Found "+strconv.Itoa(len(files))+" fig file(s).", logutil.Info)

	w("Parsing "+strconv.Itoa(len(files))+" fig file(s)...", logutil.Info)
	fFiles := []figFile{}
	wIn()
	for i, file := range files {
		w("Parsing fig file #"+strconv.Itoa(i+1)+": "+file, logutil.Verbose)
		fFile, err := parseFigFile(file)
		if err != nil {
			w("Error encountered while parsing fig file: "+err.Error(), logutil.Error)
			return err
		}

		fFiles = append(fFiles, fFile)
	}
	wOut()
	w("Parsed "+strconv.Itoa(len(fFiles))+" fig files.", logutil.Info)

	w("Looping through "+strconv.Itoa(len(fFiles))+" parsed fig files...", logutil.Info)
	wIn()
	for i, fFile := range fFiles {
		w("Extracting Dockerfile paths from fig file #"+strconv.Itoa(i+1)+"...", logutil.Verbose)
		dFiles, err := getDockerFiles(fFile)
		if err != nil {
			w("Error encountered while extracting Dockerfile paths from fig file: "+err.Error(), logutil.Error)
			return err
		}
		w("Extracted "+strconv.Itoa(len(dFiles))+" Dockerfile paths.", logutil.Verbose)

		w("Looping through "+strconv.Itoa(len(dFiles))+" Dockerfiles...", logutil.Verbose)
		wIn()
		for j, dFile := range dFiles {
			w("Processing Dockerfile #"+strconv.Itoa(j+1)+"...", logutil.Verbose)

			w("Building Dockerfile "+dFile.File+"...", logutil.Verbose)
			tag, err := buildDockerFile(dFile.File, version, writer.CreateChild())
			if err != nil {
				w("Error encountered while building Dockerfile: "+err.Error(), logutil.Error)
				return err
			}
			w("Built Dockerfile.", logutil.Verbose)

			w("Pushing Docker image "+tag+"...", logutil.Verbose)
			err = pushDockerImage(tag, writer.CreateChild())
			w("Pushed Docker image "+tag+".", logutil.Verbose)

			w("Replacing \"build\" node with appropriate \"image\" node in Fig file...", logutil.Verbose)
			err = updateFigFileWithDockerImage(fFile, dFile, tag)
			if err != nil {
				w("Error encountered while replacing \"build\" node with \"image\" node: "+err.Error(), logutil.Error)
				return err
			}
			w("Replaced \"build\" node with \"image\" node in Fig file.", logutil.Verbose)
		}
		wOut()
		w("Saving build artifact...", logutil.Verbose)
		err = saveBuildArtifact(uuid, fFile)
		w("Saved build artifact.", logutil.Verbose)
	}
	wOut()
	w("Done looping through "+strconv.Itoa(len(fFiles))+" Fig files.", logutil.Info)

	return nil
}

func findFigFiles(repoDir string) (fFiles []string, err error) {
	fFiles = []string{}
	walk := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.EqualFold(info.Name(), "fig.yml") {
			fFiles = append(fFiles, path)
		}

		return nil
	}

	err = filepath.Walk(repoDir, walk)
	if err != nil {
		return nil, err
	}

	return fFiles, nil
}

func parseFigFile(file string) (fFile figFile, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return figFile{}, err
	}

	config := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return figFile{}, err
	}

	return figFile{File: file, Config: config}, nil
}

func getDockerFiles(fFile figFile) (dFiles []dockerFile, err error) {
	dFiles = []dockerFile{}
	for k, v := range fFile.Config {
		if m, ok := v.(map[interface{}]interface{}); ok {
			if buildPath, ok := m["build"]; ok {
				file := filepath.Join(filepath.Dir(fFile.File), buildPath.(string), "Dockerfile")
				dFiles = append(dFiles, dockerFile{FigService: k.(string), File: file})
			}
		}
	}

	return dFiles, nil
}

func updateFigFileWithDockerImage(fFile figFile, dFile dockerFile, tag string) (err error) {
	serviceConfig := fFile.Config[dFile.FigService]

	if m, ok := serviceConfig.(map[interface{}]interface{}); ok {
		delete(m, "build")
		m["image"] = tag
		fFile.Config[dFile.FigService] = m

		return nil
	}

	return errors.New("Could not interpret Fig file.")
}
