package task

import (
	"critical-path-analysis-api/internal/domain/task"
	"critical-path-analysis-api/internal/tests"
)

type taskStorage struct {
	//db client
}

func NewStorage() task.Storage {
	return &taskStorage{}
}

func (ts *taskStorage) GetOne(id int) (*task.Task, error) {
	tasks, err := tests.TestRepository.GetOne(id)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
func (ts *taskStorage) GetAll() *[]task.Task {
	return tests.TestRepository.GetAll()
}
func (ts *taskStorage) Create(tasks *[]task.Task) *[]task.Task {
	return tests.TestRepository.Create(tasks)
}
func (ts *taskStorage) Delete(id int) error {
	return nil
}
func (ts *taskStorage) Update(task *task.Task) error {
	return nil
}
