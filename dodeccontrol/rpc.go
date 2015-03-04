package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

func rpcListen(port int) (err error) {
	err = rpc.RegisterName("Logger", &RPCLogger{})
	if err != nil {
		return err
	}

	err = rpc.RegisterName("ServiceRegistry", &RPCServiceRegistry{})
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

type RPCLogger struct{}

func (rpcL *RPCLogger) Write(log string, success *bool) (err error) {
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

type RegisterServiceArgs struct {
	Service  string
	Endpoint string
}

type RPCServiceRegistry struct{}

func (rpcS *RPCServiceRegistry) Register(args RegisterServiceArgs, success *bool) (err error) {
	switch args.Service {
	case "build":
		buildAddr = args.Endpoint
	case "deploy":
		deployAddr = args.Endpoint
	default:
		return errors.New("Attempted to register unrecognized service: \"" + args.Service + "\"")
	}

	*success = true
	return nil
}
