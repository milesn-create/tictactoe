package dsrepository

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/datasource/dsmapper"
	"tictactoe/datasource/dsmodel"
	"tictactoe/domain/domainerror"
	"tictactoe/domain/domainmodel"
	"tictactoe/domain/domainrepository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) domainrepository.UserRepository {
	return &userRepository{pool: pool}
}

const (
	getUserByLoginQuery = `SELECT  id, login, password_hash
					   FROM users
					   WHERE login = $1`
	createUserQuery = `INSERT INTO users(id, login, password_hash)
					   VALUES ($1, $2, $3)`
	getUserByIDQuery = `SELECT id, login,password_hash
				   		FROM users 
				   		WHERE id = $1`
)

func (r *userRepository) GetByLogin(ctx context.Context, login string) (domainmodel.User, error) {
	rows, err := r.pool.Query(ctx, getUserByLoginQuery, login)
	if err != nil {
		return domainmodel.User{}, fmt.Errorf("failed query get user by login: %w", err)
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dsmodel.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainmodel.User{}, domainerror.ErrUserNotFound
		}
		return domainmodel.User{}, fmt.Errorf("get by login: %w", err)
	}
	return dsmapper.UserToDomain(user), nil
}

func (r *userRepository) Create(ctx context.Context, user domainmodel.User) error {
	userDatasource := dsmapper.UserToDatasource(user)
	_, err := r.pool.Exec(ctx, createUserQuery, userDatasource.ID, userDatasource.Login, userDatasource.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domainerror.ErrUserExists
		}
		return fmt.Errorf("failed query create user: %w", err)
	}
	return nil
}
func (r *userRepository) GetByID(ctx context.Context, userID uuid.UUID) (domainmodel.User, error) {
	rows, err := r.pool.Query(ctx, getUserByIDQuery, userID)
	if err != nil {
		return domainmodel.User{}, fmt.Errorf("failed query get user by id: %w", err)
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dsmodel.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainmodel.User{}, domainerror.ErrUserNotFound
		}
		return domainmodel.User{}, fmt.Errorf("get by id: %w", err)
	}
	return dsmapper.UserToDomain(user), nil
}
