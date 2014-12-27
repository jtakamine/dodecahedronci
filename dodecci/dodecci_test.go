package main

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	var err error
	var cmd *exec.Cmd

	cmd = createCmd("fig", "kill")
	cmd.Dir = ".."

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	cmd = createCmd("fig", "rm", "--force")
	cmd.Dir = ".."

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	cmd = createCmd("fig", "build")
	cmd.Dir = ".."

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	cmd = createCmd("fig", "up", "-d")
	cmd.Dir = ".."

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)

	testWebhook(t)
}

func testWebhook(t *testing.T) {
	var err error
	var req *http.Request
	var resp *http.Response

	payload := "{\"repository\":{\"id\":1234567,\"ssh_url\":\"git@github.com:jtakamine/dodecahedronci.git\", \"clone_url\":\"https://github.com/progrium/logspout.git\"}}"
	req, err = http.NewRequest("POST", "http://localhost:8000", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}
	resp.Body.Close()
}

func createCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd
}
