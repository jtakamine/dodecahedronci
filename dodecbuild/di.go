package main

import (
	dodecregistry_API "github.com/jtakamine/dodecahedronci/dodecregistry/api"
	"github.com/jtakamine/dodecahedronci/logutil"
	"gopkg.in/yaml.v2"
	"os/exec"
	"path/filepath"
	"strings"
)

var buildDockerFile = func(dFile string, version string, writer *logutil.Writer) (tag string, err error) {
	//TODO: include private registry specification (and user?)
	dockerRegistryUrl := "registry:5000"
	dockerUser := "jtakamine"

	registryPrefix := ""
	if dockerRegistryUrl != "" {
		dockerRegistryUrl = strings.TrimPrefix(dockerRegistryUrl, "http://")
		dockerRegistryUrl = strings.TrimPrefix(dockerRegistryUrl, "https://")
		registryPrefix = dockerRegistryUrl + "/"
	}
	userPrefix := ""
	if dockerUser != "" {
		userPrefix = dockerUser + "/"
	}
	versionSuffix := ""
	if version != "" {
		versionSuffix = ":" + version
	}
	tag = registryPrefix + userPrefix + generateRandID(8) + versionSuffix

	cmd := exec.Command("docker", "build", "-t", tag, ".")
	cmd.Dir = filepath.Dir(dFile)
	cmd.Stdout = writer.CreateWriter(logutil.Verbose)
	cmd.Stderr = writer.CreateWriter(logutil.Error)

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return tag, nil
}

var saveBuild = func(app string, version string, fFile figFile) (err error) {
	data, err := yaml.Marshal(fFile.Config)
	if err != nil {
		return err
	}
	artifactStr := string(data)

	build := dodecregistry_API.Build{Artifact: artifactStr}

	err = dodecregistry_API.PostBuild(app, version, build, "http://dodecregistry:8000/")
	if err != nil {
		return err
	}

	return nil
}

var pushDockerImage = func(tag string, writer *logutil.Writer) (err error) {
	cmd := exec.Command("docker", "push", tag)
	cmd.Stdout = writer.CreateWriter(logutil.Verbose)
	cmd.Stderr = writer.CreateWriter(logutil.Error)

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
