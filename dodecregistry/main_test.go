package main

import (
	"encoding/json"
	"github.com/jtakamine/dodecahedronci/testutils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testutils.FigBuild(t)
	testutils.FigUp(t)
	defer testutils.FigKillAndRm(t)

	testPostAndGetBuild("myapp", "1.2.0.345", "app:\n  image: scratch", "", "http://localhost:8000", t)
}

func TestMainShort(t *testing.T) {
	testutils.GoInstall(t)
	process := dodecregistry(t)
	defer process.Kill()

	testPostAndGetBuild("myapp", "1.2.0.345", "app:\n  image: scratch", "asdf", "http://localhost:8000", t)
}

func dodecregistry(t *testing.T) (p *os.Process) {
	var cmd *exec.Cmd
	go func() {
		cmd = testutils.CreateCmd("dodecregistry", "--port", "8000")

		err := cmd.Run()
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(500 * time.Millisecond)

	p = cmd.Process
	return p
}

func testPostAndGetBuild(app string, version string, artifact string, dockerRegistryUrl string, targetUrl string, t *testing.T) {
	testPostBuild(app, version, artifact, dockerRegistryUrl, targetUrl, t)
	build := testGetBuild(app, version, targetUrl, t)

	if build.DockerRegistryUrl != dockerRegistryUrl || build.Artifact != artifact {
		log.Printf("Failed: The retrieved build does not match the build that was posted.\n")
		t.Fail()
	}
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

func testGetBuild(app string, version string, targetUrl string, t *testing.T) (build *dodecBuild) {
	vals := url.Values{"app": {app}, "version": {version}}
	url := targetUrl + "?" + vals.Encode()

	resp, err := http.Get(url)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	log.Printf("get build response: %v\n", string(respBody))

	build = &dodecBuild{}
	err = json.Unmarshal(respBody, build)
	if err != nil {
		t.Error(err)
	}

	return build
}
