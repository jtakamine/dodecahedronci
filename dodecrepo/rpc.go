package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
)

type RPCBuildRepo struct{}

func (rpcB *RPCBuildRepo) Save(b Build, success *bool) (err error) {
	c, err := getConnStr()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", c)
	if err != nil {
		return err
	}

	sqlFmt := `
INSERT INTO task(task_type_id, uuid, artifact)
	SELECT TOP 1 tt.id, '%s', '%s'
	FROM task_type tt
	WHERE tt.code = 'build'
`
	err := db.Exec("INSERT INTO task()")

	*success = true
	return nil
}

func (rpcB *RPCBuildRepo) Get(uuid string, b *Build) (err error) {
	*build, err = getBuild(uuid)

	return nil
}
