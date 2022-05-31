package task

import (
	"context"
	"critical-path-analysis-api/internal/domain/task"
	"critical-path-analysis-api/pkg/client/postgresql"
	"critical-path-analysis-api/pkg/cpa"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

const (
	selectAllQuery      = `SELECT id, name, start_date, end_date FROM t_task`
	selectOneQuery      = `SELECT id, name, start_date, end_date FROM t_task WHERE id = $1`
	callStoredProcedure = `CALL sp_insert_merge($1, $2, $3, $4)`
	deleteQuery         = `DELETE FROM t_task WHERE id = $1`
)

type taskStorage struct {
	client postgresql.Client
}

func NewStorage(postgresqlClient postgresql.Client) task.Storage {
	return &taskStorage{client: postgresqlClient}
}

func (ts *taskStorage) GetOne(id int) (*task.Task, error) {
	var t task.Task

	err := ts.client.QueryRow(context.TODO(), selectOneQuery, id).Scan(&t.Id, &t.Name, &t.StartDate, &t.EndDate)
	if err != nil {
		return nil, errors.New("there is no task with this id in database")
	}
	return &t, nil
}

func (ts *taskStorage) GetAll() *[]task.Task {
	rows, _ := ts.client.Query(context.Background(), selectAllQuery)

	tasks := make([]task.Task, 0)

	for rows.Next() {
		var t task.Task
		err := rows.Scan(&t.Id, &t.Name, &t.StartDate, &t.EndDate)
		if err != nil {
			fmt.Println("*** rows.Scan error: ", err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("*** iteration error: ", err)
	}
	return &tasks
}

func (ts *taskStorage) Create(newTasks *[]task.Task) *[]task.Task {
	/*
		For critical path calculation it is necessary to get all project
		tasks from the database and overwrite changed start
		and end dates, as well as add new ones.
	*/
	createdTasks := make([]task.Task, 0)
	oldTasks := ts.GetAll()

	for _, newT := range *newTasks {
		index := slices.IndexFunc(*oldTasks, func(t task.Task) bool { return t.Id == newT.Id })
		if index == -1 {
			createdTasks = append(createdTasks, newT)
		} else {
			(*oldTasks)[index] = newT
		}
	}

	for _, t := range createdTasks {
		*oldTasks = append(*oldTasks, t)
	}

	cpa.Arrange(oldTasks)

	for _, t := range *oldTasks {
		err := ts.client.QueryRow(context.TODO(), callStoredProcedure, t.Id, t.Name, t.StartDate, t.EndDate).Scan()
		if err != nil {
			pgError, ok := err.(*pgconn.PgError)
			if ok {
				fmt.Println("*** message err: ", pgError.Message)
			}
		}
	}
	return oldTasks
}

func (ts *taskStorage) Delete(id int) (*[]task.Task, error) {
	var task task.Task
	task.Id = id

	err := ts.client.QueryRow(context.TODO(), deleteQuery, task.Id).Scan()
	if err != nil {
		pgError, ok := err.(*pgconn.PgError)
		if ok {
			fmt.Println("database err: ", pgError.Message)
		}
	}
	allTasks := ts.GetAll()
	cpa.Arrange(allTasks)
	return allTasks, nil
}

func (ts *taskStorage) Update(task *task.Task) (*[]task.Task, error) {
	err := ts.client.QueryRow(context.TODO(), callStoredProcedure, (*task).Id, (*task).Name, (*task).StartDate, (*task).EndDate).Scan()
	if err != nil {
		pgError, ok := err.(*pgconn.PgError)
		if ok {
			fmt.Println("*** message err: ", pgError.Message)
		}
	}
	allTasks := ts.GetAll()
	cpa.Arrange(allTasks)
	return allTasks, nil
}
