package repository

import (
	"context"
	"tictactoe/domain"

	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Game, error)
	Save(ctx context.Context, game domain.Game) error
}
