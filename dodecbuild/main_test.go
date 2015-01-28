package main

import (
	"bytes"
	"fmt"
	dodecpubsub_API "github.com/jtakamine/dodecahedronci/dodecpubsub/api"
	"github.com/jtakamine/dodecahedronci/testutil"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testutil.FigBuild(t)
	testutil.FigUp(t)
	defer testutil.FigKillAndRm(t)

	testSubscribeLogs(t, "localhost:8000")
	testWebhook(t, "https://github.com/progrium/logspout.git", "http://localhost:8002")
}

func TestMainShort(t *testing.T) {
	if !testing.Short() {
		t.Skip()
	}

	parseArgs = func() (port int) {
		return 8002
	}

	saveBuild = func(app string, version string, fFile figFile, dockerRegistryUrl string) (err error) {
		fmt.Printf("Saved build. Fig file:%v\n", fFile.Config)
		return nil
	}

	log = func(msg string, lType logType) (err error) {
		fmt.Printf("Published log (level %v): %s\n", lType, msg)
		return nil
	}

	go main()
	time.Sleep(500 * time.Millisecond)

	testWebhook(t, "https://github.com/jtakamine/dodecahedronci.git", "http://localhost:8002")
	testWebhook(t, "https://github.com/Leland-Takamine/testtarget.git", "http://localhost:8002")
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

func testSubscribeLogs(t *testing.T, address string) {
	subChan, err := dodecpubsub_API.Subscribe(strconv.Itoa(int(verboseLogType)), address)
	if err != nil {
		t.Error(err)
	}

	go func() {
		for msg := range subChan {
			fmt.Println("Received message: " + msg)
		}
	}()
}
