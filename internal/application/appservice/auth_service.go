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
	users        domainservice.UserService
	jwtProvider  domainservice.JwtProvider
	refreshToken domainservice.RefreshTokenService
}

func NewAuthService(users domainservice.UserService, jwt domainservice.JwtProvider, refreshToken domainservice.RefreshTokenService) domainservice.AuthService {
	return &authService{users: users, jwtProvider: jwt, refreshToken: refreshToken}
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
	accessToken, err := s.jwtProvider.GenerateAccessToken(user.ID)
	if err != nil {
		return domainmodel.JwtResponse{}, fmt.Errorf("generate access toke by authorize: %w", err)
	}
	refreshToken, err := s.jwtProvider.GenerateRefreshToken(user.ID)
	if err != nil {
		return domainmodel.JwtResponse{}, fmt.Errorf("generate refresh token by authorize: %w", err)
	}
	jtiStr, err := s.jwtProvider.GetJtiFromRefreshToken(refreshToken)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	jti, err := uuid.Parse(jtiStr)
	if err != nil {
		return domainmodel.JwtResponse{}, domainerror.ErrInvalidJti
	}
	err = s.refreshToken.SaveNewJti(ctx, user.ID, jti)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	return domainmodel.JwtResponse{Type: "Bearer", AccessToken: accessToken, RefreshToken: refreshToken}, nil

}

func (s *authService) UpdateAccessToken(ctx context.Context, request domainmodel.RefreshJwtRequest) (domainmodel.JwtResponse, error) {
	return s.rotateTokens(ctx, request.RefreshToken)
}
func (s *authService) UpdateRefreshToken(ctx context.Context, request domainmodel.RefreshJwtRequest) (domainmodel.JwtResponse, error) {
	return s.rotateTokens(ctx, request.RefreshToken)

}
func (s *authService) rotateTokens(ctx context.Context, refreshToken string) (domainmodel.JwtResponse, error) {
	err := s.jwtProvider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	userID, err := s.jwtProvider.GetUUIDFromToken(refreshToken)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	jtiStr, err := s.jwtProvider.GetJtiFromRefreshToken(refreshToken)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	jti, err := uuid.Parse(jtiStr)
	if err != nil {
		return domainmodel.JwtResponse{}, domainerror.ErrInvalidJti
	}
	err = s.refreshToken.MarkUsedByJti(ctx, jti)
	if err != nil {
		if errors.Is(err, domainerror.ErrRefreshTokenUsed) {
			if revokeErr := s.refreshToken.RevokeAllByUserID(ctx, userID); revokeErr != nil {
				return domainmodel.JwtResponse{}, revokeErr
			}
		}
		return domainmodel.JwtResponse{}, err
	}
	newAccessToken, err := s.jwtProvider.GenerateAccessToken(userID)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	newRefreshToken, err := s.jwtProvider.GenerateRefreshToken(userID)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	newJtiStr, err := s.jwtProvider.GetJtiFromRefreshToken(newRefreshToken)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	newJti, err := uuid.Parse(newJtiStr)
	if err != nil {
		return domainmodel.JwtResponse{}, domainerror.ErrInvalidJti
	}
	err = s.refreshToken.SaveNewJti(ctx, userID, newJti)
	if err != nil {
		return domainmodel.JwtResponse{}, err
	}
	return domainmodel.JwtResponse{Type: "Bearer", AccessToken: newAccessToken, RefreshToken: newRefreshToken}, nil
}
