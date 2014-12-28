package main

import (
	"encoding/json"
	"github.com/jtakamine/dodecahedronci/testutils"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testutils.FigBuild(t)
	testutils.FigUp(t)
	defer testutils.FigKillAndRm(t)

	testPostBuild("myapp", "1.2.0.345", "app:\n  image: scratch", "", "http://localhost:8000", t)
}

func testPostBuild(app string, version string, artifact string, dockerRegistryUrl string, targetUrl string, t *testing.T) {
	var err error

	reqObj := postBuildReq{App: app, Version: version, Artifact: artifact, DockerRegistryUrl: dockerRegistryUrl}

	reqData, err := json.Marshal(reqObj)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", targetUrl, strings.NewReader(string(reqData)))
	if err != nil {
		t.Error(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	log.Printf("%v", string(respBody))
}
