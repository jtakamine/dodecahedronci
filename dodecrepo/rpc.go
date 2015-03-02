package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	// "encoding/json"
	_ "github.com/lib/pq"
)

func rpcListen(port int) (err error) {
	err = rpc.RegisterName("BuildRepo", &RPCBuildRepo{})
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

type RPCBuildRepo struct{}

func (rpcB *RPCBuildRepo) Save(b Build, success *bool) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	err = saveBuild(b, c)
	if err != nil {
		return err
	}

	*success = true
	return nil
}

func (rpcB *RPCBuildRepo) Get(uuid string, b *Build) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	*b, err = getBuild(uuid, c)
	if err != nil {
		return err
	}

	return nil
}
