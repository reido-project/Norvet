package profile

import (
	"api/modules/accounts/_shared"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]_shared.User, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByID(id uint) (*_shared.User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GetByUsername(username string) (*_shared.User, error) {
	return s.repo.FindByUsername(username)
}

func (s *Service) GetByEmail(email string) (*_shared.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *Service) Update(id uint, req *_shared.UpdateUserRequest) (*_shared.User, error) {
	user := &_shared.User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
	}

	err := s.repo.Update(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Delete(id uint) error {
	return s.repo.Delete(id)
}
