package handlers

import (
	"net/http"
	"encoding/json"
	"log"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	req := &gitHubReq{}

	err := decoder.Decode(req)

	if err != nil {
		log.Panicf("Could not parse JSON")
	}
	
	log.Println(req)
}

type gitHubReq struct {
	Repository gitHubRepo
}

type gitHubRepo struct {
	Ssh_url string
}
