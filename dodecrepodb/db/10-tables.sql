CREATE TABLE IF NOT EXISTS task_type (
	id		serial PRIMARY KEY,
	code		varchar(64) NOT NULL UNIQUE,
	description	text NOT NULL
);

CREATE TABLE IF NOT EXISTS task_attribute_type (
	id		serial PRIMARY KEY,
	task_type_id	integer NOT NULL REFERENCES task_type (id),
	code		varchar(64) NOT NULL UNIQUE,
	description	text NOT NULL
);

CREATE TABLE IF NOT EXISTS task_artifact_type (
	id		serial PRIMARY KEY,
	code		varchar(64) NOT NULL UNIQUE,
	description	text NOT NULL
);


CREATE TABLE IF NOT EXISTS application (
	id 		serial PRIMARY KEY,
	name		varchar(64) NOT NULL UNIQUE,
	description	text NOT NULL
);

CREATE TABLE IF NOT EXISTS task (
	id 		serial PRIMARY KEY,
	parent_id	integer REFERENCES task (id),
	task_type_id 	integer NOT NULL REFERENCES task_type (id), 
	application_id	integer NOT NULL REFERENCES application (id),
	uuid 		varchar(64) NOT NULL UNIQUE,
	started		timestamp NOT NULL DEFAULT now(),
	completed	timestamp,
	success		boolean,
	node_path	ltree
);

CREATE TABLE IF NOT EXISTS task_log (
	id 		bigserial PRIMARY KEY,
	task_id		integer NOT NULL REFERENCES task (id),
	severity 	smallint NOT NULL,
	message 	text NOT NULL
);

CREATE TABLE IF NOT EXISTS task_artifact (
	id			serial PRIMARY KEY,
	task_artifact_type_id 	integer NOT NULL REFERENCES task_artifact_type (id),
	task_id			integer NOT NULL REFERENCES task (id),
	artifact		text NOT NULL
);

CREATE TABLE IF NOT EXISTS task_attribute (
	id 			serial PRIMARY KEY,
	task_id 		integer NOT NULL REFERENCES task (id),
	task_attribute_type_id 	integer NOT NULL REFERENCES task_attribute_type (id),
	value 			text,
	UNIQUE (task_id, task_attribute_type_id)
);

