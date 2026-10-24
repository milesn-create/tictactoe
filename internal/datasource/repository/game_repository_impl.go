package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"tictactoe/datasource/mapper"
	"tictactoe/datasource/model"
	"tictactoe/domain"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type gameRepository struct {
	pool *pgxpool.Pool
}

func NewGameRepository(pool *pgxpool.Pool) domain.GameRepository {
	return &gameRepository{pool: pool}
}

const (
	saveGameQuery = `INSERT INTO games(id, board)
					 VALUES ($1, $2)
					 ON CONFLICT (id)
					 DO UPDATE SET board= EXCLUDED.board`
	getGameQuery = `SELECT board 
					FROM games 
					WHERE id = $1`
)

func (r *gameRepository) Save(ctx context.Context, game domainmodel.Game) error {
	saveGame := mapper.ToDatasource(game)
	board, err := json.Marshal(saveGame.Board)
	if err != nil {
		return fmt.Errorf("marshal board: %w", err)
	}
	_, err = r.pool.Exec(ctx, saveGameQuery, saveGame.ID, board)
	if err != nil {
		return fmt.Errorf("save game: %w", err)
	}
	return nil
}

func (r *gameRepository) Get(ctx context.Context, id uuid.UUID) (domainmodel.Game, error) {
	var board []byte
	err := r.pool.QueryRow(ctx, getGameQuery, id).Scan(&board)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainmodel.Game{}, domainerror.ErrGameNotFound
		}
		return domainmodel.Game{}, fmt.Errorf("get game : %w", err)
	}
	game := model.Game{ID: id}
	err = json.Unmarshal(board, &game.Board)
	if err != nil {
		return domainmodel.Game{}, fmt.Errorf("unmarshal board: %w", err)
	}
	return mapper.ToDomain(game)
}
