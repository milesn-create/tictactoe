package domain

import (
	"context"
	"tictactoe/domain/domainmodel"
)

type Service interface {
	NextMove(game domainmodel.Game) (domainmodel.Game, error)
	ValidateBoard(ctx context.Context, game domainmodel.Game) error
	IsGameOver(game domainmodel.Game) (bool, domainmodel.Cell)
}
