package main

import (
	"strconv"
)

type version struct {
	Major    int
	Minor    int
	Revision int
	Build    int
}

var versions = make(map[string]version)

func getNextVersion(app string) (v string) {
	if vers, ok := versions[app]; ok {
		vers.Build = vers.Build + 1
		versions[app] = vers
		return v
	}
	versions[app] = version{Major: 0, Minor: 0, Revision: 0, Build: 1}

	vers := versions[app]

	return strconv.Itoa(vers.Major) + "." + strconv.Itoa(vers.Minor) + "." + strconv.Itoa(vers.Revision) + "." + strconv.Itoa(vers.Build)
}
