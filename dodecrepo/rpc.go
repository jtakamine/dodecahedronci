package main

import (
	_ "github.com/lib/pq"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"time"
)

func rpcListen(port int) (err error) {
	err = rpc.RegisterName("BuildRepo", &RPCBuildRepo{})
	if err != nil {
		return err
	}

	err = rpc.RegisterName("ArtifactRepo", &RPCArtifactRepo{})
	if err != nil {
		return err
	}

	err = rpc.RegisterName("AppRepo", &RPCApplicationRepo{})
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

	return nil
}

type Application struct {
	Name        string
	Description string
}

type Build struct {
	UUID    string
	AppName string
	Version string
}

type TaskCompletionInfo struct {
	UUID    string
	Success bool
}

type BuildDetails struct {
	Build
	Started   time.Time
	Completed time.Time
	Success   bool
	Artifact  string
}

type Artifact struct {
	Artifact  string
	Type      string
	BuildUUID string
}

type RPCApplicationRepo struct{}

func (rpcA *RPCApplicationRepo) Save(a Application, success *bool) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	err = saveApplication(a.Name, a.Description, c)
	if err != nil {
		return err
	}

	*success = true
	return nil
}

type RPCBuildRepo struct{}

func (rpcB *RPCBuildRepo) Save(b Build, success *bool) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	err = saveBuild(b.UUID, b.AppName, b.Version, c)
	if err != nil {
		return err
	}

	*success = true
	return nil
}

func (rpcB *RPCBuildRepo) RecordCompletion(tci TaskCompletionInfo, success *bool) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	err = recordCompletion(tci.UUID, tci.Success, c)
	if err != nil {
		return err
	}

	*success = true
	return nil
}

func (rpcB *RPCBuildRepo) Get(uuid string, b *BuildDetails) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	*b, err = getBuild(uuid, c)
	if err != nil {
		return err
	}

	return nil
}

func (rpcB *RPCBuildRepo) GetAll(appName string, bs *[]Build) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	*bs, err = getBuilds(appName, c)
	if err != nil {
		return err
	}

	return nil
}

type RPCArtifactRepo struct{}

func (rpcA *RPCArtifactRepo) Save(a Artifact, success *bool) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	err = saveArtifact(a.Artifact, a.BuildUUID, a.Type, c)
	if err != nil {
		return err
	}

	*success = true
	return nil
}
