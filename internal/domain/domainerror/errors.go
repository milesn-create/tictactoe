package domainerror

import "errors"

var (
	ErrGameNotFound       = errors.New("game not found")
	ErrBoardModified      = errors.New("previous moves were modified")
	ErrInvalidMovesCount  = errors.New("exactly one new move is required")
	ErrGameOver           = errors.New("game over")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("login already exists")
	ErrInvalidCredentials = errors.New("incorrect password or login")
)
