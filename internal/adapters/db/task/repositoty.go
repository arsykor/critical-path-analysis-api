package task

import (
	"critical-path-analysis-api/internal/domain/task"
	"critical-path-analysis-api/internal/tests"
)

type taskRepository struct{}

func NewRepository() task.Storage {
	return &taskRepository{}
}

func (ts *taskRepository) GetOne(id int) (*task.Task, error) {
	tasks, err := tests.TestRepository.GetOne(id)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
func (ts *taskRepository) GetAll() *[]task.Task {
	return tests.TestRepository.GetAll()
}
func (ts *taskRepository) Create(tasks *[]task.Task) *[]task.Task {
	return tests.TestRepository.Create(tasks)
}
func (ts *taskRepository) Delete(id int) (*[]task.Task, error) {
	return nil, nil
}
func (ts *taskRepository) Update(task *task.Task) (*[]task.Task, error) {
	return nil, nil
}
