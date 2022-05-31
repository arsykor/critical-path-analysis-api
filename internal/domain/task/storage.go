package task

type Storage interface {
	GetOne(id int) (*Task, error)
	GetAll() *[]Task
	Create(task *[]Task) *[]Task
	Delete(id int) (*[]Task, error)
	Update(task *Task) (*[]Task, error)
}
