package main

import (
	"encoding/json"
)

func parseGitHubRequest(data []byte) (repoUrl string, appName string, err error) {
	req := &struct {
		Repository struct {
			Name      string
			Clone_url string
		}
	}{}

	err = json.Unmarshal(data, req)
	if err != nil {
		return "", "", err
	}

	return req.Repository.Clone_url, req.Repository.Name, nil
}
