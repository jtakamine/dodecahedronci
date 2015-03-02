-- From http://bit.ly/1EUAGyc
CREATE OR REPLACE FUNCTION trig_update_task_node_path() RETURNS trigger AS
$$
BEGIN
	IF TG_OP = 'UPDATE' THEN
		IF COALESCE(OLD.parent_id,0) != COALESCE(NEW.parent_id,0) THEN
			-- update all nodes that are children of this one including this one
            		UPDATE task SET node_path = get_calculated_task_node_path(id)
				WHERE OLD.node_path @> task.node_path;
		END IF;
	ELSIF TG_OP = 'INSERT' THEN
		UPDATE task SET node_path = get_calculated_task_node_path(NEW.id)
			WHERE task.id = NEW.id;
	END IF;

	RETURN NEW;
END
$$
LANGUAGE 'plpgsql' VOLATILE;

CREATE TRIGGER trig_update_task_node_path AFTER INSERT OR UPDATE OF parent_id
   ON task FOR EACH ROW
   EXECUTE PROCEDURE trig_update_task_node_path();
