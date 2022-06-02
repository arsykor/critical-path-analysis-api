package task

import (
	"context"
	"critical-path-analysis-api/internal/domain/task"
	"critical-path-analysis-api/pkg/client/postgresql"
	"critical-path-analysis-api/pkg/cpa"
	"golang.org/x/exp/slices"
	"sync"
)

const (
	selectAllQuery           = `SELECT id, name, start_date, end_date FROM t_task`
	selectOneQuery           = `SELECT id, name, start_date, end_date FROM t_task WHERE id = $1`
	selectPredecessors       = `SELECT id_predecessor FROM t_task_predecessor WHERE id_task = $1`
	callSPInsertTasks        = `CALL sp_insert_merge($1, $2, $3, $4)`
	callSPInsertPredecessors = `CALL sp_predecessors($1, $2);`
	callSPDeleteTask         = `CALL sp_delete_tasks($1);`
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
	if err != nil && err.Error() != "no rows in result set" {
		return nil, err
	}

	rowsPred, err := ts.client.Query(context.TODO(), selectPredecessors, t.Id)
	if err != nil {
		return nil, err
	}
	for rowsPred.Next() {
		var pred int
		err = rowsPred.Scan(&pred)
		if err != nil {
			return nil, err
		}
		t.Predecessors = append(t.Predecessors, pred)
	}
	return &t, nil
}

func (ts *taskStorage) GetAll() (*[]task.Task, error) {
	tasks := make([]task.Task, 0)

	rowsTasks, err := ts.client.Query(context.TODO(), selectAllQuery)
	if err != nil {
		return nil, err
	}

	for rowsTasks.Next() {
		var t task.Task
		err := rowsTasks.Scan(&t.Id, &t.Name, &t.StartDate, &t.EndDate)
		if err != nil {
			return nil, err
		}

		rowsPred, err := ts.client.Query(context.TODO(), selectPredecessors, t.Id)
		if err != nil {
			return nil, err
		}
		for rowsPred.Next() {
			var pred int
			err = rowsPred.Scan(&pred)
			if err != nil {
				return nil, err
			}
			t.Predecessors = append(t.Predecessors, pred)
		}
		if t.Predecessors == nil {
			t.Predecessors = make([]int, 0)
		}

		tasks = append(tasks, t)
	}

	return &tasks, nil
}

func (ts *taskStorage) Create(newTasks *[]task.Task) (*[]task.Task, error) {
	/*
		For critical path calculation it is necessary to get all project
		tasks from the database and overwrite changed start
		and end dates, as well as add new ones.
	*/
	createdTasks := make([]task.Task, 0)
	oldTasks, err := ts.GetAll()
	if err != nil {
		return nil, err
	}

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

	err = cpa.Arrange(oldTasks)
	if err != nil {
		return nil, err
	}

	err = ts.addTasksAsync(oldTasks)
	if err != nil {
		return nil, err
	}

	return oldTasks, nil
}

func (ts *taskStorage) Delete(id int) (*[]task.Task, error) {
	var task task.Task
	task.Id = id

	err := ts.client.QueryRow(context.TODO(), callSPDeleteTask, task.Id).Scan()
	if err != nil && err.Error() != "no rows in result set" {
		return nil, err
	}
	allTasks, err := ts.GetAll()
	if err != nil {
		return nil, err
	}

	err = cpa.Arrange(allTasks)
	if err != nil {
		return nil, err
	}
	return allTasks, nil
}

func (ts *taskStorage) Update(newTask *task.Task) (*[]task.Task, error) {

	oldTasks, err := ts.GetAll()
	if err != nil {
		return nil, err
	}

	index := slices.IndexFunc(*oldTasks, func(t task.Task) bool { return t.Id == newTask.Id })
	if index == -1 {
		*oldTasks = append(*oldTasks, *newTask)
	} else {
		(*oldTasks)[index] = *newTask
	}

	err = cpa.Arrange(oldTasks)
	if err != nil {
		return nil, err
	}

	err = ts.addTasksAsync(oldTasks)
	if err != nil {
		return nil, err
	}

	return oldTasks, nil
}

func (ts *taskStorage) writePredecessors(tasks *[]task.Task, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, t := range *tasks {
		if len(t.Predecessors) > 0 {
			err := ts.client.QueryRow(context.TODO(), callSPInsertPredecessors, t.Id, t.Predecessors).Scan()
			if err != nil && err.Error() != "no rows in result set" {
				select {
				case <-errCh:
					return
				default:
					errCh <- err
				}
				close(errCh)
			}
		}
	}
}

func (ts *taskStorage) writeTasks(tasks *[]task.Task, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, t := range *tasks {
		err := ts.client.QueryRow(context.TODO(), callSPInsertTasks, t.Id, t.Name, t.StartDate, t.EndDate).Scan()
		if err != nil && err.Error() != "no rows in result set" {
			select {
			case <-errCh:
				return
			default:
				errCh <- err
			}
			close(errCh)
		}
	}
}

func (ts *taskStorage) addTasksAsync(tasks *[]task.Task) error {
	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(2)
	go ts.writePredecessors(tasks, errCh, &wg)
	go ts.writeTasks(tasks, errCh, &wg)
	wg.Wait()

	if len(errCh) != 0 {
		err := <-errCh
		if err != nil {
			return err
		}
	}

	return nil
}
