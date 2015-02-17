package main

import (
	"fmt"
	"github.com/jtakamine/dodecahedronci/logutil"
	"net"
	"net/rpc/jsonrpc"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	parseArgs = func() (port int) {
		return 8002
	}

	pushDockerImage = func(tag string, writer *logutil.Writer) (err error) {
		fmt.Printf("**Mocked: push docker image %s to registry.\n", tag)
		return nil
	}

	saveBuild = func(app string, version string, fFile figFile) (err error) {
		fmt.Printf("**Mocked: Saved build. Fig file:%v\n", fFile.Config)
		return nil
	}

	rpcRegisterService = func() (err error) {
		fmt.Println("**Mocked: Registered service.")
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
	c := jsonrpc.NewClient(conn)

	args := &BuildArgs{
		RepoUrl: repoUrl,
	}

	var buildID string
	err = c.Call("Build.Execute", args, &buildID)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("buildID=%s\n", buildID)
}
