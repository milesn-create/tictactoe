package appservice

import (
	"context"
	"tictactoe/domain/domainrepository"
	"tictactoe/domain/domainservice"

	"github.com/google/uuid"
)

type refreshTokenService struct {
	repo domainrepository.RefreshTokenRepository
}

func NewRefreshTokenService(repo domainrepository.RefreshTokenRepository) domainservice.RefreshTokenService {
	return &refreshTokenService{repo: repo}
}
func (s *refreshTokenService) SaveNewJti(ctx context.Context, userID uuid.UUID, jti uuid.UUID) error {
	return s.repo.SaveNewJti(ctx, userID, jti)
}
func (s *refreshTokenService) MarkUsedByJti(ctx context.Context, jti uuid.UUID) error {
	return s.repo.MarkUsedByJti(ctx, jti)
}
func (s *refreshTokenService) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	return s.repo.RevokeAllByUserID(ctx, userID)
}
