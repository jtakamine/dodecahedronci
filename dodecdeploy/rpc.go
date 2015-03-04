package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

func rpcListen(port int) (err error) {
	err = rpc.RegisterName("Deployer", &RPCDeployer{})
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

	return nil
}

type RPCDeployer struct{}

func (*RPCDeployer) Execute(buildUUID string, uuid *string) (err error) {
	*uuid = generateRandID(8)

	a, err := rpcGetBuildArtifact(buildUUID)
	if err != nil {
		return err
	}

	err = deploy(a)
	if err != nil {
		return err
	}

	return nil
}
