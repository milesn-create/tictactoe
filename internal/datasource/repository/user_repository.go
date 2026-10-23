package repository

import (
	"context"
	"tictactoe/domain"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (domain.User, error)
	Create(ctx context.Context, user domain.User) error
}
