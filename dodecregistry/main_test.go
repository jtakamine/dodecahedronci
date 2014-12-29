package main

import (
	"github.com/jtakamine/dodecahedronci/dodecregistry/api"
	"github.com/jtakamine/dodecahedronci/testutils"
	"log"
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

	testPostAndGetBuild("myapp", "1.2.0.345", "app:\n  image: scratch", "", "http://localhost:8000", t)
}

func TestMainShort(t *testing.T) {
	testutils.GoInstall(t)
	process := dodecregistry(t)
	defer process.Kill()

	testPostAndGetBuild("myapp", "1.2.0.345", "app:\n  image: scratch", "asdf", "http://localhost:8000", t)
}

func dodecregistry(t *testing.T) (p *os.Process) {
	var cmd *exec.Cmd
	go func() {
		cmd = testutils.CreateCmd("dodecregistry", "--port", "8000")

		err := cmd.Run()
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(500 * time.Millisecond)

	p = cmd.Process
	return p
}

func testPostAndGetBuild(app string, version string, artifact string, dockerRegistryUrl string, targetUrl string, t *testing.T) {
	var err error

	build_post := api.Build{Artifact: artifact, DockerRegistryUrl: dockerRegistryUrl}
	err = api.PostBuild(app, version, build_post, targetUrl)
	if err != nil {
		t.Error(err)
	}

	build_get, err := api.GetBuild(app, version, targetUrl)
	if err != nil {
		t.Error(err)
	}

	if build_post != build_get {
		log.Printf("Failed: The retrieved build does not match the build that was posted.\n")
		t.Fail()
	}
}
