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
	saveGameQuery = `INSERT INTO games(id, board,state,turn_player_id,winner_id,vs_computer,cross_player_id,nought_player_id,created_at)
					 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
					 ON CONFLICT (id)
					 DO UPDATE SET board = EXCLUDED.board,state = EXCLUDED.state,turn_player_id = EXCLUDED.turn_player_id,winner_id = EXCLUDED.winner_id,vs_computer = EXCLUDED.vs_computer,cross_player_id = EXCLUDED.cross_player_id,nought_player_id = EXCLUDED.nought_player_id `
	getGameQuery = `SELECT id, board,state,turn_player_id,winner_id,vs_computer,cross_player_id,nought_player_id,created_at
					FROM games 
					WHERE id = $1`
	getAvailableQuery = `SELECT id, board,state,turn_player_id,winner_id,vs_computer,cross_player_id,nought_player_id,created_at
						FROM games 
						WHERE state = $1 AND vs_computer= false AND cross_player_id <> $2`
	getCompletedQuery = `SELECT id, board,state,turn_player_id,winner_id,vs_computer,cross_player_id,nought_player_id,created_at
						FROM games
						WHERE (cross_player_id=$1 OR nought_player_id=$1) AND (state = $2 OR state = $3)
						ORDER BY created_at DESC NULLS LAST`
	getLeaderboardQuery = `WITH 
							results AS (
								SELECT 
						   			cross_player_id AS user_id,
						   			CASE WHEN cross_player_id = winner_id THEN 1 ELSE 0 END AS is_win
						   		FROM games
						   		WHERE state IN (2,3) AND cross_player_id IS NOT NULL
						   		UNION ALL
						   		SELECT 
						   			nought_player_id AS user_id,
									CASE WHEN nought_player_id = winner_id THEN 1 ELSE 0 END AS is_win
						   		FROM games
						   		WHERE state IN(2, 3) AND nought_player_id IS NOT NULL),
							stats AS(
								SELECT
								 	user_id, 
								 	sum(is_win) AS wins , 
								 	count(*) AS total
								FROM results
								GROUP BY user_id)
						SELECT 
							user_id,
							users.login AS login,
							wins::float8 / total AS win_rate
						FROM stats
							JOIN users ON users.id = stats.user_id 
						ORDER BY win_rate DESC LIMIT $1`
)

func (r *gameRepository) Save(ctx context.Context, game domainmodel.Game) error {
	saveGame := dsmapper.GameToDatasource(game)
	board, err := json.Marshal(saveGame.Board)
	if err != nil {
		return fmt.Errorf("marshal board: %w", err)
	}
	_, err = r.pool.Exec(ctx, saveGameQuery, saveGame.ID, board, saveGame.State, saveGame.TurnPlayerID, saveGame.WinnerID, saveGame.VsComputer, saveGame.CrossPlayerID, saveGame.NoughtPlayerID, saveGame.CreatedAt)
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
	return dsmapper.GameToDomain(game)
}
func (r *gameRepository) GetAvailable(ctx context.Context, userID uuid.UUID) ([]domainmodel.Game, error) {
	return r.queryGetGames(ctx, getAvailableQuery, int(domainmodel.StateWaitingPlayers), userID)
}

func (r *gameRepository) GetCompleted(ctx context.Context, userID uuid.UUID) ([]domainmodel.Game, error) {
	return r.queryGetGames(ctx, getCompletedQuery, userID, int(domainmodel.StateDraw), int(domainmodel.StateWon))
}

func (r *gameRepository) queryGetGames(ctx context.Context, query string, args ...any) ([]domainmodel.Game, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return []domainmodel.Game{}, fmt.Errorf("failed  getQuery: %w", err)
	}
	games, err := pgx.CollectRows(rows, pgx.RowToStructByName[dsmodel.Game])
	if err != nil {
		return []domainmodel.Game{}, fmt.Errorf("get games: %w", err)
	}
	domainGames := make([]domainmodel.Game, 0, len(games))
	for i := range len(games) {
		domainGame, err := dsmapper.GameToDomain(games[i])
		if err != nil {
			return []domainmodel.Game{}, fmt.Errorf("map game to domain: %w", err)
		}
		domainGames = append(domainGames, domainGame)
	}
	return domainGames, nil

}
func (r *gameRepository) GetLeaderboard(ctx context.Context, limit int) ([]domainmodel.LeaderboardEntry, error) {
	rows, err := r.pool.Query(ctx, getLeaderboardQuery, limit)
	if err != nil {
		return []domainmodel.LeaderboardEntry{}, fmt.Errorf("failed getLeaderboardQuery: %w", err)
	}
	entries, err := pgx.CollectRows(rows, pgx.RowToStructByName[dsmodel.LeaderboardRow])
	if err != nil {
		return []domainmodel.LeaderboardEntry{}, fmt.Errorf("get leaderboard: %w", err)
	}
	return dsmapper.LeaderboardListToDomain(entries), nil
}
