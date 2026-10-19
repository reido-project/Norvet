package user

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]User, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByID(id uint) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GetByUsername(username string) (*User, error) {
	return s.repo.FindByUsername(username)
}

func (s *Service) GetByEmail(email string) (*User, error) {
	return s.repo.FindByEmail(email)
}

func (s *Service) Create(req *CreateUserRequest) (*User, error) {
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Update(id uint, req *UpdateUserRequest) (*User, error) {
	user := &User{
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
