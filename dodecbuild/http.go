package main

import (
	"io/ioutil"
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
		panic("Error reading request body: " + err.Error())
	}

	//Eventually, take a look at the header/body to determine which handler to use.  For now assume it's a github request
	repoUrl, err := parseGitHubRequest(data)
	if err != nil {
		panic("Error parsing GitHub request: " + err.Error())
	}

	repoDir, err := cloneOrUpdateGitRepo(repoUrl)
	if err != nil {
		panic("Error cloning or updating git repo: " + err.Error())
	}

	err = build(repoDir, "myAwesomeApp", "http://localhost:8080")
	if err != nil {
		panic("Error building Docker images: " + err.Error())
	}
}
