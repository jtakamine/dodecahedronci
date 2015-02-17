package main

import (
	"github.com/jtakamine/dodecahedronci/logutil"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type SaveBuildArgs struct {
	UUID     string
	AppName  string
	Version  string
	Artifact string
}

type RPCBuild struct{}

func (rpcB *RPCBuild) Save(args SaveBuildArgs, success *bool) (err error) {
	b := Build{
		UUID:     args.UUID,
		AppName:  args.AppName,
		Version:  args.Version,
		Artifact: args.Artifact,
	}

	err := addBuild(b)
	if err != nil {
		return err
	}

	*sucess = true
	return nil
}

func (rpcB *RPCBuild) Get(uuid string, build *Build) (err error) {
	*build, err = getBuild(uuid)
	if err != nil {
		return err
	}
	return nil
}
