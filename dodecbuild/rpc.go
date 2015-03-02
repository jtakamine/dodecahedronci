package main

import (
	"github.com/jtakamine/dodecahedronci/logutil"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

type ExecuteBuildArgs struct {
	RepoUrl string
	AppName string
}

type RPCBuilder struct{}

func (b *RPCBuilder) Execute(args ExecuteBuildArgs, uuid *string) (err error) {
	*uuid = generateRandID(8)
	writer := logutil.NewWriter("build", *uuid)

	repoDir, err := cloneOrUpdateGitRepo(args.RepoUrl, writer)
	if err != nil {
		return err
	}

	err = build(repoDir, *uuid, args.AppName, writer)
	if err != nil {
		return err
	}

	return nil
}

func rpcListen(port int) (err error) {
	err = rpc.RegisterName("Build", &RPCBuilder{})
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
