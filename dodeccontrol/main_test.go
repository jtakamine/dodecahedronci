package main

import (
	"fmt"
	"github.com/jtakamine/dodecahedronci/testutil"
	"net"
	"net/rpc/jsonrpc"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testutil.FigBuild(t)
	testutil.FigUp(t)
	defer testutil.FigKillAndRm(t)
}

func TestMainShort(t *testing.T) {
	if !testing.Short() {
		t.Skip()
	}

	parseArgs = func() (port int, rpcPort int) {
		return 8000, 9000
	}

	go main()
	time.Sleep(500 * time.Millisecond)

	testRPCExecute(t, "my message", "localhost:9000")
}

func testRPCExecute(t *testing.T, msg string, addr string) {
	tStamp := "2015-02-14T15:04:05+07:00"
	logType := 0
	taskID := "mytaskID"
	src := "mysource"

	log := fmt.Sprintf("[%s][%s][%d] %s\t|%s", src, taskID, logType, tStamp, msg)

	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		t.Error(err)
		return
	}
	c := jsonrpc.NewClient(conn)

	var success bool

	err = c.Call("Stdin.Write", log, &success)
	if err != nil {
		t.Error(err)
		return
	}

	ls, ok := logs[src]
	if !ok {
		t.Errorf("Could not find log lists for source \"%s\"", src)
		return
	}

	l, ok := ls[taskID]
	if !ok {
		t.Errorf("Could not find log list for task id \"%s\"", taskID)
		return
	}

	found := false
	for _, e := range l {
		if e.Msg == msg && e.Type == logType {
			found = true
			break
		}
	}
	fmt.Printf("log list: %v\n", logs)
	if !found {
		t.Errorf("Sent message \"%s\", but could not find it list", msg)
		return
	}
}
