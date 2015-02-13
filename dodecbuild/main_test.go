package main

import (
	"bytes"
	"fmt"
	"github.com/jtakamine/dodecahedronci/testutil"
	"net"
	"net/http"
	"net/rpc"
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

	testRPCExecute(t, "https://github.com/progrium/logspout.git", "localhost:8002")
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

	go main()
	time.Sleep(500 * time.Millisecond)

	testRPCExecute(t, "https://github.com/jtakamine/dodecahedronci.git", "localhost:8002")
	testRPCExecute(t, "https://github.com/Leland-Takamine/testtarget.git", "localhost:8002")
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
		return
	}
	resp.Body.Close()
}

func testRPCExecute(t *testing.T, repoUrl string, addr string) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		t.Error(err)
		return
	}
	c := rpc.NewClient(conn)

	var buildID string
	err = c.Call("Build.Execute", repoUrl, &buildID)
	if err != nil {
		t.Error(err)
		return
	}
}
