package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

type RPCLog struct{}

func (rpcL *RPCLog) Write(log string, success *bool) (err error) {
	msg, _, logType, taskID, src := parseLog(log)

	ls, ok := logs[src]
	if !ok {
		ls = make(map[string][]LogEntry)
	}

	l, ok := ls[taskID]
	if !ok {
		l = []LogEntry{}
	}

	l = append(l, LogEntry{Type: logType, Msg: msg})
	ls[taskID] = l
	logs[src] = ls

	fmt.Printf("logged: %s\n", msg)

	return nil
}

type RegisterArgs struct {
	Service  string
	Endpoint string
}
type RPCService struct{}

func (rpcS *RPCService) Register(args RegisterArgs, success *bool) (err error) {
	switch args.Service {
	case "build":
		buildAddr = args.Endpoint
	}

	*success = true
	return nil
}

func rpcListen(port int) (err error) {
	err = rpc.RegisterName("Log", &RPCLog{})
	if err != nil {
		return err
	}

	err = rpc.RegisterName("Service", &RPCService{})
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
