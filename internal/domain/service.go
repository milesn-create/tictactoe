package domain

import "context"

type Service interface {
	NextMove(game Game) (Game, error)
	ValidateBoard(ctx context.Context, game Game) error
	IsGameOver(game Game) (bool, Cell)
}
