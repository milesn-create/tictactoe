package dsrepository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"tictactoe/datasource/dsmapper"
	"tictactoe/datasource/dsmodel"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"
	"tictactoe/domain/domainrepository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type gameRepository struct {
	pool *pgxpool.Pool
}

func NewGameRepository(pool *pgxpool.Pool) domainrepository.GameRepository {
	return &gameRepository{pool: pool}
}

const (
	saveGameQuery = `INSERT INTO games(id, board,state,turn_player_id,winner_id,vs_computer,cross_player_id,nought_player_id)
					 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
					 ON CONFLICT (id)
					 DO UPDATE SET board = EXCLUDED.board,state = EXCLUDED.state,turn_player_id = EXCLUDED.turn_player_id,winner_id = EXCLUDED.winner_id,vs_computer = EXCLUDED.vs_computer,cross_player_id = EXCLUDED.cross_player_id,nought_player_id = EXCLUDED.nought_player_id `
	getGameQuery = `SELECT id, board,state,turn_player_id,winner_id,vs_computer,cross_player_id,nought_player_id
					FROM games 
					WHERE id = $1`
	getAvailableQuery = `SELECT id, board,state,turn_player_id,winner_id,vs_computer,cross_player_id,nought_player_id
						FROM games 
						WHERE state = $1 AND vs_computer= false AND cross_player_id <> $2`
)

func (r *gameRepository) Save(ctx context.Context, game domainmodel.Game) error {
	saveGame := dsmapper.ToDatasource(game)
	board, err := json.Marshal(saveGame.Board)
	if err != nil {
		return fmt.Errorf("marshal board: %w", err)
	}
	_, err = r.pool.Exec(ctx, saveGameQuery, saveGame.ID, board, saveGame.State, saveGame.TurnPlayerID, saveGame.WinnerID, saveGame.VsComputer, saveGame.CrossPlayerID, saveGame.NoughtPlayerID)
	if err != nil {
		return fmt.Errorf("save game: %w", err)
	}
	return nil
}

func (r *gameRepository) Get(ctx context.Context, id uuid.UUID) (domainmodel.Game, error) {
	rows, err := r.pool.Query(ctx, getGameQuery, id)
	if err != nil {
		return domainmodel.Game{}, fmt.Errorf("failed getGameQuery: %w", err)

	}
	game, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dsmodel.Game])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainmodel.Game{}, domainerror.ErrGameNotFound
		}
		return domainmodel.Game{}, fmt.Errorf("get game : %w", err)
	}
	return dsmapper.ToDomain(game)
}
func (r *gameRepository) GetAvailable(ctx context.Context, userID uuid.UUID) ([]domainmodel.Game, error) {
	rows, err := r.pool.Query(ctx, getAvailableQuery, int(domainmodel.StateWaitingPlayers), userID)
	if err != nil {
		return []domainmodel.Game{}, fmt.Errorf("failed getAvailableQuery: %w", err)
	}
	games, err := pgx.CollectRows(rows, pgx.RowToStructByName[dsmodel.Game])
	if err != nil {
		return []domainmodel.Game{}, fmt.Errorf("get available game: %w", err)
	}
	domainGames := make([]domainmodel.Game, 0, len(games))
	for i := range len(games) {
		domainGame, err := dsmapper.ToDomain(games[i])
		if err != nil {
			return []domainmodel.Game{}, fmt.Errorf("map game to domain: %w", err)
		}
		domainGames = append(domainGames, domainGame)
	}
	return domainGames, nil

}
