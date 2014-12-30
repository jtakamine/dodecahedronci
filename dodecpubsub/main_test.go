package main

import (
	"github.com/jtakamine/dodecahedronci/dodecpubsub/api"
	"github.com/jtakamine/dodecahedronci/testutil"
	"testing"
)

func TestMain(t *testing.T) {
	testutil.FigBuild(t)
	testutil.FigUp(t)
	defer testutil.FigKillAndRm(t)

	testSubscribeAndPublish(t)
}

func testSubscribeAndPublish(t *testing.T) {
	var err error

	err = api.Subscribe("asdf")
	if err != nil {
		t.Error(err)
	}
}
