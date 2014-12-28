package testutils

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func FigBuild(t *testing.T) {
	cmd := CreateCmd("fig", "build")
	cmd.Dir = ".."

	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}
}

func FigUp(t *testing.T) {
	cmd := CreateCmd("fig", "up", "-d")
	cmd.Dir = ".."

	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
}

func FigKillAndRm(t *testing.T) {
	var err error
	var cmd *exec.Cmd

	cmd = CreateCmd("fig", "kill")
	cmd.Dir = ".."

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	cmd = CreateCmd("fig", "rm", "--force")
	cmd.Dir = ".."

	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}
}

func CreateCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd
}
