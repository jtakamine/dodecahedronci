package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

/*
	GET /info
	POST /build
	POST /build/github
	GET /build/(id)
	GET /build/(id)/logs
	POST /release
	GET /release/(id)
	POST /deploy
	GET /deploy/(id)
	GET /deploy/(id)/logs
*/
func httpListen(port int) {
	r := mux.NewRouter()

	r.HandleFunc("/info", handleGetInfo).Methods("GET")

	r.HandleFunc("/builds", handleGetBuilds).Methods("GET")
	r.HandleFunc("/builds", handlePostBuild).Methods("POST")
	r.HandleFunc("/builds/{id}", handleGetBuild).Methods("GET")
	r.HandleFunc("/github/builds", handlePostGitHubBuild).Methods("POST")

	r.HandleFunc("/releases", handleGetReleases).Methods("GET")
	r.HandleFunc("/releases", handlePostRelease).Methods("POST")
	r.HandleFunc("/releases/{id}", handleGetRelease).Methods("GET")

	r.HandleFunc("/deploys", handleGetDeploys).Methods("GET")
	r.HandleFunc("/deploys", handlePostDeploy).Methods("POST")
	r.HandleFunc("/deploys/{id}", handleGetDeploy).Methods("GET")

	r.HandleFunc("/{entity}/{id}/logs", handleGetLogs).Methods("GET")
	r.HandleFunc("/{entity}/{id}/logs/stream", handleStreamLogs).Methods("UPGRADE")

	http.Handle("/", r)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func handleGetInfo(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}

func handleGetBuilds(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}
func handlePostBuild(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}
func handleGetBuild(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}

func handlePostGitHubBuild(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic("Error reading request body: " + err.Error())
	}

	repoUrl, appName, description, err := parseGitHubRequest(data)
	if err != nil {
		panic("Error parsing GitHub request: " + err.Error())
	}

	fmt.Printf("appname = %s\n", appName)

	err = rpcAddApplication(appName, description)
	if err != nil {
		panic("Error adding application: " + err.Error())
	}

	err = rpcExecuteBuild(repoUrl, appName)
	if err != nil {
		panic("Error executing RPC Build Execute: " + err.Error())
	}
}

func handleGetReleases(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}
func handlePostRelease(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}
func handleGetRelease(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}

func handleGetDeploys(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}
func handlePostDeploy(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}
func handleGetDeploy(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}

func handleGetLogs(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}
func handleStreamLogs(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}
