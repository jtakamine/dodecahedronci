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

func (v *version) String() (s string) {
	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Revision) + "." + strconv.Itoa(v.Build)
}

var versions = make(map[string]version)

func getNextVersion(app string) (v string) {
	if vers, ok := versions[app]; ok {
		vers.Build = vers.Build + 1
		versions[app] = vers
		v = vers.String()

		return v
	}

	versions[app] = version{Major: 0, Minor: 0, Revision: 0, Build: 1}
	vers := versions[app]
	v = vers.String()

	return v
}
