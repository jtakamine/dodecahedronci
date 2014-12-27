package main

import (
	"bufio"
	"github.com/jtakamine/dodecahedronci/config"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildDockerImages(repoDir string) (err error) {
	dFiles, err := getDockerfiles(repoDir)
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

func getDockerfiles(dir string) (dFiles []string, err error) {
	dFiles = []string{}
	walk := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), "Dockerfile") {
			dFiles = append(dFiles, path)
		}

		return err
	}

	err = filepath.Walk(dir, walk)
	if err != nil {
		return nil, err
	}

	return dFiles, nil
}

func getImageNameHint(dockerFile string) (hint string, err error) {
	hint = "builtbydodecci" //default image name hint
	hintPrefix := "#imagenamehint:"

	file, err := os.Open(dockerFile)
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
