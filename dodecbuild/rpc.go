package main

import (
	"github.com/jtakamine/dodecahedronci/logutil"
)

type BuildArgs struct {
	RepoUrl    string
	DockerUser string
}

type Build struct{}

func (b *Build) Execute(args BuildArgs, buildID *string) (err error) {
	*buildID = generateBuildID()
	writer := logutil.NewWriter("build", *buildID)

	repoDir, err := cloneOrUpdateGitRepo(args.RepoUrl, writer)
	if err != nil {
		return err
	}

	err = build(repoDir, "myAwesomeApp", "http://localhost:8080", args.DockerUser, writer)
	if err != nil {
		return err
	}

	return nil
}
