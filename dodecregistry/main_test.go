package main

import (
	"github.com/jtakamine/dodecahedronci/dodecregistry/api"
	"github.com/jtakamine/dodecahedronci/utils/testutil"
	"log"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testutil.FigBuild(t)
	testutil.FigUp(t)
	defer testutil.FigKillAndRm(t)

	testPostAndGetBuild("myapp", "1.2.0.345", "app:\n  image: scratch", "", "http://localhost:8001", t)
}

func TestMainShort(t *testing.T) {
	parseArgs = func() (port int) {
		return 8001
	}

	go main()
	time.Sleep(500 * time.Millisecond)

	testPostAndGetBuild("myapp", "1.2.0.345", "app:\n  image: scratch", "asdf", "http://localhost:8001", t)
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
