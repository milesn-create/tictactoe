package domainrule

import (
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"
)

func findWinner(board domainmodel.Board) domainmodel.Cell {
	for i := 0; i < 3; i++ {
		if board[i][0] == board[i][1] && board[i][1] == board[i][2] && board[i][2] != domainmodel.Empty {
			return board[i][0]
		}
		if board[0][i] == board[1][i] && board[1][i] == board[2][i] && board[2][i] != domainmodel.Empty {
			return board[0][i]
		}
	}
	if (board[0][0] == board[1][1] && board[1][1] == board[2][2] && board[0][0] != domainmodel.Empty) ||
		(board[0][2] == board[1][1] && board[1][1] == board[2][0] && board[0][2] != domainmodel.Empty) {
		return board[1][1]
	}
	return domainmodel.Empty
}
func isBoardFull(board domainmodel.Board) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domainmodel.Empty {
				return false
			}
		}
	}
	return true

}
func minimax(board domainmodel.Board, isComputerTurn bool) int {
	winner := findWinner(board)
	if winner == domainmodel.Cross {
		return -1
	}
	if winner == domainmodel.Nought {
		return 1
	}
	if isBoardFull(board) {
		return 0
	}
	if isComputerTurn {
		best := -2
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == domainmodel.Empty {
					board[i][j] = domainmodel.Nought
					score := minimax(board, false)
					board[i][j] = domainmodel.Empty
					if score > best {
						best = score
					}

				}
			}
		}
		return best
	}
	best := 2
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domainmodel.Empty {
				board[i][j] = domainmodel.Cross
				score := minimax(board, true)
				board[i][j] = domainmodel.Empty
				if score < best {
					best = score
				}

			}
		}
	}
	return best
}
func NextMove(game domainmodel.Game) (domainmodel.Game, error) {
	over, _ := IsGameOver(game)
	if over {
		return domainmodel.Game{}, domainerror.ErrGameOver
	}
	bestI := -1
	bestJ := -1
	bestScore := -2
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if game.Board[i][j] == domainmodel.Empty {
				game.Board[i][j] = domainmodel.Nought
				score := minimax(game.Board, false)
				game.Board[i][j] = domainmodel.Empty
				if score > bestScore {
					bestScore = score
					bestI = i
					bestJ = j
				}
			}
		}
	}
	if bestI == -1 {
		return domainmodel.Game{}, domainerror.ErrBoardFull
	}
	game.Board[bestI][bestJ] = domainmodel.Nought
	return game, nil
}

func IsGameOver(game domainmodel.Game) (bool, domainmodel.Cell) {
	winner := findWinner(game.Board)
	if winner != domainmodel.Empty {
		return true, winner
	}
	return isBoardFull(game.Board), domainmodel.Empty
}
func ValidateMove(oldGame, newGame domainmodel.Game) error {
	if over, _ := IsGameOver(oldGame); over {
		return domainerror.ErrGameOver
	}
	newMoves := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if oldGame.Board[i][j] != domainmodel.Empty && oldGame.Board[i][j] != newGame.Board[i][j] {
				return domainerror.ErrBoardModified
			}

			if oldGame.Board[i][j] == domainmodel.Empty && newGame.Board[i][j] != domainmodel.Empty {
				if newGame.Board[i][j] != domainmodel.Cross {
					return domainerror.ErrBoardModified
				}
				newMoves++
			}
		}
	}
	if newMoves != 1 {
		return domainerror.ErrInvalidMovesCount
	}
	return nil

}
