package service

import (
	"context"
	"expense_tracker/internal/model"
	"expense_tracker/internal/repository"
	"fmt"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserServise(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) UpdateUsername(ctx context.Context, userID int, input *model.UpdateUsernameInput) (*model.User, error) {
	if input.Username == "" {
		return nil, fmt.Errorf("service/user: username can't be empty")
	}

	exists, err := s.userRepository.IsExistsUser(ctx, userID)
	if !exists {
		return nil, fmt.Errorf("service/auth: user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("service/auth: can't found this user: %w", err)
	}

	updated, err := s.userRepository.UpdateUsername(ctx, userID, input)
	if err != nil {
		return nil, fmt.Errorf("service/user: can't update username: %w", err)
	}
	return updated, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	if err := s.userRepository.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("service/user: can't delete user: %w", err)
	}
	return nil
}

func (s *UserService) GetUserProfile(ctx context.Context, userID int) (*model.User, error) {
	user, err := s.userRepository.GetUserById(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service/user: can't get user profile: %w", err)
	}

	user.Password = ""
	return user, nil
}
