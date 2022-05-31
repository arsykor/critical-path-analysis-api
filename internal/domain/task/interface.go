package task

type Service interface {
	GetById(int) (*Task, error)
	GetAll() *[]Task
	Create(task *[]Task) *[]Task
	Delete(id int) (*[]Task, error)
	Update(task *Task) (*[]Task, error)
}
