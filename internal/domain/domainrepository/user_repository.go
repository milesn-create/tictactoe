package domainrepository

import (
	"context"
	"tictactoe/domain/domainmodel"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (domainmodel.User, error)
	Create(ctx context.Context, user domainmodel.User) error
}
