package main

import (
	"encoding/json"
	"net/http"
)

type gitHubReq struct {
	Repository struct {
		Id        int
		Ssh_url   string
		Clone_url string
	}
}

func parseGitHubRequest(w http.ResponseWriter, r *http.Request) (repoUrl string, err error) {
	req := &gitHubReq{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(req)
	if err != nil {
		return "", err
	}

	return req.Repository.Clone_url, nil
}
