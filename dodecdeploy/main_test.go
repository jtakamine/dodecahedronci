package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	artifact := `
myapp:
   image: golang
   command: ping google.com
`

	err := deploy(artifact)
	if err != nil {
		t.Error(err)
		return
	}

	err = killDeployments()
	if err != nil {
		t.Error(err)
		return
	}

	err = deploy(artifact)
	if err != nil {
		t.Error(err)
		return
	}
}
