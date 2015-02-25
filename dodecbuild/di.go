package main

import (
	"errors"
	"fmt"
	"github.com/jtakamine/dodecahedronci/logutil"
	"gopkg.in/yaml.v2"
	"net"
	"net/rpc/jsonrpc"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

var saveBuild = func(appName string, version string, fFile figFile) (err error) {
	data, err := yaml.Marshal(fFile.Config)
	if err != nil {
		return err
	}
	artifactStr := string(data)

	build := dodecregistry_API.Build{Artifact: artifactStr}

	err = dodecregistry_API.PostBuild(appName, version, build, "http://dodecrepo:8000/")
	if err != nil {
		return err
	}

	return nil
}

var pushDockerImage = func(tag string, writer *logutil.Writer) (err error) {
	cmd := exec.Command("docker", "push", tag)
	cmd.Stdout = writer.CreateWriter(logutil.Verbose)
	cmd.Stderr = writer.CreateWriter(logutil.Error)

	fmt.Printf("pushing docker image: %s", tag)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

var rpcRegisterService = func() (err error) {
	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		return err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}
	if len(addrs) == 0 {
		return errors.New("Could not find IP address associated with eth0 device")
	}

	ip := strings.Split(addrs[0].String(), "/")[0]

	conn, err := net.DialTimeout("tcp", "dodeccontrol:9000", time.Second)
	if err != nil {
		return err
	}
	c := jsonrpc.NewClient(conn)

	args := &struct {
		Service  string
		Endpoint string
	}{
		Service:  "build",
		Endpoint: ip + ":9000",
	}

	var success bool
	err = c.Call("ServiceRegistry.Register", args, &success)
	if err != nil {
		return err
	}

	return nil
}
