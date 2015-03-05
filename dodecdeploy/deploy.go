package main

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/jtakamine/dodecahedronci/logutil"
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

func deploy(buildArtifact string, writer *logutil.Writer) (err error) {
	w := writer.WriteType

	w("Creating temp directory to write fig file...", logutil.Info)
	tempDir, err := ioutil.TempDir("", "dodec")
	if err != nil {
		return err
	}
	w("Created temp directory "+tempDir+" for fig file.", logutil.Info)

	tempFile := tempDir + "/fig.yml"

	w("Writing fig file to temp directory...", logutil.Info)
	err = ioutil.WriteFile(tempFile, []byte(buildArtifact), 0644)
	if err != nil {
		return err
	}
	w("Fig file written to temp directory.", logutil.Info)

	err = os.Setenv("FIG_FILE", tempFile)
	if err != nil {
		return err
	}

	w("Running fig up...", logutil.Info)
	cmd := exec.Command("fig", "up", "-d", "--allow-insecure-ssl")
	cmd.Stderr = writer.CreateChild().CreateWriter(logutil.Error)
	cmd.Stdout = writer.CreateChild().CreateWriter(logutil.Verbose)

	err = cmd.Run()
	if err != nil {
		return err
	}
	w("Fig up complete.", logutil.Info)

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
