package repository

import (
	"errors"
	"tictactoe/datasource/mapper"
	"tictactoe/datasource/model"
	"tictactoe/datasource/storage"
	"tictactoe/domain"

	"github.com/google/uuid"
)

var ErrInvalidStorageData = errors.New("invalid data in storage")

type gameRepository struct {
	storage *storage.Storage
}

func NewRepository(storage *storage.Storage) Repository {
	return &gameRepository{storage: storage}
}

func (r *gameRepository) Save(game domain.Game) error {
	saveGame := mapper.ToDatasource(game)

	r.storage.Games.Store(game.ID, saveGame)
	return nil
}

func (r *gameRepository) Get(id uuid.UUID) (domain.Game, error) {
	value, ok := r.storage.Games.Load(id)
	if !ok {
		return domain.Game{}, domain.ErrGameNotFound
	}
	game, ok := value.(model.Game)
	if !ok {
		return domain.Game{}, ErrInvalidStorageData
	}
	return mapper.ToDomain(game)
}
