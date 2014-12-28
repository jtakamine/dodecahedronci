package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type deployReq struct {
	string App
	string Version
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
	var err error

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicf("Error reading request body: %v\n", err)
	}

	req := &deployReq{}
	err = json.Unmarshal(data, req)
	if err != nil {
		log.Panicf("Error parsing json request: %v\n", err)
	}

	fmt.Fprint(w, "build successful\n")
}
