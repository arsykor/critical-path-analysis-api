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

/* 3. Create t_task_predecessor table - predecessors of tasks (required for critical path analysis). */

CREATE TABLE IF NOT EXISTS t_task_predecessor
(
    id bigint NOT NULL,
    id_task bigint,
    id_predecessor bigint,
    CONSTRAINT t_task_predecessor_pkey PRIMARY KEY (id)
    )

/* 4. Stored procedure for inserting all tasks.  New tasks from the array are added, old ones are updated. */

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