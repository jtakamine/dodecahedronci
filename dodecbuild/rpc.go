package main

import (
	"github.com/jtakamine/dodecahedronci/logutil"
)

type BuildArgs struct {
	RepoUrl string
}

type Build struct{}

func (b *Build) Execute(args BuildArgs, buildID *string) (err error) {
	*buildID = generateRandID(8)
	writer := logutil.NewWriter("build", *buildID)

	repoDir, err := cloneOrUpdateGitRepo(args.RepoUrl, writer)
	if err != nil {
		return err
	}

	err = build(repoDir, "myAwesomeApp", writer)
	if err != nil {
		return err
	}

	return nil
}
