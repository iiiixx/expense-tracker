package service

import "expense_tracker/internal/repository"

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserServise(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}
