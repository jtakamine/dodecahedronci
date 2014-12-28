package main

import (
	"encoding/json"
)

type gitHubReq struct {
	Repository struct {
		Id        int
		Ssh_url   string
		Clone_url string
	}
}

func parseGitHubRequest(data []byte) (repoUrl string, err error) {
	req := &gitHubReq{}
	err = json.Unmarshal(data, req)
	if err != nil {
		return "", err
	}

	return req.Repository.Clone_url, nil
}
