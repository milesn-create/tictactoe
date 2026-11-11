package dsrepository

import (
	"context"
	"fmt"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainrepository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	saveNewJtiQuery = `INSERT INTO refresh_tokens(jti, user_id)
					   VALUES ($1, $2)`
	markUsedByJtiQuery     = `UPDATE refresh_tokens SET used_at = now() WHERE jti = $1 AND used_at IS NULL`
	revokeAllByUserIDQuery = `DELETE FROM refresh_tokens WHERE user_id = $1`
)

type refreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) domainrepository.RefreshTokenRepository {
	return &refreshTokenRepository{pool: pool}
}
func (r *refreshTokenRepository) SaveNewJti(ctx context.Context, userID uuid.UUID, jti uuid.UUID) error {
	_, err := r.pool.Exec(ctx, saveNewJtiQuery, jti, userID)
	if err != nil {
		return fmt.Errorf("failed query saveNewJti: %w", err)
	}
	return nil
}
func (r *refreshTokenRepository) MarkUsedByJti(ctx context.Context, jti uuid.UUID) error {
	tags, err := r.pool.Exec(ctx, markUsedByJtiQuery, jti)
	if err != nil {
		return fmt.Errorf("failed query markUsedByJti: %w", err)
	}
	if tags.RowsAffected() == 0 {
		return domainerror.ErrRefreshTokenUsed
	}
	return nil
}
func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, revokeAllByUserIDQuery, userID)
	if err != nil {
		return fmt.Errorf("failed query revokeAllByUserID: %w", err)
	}
	return nil
}
