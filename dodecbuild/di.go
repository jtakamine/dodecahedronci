package main

import (
	"errors"
	"fmt"
	"github.com/jtakamine/dodecahedronci/logutil"
	"gopkg.in/yaml.v2"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var buildDockerFile = func(dFile string, version string, writer *logutil.Writer) (tag string, err error) {
	dockerRegistryAddr := os.Getenv("DOCKER_REGISTRYADDR")
	if dockerRegistryAddr == "" {
		return "", errors.New("Missing environment variable: \"DOCKER_REGISTRYADDR\"")
	}

	registryPrefix := ""
	if dockerRegistryAddr != "" {
		registryPrefix = dockerRegistryAddr + "/"
	}

	dockerUser := "" //For now, do not provide a user specification
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

var saveBuild = func(uuid string, appName string) (version string, err error) {
	build := struct {
		UUID    string
		AppName string
	}{
		UUID:    uuid,
		AppName: appName,
	}

	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return "", errors.New("Missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return "", err
	}
	c := jsonrpc.NewClient(conn)

	err = c.Call("BuildRepo.Save", build, &version)
	if err != nil {
		return "", err
	}

	return version, nil
}

var saveBuildArtifact = func(uuid string, fFile figFile) (err error) {
	fData, err := yaml.Marshal(fFile.Config)
	if err != nil {
		return err
	}
	fYml := string(fData)

	args := struct {
		Artifact  string
		BuildUUID string
	}{
		Artifact:  fYml,
		BuildUUID: uuid,
	}

	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return errors.New("Missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return err
	}
	c := jsonrpc.NewClient(conn)

	var success bool
	err = c.Call("ArtifactRepo.Save", args, &success)
	if err != nil {
		return err
	}

	return nil
}

var recordCompletion = func(uuid string, success bool) (err error) {
	addr := os.Getenv("DODEC_REPOADDR")
	if addr == "" {
		return errors.New("Missing environment variable: DODEC_REPOADDR")
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return err
	}
	c := jsonrpc.NewClient(conn)

	args := struct {
		UUID    string
		Success bool
	}{
		UUID:    uuid,
		Success: success,
	}

	err = c.Call("TaskRepo.RecordCompletion", args, &success)
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
	controlAddr := os.Getenv("DODEC_CONTROLADDR")
	if controlAddr == "" {
		return errors.New("Missing environment variable: DODEC_CONTROLADDR")
	}

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

	conn, err := net.DialTimeout("tcp", controlAddr, time.Second)
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
