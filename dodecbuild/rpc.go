package main

import (
	"github.com/jtakamine/dodecahedronci/logutil"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

type BuildArgs struct {
	RepoUrl string
}

type RPCBuild struct{}

func (b *RPCBuild) Execute(args BuildArgs, buildID *string) (err error) {
	*buildID = generateRandID(8)
	writer := logutil.NewWriter("build", *buildID)

	repoDir, err := cloneOrUpdateGitRepo(args.RepoUrl, writer)
	if err != nil {
		return err
	}

	err = build(repoDir, "myAwesomeApp", writer)
	if err != nil {
		return err
	}

	return nil
}

func rpcListen(port int) (err error) {
	err = rpc.RegisterName("Build", &RPCBuild{})
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
