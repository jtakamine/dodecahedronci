package handlers

import (
	"net/http"
	"encoding/json"
	"log"
	"bytes"
	"os/exec"
	"os"
	"strconv"
	"path/filepath"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	req := &gitHubReq{}

	err := decoder.Decode(req)

	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		log.Panicf("Could not parse JSON: %v\n", err)
		return
	}

	requestBuild(req.Repository.Id, req.Repository.Ssh_url)
}

type gitHubReq struct {
	Repository gitHubRepo
}

type gitHubRepo struct {
	Id int
	Ssh_url string
}

//Ideally, this method should send a build request
// to some queue which is asynchronously consumed by a separate
// build server.  For now, we execute the build synchronously
// and on the same server.
func requestBuild(repoId int, repoUrl string) {
	log.Printf("Triggering build for repo with url: %v\n", repoUrl)

	repoDir := cloneOrUpdateGitRepo(repoId, repoUrl)
	buildDockerImages(repoDir)
}

func cloneOrUpdateGitRepo(repoId int, repoUrl string) string {
	dir := "/var/lib/dodecci/" + strconv.Itoa(repoId)

	var cmd *exec.Cmd

	if fInfo, err := os.Stat(dir); os.IsNotExist(err) || !fInfo.IsDir() {
		cmd = exec.Command("git", "clone", repoUrl, dir)
	} else {
		cmd = exec.Command("git", "pull", repoUrl)
		cmd.Dir = dir
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Panicf("Error running git operation: %v\n", err)
	}

	return dir
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

		cmd := exec.Command("docker", "build", "-t", "jtakamine/autobuild", ".")
		cmd.Dir = filepath.Dir(dFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Panicf("Error building Dockerfile: %v\n", err)
		}

		cmd = exec.Command("docker", "push", "jtakamine/autobuild")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			log.Panicf("Error pushing Docker image: %v\n", err)
		}
	}
}
