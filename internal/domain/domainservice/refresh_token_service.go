package domainservice

import (
	"context"

	"github.com/google/uuid"
)

type RefreshTokenService interface {
	SaveNewJti(ctx context.Context, userID uuid.UUID, jti uuid.UUID) error
	MarkUsedByJti(ctx context.Context, jti uuid.UUID) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
}
