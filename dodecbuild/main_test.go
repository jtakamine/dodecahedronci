package main

import (
	"bytes"
	"github.com/jtakamine/dodecahedronci/testutils"
	"net/http"
	"os"
	"os/exec"
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
	testutils.GoInstall(t)
	process := dodecbuild(t)
	defer process.Kill()

	testWebhook(t, "https://github.com/jtakamine/dodecahedronci.git", "http://localhost:8001")
	testWebhook(t, "https://github.com/Leland-Takamine/testtarget.git", "http://localhost:8001")
}

func dodecbuild(t *testing.T) (p *os.Process) {
	var cmd *exec.Cmd
	go func() {
		cmd = testutils.CreateCmd("dodecbuild", "--port", "8001")

		err := cmd.Run()
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(500 * time.Millisecond)

	p = cmd.Process
	return p
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
