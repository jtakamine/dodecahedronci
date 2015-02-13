package main

import (
	"github.com/jtakamine/dodecahedronci/logutil"
)

type Build struct{}

func (b *Build) Execute(repoUrl string, buildID *string) (err error) {
	id := generateBuildID()
	buildID = &id
	writer := logutil.NewWriter(id)

	repoDir, err := cloneOrUpdateGitRepo(repoUrl, writer)
	if err != nil {
		return err
	}

	err = build(repoDir, "myAwesomeApp", "http://localhost:8080", writer)
	if err != nil {
		return err
	}

	return nil
}
