package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type gitHubReq struct {
	Repository struct {
		Id        int
		Ssh_url   string
		Clone_url string
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

	repoDir := cloneOrUpdateGitRepo(req.Repository.Id, req.Repository.Clone_url)
	buildDockerImages(repoDir)
}
