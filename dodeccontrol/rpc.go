package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
)

type LogEntry struct {
	Type int
	Msg  string
}

type RPCLog struct{}

func (rpcL *RPCLog) Write(log string, success *bool) (err error) {
	msg, logType, taskID, src := parseLog(log)

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

func parseLog(log string) (msg string, logType int, taskID string, src string) {
	msg = log
	logType = 0
	taskID = "default"
	src = "default"
	whitespace := "\t "

	msg = strings.TrimLeft(msg, whitespace)

	if !strings.HasPrefix(msg, "[") {
		return msg, logType, taskID, src
	}

	i := strings.Index(msg, "|")
	if i < 0 {
		return msg, logType, taskID, src
	}

	pre := strings.Split(msg, "|")[0]
	if strings.Count(pre, "[") != 3 || strings.Count(pre, "]") != 3 {
		return msg, logType, taskID, src
	}

	pre = strings.TrimRight(pre, whitespace)
	if !strings.HasSuffix(pre, "]") {
		return msg, logType, taskID, src
	}

	parts := strings.Split(pre, "][")
	if len(parts) != 3 {
		return msg, logType, taskID, src
	}

	logType, err := strconv.Atoi(strings.TrimSuffix(parts[2], "]"))
	if err != nil {
		return msg, logType, taskID, src
	}

	src = strings.TrimPrefix(parts[0], "[")
	taskID = parts[1]

	msg = strings.Join(strings.Split(msg, "|")[1:], "|")

	return msg, logType, taskID, src
}

func rpcListen(port int) (err error) {
	err = rpc.RegisterName("Stdin", &RPCLog{})
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
