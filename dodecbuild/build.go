package main

import (
	"bufio"
	"errors"
	"github.com/jtakamine/dodecahedronci/configutil"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type logType int

const (
	verboseLogType logType = iota
	infoLogType
	warnLogType
	errorLogType
)

type figFile struct {
	File   string
	Config map[interface{}]interface{}
}

type dockerFile struct {
	FigService string
	File       string
}

func build(repoDir string, app string, dockerRegistryUrl string) (err error) {
	logPub("Retrieving next version number...", infoLogType)
	version := getNextVersion(app)
	logPub("Retrieved version number: "+version, infoLogType)

	logPub("Searching for fig files in "+repoDir+"...", infoLogType)
	files, err := findFigFiles(repoDir)
	if err != nil {
		logPub("Error encountered while searching for fig files: "+err.Error(), errorLogType)
		return err
	}
	logPub("Found "+strconv.Itoa(len(files))+" fig file(s).", infoLogType)

	logPub("Parsing "+strconv.Itoa(len(files))+" fig file(s)...", infoLogType)
	fFiles := []figFile{}
	for i, file := range files {
		logPub("Parsing fig file #"+strconv.Itoa(i+1)+": "+file, verboseLogType)
		fFile, err := parseFigFile(file)
		if err != nil {
			logPub("Error encountered while parsing fig file: "+err.Error(), errorLogType)
			return err
		}

		fFiles = append(fFiles, fFile)
	}
	logPub("Parsed "+strconv.Itoa(len(fFiles))+" fig files.", infoLogType)

	logPub("Looping through "+strconv.Itoa(len(fFiles))+" parsed fig files...", infoLogType)
	for i, fFile := range fFiles {
		logPub("Extracting Dockerfile paths from fig file #"+strconv.Itoa(i+1)+"...", verboseLogType)
		dFiles, err := getDockerFiles(fFile)
		if err != nil {
			logPub("Error encountered while extracting Dockerfile paths from fig file: "+err.Error(), errorLogType)
			return err
		}
		logPub("Extracted "+strconv.Itoa(len(dFiles))+" Dockerfile paths.", verboseLogType)

		logPub("Looping through "+strconv.Itoa(len(dFiles))+" Dockerfiles...", verboseLogType)
		for j, dFile := range dFiles {
			logPub("Processing Dockerfile #"+strconv.Itoa(j+1)+"...", verboseLogType)
			logPub("Retrieving Docker repository name based on Dockerfile path...", verboseLogType)
			repo, err := getDockerRepo(dFile.File)
			if err != nil {
				logPub("Error encountered while retrieving Docker repository name: "+err.Error(), errorLogType)
				return err
			}
			logPub("Retrieved Docker repository name: "+repo, verboseLogType)

			logPub("Generating Docker image tag...", verboseLogType)
			tag := getDockerTag(dockerRegistryUrl, configutil.Get("DODEC_DOCKER_USER"), repo, version)
			logPub("Generated Docker image tag: "+tag, verboseLogType)

			logPub("Building Dockerfile "+dFile.File+"...", verboseLogType)
			err = buildDockerFile(dFile.File, tag)
			if err != nil {
				logPub("Error encountered while building Dockerfile: "+err.Error(), errorLogType)
				return err
			}
			logPub("Built Dockerfile.", verboseLogType)

			logPub("Replacing \"build\" node with appropriate \"image\" node in Fig file...", verboseLogType)
			err = updateFigFileWithDockerImage(fFile, dFile, tag)
			if err != nil {
				logPub("Error encountered while replacing \"build\" node with \"image\" node: "+err.Error(), errorLogType)
				return err
			}
			logPub("Replaced \"build\" node with \"image\" node in Fig file.", verboseLogType)
		}

		logPub("Posting the build to the Dodec Registry...", verboseLogType)
		err = saveBuild(app, version, fFile, dockerRegistryUrl)
		if err != nil {
			logPub("Error encountered while posting the build to the Dodec Registry: "+err.Error(), errorLogType)
			return err
		}
		logPub("Posted the build to the Dodec Registry.", verboseLogType)
	}
	logPub("Done looping through "+strconv.Itoa(len(fFiles))+" Fig files.", infoLogType)

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

func getDockerRepo(dFile string) (repo string, err error) {
	dir := filepath.Dir(dFile)
	parts := strings.Split(dir, "/")

	//default repo name is the name of the directory containing the Dockerfile
	repo = parts[len(parts)-1]

	repoHint := "#repoHint:"

	file, err := os.Open(dFile)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, repoHint) {
			repo = strings.TrimPrefix(line, repoHint)
			return repo, nil
		}
	}

	err = scanner.Err()
	if err != nil {
		return "", err
	}

	return repo, nil
}

func getDockerTag(dockerRegistryUrl string, dockerUser string, dockerRepo string, version string) (tag string) {
	registryPrefix := ""
	if dockerRegistryUrl != "" {
		dockerRegistryUrl = strings.TrimPrefix(dockerRegistryUrl, "http://")
		dockerRegistryUrl = strings.TrimPrefix(dockerRegistryUrl, "https://")
		registryPrefix = dockerRegistryUrl + "/"
	}

	userPrefix := ""
	if dockerUser != "" {
		userPrefix = dockerUser + "/"
	}

	versionSuffix := ""
	if version != "" {
		versionSuffix = ":" + version
	}

	tag = registryPrefix + userPrefix + dockerRepo + versionSuffix
	return tag
}

func buildDockerFile(dFile string, tag string) (err error) {
	cmd := exec.Command("docker", "build", "-t", tag, ".")
	cmd.Dir = filepath.Dir(dFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
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
