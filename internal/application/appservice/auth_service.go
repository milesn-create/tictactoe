package appservice

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"
	"tictactoe/domain/domainservice"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	users    domainservice.UserService
	provider domainservice.JwtProvider
}

func NewAuthService(users domainservice.UserService, provider domainservice.JwtProvider) domainservice.AuthService {
	return &authService{users: users, provider: provider}
}
func (s *authService) SignUp(ctx context.Context, request domainmodel.SignUpRequest) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed generate hash from password: %w", err)
	}
	id := uuid.New()
	user := domainmodel.User{ID: id, Login: request.Login, PasswordHash: string(passwordHash)}
	err = s.users.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed create user: %w", err)
	}
	return nil
}
func (s *authService) Authorize(ctx context.Context, request domainmodel.JwtRequest) (domainmodel.JwtResponse, error) {
	user, err := s.users.GetByLogin(ctx, request.Login)
	if err != nil {
		if errors.Is(err, domainerror.ErrUserNotFound) {
			return domainmodel.JwtResponse{}, domainerror.ErrInvalidCredentials
		}
		return domainmodel.JwtResponse{}, fmt.Errorf("failed get by login: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return domainmodel.JwtResponse{}, domainerror.ErrInvalidCredentials
		}
		return domainmodel.JwtResponse{}, fmt.Errorf("compare hash and password: %w", err)
	}
	accessToken, err := s.provider.GenerateAccessToken(user.ID)
	if err != nil {
		return domainmodel.JwtResponse{}, fmt.Errorf("generate access toke by authorize: %w", err)
	}
	refreshToken, err := s.provider.GenerateRefreshToken(user.ID)
	if err != nil {
		return domainmodel.JwtResponse{}, fmt.Errorf("generate refresh toke by authorize: %w", err)
	}
	return domainmodel.JwtResponse{Type: "Bearer", AccessToken: accessToken, RefreshToken: refreshToken}, nil

}

func (s *authService) UpdateAccessToken(request domainmodel.RefreshJwtRequest) (domainmodel.JwtResponse, error) {
	return s.rotateTokens(request.RefreshToken)
}
func (s *authService) UpdateRefreshToken(request domainmodel.RefreshJwtRequest) (domainmodel.JwtResponse, error) {
	return s.rotateTokens(request.RefreshToken)

}
func (s *authService) rotateTokens(refreshToken string) (domainmodel.JwtResponse, error) {
	err := s.provider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	userID, err := s.provider.GetUUIDFromToken(refreshToken)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	newAccessToken, err := s.provider.GenerateAccessToken(userID)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	newRefreshToken, err := s.provider.GenerateRefreshToken(userID)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	return domainmodel.JwtResponse{Type: "Bearer", AccessToken: newAccessToken, RefreshToken: newRefreshToken}, nil

}
