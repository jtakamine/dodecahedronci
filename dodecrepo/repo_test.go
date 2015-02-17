package main

import (
	"github.com/jtakamine/dodecahedronci/dodecregistry/api"
	"testing"
)

var fakeFigFile = `
dodecbuild:
        build: ./dodecbuild
        ports:
                - "8000:8000"
        environment:
                DODEC_HOME:
                DODEC_GITHUB_USER:
                DODEC_GITHUB_PASSWORD:
                DODEC_DOCKER_USER:
                DODEC_DOCKER_PASSWORD:
                DODEC_DOCKER_EMAIL:
        privileged: true
`

func TestAddPackage(t *testing.T) {
	app := "exampleApp"
	version := "3.2.0.8732"
	artifact := fakeFigFile
	dockerRegistryUrl := ""

	build_add := api.Build{Artifact: artifact, DockerRegistryUrl: dockerRegistryUrl}

	err := addBuild(app, version, build_add)
	if err != nil {
		t.Error(err)

	}

	build_get, err := getBuild(app, version)
	if err != nil {
		t.Error(err)
	}
	if build_add != build_get {
		t.Fail()
	}
}
