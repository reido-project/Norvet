package auth

import (
	"api/libs/nvt/nvtdberr"
	"api/libs/nvt/nvterr"
	"api/libs/nvt/nvthash"
	"errors"
	"time"
)

type service struct {
	repo *repository
}

func newService(repo *repository) *service {
	return &service{repo}
}

func (s *service) register(req *registerRequest) (*registerResponse, error) {
	hashedPassword, err := nvthash.BcryptHash(req.Password)

	if err != nil {
		return nil, err
	}

	user := &userModel{
		Username:    req.Username,
		ProfileName: req.ProfileName,
		Email:       req.Email,
		Password:    hashedPassword,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.create(user); err != nil {
		err, info := nvtdberr.Identify(err)

		if errors.Is(err, nvterr.ErrDuplicateValue) {
			return nil, nvterr.NewAlreadyRegistered(info.Column)
		}

		return nil, err
	}

	response := &registerResponse{
		ID:          user.ID,
		Username:    user.Username,
		ProfileName: user.ProfileName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
	}

	return response, nil
}
