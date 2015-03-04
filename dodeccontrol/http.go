package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func httpListen(port int) {
	r := mux.NewRouter()

	r.HandleFunc("/info", handleGetInfo).Methods("GET")

	r.HandleFunc("/builds", handleGetBuilds).Methods("GET")
	r.HandleFunc("/builds", handlePostBuild).Methods("POST")
	r.HandleFunc("/builds/{id}", handleGetBuild).Methods("GET")

	r.HandleFunc("/github/builds", handlePostGitHubBuild).Methods("POST")

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
	appName := r.URL.Query().Get("appname")

	bs, err := rpcGetBuilds(appName)
	if err != nil {
		panic("Error getting builds: " + err.Error())
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(bs)
	if err != nil {
		panic("Error encoding response: " + err.Error())
	}
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

	app, err := rpcGetApplication(appName)
	if app.Name == "" {
		err = rpcAddApplication(appName, description)
		if err != nil {
			panic("Error adding application: " + err.Error())
		}
	}

	buildUUID, err := rpcExecuteBuild(repoUrl, appName)
	if err != nil {
		panic("Error executing RPC Build Execute: " + err.Error())
	}

	go func() {
		for {
			b, err := rpcGetBuild(buildUUID)
			if err != nil {
				panic(err)
			}

			zeroTime := time.Time{}
			if !b.Completed.Equal(zeroTime) {
				break
			}

			time.Sleep(time.Second * 5)
		}

		deployUUID, err := rpcExecuteDeploy(buildUUID)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "{buildUUID: \"%s\"; deployUUID: \"%s\"}", buildUUID, deployUUID)
	}()
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
