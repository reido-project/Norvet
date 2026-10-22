package auth

import (
	"api/libs/hashing"
	"api/modules/accounts/_shared"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(req *_shared.CreateUserRequest) (*_shared.User, error) {
	hashedPassword, err := hashing.BcryptHash(req.Password)

	if err != nil {
		return nil, err
	}

	user := &_shared.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
