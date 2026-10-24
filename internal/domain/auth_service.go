package domain

import (
	"context"
	"tictactoe/domain/domainmodel"

	"github.com/google/uuid"
)

type AuthService interface {
	SignUp(ctx context.Context, request domainmodel.SignUpRequest) error
	Authorize(ctx context.Context, login, password string) (uuid.UUID, error)
}
