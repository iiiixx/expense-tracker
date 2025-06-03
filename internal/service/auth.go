package service

import (
	"context"
	"expense_tracker/internal/model"
	"expense_tracker/internal/repository"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthService provides methods for authentication and authorization operations.
type AuthService struct {
	userRepository *repository.UserRepository
	jwtSecret      string
	tokenExpiry    time.Duration
}

// NewAuthService create an instance of AuthService.
func NewAuthService(userRepository *repository.UserRepository, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
		tokenExpiry:    tokenExpiry,
	}
}

// Register creates a new user account with hashed password.
func (s *AuthService) Register(ctx context.Context, user *model.User) error {
	exists, err := s.userRepository.IsExistsUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("service/auth: can`t check existence of user: %w", err)
	}
	if exists {
		return fmt.Errorf("service/auth: user is already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("service/auth: can't hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	if err := s.userRepository.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("service/auth: can't registrate user")
	}
	return nil
}

// Login authenticates a user and generates a JWT token.
func (s *AuthService) Login(ctx context.Context, input *model.LoginInput) (string, error) {
	user, err := s.userRepository.GetUserByName(ctx, input.Username)
	if err != nil {
		return "", fmt.Errorf("service/auth: wrong username: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", fmt.Errorf("service/auth: wrong password: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
	})

	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("service/auth: failed to sign token: %w", err)
	}
	return signedToken, nil
}

// ValidateToken verifies a JWT token and extracts user ID.
func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("service/auth: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("service/auth: can't validate token %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}
	return 0, fmt.Errorf("service/auth: invalid token")
}
