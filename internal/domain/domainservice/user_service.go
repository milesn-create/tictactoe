package domainservice

import (
	"context"
	"tictactoe/domain/domainmodel"

	"github.com/google/uuid"
)

type UserService interface {
	GetByLogin(ctx context.Context, login string) (domainmodel.User, error)
	Create(ctx context.Context, user domainmodel.User) error
	GetByID(ctx context.Context, userID uuid.UUID) (domainmodel.User, error)
}
