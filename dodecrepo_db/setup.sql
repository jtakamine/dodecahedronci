-- Tables

CREATE TABLE IF NOT EXISTS task_type (
	id		serial PRIMARY KEY,
	code		varchar(128) NOT NULL UNIQUE,
	description	text NOT NULL
);

CREATE TABLE IF NOT EXISTS task (
	id 		serial PRIMARY KEY,
	task_type_id 	integer NOT NULL REFERENCES task_type (id), 
	uuid 		varchar(128) NOT NULL UNIQUE,
	artifact 	json NOT NULL
);

CREATE TABLE IF NOT EXISTS task_log (
	id 		bigserial PRIMARY KEY,
	task_id		integer NOT NULL REFERENCES task (id),
	severity 	smallint NOT NULL,
	message 	text NOT NULL
);

CREATE TABLE IF NOT EXISTS task_attribute_type (
	id		serial PRIMARY KEY,
	task_type_id	integer NOT NULL REFERENCES task_type (id),
	code		varchar(128) NOT NULL UNIQUE,
	description	text NOT NULL
);

CREATE TABLE IF NOT EXISTS task_attribute (
	id 			serial PRIMARY KEY,
	task_id 		integer NOT NULL REFERENCES task (id),
	task_attribute_type_id 	integer NOT NULL REFERENCES task_attribute_type (id),
	value 			text
);


-- Seed Data
INSERT INTO task_type (code, description)
	VALUES ('build', 'Dodec build task'),
	VALUES ('deploy', 'Dodec deploy task'),
	VALUES ('test', 'Dodec test task');
