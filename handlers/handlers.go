package handlers

import (
	"net/http"
	"encoding/json"
	"log"
	"bytes"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	req := &gitHubReq{}

	err := decoder.Decode(req)

	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body) 
		log.Panicf("Could not parse JSON: %v\n", err)
	} else {
		log.Printf("Parsed JSON: %v\n", req)
	}
}

type gitHubReq struct {
	Repository gitHubRepo
}

type gitHubRepo struct {
	Ssh_url string
}
