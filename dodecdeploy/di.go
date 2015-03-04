package main

import (
	"errors"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"strings"
	"time"
)

var rpcRegisterService = func() (err error) {
	controlAddr := os.Getenv("DODEC_CONTROLADDR")
	if controlAddr == "" {
		return errors.New("Missing environment variable: DODEC_CONTROLADDR")
	}

	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		return err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}
	if len(addrs) == 0 {
		return errors.New("Could not find IP address associated with eth0 device")
	}

	ip := strings.Split(addrs[0].String(), "/")[0]

	conn, err := net.DialTimeout("tcp", controlAddr, time.Second)
	if err != nil {
		return err
	}
	c := jsonrpc.NewClient(conn)

	args := &struct {
		Service  string
		Endpoint string
	}{
		Service:  "deploy",
		Endpoint: ip + ":9000",
	}

	var success bool
	err = c.Call("ServiceRegistry.Register", args, &success)
	if err != nil {
		return err
	}

	return nil
}

var rpcGetBuildArtifact = func(buildUUID string) (artifact string, err error) {
	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return "", errors.New("Missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return "", err
	}

	c := jsonrpc.NewClient(conn)

	b := struct {
		UUID      string
		AppName   string
		Version   string
		Started   time.Time
		Completed time.Time
		Success   bool
		Artifact  string
	}{}

	err = c.Call("BuildRepo.Get", buildUUID, &b)
	if err != nil {
		return "", err
	}

	return b.Artifact, nil
}

var rpcSaveDeploy = func(buildUUID string, deployUUID string, appName string) (err error) {
	_ = struct {
		deployUUID string
		buildUUID  string
		AppName    string
	}{
		deployUUID: deployUUID,
		buildUUID:  buildUUID,
		AppName:    appName,
	}

	return nil
}
