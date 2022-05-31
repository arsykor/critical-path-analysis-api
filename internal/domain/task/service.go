package task

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) GetById(id int) (*Task, error) {
	tasks, err := s.storage.GetOne(id)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *service) GetAll() *[]Task {
	return s.storage.GetAll()
}

func (s *service) Create(tasks *[]Task) *[]Task {
	return s.storage.Create(tasks)
}

func (s *service) Delete(id int) (*[]Task, error) {
	tasks, err := s.storage.Delete(id)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *service) Update(task *Task) (*[]Task, error) {
	tasks, err := s.storage.Update(task)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
