package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"time"
)

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

	a := Application{
		Name:        name,
		Description: description,
	}

	var success bool
	err = c.Call("AppRepo.Save", a, &success)
	if err != nil {
		return err
	}

	return nil
}

var rpcGetApplication = func(name string) (a Application, err error) {
	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return Application{}, errors.New("Missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return Application{}, err
	}

	c := jsonrpc.NewClient(conn)

	err = c.Call("AppRepo.Get", name, &a)
	if err != nil {
		return Application{}, err
	}

	return a, nil
}

var rpcExecuteBuild = func(repoUrl string, appName string) (uuid string, err error) {
	conn, err := net.DialTimeout("tcp", buildAddr, time.Second)
	if err != nil {
		return "", err
	}

	c := jsonrpc.NewClient(conn)

	args := struct {
		RepoUrl string
		AppName string
	}{
		RepoUrl: repoUrl,
		AppName: appName,
	}

	err = c.Call("Builder.Execute", args, &uuid)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

var rpcGetBuild = func(uuid string) (b BuildDetails, err error) {
	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return BuildDetails{}, errors.New("Missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return BuildDetails{}, err
	}

	c := jsonrpc.NewClient(conn)

	err = c.Call("BuildRepo.Get", uuid, &b)
	if err != nil {
		return BuildDetails{}, err
	}

	return b, nil
}

var rpcGetBuilds = func(appName string) (bs []Build, err error) {
	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return nil, errors.New("missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return nil, fmt.Errorf("net dial: %s", err.Error())
	}

	c := jsonrpc.NewClient(conn)

	bs = make([]Build, 1)
	err = c.Call("BuildRepo.GetAll", appName, &bs)
	if err != nil {
		return nil, fmt.Errorf("rpc call: %s", err.Error())
	}

	return bs, nil
}

var rpcExecuteDeploy = func(buildUUID string) (deployUUID string, err error) {
	conn, err := net.DialTimeout("tcp", deployAddr, time.Second)
	if err != nil {
		return "", err
	}

	c := jsonrpc.NewClient(conn)

	err = c.Call("Deployer.Execute", buildUUID, &deployUUID)
	if err != nil {
		return "", err
	}

	return deployUUID, nil
}

var rpcGetDeploy = func(uuid string) (d DeployDetails, err error) {
	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return DeployDetails{}, errors.New("Missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return DeployDetails{}, err
	}

	c := jsonrpc.NewClient(conn)

	err = c.Call("DeployRepo.Get", uuid, &d)
	if err != nil {
		return DeployDetails{}, err
	}

	return d, nil
}

var rpcGetDeploys = func(appName string) (ds []Deploy, err error) {
	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return nil, errors.New("Missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return nil, err
	}

	c := jsonrpc.NewClient(conn)

	err = c.Call("DeployRepo.GetAll", appName, &ds)
	if err != nil {
		return nil, err
	}

	return ds, nil
}

var rpcSaveLog = func(l Log) (err error) {
	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return errors.New("Missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return err
	}

	c := jsonrpc.NewClient(conn)

	var success bool
	err = c.Call("LogRepo.Save", l, &success)
	if err != nil {
		return err
	}

	return nil
}

var rpcGetLogs = func(taskUUID string, severity int, startID int64) (ls []Log, err error) {
	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return nil, errors.New("Missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return nil, err
	}

	c := jsonrpc.NewClient(conn)

	args := struct {
		TaskUUID string
		Severity int
		StartID  int64
	}{
		TaskUUID: taskUUID,
		Severity: severity,
		StartID:  startID,
	}

	err = c.Call("LogRepo.GetAll", args, &ls)
	if err != nil {
		return nil, err
	}

	return ls, nil
}
