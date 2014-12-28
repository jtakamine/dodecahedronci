package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type postBuildReq struct {
	App               string
	Version           string
	Artifact          string
	DockerRegistryUrl string
}

func ListenAndServe(addr string) (err error) {
	http.HandleFunc("/", httpHandle)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}

	return nil
}

func httpHandle(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		err := r.ParseForm()
		if err != nil {
			log.Panicf("Error parsing request form: %v\n", err)
		}

		app := r.Form.Get("app")
		version := r.Form.Get("version")

		build, err := getBuild(app, version)
		if err != nil {
			log.Panicf("Error getting build: %v\n", err)
		}

		data, err := json.Marshal(build)
		if err != nil {
			log.Panicf("Error serializing build: %v\n", err)
		}

		_, err = w.Write(data)
		if err != nil {
			log.Panicf("Error writing build response: %v\n", err)
		}
	case r.Method == "POST":
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Panicf("Error reading request body: %v\n", err)
		}

		req := &postBuildReq{}
		err = json.Unmarshal(data, req)
		if err != nil {
			log.Panicf("Error parsing json request: %v\n", err)
		}

		err = addBuild(req.App, req.Version, dodecBuild{Artifact: req.Artifact, DockerRegistryUrl: req.DockerRegistryUrl})
		if err != nil {
			log.Panicf("Error adding package: %v\n", err)
		}

		fmt.Fprint(w, "add successful\n")
	}
}
