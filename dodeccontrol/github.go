package main

import (
	"encoding/json"
)

func parseGitHubRequest(data []byte) (repoUrl string, err error) {
	req := &struct {
		Repository struct {
			Id        int
			Ssh_url   string
			Clone_url string
		}
	}{}

	err = json.Unmarshal(data, req)
	if err != nil {
		return "", err
	}

	return req.Repository.Clone_url, nil
}
