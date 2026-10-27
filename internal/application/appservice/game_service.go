package appservice

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"
	"tictactoe/domain/domainrepository"
	"tictactoe/domain/domainrule"
	"tictactoe/domain/domainservice"

	"github.com/google/uuid"
)

type gameService struct {
	repo domainrepository.GameRepository
}

func NewGameService(repo domainrepository.GameRepository) domainservice.Service {
	return &gameService{repo: repo}

}

func (s *gameService) NextMove(game domainmodel.Game) (domainmodel.Game, error) {
	return domainrule.NextMove(game)
}

func (s *gameService) ValidateBoard(ctx context.Context, game domainmodel.Game) error {
	oldGame, err := s.repo.Get(ctx, game.ID)
	if err != nil {
		if !errors.Is(err, domainerror.ErrGameNotFound) {
			return err
		}
		oldGame = domainmodel.Game{ID: game.ID}
	}
	return domainrule.ValidateMove(oldGame, game)
}
func (s *gameService) IsGameOver(game domainmodel.Game) (bool, domainmodel.Cell) {
	return domainrule.IsGameOver(game)
}
func (s *gameService) CreateGame(ctx context.Context, userID uuid.UUID, vsComputer bool) (domainmodel.Game, error) {
	var game domainmodel.Game
	game.ID = uuid.New()
	game.CrossPlayerID = userID
	game.VsComputer = vsComputer
	game = domainrule.UpdateState(game)
	err := s.repo.Save(ctx, game)
	if err != nil {
		return domainmodel.Game{}, fmt.Errorf("failed create game by save: %w", err)
	}
	return game, nil
}
