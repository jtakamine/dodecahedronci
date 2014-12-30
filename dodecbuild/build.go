package main

import (
	"bufio"
	"errors"
	"github.com/jtakamine/dodecahedronci/dodecregistry/api"
	"github.com/jtakamine/dodecahedronci/utils/configutil"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	//Retrieve next version number
	version := getNextVersion(app)

	//Find all Fig files in repoDir
	files, err := findFigFiles(repoDir)
	if err != nil {
		return err
	}

	//Parse Fig files
	fFiles := []figFile{}
	for _, file := range files {
		fFile, err := parseFigFile(file)
		if err != nil {
			return err
		}

		fFiles = append(fFiles, fFile)
	}

	//Loop through Fig files
	for _, fFile := range fFiles {
		dFiles, err := getDockerFiles(fFile)
		if err != nil {
			return err
		}

		//Loop through Dockerfiles
		for _, dFile := range dFiles {
			repo, err := getDockerRepo(dFile.File)
			if err != nil {
				return err
			}

			tag := getDockerTag(dockerRegistryUrl, configutil.Get("DODEC_DOCKER_USER"), repo, version)

			//Build Dockerfile
			err = buildDockerFile(dFile.File, tag)
			if err != nil {
				return err
			}

			//In Fig file, replace "build" node with appropriate "image" node
			err = updateFigFileWithDockerImage(fFile, dFile, tag)
			if err != nil {
				return err
			}
		}

		//Post the build to the dodec registry
		err = postBuildToDodecRegistry(app, version, fFile, dockerRegistryUrl)
		if err != nil {
			return err
		}
	}

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
	log.Printf("building %v\n", tag)

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

var postBuildToDodecRegistry = func(app string, version string, fFile figFile, dockerRegistryUrl string) (err error) {
	data, err := yaml.Marshal(fFile.Config)
	if err != nil {
		return err
	}
	artifact := string(data)

	build := api.Build{Artifact: artifact, DockerRegistryUrl: dockerRegistryUrl}

	err = api.PostBuild(app, version, build, "http://dodecregistry:8000/")
	if err != nil {
		return err
	}

	return nil
}
