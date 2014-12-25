package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/jtakamine/dodecahedronci/config"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type gitHubReq struct {
	Repository struct {
		Id      int
		Ssh_url string
	}
}

func gitHubHandle(w http.ResponseWriter, r *http.Request) {
	req := &gitHubReq{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(req)

	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		log.Panicf("Could not parse JSON: %v\n", err)
	}

	repoDir := cloneOrUpdateGitRepo(req.Repository.Id, req.Repository.Ssh_url)
	buildDockerImages(repoDir)
}

func cloneOrUpdateGitRepo(repoId int, repoUrl string) string {
	dir := strings.TrimSuffix(config.Get("DODEC_HOME"), "/") + "/" + strconv.Itoa(repoId)

	var cmd *exec.Cmd

	if fInfo, err := os.Stat(dir); os.IsNotExist(err) || !fInfo.IsDir() {
		log.Println("Cloning git repo.")
		cmd = exec.Command("git", "clone", repoUrl, dir)
	} else {
		log.Println("Pulling git repo.")
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
