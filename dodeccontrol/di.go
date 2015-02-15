package main

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

var rpcExecuteBuild = func(repoUrl string) (err error) {
	fmt.Println("here3")
	conn, err := net.DialTimeout("tcp", "dodecbuild:9000", time.Second)
	if err != nil {
		return err
	}
	c := rpc.NewClient(conn)

	args := &struct {
		RepoUrl string
	}{
		RepoUrl: repoUrl,
	}

	var buildID string
	err = c.Call("Build.Execute", args, &buildID)
	if err != nil {
		return err
	}

	return nil
}
