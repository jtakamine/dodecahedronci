package handlers

import (
	"net/http"
	"encoding/json"
	"log"
	"bytes"
	"os/exec"
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

	requestBuild(req.Repository.Ssh_url)
}

type gitHubReq struct {
	Repository gitHubRepo
}

type gitHubRepo struct {
	Ssh_url string
}

//Ideally, this method should send an a build request
// to some queue which is asynchronously consumed by a separate
// build server.  For now, we execute the build synchronously
// on the same server.
func requestBuild(url string) {
	log.Printf("Triggering build for repo with url: %v\n", url)

	dir := "/var/lib/dodecci/"

	cmd := exec.Command("git", "clone", url, dir + "dodecci/")

	//Right now this causes an error. Need to figure out ssh keys, or do username/password with https github url.
	err := cmd.Run()
	if err != nil {
		log.Panicf("Error running git clone: %v\n", err)
	}
}
