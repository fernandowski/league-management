package service

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) RegisterUser() {
	// registration logic.
}
