package domain

import (
	"context"

	"github.com/google/uuid"
)

type AuthService interface {
	SignUp(ctx context.Context, request SignUpRequest) error
	Authorize(ctx context.Context, login, password string) (uuid.UUID, error)
}
