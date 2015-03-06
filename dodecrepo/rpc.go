package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

func rpcListen(port int) (err error) {
	err = rpc.RegisterName("AppRepo", &RPCApplicationRepo{})
	if err != nil {
		return err
	}

	err = rpc.RegisterName("BuildRepo", &RPCBuildRepo{})
	if err != nil {
		return err
	}

	err = rpc.RegisterName("DeployRepo", &RPCDeployRepo{})
	if err != nil {
		return err
	}

	err = rpc.RegisterName("TaskRepo", &RPCTaskRepo{})
	if err != nil {
		return err
	}

	err = rpc.RegisterName("ArtifactRepo", &RPCArtifactRepo{})
	if err != nil {
		return err
	}

	err = rpc.RegisterName("LogRepo", &RPCLogRepo{})

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

type RPCApplicationRepo struct{}

func (*RPCApplicationRepo) Save(a Application, success *bool) (err error) {
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

func (*RPCApplicationRepo) Get(name string, a *Application) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	*a, err = getApplication(name, c)
	if err != nil {
		return err
	}

	return nil
}

type RPCBuildRepo struct{}

func (*RPCBuildRepo) Save(b Build, version *string) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	*version, err = saveBuild(b.UUID, b.AppName, c)
	if err != nil {
		return err
	}

	return nil
}

func (*RPCBuildRepo) Get(uuid string, b *BuildDetails) (err error) {
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

func (*RPCBuildRepo) GetAll(appName string, bs *[]Build) (err error) {
	c, err := getConnStr()
	if err != nil {
		return fmt.Errorf("getting connection string: %s", err.Error())
	}

	*bs, err = getBuilds(appName, c)
	if err != nil {
		return fmt.Errorf("getting builds from database: %s", err.Error())
	}

	return nil
}

type RPCDeployRepo struct{}

func (*RPCDeployRepo) Save(d Deploy, success *bool) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	err = saveDeploy(d.UUID, d.BuildUUID, d.AppName, c)
	if err != nil {
		return err
	}

	return nil
}

func (*RPCDeployRepo) Get(uuid string, d *DeployDetails) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	*d, err = getDeploy(uuid, c)
	if err != nil {
		return err
	}

	return nil
}

func (*RPCDeployRepo) GetAll(appName string, ds *[]Deploy) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	*ds, err = getDeploys(appName, c)
	if err != nil {
		return err
	}

	return nil
}

type RPCTaskRepo struct{}

func (*RPCTaskRepo) RecordCompletion(tci TaskCompletionInfo, success *bool) (err error) {
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

type RPCArtifactRepo struct{}

func (*RPCArtifactRepo) Save(a Artifact, success *bool) (err error) {
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

type RPCLogRepo struct{}

func (*RPCLogRepo) Save(l Log, success *bool) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	err = saveLog(l.TaskUUID, l.Message, l.Severity, l.Created, c)
	if err != nil {
		return err
	}

	*success = true
	return nil
}

func (*RPCLogRepo) GetAll(lq LogQuery, ls *[]Log) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	*ls, err = getLogs(lq.TaskUUID, lq.Severity, c)
	if err != nil {
		return err
	}

	return nil
}
