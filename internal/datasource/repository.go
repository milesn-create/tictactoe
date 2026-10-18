package datasource

import (
	"tictactoe/domain"

	"github.com/google/uuid"
)

type Repository interface {
	Get(id uuid.UUID) (domain.Game, error)
	Save(game domain.Game) error
}
