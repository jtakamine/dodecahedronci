package main

import (
	"fmt"
	"github.com/jtakamine/dodecahedronci/dodecpubsub/api"
	"github.com/jtakamine/dodecahedronci/testutil"
	"strconv"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	testutil.FigBuild(t)
	testutil.FigUp(t)
	defer testutil.FigKillAndRm(t)

	testSubscribeAndPublish("myChannel", "localhost:8000", t)
}

func testSubscribeAndPublish(channel string, address string, t *testing.T) {
	var err error

	subChan, err := api.Subscribe(channel, address)
	if err != nil {
		t.Error(err)
	}

	iterations := 25

	for i := 1; i <= iterations; i++ {
		msg := createMessage(i)
		err = api.Publish(msg, channel, address)
		if err != nil {
			t.Error(err)
		}
	}

	timeout := 500
	msgs := make(map[string]struct{})
	for i := 1; i <= iterations; i++ {
		select {
		case msg, ok := <-subChan:
			if !ok {
				t.Error("Subscription channel closed unexpectedly")
			}
			msgs[msg] = struct{}{}
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			t.Error("Receive timed out after " + strconv.Itoa(timeout) + "ms")
		}
	}

	if len(msgs) != iterations {
		t.Errorf("Received %v messages, but %v messages were published", len(msgs), iterations)
	}
}

func createMessage(i int) (msg string) {
	return "My Message #" + strconv.Itoa(i)
}
