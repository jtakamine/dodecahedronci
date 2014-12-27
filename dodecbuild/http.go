package main

import (
	"fmt"
	"log"
	"net/http"
)

func httpHandle(w http.ResponseWriter, r *http.Request) {
	//Eventually, take a look at the header/body to determine which handler to use.  For now assume it's a github request
	repoUrl, err := parseGitHubRequest(w, r)
	if err != nil {
		log.Panicf("Error parsing GitHub request: %v\n", err)
	}

	repoDir, err := cloneOrUpdateGitRepo(repoUrl)
	if err != nil {
		log.Panicf("Error cloning or updating git repo: %v\n", err)
	}

	err = buildDockerImages(repoDir)
	if err != nil {
		log.Panicf("Error building Docker images: %v\n", err)
	}

	fmt.Fprint(w, "build successful\n")
}
