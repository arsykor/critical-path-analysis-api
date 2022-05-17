package task

type Service interface {
	GetById(int) (*Task, error)
	GetAll() *[]Task
	Create(task *[]Task) *[]Task
	Delete(id int)
	Update(task *Task)
}
