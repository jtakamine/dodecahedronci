-- task_type

INSERT INTO task_type (code, description)
	SELECT 'build', 'Dodec build task'
	WHERE NOT EXISTS (SELECT 1 FROM task_type WHERE code = 'build');

INSERT INTO task_type (code, description)
	SELECT 'deploy', 'Dodec deploy task'
	WHERE NOT EXISTS (SELECT 1 FROM task_type WHERE code = 'deploy');

INSERT INTO task_type (code, description)
	SELECT 'test', 'Dodec test task'
	WHERE NOT EXISTS (SELECT 1 FROM task_type WHERE code = 'test');


-- task_attribute_type

INSERT INTO task_attribute_type (task_type_id, code, description)
	SELECT id, 'version', 'Version number'
	FROM task_type
	WHERE code = 'build'
	AND NOT EXISTS (SELECT 1 FROM task_attribute_type WHERE code = 'version');


-- task_artifact_type

INSERT INTO task_artifact_type (code, description)
	SELECT 'build_artifact', 'Build artifact (fig.yml)'
	WHERE NOT EXISTS (SELECT 1 FROM task_artifact_type WHERE code = 'build_artifact');
