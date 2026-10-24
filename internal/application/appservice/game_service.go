package appservice

import (
	"context"
	"errors"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"
	"tictactoe/domain/domainrepository"
	"tictactoe/domain/domainrule"
	"tictactoe/domain/domainservice"
)

type GameService struct {
	repo domainrepository.GameRepository
}

func NewGameService(repo domainrepository.GameRepository) domainservice.Service {
	return &GameService{repo: repo}

}

func (s *GameService) NextMove(game domainmodel.Game) (domainmodel.Game, error) {
	return domainrule.NextMove(game)
}

func (s *GameService) ValidateBoard(ctx context.Context, game domainmodel.Game) error {
	oldGame, err := s.repo.Get(ctx, game.ID)
	if err != nil {
		if !errors.Is(err, domainerror.ErrGameNotFound) {
			return err
		}
		oldGame = domainmodel.Game{ID: game.ID}
	}
	if over, _ := s.IsGameOver(oldGame); over {
		return domainerror.ErrGameOver
	}
	newMoves := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if oldGame.Board[i][j] != domainmodel.Empty && oldGame.Board[i][j] != game.Board[i][j] {
				return domainerror.ErrBoardModified
			}

			if oldGame.Board[i][j] == domainmodel.Empty && game.Board[i][j] != domainmodel.Empty {
				if game.Board[i][j] != domainmodel.Cross {
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
func (s *GameService) IsGameOver(game domainmodel.Game) (bool, domainmodel.Cell) {
	return domainrule.IsGameOver(game)
}
