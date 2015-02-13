package main

import (
	dodecregistry_API "github.com/jtakamine/dodecahedronci/dodecregistry/api"
	"gopkg.in/yaml.v2"
)

var saveBuild = func(app string, version string, fFile figFile, dockerRegistryUrl string) (err error) {
	data, err := yaml.Marshal(fFile.Config)
	if err != nil {
		return err
	}
	artifact := string(data)

	build := dodecregistry_API.Build{Artifact: artifact, DockerRegistryUrl: dockerRegistryUrl}

	err = dodecregistry_API.PostBuild(app, version, build, "http://dodecregistry:8000/")
	if err != nil {
		return err
	}

	return nil
}
