package service

import (
	"context"
	"tictactoe/domain"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}
func (s *userService) Create(ctx context.Context, user domain.User) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	return s.repo.GetByLogin(ctx, login)
}
