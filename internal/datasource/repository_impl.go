package datasource

import (
	"errors"
	"tictactoe/domain"

	"github.com/google/uuid"
)

var ErrGameNotFound = errors.New("game not found")
var ErrInvalidStorageData = errors.New("invalid data in storage")

type gameRepository struct {
	storage *Storage
}

func NewRepository(storage *Storage) Repository {
	return &gameRepository{storage: storage}
}

func (r *gameRepository) Save(game domain.Game) error {
	saveGame := toDatasource(game)

	r.storage.games.Store(game.ID, saveGame)
	return nil
}

func (r *gameRepository) Get(id uuid.UUID) (domain.Game, error) {
	value, ok := r.storage.games.Load(id)
	if !ok {
		return domain.Game{}, ErrGameNotFound
	}
	game, ok := value.(Game)
	if !ok {
		return domain.Game{}, ErrInvalidStorageData
	}
	return toDomain(game)
}
