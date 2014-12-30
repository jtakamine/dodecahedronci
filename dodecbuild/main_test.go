package main

import (
	"bytes"
	"github.com/jtakamine/dodecahedronci/testutils"
	"net/http"
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

	testWebhook(t, "https://github.com/progrium/logspout.git", "http://localhost:8001")
}

func TestMainShort(t *testing.T) {
	parseArgs = func() (port int) {
		return 8001
	}

	postBuildToDodecRegistry = func(app string, version string, fFile figFile, dockerRegistryUrl string) (err error) {
		return nil
	}

	go main()
	time.Sleep(500 * time.Millisecond)

	testWebhook(t, "https://github.com/jtakamine/dodecahedronci.git", "http://localhost:8001")
	testWebhook(t, "https://github.com/Leland-Takamine/testtarget.git", "http://localhost:8001")
}

func testWebhook(t *testing.T, cloneUrl string, targetUrl string) {
	var err error
	var req *http.Request
	var resp *http.Response

	payload := "{\"repository\":{\"id\":1234567,\"ssh_url\":\"git@github.com:jtakamine/dodecahedronci.git\", \"clone_url\":\"" + cloneUrl + "\"}}"
	req, err = http.NewRequest("POST", targetUrl, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}
	resp.Body.Close()
}
