package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	CreateUser(username, passwordHash string) (User, error)
	GetUserByID(userID uuid.UUID) (User, error)
	GetUserByUsername(username string) (User, error)
}

type service struct {
	ctx     context.Context
	queries *Queries
}

func NewRepository(ctx context.Context, db DBTX) Repository {
	return &service{
		ctx:     ctx,
		queries: New(db),
	}
}

func (s service) CreateUser(username, passwordHash string) (User, error) {
	newUser, err := s.queries.CreateUser(
		s.ctx,
		CreateUserParams{
			username,
			passwordHash,
		},
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return User{}, ErrUserAlreadyExists
	}
	return newUser, nil
}

func (s service) GetUserByID(userID uuid.UUID) (User, error) {
	return s.queries.GetUserByID(s.ctx, userID)
}

func (s service) GetUserByUsername(username string) (User, error) {
	user, err := s.queries.GetUserByUsername(s.ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}
	return user, nil
}
