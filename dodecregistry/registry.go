package main

import (
	"errors"
)

type dodecBuild struct {
	Artifact          string
	DockerRegistryUrl string
}

var registry = make(map[string]map[string]dodecBuild)

func addBuild(app string, version string, pkg dodecBuild) (err error) {
	if _, ok := registry[app]; !ok {
		registry[app] = make(map[string]dodecBuild)
	}
	registry[app][version] = pkg

	return nil
}

func getBuild(app string, version string) (build dodecBuild, err error) {
	if m, ok := registry[app]; ok {
		if build, ok := m[version]; ok {
			return build, nil
		}

		errStr := "No build found for version" + version + " of application \"" + app + "\""
		return dodecBuild{}, errors.New(errStr)
	}

	errStr := "No application found called \"" + app + "\""
	return dodecBuild{}, errors.New(errStr)
}
