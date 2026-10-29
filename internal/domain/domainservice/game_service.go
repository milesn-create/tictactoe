package domainservice

import (
	"context"
	"tictactoe/domain/domainmodel"

	"github.com/google/uuid"
)

type GameService interface {
	CreateGame(ctx context.Context, userID uuid.UUID, vsComputer bool) (domainmodel.Game, error)
	GetAvailable(ctx context.Context, userID uuid.UUID) ([]domainmodel.Game, error)
	JoinGame(ctx context.Context, gameID uuid.UUID, userID uuid.UUID) (domainmodel.Game, error)
	MakeMove(ctx context.Context, gameID uuid.UUID, userID uuid.UUID, board domainmodel.Board) (domainmodel.Game, error)
	GetGame(ctx context.Context, gameID uuid.UUID, userID uuid.UUID) (domainmodel.Game, error)
}
