package main

import (
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

	payload := "{\"repository\":{\"id\":1234567,\"ssh_url\":\"git@github.com:jtakamine/dodecahedronci.git\", \"clone_url\":\"https://github.com/progrium/logspout.git\"}}"
	cmd = exec.Command("curl", "-H", "Content-Type: application/json", "-d", payload, "http://localhost:8000")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}
}
