package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func generateRandID(length int) string {
	id := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		panic(err)
	}
	return hex.EncodeToString(id)
}

func deploy(buildArtifact string) (err error) {
	tempDir, err := ioutil.TempDir("", "dodec")
	if err != nil {
		return err
	}

	tempFile := tempDir + "/fig.yml"

	err = ioutil.WriteFile(tempFile, []byte(buildArtifact), 0644)
	if err != nil {
		return err
	}

	os.Setenv("FIG_FILE", tempFile)

	cmd := exec.Command("fig", "up", "-d", "--allow-insecure-ssl")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 2)

	return nil
}

func killDeployments() (err error) {
	cmd := exec.Command("fig", "kill")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
