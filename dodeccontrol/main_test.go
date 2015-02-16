package main

import (
	"bytes"
	"fmt"
	"github.com/jtakamine/dodecahedronci/testutil"
	"net"
	"net/http"
	"net/rpc/jsonrpc"
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

	time.Sleep(time.Second * 2)

	testWebhook(t, "https://github.com/progrium/logspout.git", "http://localhost:8080/github/builds")

	time.Sleep(time.Second * 15)
}

func TestMainShort(t *testing.T) {
	if !testing.Short() {
		t.Skip()
	}

	parseArgs = func() (port int, rpcPort int) {
		return 8000, 9000
	}

	rpcExecuteBuild = func(repoUrl string) (err error) {
		fmt.Printf("***Mocked: RPC Execute Build. Repo Url: %s\n", repoUrl)
		return nil
	}

	go main()
	time.Sleep(500 * time.Millisecond)

	testRPCExecute(t, "my message", "localhost:9000")
	testWebhook(t, "https://github.com/jtakamine/dodecahedronci.git", "http://localhost:8000/github/builds")
	testWebhook(t, "https://github.com/Leland-Takamine/testtarget.git", "http://localhost:8000/github/builds")
}

func testRPCExecute(t *testing.T, msg string, addr string) {
	tStamp := "2015-02-14T15:04:05+07:00"
	logType := 0
	taskID := "mytaskID"
	src := "mysource"

	log := fmt.Sprintf("[%s][%s][%d] %s\t|%s", src, taskID, logType, tStamp, msg)

	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		t.Error(err)
		return
	}
	c := jsonrpc.NewClient(conn)

	var success bool

	err = c.Call("Log.Write", log, &success)
	if err != nil {
		t.Error(err)
		return
	}

	ls, ok := logs[src]
	if !ok {
		t.Errorf("Could not find log lists for source \"%s\"", src)
		return
	}

	l, ok := ls[taskID]
	if !ok {
		t.Errorf("Could not find log list for task id \"%s\"", taskID)
		return
	}

	found := false
	for _, e := range l {
		if e.Msg == msg && e.Type == logType {
			found = true
			break
		}
	}
	fmt.Printf("log list: %v\n", logs)
	if !found {
		t.Errorf("Sent message \"%s\", but could not find it list", msg)
		return
	}
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
