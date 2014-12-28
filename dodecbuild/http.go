package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func ListenAndServe(addr string) (err error) {
	http.HandleFunc("/", httpHandle)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}

	return nil
}

func httpHandle(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicf("Error reading request body: %v\n", err)
	}

	//Eventually, take a look at the header/body to determine which handler to use.  For now assume it's a github request
	repoUrl, err := parseGitHubRequest(data)
	if err != nil {
		log.Panicf("Error parsing GitHub request: %v\n", err)
	}

	repoDir, err := cloneOrUpdateGitRepo(repoUrl)
	if err != nil {
		log.Panicf("Error cloning or updating git repo: %v\n", err)
	}

	err = build(repoDir, "myAwesomeApp", "http://localhost:8080")
	if err != nil {
		log.Panicf("Error building Docker images: %v\n", err)
	}

	fmt.Fprint(w, "build successful\n")
}
