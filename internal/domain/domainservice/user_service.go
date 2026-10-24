package domainservice

import (
	"context"
	"tictactoe/domain/domainmodel"
)

type UserService interface {
	GetByLogin(ctx context.Context, login string) (domainmodel.User, error)
	Create(ctx context.Context, user domainmodel.User) error
}
