package domain

import "errors"

var (
	ErrGameNotFound      = errors.New("game not found")
	ErrBoardModified     = errors.New("previous moves were modified")
	ErrInvalidMovesCount = errors.New("exactly one new move is required")
	ErrGameOver          = errors.New("game over")
)
