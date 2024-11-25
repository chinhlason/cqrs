package command

type IService interface {
	InsertUser(name, email string) error
	UpdateUser(id int64, name, email string) error

	InsertOrder(userId int64, product string) error
	UpdateOrder(orderId int64, product string) error
}

type Service struct {
	repository IRepository
}

func NewService(repository IRepository) IService {
	return &Service{repository}
}

func (s *Service) InsertUser(name, email string) error {
	return s.repository.InsertUser(name, email)
}

func (s *Service) UpdateUser(id int64, name, email string) error {
	return s.repository.UpdateUser(id, name, email)
}

func (s *Service) InsertOrder(userId int64, product string) error {
	return s.repository.InsertOrder(userId, product)
}

func (s *Service) UpdateOrder(orderId int64, product string) error {
	return s.repository.UpdateOrder(orderId, product)
}
