package appservice

import (
	"context"
	"tictactoe/domain/domainmodel"
	"tictactoe/domain/domainrepository"
	"tictactoe/domain/domainservice"

	"github.com/google/uuid"
)

type userService struct {
	repo domainrepository.UserRepository
}

func NewUserService(repo domainrepository.UserRepository) domainservice.UserService {
	return &userService{repo: repo}
}
func (s *userService) Create(ctx context.Context, user domainmodel.User) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) GetByLogin(ctx context.Context, login string) (domainmodel.User, error) {
	return s.repo.GetByLogin(ctx, login)
}
func (s *userService) GetByID(ctx context.Context, userID uuid.UUID) (domainmodel.User, error) {
	return s.repo.GetByID(ctx, userID)
}
