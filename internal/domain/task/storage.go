package task

type Storage interface {
	GetOne(id int) (*Task, error)
	GetAll() *[]Task
	Create(task *[]Task) *[]Task
	Delete(id int) error
	Update(task *Task) error
}
