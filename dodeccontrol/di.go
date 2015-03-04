package main

import (
	"errors"
	"net"
	"net/rpc/jsonrpc"
	"os"
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
		return err
	}

	return nil
}

var rpcAddApplication = func(name string, description string) (err error) {
	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return errors.New("Missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return err
	}

	c := jsonrpc.NewClient(conn)

	args := struct {
		Name        string
		Description string
	}{
		Name:        name,
		Description: description,
	}

	var success bool
	err = c.Call("AppRepo.Save", args, &success)
	if err != nil {
		return err
	}

	return nil
}
