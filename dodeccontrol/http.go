package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
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
	vals := r.URL.Query()
	deployStr := vals.Get("deploy")

	var err error
	var deploy bool
	if deployStr != "" {
		deploy, err = strconv.ParseBool(deployStr)
		if err != nil {
			panic("Error parsing \"deploy\" query argument: " + err.Error())
		}
	}

	reqBody := struct {
		RepoUrl string
		AppName string
	}{}

	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&reqBody)
	if err != nil {
		panic("Error decoding request body: " + err.Error())
	}

	uuid, err := execBuild(reqBody.RepoUrl, reqBody.AppName, deploy)
	if err != nil {
		panic("Error executing RPC Build: " + err.Error())
	}

	fmt.Fprintf(w, "{\"UUID\": \"%s\"}", uuid)
}

func handleGetBuild(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	b, err := rpcGetBuild(vars["id"])
	if err != nil {
		panic("Error getting build: " + err.Error())
	}

	enc := json.NewEncoder(w)

	if b.UUID != "" {
		err = enc.Encode(b)
	} else {
		err = enc.Encode(struct{}{})
	}
	if err != nil {
		panic("Error encoding response: " + err.Error())
	}
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

	uuid, err := execBuild(repoUrl, appName, true)
	if err != nil {
		panic("Error executing RPC Build: " + err.Error())
	}

	fmt.Fprintf(w, "{\"UUID:\" \"%s\"}", uuid)
}

func handleGetDeploys(w http.ResponseWriter, r *http.Request) {
	appName := r.URL.Query().Get("appname")

	ds, err := rpcGetDeploys(appName)
	if err != nil {
		panic("Error getting deploys: " + err.Error())
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(ds)
	if err != nil {
		panic("Error encoding response: " + err.Error())
	}
}
func handlePostDeploy(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}
func handleGetDeploy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	d, err := rpcGetDeploy(vars["id"])
	if err != nil {
		panic("Error getting build: " + err.Error())
	}

	enc := json.NewEncoder(w)

	if d.UUID != "" {
		err = enc.Encode(d)
	} else {
		err = enc.Encode(struct{}{})
	}
	if err != nil {
		panic("Error encoding response: " + err.Error())
	}
}

func handleGetLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sevStr := r.URL.Query().Get("severity")
	startIDStr := r.URL.Query().Get("startid")

	var sev int
	var startID int64
	var err error
	if sevStr != "" {
		sev, err = strconv.Atoi(sevStr)
		if err != nil {
			panic("Error parsing severity from query string: " + err.Error())
		}
	}
	if startIDStr != "" {
		startID32, err := strconv.Atoi(startIDStr)
		if err != nil {
			panic("Error parsing startID from query string: " + err.Error())
		}
		startID = int64(startID32)
	}

	ls, err := rpcGetLogs(vars["id"], sev, startID)
	if err != nil {
		panic("Error getting logs: " + err.Error())
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(ls)
	if err != nil {
		panic("Error encoding response: " + err.Error())
	}
}

func handleStreamLogs(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented!")
}
