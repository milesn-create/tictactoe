package domainrule

import (
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"

	"github.com/google/uuid"
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
func minimax(board domainmodel.Board, turn domainmodel.Cell, computerBadge, playerBadge domainmodel.Cell) int {
	winner := findWinner(board)
	if winner == playerBadge {
		return -1
	}
	if winner == computerBadge {
		return 1
	}
	if isBoardFull(board) {
		return 0
	}
	if turn == computerBadge {
		best := -2
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == domainmodel.Empty {
					board[i][j] = computerBadge
					score := minimax(board, playerBadge, computerBadge, playerBadge)
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
				board[i][j] = playerBadge
				score := minimax(board, computerBadge, computerBadge, playerBadge)
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
	computerBadge := nextBadge(game.Board)
	playerBadge := -computerBadge
	bestI := -1
	bestJ := -1
	bestScore := -2
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if game.Board[i][j] == domainmodel.Empty {
				game.Board[i][j] = computerBadge
				score := minimax(game.Board, playerBadge, computerBadge, playerBadge)
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
	game.Board[bestI][bestJ] = computerBadge
	return game, nil
}

func IsGameOver(game domainmodel.Game) (bool, domainmodel.Cell) {
	winner := findWinner(game.Board)
	if winner != domainmodel.Empty {
		return true, winner
	}
	return isBoardFull(game.Board), domainmodel.Empty
}
func nextBadge(board domainmodel.Board) domainmodel.Cell {
	countX := 0
	countO := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domainmodel.Cross {
				countX++
			}
			if board[i][j] == domainmodel.Nought {
				countO++
			}
		}
	}
	if countO == countX {
		return domainmodel.Cross
	}
	return domainmodel.Nought
}
func playerByBadge(game domainmodel.Game, badge domainmodel.Cell) uuid.UUID {
	if badge == domainmodel.Cross {
		return game.CrossPlayerID
	}
	if badge == domainmodel.Nought {
		return game.NoughtPlayerID
	}
	return uuid.Nil
}
func UpdateState(game domainmodel.Game) domainmodel.Game {
	if game.State == domainmodel.StateArchived {
		return game
	}
	game.TurnPlayerID = uuid.Nil
	game.WinnerID = uuid.Nil
	if !game.VsComputer && (game.CrossPlayerID == uuid.Nil || game.NoughtPlayerID == uuid.Nil) {
		game.State = domainmodel.StateWaitingPlayers
	} else if w := findWinner(game.Board); w != domainmodel.Empty {
		game.State = domainmodel.StateWon
		game.WinnerID = playerByBadge(game, w)
	} else if isBoardFull(game.Board) {
		game.State = domainmodel.StateDraw
	} else {
		game.State = domainmodel.StatePlayerTurn
		game.TurnPlayerID = playerByBadge(game, nextBadge(game.Board))
	}
	return game
}
func ValidateMove(oldGame, newGame domainmodel.Game) error {
	if oldGame.State == domainmodel.StateArchived {
		return domainerror.ErrGameOver
	}
	if over, _ := IsGameOver(oldGame); over {
		return domainerror.ErrGameOver
	}
	newMoves := 0
	expected := nextBadge(oldGame.Board)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if oldGame.Board[i][j] != domainmodel.Empty && oldGame.Board[i][j] != newGame.Board[i][j] {
				return domainerror.ErrBoardModified
			}

			if oldGame.Board[i][j] == domainmodel.Empty && newGame.Board[i][j] != domainmodel.Empty {
				if newGame.Board[i][j] != expected {
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
func JoinGame(game domainmodel.Game, userID uuid.UUID) (domainmodel.Game, error) {
	if game.State != domainmodel.StateWaitingPlayers {
		return domainmodel.Game{}, domainerror.ErrGameNotJoinable
	}
	if game.CrossPlayerID == userID {
		return domainmodel.Game{}, domainerror.ErrJoinOwnGame
	}
	game.NoughtPlayerID = userID
	return UpdateState(game), nil
}
func ValidateTurn(game domainmodel.Game, userID uuid.UUID) error {
	if userID != game.CrossPlayerID && userID != game.NoughtPlayerID {
		return domainerror.ErrNotParticipant
	}
	if game.State == domainmodel.StateWon || game.State == domainmodel.StateDraw || game.State == domainmodel.StateArchived {
		return domainerror.ErrGameOver
	}
	if game.State == domainmodel.StateWaitingPlayers {
		return domainerror.ErrGameNotStarted
	}
	if userID != game.TurnPlayerID {
		return domainerror.ErrNotYourTurn
	}
	return nil
}
