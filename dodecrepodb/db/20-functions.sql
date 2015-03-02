-- From http://bit.ly/1EUAGyc
CREATE OR REPLACE FUNCTION get_calculated_task_node_path(param_task_id integer)
	RETURNS ltree AS
$$
SELECT CASE WHEN t.parent_id IS NULL THEN t.id::text::ltree
	ELSE get_calculated_task_node_path(t.parent_id) || t.parent_id::text
	END
FROM task AS t 
WHERE t.id = $1;
$$
LANGUAGE sql;
