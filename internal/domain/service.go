package domain

type Service interface {
	NextMove(game Game) (Game, error)
	ValidateBoard(game Game) error
	IsGameOver(game Game) (bool, Cell)
}
