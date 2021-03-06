package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"time"
)

var database *sql.DB

func getDB(connStr string) (db *sql.DB, err error) {
	if database == nil {
		database, err = sql.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}
	}

	db = database
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

func getApplication(name string, connStr string) (a Application, err error) {
	db, err := getDB(connStr)
	if err != nil {
		return Application{}, err
	}

	s := `
SELECT name, description
FROM application
WHERE name = $1
`

	st, err := db.Prepare(s)
	if err != nil {
		return Application{}, err
	}

	err = st.QueryRow(name).Scan(&a.Name, &a.Description)

	if err != nil && err != sql.ErrNoRows {
		return Application{}, err
	}

	return a, nil
}

func saveBuild(uuid string, appName string, connStr string) (version string, err error) {
	db, err := getDB(connStr)
	if err != nil {
		return "", err
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
	SELECT $1, x.tat_id, x.ver
	FROM (
		SELECT y.tat_id, array_to_string(ARRAY[y.vArr[1], y.vArr[2], y.vArr[3], y.vArr[4]+1], '.') as ver
		FROM (
			SELECT tat.id as tat_id, regexp_split_to_array(ta.value, E'\\.')::int[] as vArr
			FROM task t
				INNER JOIN application a on a.id = t.application_id
				INNER JOIN task_attribute ta on ta.task_id = t.id
				INNER JOIN task_attribute_type tat on tat.id = ta.task_attribute_type_id
			WHERE tat.code = 'version'
				AND a.name = $2
		) y

		UNION

		SELECT tat.id as tat_id, '0.0.0.1' as ver
		FROM task_attribute_type tat
		WHERE tat.code = 'version'
	) x
	ORDER BY x.ver desc
	LIMIT 1
	RETURNING value
`

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	st1, err := db.Prepare(s1)
	if err != nil {
		return "", err
	}

	st2, err := db.Prepare(s2)
	if err != nil {
		return "", err
	}

	var id int
	err = st1.QueryRow(uuid, appName).Scan(&id)
	if err != nil {
		return "", err
	}

	err = st2.QueryRow(id, appName).Scan(&version)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return version, nil
}

func getBuild(uuid string, connStr string) (b BuildDetails, err error) {
	db, err := getDB(connStr)
	if err != nil {
		return BuildDetails{}, err
	}

	s := `
SELECT t.uuid, a.name, tatt.value, t.started, t.completed, t.success, tart.artifact
FROM task t
	LEFT OUTER JOIN task_type tt on tt.id = t.task_type_id
	LEFT OUTER JOIN task_attribute tatt on tatt.task_id = t.id
	LEFT OUTER JOIN task_attribute_type tat on tat.id = tatt.task_attribute_type_id
	LEFT OUTER JOIN task_artifact tart on tart.task_id = t.id
	LEFT OUTER JOIN application a on a.id = t.application_id
WHERE t.uuid = $1
	AND tat.code = 'version'
	AND tt.code = 'build'
`

	st, err := db.Prepare(s)
	if err != nil {
		return BuildDetails{}, err
	}

	nCompleted := pq.NullTime{}
	nSuccess := sql.NullBool{}
	nArtifact := sql.NullString{}

	err = st.QueryRow(uuid).Scan(&b.UUID, &b.AppName, &b.Version, &b.Started, &nCompleted, &nSuccess, &nArtifact)
	if err != nil && err != sql.ErrNoRows {
		return BuildDetails{}, err
	}

	if nCompleted.Valid {
		b.Completed = nCompleted.Time
	}
	if nSuccess.Valid {
		b.Success = nSuccess.Bool
	}
	if nArtifact.Valid {
		b.Artifact = nArtifact.String
	}

	return b, nil
}

func getBuilds(appName string, connStr string) (builds []Build, err error) {
	db, err := getDB(connStr)
	if err != nil {
		return nil, fmt.Errorf("getting database: %s", err.Error())
	}

	s := `
SELECT t.uuid, a.name, tatt.value
FROM task t
	LEFT OUTER JOIN task_type tt on tt.id = t.task_type_id
	LEFT OUTER JOIN task_attribute tatt on tatt.task_id = t.id
	LEFT OUTER JOIN task_attribute_type tat on tat.id = tatt.task_attribute_type_id
	LEFT OUTER JOIN application a on a.id = t.application_id
WHERE (a.name = $1 OR '' = $1)
	AND tt.code = 'build'
ORDER BY t.id
`

	st, err := db.Prepare(s)
	if err != nil {
		return nil, fmt.Errorf("preparing sql statement: %s", err.Error())
	}

	rows, err := st.Query(appName)
	if err != nil {
		return nil, fmt.Errorf("executing sql query: %s", err.Error())
	}

	builds = make([]Build, 0)
	for rows.Next() {
		b := Build{}
		err = rows.Scan(&b.UUID, &b.AppName, &b.Version)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %s", err.Error())
		}

		builds = append(builds, b)
	}

	return builds, nil
}

func saveDeploy(deployUUID string, buildUUID string, appName string, connStr string) (err error) {
	db, err := getDB(connStr)
	if err != nil {
		return err
	}

	fmt.Printf("deployUUID: %s; buildUUID: %s; appName: %s;\n", deployUUID, buildUUID, appName)

	s := `
INSERT INTO task(parent_id, task_type_id, application_id, uuid)
	SELECT t.id, tt.id, a.id, $1
	FROM task t
		CROSS JOIN task_type tt
		CROSS JOIN application a
	WHERE 
		t.uuid = $2
		AND tt.code = 'deploy'
		AND a.name = $3
`

	st, err := db.Prepare(s)
	if err != nil {
		return err
	}

	_, err = st.Exec(deployUUID, buildUUID, appName)
	if err != nil {
		return err
	}

	return nil
}

func getDeploy(uuid string, connStr string) (d DeployDetails, err error) {
	db, err := getDB(connStr)
	if err != nil {
		return DeployDetails{}, err
	}

	s := `
SELECT t.uuid, a.name, t.started, t.completed, t.success, t2.uuid
FROM task t
	LEFT OUTER JOIN task_type tt on tt.id = t.task_type_id
	LEFT OUTER JOIN task t2 on t2.id = t.parent_id
	LEFT OUTER JOIN application a on a.id = t.application_id
WHERE t.uuid = $1
	AND tt.code = 'deploy'
`

	st, err := db.Prepare(s)
	if err != nil {
		return DeployDetails{}, err
	}

	nCompleted := pq.NullTime{}
	nSuccess := sql.NullBool{}

	err = st.QueryRow(uuid).Scan(&d.UUID, &d.AppName, &d.Started, &nCompleted, &nSuccess, &d.BuildUUID)
	if err != nil && err != sql.ErrNoRows {
		return DeployDetails{}, err
	}

	if nCompleted.Valid {
		d.Completed = nCompleted.Time
	}
	if nSuccess.Valid {
		d.Success = nSuccess.Bool
	}

	return d, nil
}

func getDeploys(appName string, connStr string) (deploys []Deploy, err error) {
	db, err := getDB(connStr)
	if err != nil {
		return nil, err
	}

	s := `
SELECT t.uuid, a.name, t2.uuid
FROM task t
	LEFT OUTER JOIN task_type tt on tt.id = t.task_type_id
	LEFT OUTER JOIN task t2 on t2.id = t.parent_id
	LEFT OUTER JOIN application a on a.id = t.application_id
WHERE (a.name = $1 OR '' = $1)
	AND tt.code = 'deploy'
ORDER BY t.id
`

	st, err := db.Prepare(s)
	if err != nil {
		return nil, err
	}

	rows, err := st.Query(appName)

	deploys = make([]Deploy, 0)
	for rows.Next() {
		d := Deploy{}
		err = rows.Scan(&d.UUID, &d.AppName, &d.BuildUUID)
		if err != nil {
			return nil, err
		}

		deploys = append(deploys, d)
	}

	return deploys, err
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

func recordCompletion(uuid string, success bool, connStr string) (err error) {
	db, err := getDB(connStr)
	if err != nil {
		return err
	}

	s := `
UPDATE task
	SET completed = now(),
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

func saveLog(taskUUID string, message string, severity int, created time.Time, connStr string) (err error) {
	db, err := getDB(connStr)
	if err != nil {
		return err
	}

	s := `
INSERT INTO task_log (task_id, severity, message, created)
	SELECT t.id, $1, $2, $3
	FROM task t
	WHERE t.uuid = $4
`

	st, err := db.Prepare(s)
	if err != nil {
		return err
	}

	_, err = st.Exec(severity, message, created, taskUUID)
	if err != nil {
		return err
	}

	return nil
}

func getLogs(taskUUID string, severity int, startID int64, connStr string) (logs []Log, err error) {
	db, err := getDB(connStr)
	if err != nil {
		return nil, err
	}

	s := `
SELECT tl.id, t.uuid, tl.message, tl.severity, tl.created
FROM task_log tl
	INNER JOIN task t on t.id = tl.task_id
WHERE tl.severity >= $1
	AND t.uuid = $2
	AND tl.id >= $3
ORDER BY tl.id
`

	st, err := db.Prepare(s)
	if err != nil {
		return nil, err
	}

	rows, err := st.Query(severity, taskUUID, startID)

	logs = make([]Log, 0)
	for rows.Next() {
		l := Log{}
		err = rows.Scan(&l.ID, &l.TaskUUID, &l.Message, &l.Severity, &l.Created)
		if err != nil {
			return nil, err
		}

		logs = append(logs, l)
	}

	return logs, nil
}
