package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func getDB(connStr string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func saveApplication(name string, description string, connStr string) (err error) {
	db, err := getDB(connStr)
	if err != nil {
		return err
	}

	s := `
INSERT INTO application (name, description)
	VALUES($1, $2)
`

	st, err := db.Prepare(s)
	if err != nil {
		return err
	}

	_, err = st.Exec(name, description)
	if err != nil {
		return err
	}

	return nil
}

func saveBuild(uuid string, appName string, version string, connStr string) (err error) {
	db, err := getDB(connStr)
	if err != nil {
		return err
	}

	s1 := `
INSERT INTO task (task_type_id, application_id, uuid)
	SELECT tt.id, a.id, $1
	FROM task_type tt
		CROSS JOIN application a
	WHERE tt.code = 'build'
		AND a.name = $2
	RETURNING id;
`

	s2 := `
INSERT INTO task_attribute (task_id, task_attribute_type_id, value)
	SELECT $1, tat.id, $2
	FROM task_attribute_type tat
	WHERE tat.code = 'version'
`

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	st1, err := db.Prepare(s1)
	if err != nil {
		return err
	}

	st2, err := db.Prepare(s2)
	if err != nil {
		return err
	}

	var id int
	err = st1.QueryRow(uuid, appName).Scan(&id)
	if err != nil {
		return err
	}

	_, err = st2.Exec(id, version)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func recordCompletion(uuid string, success bool, connStr string) (err error) {
	db, err := getDB(connStr)
	if err != nil {
		return err
	}

	s := `
UPDATE task
	SET completed = true,
	success = $1
WHERE uuid = $2
`

	st, err := db.Prepare(s)
	if err != nil {
		return err
	}

	_, err = st.Exec(success, uuid)
	if err != nil {
		return err
	}

	return nil
}

func getBuild(uuid string, connStr string) (b BuildDetails, err error) {
	db, err := getDB(connStr)
	if err != nil {
		return BuildDetails{}, err
	}

	s := `
SELECT t.uuid, a.name, tatt.value, t.started, t.completed, t.success, tart.artifact
FROM task t
	LEFT OUTER JOIN task_attribute tatt on tatt.task_id = t.id
	LEFT OUTER JOIN task_attribute_type tat on tat.id = tatt.task_attribute_type_id
	LEFT OUTER JOIN task_artifact tart on tart.task_id = t.id
	LEFT OUTER JOIN application a on a.id = t.application_id
WHERE uuid = $1
	AND tat.code = 'version'
`

	st, err := db.Prepare(s)
	if err != nil {
		return BuildDetails{}, err
	}

	err = st.QueryRow(uuid).Scan(&b.UUID, &b.AppName, &b.Version, &b.Started, &b.Completed, &b.Success, &b.Artifact)
	if err != nil {
		return BuildDetails{}, err
	}

	return b, nil
}

func getBuilds(appName string, connStr string) (builds []Build, err error) {
	db, err := getDB(connStr)

	s := `
SELECT t.uuid, a.name, tatt.value
FROM task t
	LEFT OUTER JOIN task_attribute tatt on tatt.task_id = t.id
	LEFT OUTER JOIN task_attribute_type tat on tat.id = tatt.task_attribute_type_id
	LEFT OUTER JOIN application a on a.id = t.application_id
WHERE a.name = $1
ORDER BY tatt.value
`

	st, err := db.Prepare(s)
	if err != nil {
		return nil, err
	}

	rows, err := st.Query(appName)

	for rows.Next() {
		b := Build{}
		err = rows.Scan(&b.UUID, &b.AppName, &b.Version)
		if err != nil {
			return nil, err
		}

		builds = append(builds, b)
	}

	return builds, err
}

func saveArtifact(artifact string, buildUUID string, artifactType string, connStr string) (err error) {
	db, err := getDB(connStr)
	if err != nil {
		return err
	}

	s := `
INSERT INTO task_artifact (task_artifact_type_id, task_id, artifact)
	SELECT tat.id, t.id, $1
	FROM task_artifact_type tat
		CROSS JOIN task t
	WHERE tat.code = $2
		AND t.UUID = $3
`

	st, err := db.Prepare(s)
	if err != nil {
		return err
	}

	_, err = st.Exec(artifact, "build_artifact", buildUUID)
	if err != nil {
		return err
	}

	return nil
}
