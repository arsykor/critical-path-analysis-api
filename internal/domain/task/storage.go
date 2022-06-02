package task

type Storage interface {
	GetOne(id int) (*Task, error)
	GetAll() (*[]Task, error)
	Create(task *[]Task) (*[]Task, error)
	Delete(id int) (*[]Task, error)
	Update(task *Task) (*[]Task, error)
}
