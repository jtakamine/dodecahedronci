package main

import (
	"github.com/jtakamine/dodecahedronci/testutil"
	"net"
	"net/rpc/jsonrpc"
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

	testRPCSaveAndGetBuild("1234uuid567", "myapp", "1.0.0.0", "app:\n\timage: scratch", "localhost:8001", t)
}

func TestMainShort(t *testing.T) {
	parseArgs = func() (port int) {
		return 8001
	}

	go main()
	time.Sleep(500 * time.Millisecond)

	testRPCSaveAndGetBuild("1234uuid567", "myapp", "1.0.0.0", "app:\n\timage: scratch", "localhost:8001", t)
}

func testRPCSaveAndGetBuild(uuid string, app string, version string, artifact string, addr string, t *testing.T) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		t.Error(err)
		return
	}
	c := jsonrpc.NewClient(conn)

	var success bool
	b_save := Build{
		UUID:     uuid,
		AppName:  app,
		Version:  version,
		Artifact: artifact,
	}

	err = c.Call("BuildRepo.Save", b_save, &success)
	if err != nil {
		t.Error(err)
		return
	}

	var b_get Build
	err = c.Call("BuildRepo.Get", uuid, &b_get)
	if err != nil {
		t.Error(err)
		return
	}

	if b_get != b_save {
		t.Errorf("Saved build did not equal retrieved build: %s <> %s\n", b_save, b_get)
	}
}

/*
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
}*/
