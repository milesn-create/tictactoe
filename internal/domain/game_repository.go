package domain

import (
	"context"

	"github.com/google/uuid"
)

type GameRepository interface {
	Get(ctx context.Context, id uuid.UUID) (Game, error)
	Save(ctx context.Context, game Game) error
}
