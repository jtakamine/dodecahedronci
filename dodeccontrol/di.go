package main

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"time"
)

var rpcExecuteBuild = func(repoUrl string) (err error) {
	conn, err := net.DialTimeout("tcp", buildAddr, time.Second)
	if err != nil {
		return err
	}
	c := jsonrpc.NewClient(conn)

	args := &struct {
		RepoUrl string
	}{
		RepoUrl: repoUrl,
	}

	var buildID string
	err = c.Call("Build.Execute", args, &buildID)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
