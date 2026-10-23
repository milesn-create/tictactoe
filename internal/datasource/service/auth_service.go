package service

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	users domain.UserService
}

func NewAuthService(users domain.UserService) domain.AuthService {
	return &authService{users: users}
}
func (s *authService) SignUp(ctx context.Context, request domain.SignUpRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed generate hash from password: %w", err)
	}
	id := uuid.New()
	user := domain.User{ID: id, Login: request.Login, PasswordHash: string(passwordHash)}
	err = s.users.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed create user: %w", err)
	}
	return nil
}
func (s *authService) Authorize(ctx context.Context, login, password string) (uuid.UUID, error) {
	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return uuid.Nil, domain.ErrInvalidCredentials
		}
		return uuid.Nil, fmt.Errorf("failed get by login: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return uuid.Nil, domain.ErrInvalidCredentials
		}
		return uuid.Nil, fmt.Errorf("compare hash and password: %w", err)
	}
	return user.ID, nil
}
