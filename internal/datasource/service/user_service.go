package service

import (
	"context"
	"tictactoe/domain"
	"tictactoe/domain/domainmodel"
	"tictactoe/domain/domainrepository"
)

type userService struct {
	repo domainrepository.UserRepository
}

func NewUserService(repo domainrepository.UserRepository) domain.UserService {
	return &userService{repo: repo}
}
func (s *userService) Create(ctx context.Context, user domainmodel.User) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) GetByLogin(ctx context.Context, login string) (domainmodel.User, error) {
	return s.repo.GetByLogin(ctx, login)
}
