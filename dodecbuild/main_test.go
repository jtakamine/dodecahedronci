package main

import (
	"fmt"
	"github.com/jtakamine/dodecahedronci/testutil"
	"net"
	"net/rpc"
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

	testRPCExecute(t, "https://github.com/progrium/logspout.git", "jtakamine", "localhost:8002")
}

func TestMainShort(t *testing.T) {
	if !testing.Short() {
		t.Skip()
	}

	parseArgs = func() (port int) {
		return 8002
	}

	saveBuild = func(app string, version string, fFile figFile, dockerRegistryUrl string) (err error) {
		fmt.Printf("Saved build. Fig file:%v\n", fFile.Config)
		return nil
	}

	go main()
	time.Sleep(500 * time.Millisecond)

	testRPCExecute(t, "https://github.com/jtakamine/dodecahedronci.git", "jtakamine", "localhost:8002")
	testRPCExecute(t, "https://github.com/Leland-Takamine/testtarget.git", "jtakamine", "localhost:8002")
}

func testRPCExecute(t *testing.T, repoUrl string, dockerUser string, addr string) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		t.Error(err)
		return
	}
	c := rpc.NewClient(conn)

	args := &BuildArgs{
		RepoUrl:    repoUrl,
		DockerUser: dockerUser,
	}

	var buildID string
	err = c.Call("Build.Execute", args, &buildID)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("buildID=%s\n", buildID)
}
