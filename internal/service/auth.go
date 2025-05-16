package service

import (
	"context"
	"expense_tracker/internal/model"
	"expense_tracker/internal/repository"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository *repository.UserRepository
	jwtSecret      string
	tokenExpiry    time.Duration
}

func NewAuthService(userRepository *repository.UserRepository, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
		tokenExpiry:    tokenExpiry,
	}
}

func (s *AuthService) Register(ctx context.Context, user *model.User) error {
	exists, err := s.userRepository.IsExistsUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("service/auth: can`t check existence of user: %w", err)
	}
	if exists {
		return fmt.Errorf("service/auth: user is already exists")
	}
	log.Printf("Пароль перед хешированием: %s", user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	log.Printf("Хеш при генерации: %s", hashedPassword)
	if err != nil {
		return fmt.Errorf("service/auth: can't hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	if err := s.userRepository.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("service/auth: can't registrate user")
	}
	return nil
}

func (s *AuthService) Login(ctx context.Context, input *model.LoginInput) (string, error) {
	user, err := s.userRepository.GetUserByName(ctx, input.Username)
	if err != nil {
		return "", fmt.Errorf("service/auth: wrong username: %w", err)
	}

	log.Printf("Хеш из базы: %s", user.Password)
	log.Printf("Введенный пароль: %s", input.Password)

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
