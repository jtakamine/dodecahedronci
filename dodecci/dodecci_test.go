package main

import (
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	cmd := exec.Command("fig", "up")
	err := cmd.Run()

	if err != nil {
		t.Error(err)
	}

}
