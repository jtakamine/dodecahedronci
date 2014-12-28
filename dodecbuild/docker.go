package main

import (
	"bufio"
	"github.com/jtakamine/dodecahedronci/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func buildDockerImages(repoDir string) (err error) {
	dFiles, err := getDFiles(repoDir)
	if err != nil {
		return err
	}

	for _, dFile := range dFiles {
		log.Printf("Building Docker file: %v\n", dFile)

		imgName, err := getImageNameHint(dFile)
		if err != nil {
			return err
		}

		cmd := exec.Command("docker", "build", "-t", config.Get("DODEC_DOCKER_USER")+"/"+imgName, ".")
		cmd.Dir = filepath.Dir(dFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func getDFiles(dir string) (dFiles []string, err error) {
	dFiles = []string{}
	walk := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.EqualFold(info.Name(), "fig.yml") {
			files, err := getDFilesFromFigYml(path)
			if err != nil {
				return err
			}
			dFiles = append(dFiles, files...)
		}

		return nil
	}

	err = filepath.Walk(dir, walk)
	if err != nil {
		return nil, err
	}

	return dFiles, nil
}

func getDFilesFromFigYml(fyml string) (dFiles []string, err error) {
	dFiles = []string{}

	data, err := ioutil.ReadFile(fyml)
	if err != nil {
		return nil, err
	}

	config := make(map[string]interface{})
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	for _, v := range config {
		if m, ok := v.(map[interface{}]interface{}); ok {
			buildPath := m["build"]
			if buildPath != nil {
				file := path.Join(path.Dir(fyml), buildPath.(string), "Dockerfile")
				dFiles = append(dFiles, file)
			}
		}
	}

	return dFiles, nil
}

func getImageNameHint(dFile string) (hint string, err error) {
	hint = "builtbydodecci" //default image name hint
	hintPrefix := "#imagenamehint:"

	file, err := os.Open(dFile)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, hintPrefix) {
			hint = strings.TrimPrefix(line, hintPrefix)
			return hint, nil
		}
	}

	err = scanner.Err()
	if err != nil {
		return "", err
	}

	return hint, nil
}
