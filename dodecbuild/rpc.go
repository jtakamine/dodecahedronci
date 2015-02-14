package main

import (
	"github.com/jtakamine/dodecahedronci/logutil"
)

type Build struct{}

func (b *Build) Execute(repoUrl string, buildID *string) (err error) {
	*buildID = generateBuildID()
	writer := logutil.NewWriter("build", *buildID)

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
