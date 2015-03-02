package main

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"time"
)

var rpcExecuteBuild = func(repoUrl string, appName string) (err error) {
	conn, err := net.DialTimeout("tcp", buildAddr, time.Second)
	if err != nil {
		return err
	}
	c := jsonrpc.NewClient(conn)

	args := struct {
		RepoUrl string
		AppName string
	}{
		RepoUrl: repoUrl,
		AppName: appName,
	}

	var uuid string
	err = c.Call("Builder.Execute", args, &uuid)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
