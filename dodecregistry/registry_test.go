package main

import (
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
	pkg := []byte(fakeFigFile)
	dockerRegistryUrl := ""

	err := addPackage(app, version, pkg, dockerRegistryUrl)
	if err != nil {
		t.Error(err)

	}

	_pkg, _dockerRegistryUrl, err := getPackage(app, version)
	if err != nil {
		t.Error(err)
	}
	if string(pkg) != string(_pkg) || dockerRegistryUrl != _dockerRegistryUrl {
		t.Fail()
	}
}
