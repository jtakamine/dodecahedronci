package main

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	var err error
	var cmd *exec.Cmd

	cmd = exec.Command("fig", "kill")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = ".."

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	cmd = exec.Command("fig", "rm", "--force")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = ".."

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	cmd = exec.Command("fig", "build")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = ".."

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	cmd = exec.Command("fig", "up", "-d")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = ".."

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

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
