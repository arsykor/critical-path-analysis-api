/* 1. Create db_tasks database */

-- DROP DATABASE IF EXISTS db_tasks;

CREATE DATABASE db_tasks
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

/* 2. Create t_task table - main table with all tasks. */

CREATE TABLE IF NOT EXISTS t_task
(
    id bigint NOT NULL,
    name character varying(20) NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    CONSTRAINT t_task_pkey PRIMARY KEY (id)
    )

/* 3. Create t_task_changes table - task deadlines shifts. */

CREATE TABLE IF NOT EXISTS t_task_changes
(
    id SERIAL NOT NULL,
    id_task bigint NOT NULL,
    start_date_old date NOT NULL,
    end_date_old date NOT NULL,
    start_date_new date NOT NULL,
    end_date_new date NOT NULL,
    deadline_delay boolean NOT NULL,
    changed_date timestamp NOT NULL,
    CONSTRAINT t_task_changes_pkey PRIMARY KEY (id)
    )

/* 4. Create t_task_predecessor table - predecessors of tasks (required for critical path analysis). */

CREATE TABLE IF NOT EXISTS t_task_predecessor
(
    id bigint NOT NULL,
    id_task bigint,
    id_predecessor bigint,
    CONSTRAINT t_task_predecessor_pkey PRIMARY KEY (id)
    )

/* 5. Stored procedure for inserting all tasks.  New tasks from the array are added, old ones are updated. */

CREATE PROCEDURE sp_insert_merge
(p_id bigint,
 p_name text,
 p_start_date date,
 p_end_date date)
    LANGUAGE plpgsql
AS $$
DECLARE
temp_row t_task%ROWTYPE;
BEGIN
SELECT * INTO temp_row FROM t_task WHERE id = p_id;
IF temp_row IS NOT NULL THEN

       IF temp_row.start_date <> p_start_date OR temp_row.end_date <> p_end_date THEN
            INSERT INTO t_task_changes
            (id_task, start_date_old, end_date_old, start_date_new, end_date_new, deadline_delay, changed_date)
            VALUES (temp_row.id, temp_row.start_date, temp_row.end_date, p_start_date, p_end_date,
                   p_end_date > temp_row.end_date, current_timestamp::date);
       END IF;

UPDATE t_task
SET id = p_id,
    name = p_name,
    start_date = p_start_date,
    end_date = p_end_date
WHERE id = P_ID;

ELSE
       INSERT INTO t_task
       (id, name, start_date, end_date)
       VALUES (p_id, p_name, p_start_date, p_end_date);
END IF;
END
$$;

/* 6. Stored procedure for inserting predecessors (connected tasks) for the tasks. */

CREATE PROCEDURE sp_predecessors
(p_id_task bigint,
 arr_pred bigint[])
    LANGUAGE plpgsql
AS $$
DECLARE
    v_pred bigint;
    arr_inserted bigint[];
BEGIN
FOR i IN array_lower(arr_pred, 1) .. array_upper(arr_pred, 1)
    LOOP
        IF NOT EXISTS (SELECT id FROM t_task_predecessor
                        WHERE id_task = p_id_task
                        AND id_predecessor = arr_pred[i])
                THEN
                INSERT INTO t_task_predecessor
                (id_task, id_predecessor)
                VALUES (p_id_task, arr_pred[i])
                RETURNING id INTO v_pred;
                    arr_inserted = array_append(arr_inserted, v_pred);

ELSE
SELECT id FROM t_task_predecessor
WHERE id_task = p_id_task
  AND id_predecessor = arr_pred[i]
    INTO v_pred;
arr_inserted = array_append(arr_inserted, v_pred);
END IF;
END LOOP;

FOR i IN array_lower(arr_inserted, 1) .. array_upper(arr_inserted, 1)
    LOOP
DELETE FROM t_task_predecessor
WHERE id != ALL(arr_inserted)
        AND id_task = p_id_task;
END LOOP;
END
$$;

/* 7. Stored procedure for deleting the tasks from t_task, t_task_changes, t_task_predecessor
    (both tasks and rows where the task is a predecessor). */

CREATE PROCEDURE sp_delete_tasks
(p_id_task bigint)
    LANGUAGE plpgsql
AS $$
BEGIN
DELETE FROM t_task WHERE id = p_id_task;
DELETE FROM t_task_changes WHERE id_task = p_id_task;
DELETE FROM t_task_predecessor WHERE id_task = p_id_task
                                  OR id_predecessor = p_id_task;
END
$$;