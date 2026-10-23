package service

import (
	"context"
	"tictactoe/datasource/repository"
	"tictactoe/domain"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) domain.UserService {
	return &userService{repo: repo}
}
func (s *userService) Create(ctx context.Context, user domain.User) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	return s.repo.GetByLogin(ctx, login)
}
