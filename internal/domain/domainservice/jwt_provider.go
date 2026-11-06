package domainservice

import (
	"github.com/google/uuid"
)

type JwtProvider interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ValidateAccessToken(token string) error
	ValidateRefreshToken(token string) error
	GetUUIDFromToken(token string) (uuid.UUID, error)
}
