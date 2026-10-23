package domain

import "context"

type UserService interface {
	GetByLogin(ctx context.Context, login string) (User, error)
	Create(ctx context.Context, user User) error
}
