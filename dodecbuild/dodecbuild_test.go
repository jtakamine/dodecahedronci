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
	if testing.Short() {
		t.Skip()
	}

	figBuild(t)
	figUp(t)
	defer figKillAndRm(t)

	time.Sleep(5 * time.Second)

	testWebhook(t, "https://github.com/progrium/logspout.git", "http://localhost:8000")
}

func figKillAndRm(t *testing.T) {
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
}

func figBuild(t *testing.T) {
	cmd := createCmd("fig", "build")
	cmd.Dir = ".."

	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}
}

func figUp(t *testing.T) {
	cmd := createCmd("fig", "up", "-d")
	cmd.Dir = ".."

	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}
}

func TestMainShort(t *testing.T) {
	var err error
	var cmd *exec.Cmd

	cmd = createCmd("go", "install")

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	var dodecbuildCmd *exec.Cmd
	go func() {
		dodecbuildCmd = createCmd("dodecbuild", "--port", "8000")

		err = dodecbuildCmd.Run()
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(500 * time.Millisecond)

	defer dodecbuildCmd.Process.Kill()

	testWebhook(t, "https://github.com/jtakamine/dodecahedronci.git", "http://localhost:8000")
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

func createCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd
}
