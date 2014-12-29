package main

import (
	"errors"
	"github.com/jtakamine/dodecahedronci/dodecregistry/api"
)

var registry = make(map[string]map[string]api.Build)

func addBuild(app string, version string, build api.Build) (err error) {
	if _, ok := registry[app]; !ok {
		registry[app] = make(map[string]api.Build)
	}
	registry[app][version] = build

	return nil
}

func getBuild(app string, version string) (build api.Build, err error) {
	if m, ok := registry[app]; ok {
		if build, ok := m[version]; ok {
			return build, nil
		}

		errStr := "No build found for version" + version + " of application \"" + app + "\""
		return api.Build{}, errors.New(errStr)
	}

	errStr := "No application found called \"" + app + "\""
	return api.Build{}, errors.New(errStr)
}
