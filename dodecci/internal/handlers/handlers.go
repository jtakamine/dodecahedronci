package handlers

import (
	"net/http"
	"log"
	"os/exec"
	"os"
	"bufio"
	"path/filepath"
	"strings"
	"github.com/jtakamine/dodecahedronci/config"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	//Eventually, take a look at the header/body to determine which handler to use.  For now assume it's a github request
	gitHubHandle(w, r)
}

func buildDockerImages(repoDir string) {
	dockerFiles := []string{}

	walk := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "Dockerfile" {
			dockerFiles = append(dockerFiles, path)
		}

		return nil
	}

	err := filepath.Walk(repoDir, walk)
	if err != nil {
		log.Panicf("Error walking the directory \"%v\": %v\n", repoDir, err)
	}

	for _,dFile := range dockerFiles {
		log.Printf("Building Docker file: %v\n", dFile)

		cmd := exec.Command("docker", "build", "-t", config.Get("DODEC_DOCKER_USER") + "/builtbydodec", ".")
		cmd.Dir = filepath.Dir(dFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Panicf("Error building Dockerfile: %v\n", err)
		}

		//Commenting out the below--do not automatically push to Docker repo
		//  for now (it's slow).  Instead, assume that dodecdeploy will run on
		//  the same server/container, so it will have access to the "local"
		//  repo which will be automatically available as soon as a Docker image
		//  is built.

		/*
		cmd = exec.Command("docker", "push", "jtakamine/autobuild")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Panicf("Error pushing Docker image: %v\n", err)
		}*/
	}
}

func getImageNameHint(dockerFile string) string {
	hintPrefix := "#imagenamehint:"

	file, err := os.Open(dockerFile)
	if err != nil {
		log.Panicf("Error opening Dockerfile: %v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, hintPrefix) {
			return strings.TrimPrefix(line, hintPrefix)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Panicf("Error reading Dockerfile: %v\n", err)
	}

	return ""
}
