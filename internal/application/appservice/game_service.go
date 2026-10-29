package appservice

import (
	"context"
	"fmt"
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
func (s *gameService) GetAvailable(ctx context.Context, userID uuid.UUID) ([]domainmodel.Game, error) {
	games, err := s.repo.GetAvailable(ctx, userID)
	if err != nil {
		return []domainmodel.Game{}, fmt.Errorf("get available to app service: %w", err)
	}
	return games, nil
}
func (s *gameService) JoinGame(ctx context.Context, gameID uuid.UUID, userID uuid.UUID) (domainmodel.Game, error) {
	game, err := s.repo.Get(ctx, gameID)
	if err != nil {
		return domainmodel.Game{}, err
	}
	gameDone, err := domainrule.JoinGame(game, userID)
	if err != nil {
		return domainmodel.Game{}, err
	}
	err = s.repo.Save(ctx, gameDone)
	if err != nil {
		return domainmodel.Game{}, fmt.Errorf("save joined game: %w", err)
	}
	return gameDone, nil

}

func (s *gameService) MakeMove(ctx context.Context, gameID uuid.UUID, userID uuid.UUID, board domainmodel.Board) (domainmodel.Game, error) {
	game, err := s.repo.Get(ctx, gameID)
	if err != nil {
		return domainmodel.Game{}, err
	}
	err = domainrule.ValidateTurn(game, userID)
	if err != nil {
		return domainmodel.Game{}, err
	}
	candidate := game
	candidate.Board = board
	err = domainrule.ValidateMove(game, candidate)
	if err != nil {
		return domainmodel.Game{}, err
	}
	game = domainrule.UpdateState(candidate)
	if game.VsComputer && game.State == domainmodel.StatePlayerTurn {
		game, err = domainrule.NextMove(game)
		if err != nil {
			return domainmodel.Game{}, fmt.Errorf("computer make move: %w", err)
		}
		game = domainrule.UpdateState(game)
	}
	err = s.repo.Save(ctx, game)
	if err != nil {
		return domainmodel.Game{}, fmt.Errorf("save game after move: %w", err)
	}
	return game, nil
}
func (s *gameService) GetGame(ctx context.Context, gameID uuid.UUID, userID uuid.UUID) (domainmodel.Game, error) {
	game, err := s.repo.Get(ctx, gameID)
	if err != nil {
		return domainmodel.Game{}, err
	}
	err = domainrule.ValidateView(game, userID)
	if err != nil {
		return domainmodel.Game{}, err
	}
	return game, nil
}
