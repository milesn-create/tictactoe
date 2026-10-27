package domainrepository

import (
	"context"
	"tictactoe/domain/domainmodel"

	"github.com/google/uuid"
)

type GameRepository interface {
	Get(ctx context.Context, id uuid.UUID) (domainmodel.Game, error)
	Save(ctx context.Context, game domainmodel.Game) error
	GetAvailable(ctx context.Context, userID uuid.UUID) ([]domainmodel.Game, error)
}
