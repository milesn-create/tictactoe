package service

import (
	"context"
	"tictactoe/domain"
	"tictactoe/domain/domainmodel"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}
func (s *userService) Create(ctx context.Context, user domainmodel.User) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) GetByLogin(ctx context.Context, login string) (domainmodel.User, error) {
	return s.repo.GetByLogin(ctx, login)
}
