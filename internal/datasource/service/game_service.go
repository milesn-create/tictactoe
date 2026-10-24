package service

import (
	"context"
	"errors"
	"tictactoe/domain"
)

var (
	ErrBoardFull = errors.New("board is full")
)

type gameService struct {
	repo domain.GameRepository
}

func NewService(repo domain.GameRepository) domain.Service {
	return &gameService{repo: repo}

}

func (s *gameService) NextMove(game domain.Game) (domain.Game, error) {
	over, _ := s.IsGameOver(game)
	if over {
		return domain.Game{}, domain.ErrGameOver
	}
	bestI := -1
	bestJ := -1
	bestScore := -2
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if game.Board[i][j] == domain.Empty {
				game.Board[i][j] = domain.Nought
				score := minimax(game.Board, false)
				game.Board[i][j] = domain.Empty
				if score > bestScore {
					bestScore = score
					bestI = i
					bestJ = j
				}
			}
		}
	}
	if bestI == -1 {
		return domain.Game{}, ErrBoardFull
	}
	game.Board[bestI][bestJ] = domain.Nought
	return game, nil
}

func (s *gameService) ValidateBoard(ctx context.Context, game domain.Game) error {
	oldGame, err := s.repo.Get(ctx, game.ID)
	if err != nil {
		if !errors.Is(err, domain.ErrGameNotFound) {
			return err
		}
		oldGame = domain.Game{ID: game.ID}
	}
	if over, _ := s.IsGameOver(oldGame); over {
		return domain.ErrGameOver
	}
	newMoves := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if oldGame.Board[i][j] != domain.Empty && oldGame.Board[i][j] != game.Board[i][j] {
				return domain.ErrBoardModified
			}

			if oldGame.Board[i][j] == domain.Empty && game.Board[i][j] != domain.Empty {
				if game.Board[i][j] != domain.Cross {
					return domain.ErrBoardModified
				}
				newMoves++
			}
		}
	}
	if newMoves != 1 {
		return domain.ErrInvalidMovesCount
	}
	return nil

}
func (s *gameService) IsGameOver(game domain.Game) (bool, domain.Cell) {
	winner := findWinner(game.Board)
	if winner != domain.Empty {
		return true, winner
	}
	return isBoardFull(game.Board), domain.Empty
}
func findWinner(board domain.Board) domain.Cell {
	for i := 0; i < 3; i++ {
		if board[i][0] == board[i][1] && board[i][1] == board[i][2] && board[i][2] != domain.Empty {
			return board[i][0]
		}
		if board[0][i] == board[1][i] && board[1][i] == board[2][i] && board[2][i] != domain.Empty {
			return board[0][i]
		}
	}
	if (board[0][0] == board[1][1] && board[1][1] == board[2][2] && board[0][0] != domain.Empty) ||
		(board[0][2] == board[1][1] && board[1][1] == board[2][0] && board[0][2] != domain.Empty) {
		return board[1][1]
	}
	return domain.Empty
}
func isBoardFull(board domain.Board) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == domain.Empty {
				return false
			}
		}
	}
	return true

}
func minimax(board domain.Board, isComputerTurn bool) int {
	winner := findWinner(board)
	if winner == domain.Cross {
		return -1
	}
	if winner == domain.Nought {
		return 1
	}
	if isBoardFull(board) {
		return 0
	}
	if isComputerTurn {
		best := -2
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == domain.Empty {
					board[i][j] = domain.Nought
					score := minimax(board, false)
					board[i][j] = domain.Empty
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
			if board[i][j] == domain.Empty {
				board[i][j] = domain.Cross
				score := minimax(board, true)
				board[i][j] = domain.Empty
				if score < best {
					best = score
				}

			}
		}
	}
	return best
}
