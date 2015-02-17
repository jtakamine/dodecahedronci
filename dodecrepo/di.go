package main

import (
	"errors"
	"github.com/lib/pq"
)

func addBuild(b Build) (err error) {
	fmt.Printf("TODO: save build to db--%s\n", b)
	return nil
}

func getBuild(app string, version string) (build Build, err error) {
	fmt.Println("TODO: get build from db")
	return api.Build{}
}
