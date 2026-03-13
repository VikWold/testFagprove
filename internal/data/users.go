package data

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"testFagprove/internal/loggingutils"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"username"`
	Password    string     `json:"password"`
	CreatedAt   *time.Time `json:"created_at"`
	LastUpdated *time.Time `json:"last_updated"`
}

type UserModel struct {
	Timeout *time.Duration
	DB      *sql.DB
}

func (u UserModel) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	logger := loggingutils.LoggerFromContext(ctx)

	stmt := `
SELECT CAST(id AS CHAR(36)), username, password, created_at, last_updated
FROM public.users
WHERE id = $1;
	`

	var user User

	ctx, cancel := context.WithTimeout(ctx, *u.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", stmt),
			slog.String("id", id.String()),
		),
	)

	logger.InfoContext(ctx, "performing query")
	err := u.DB.QueryRowContext(ctx, stmt, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.LastUpdated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.InfoContext(ctx, "no rows found")
			return nil, ErrRecordNotFound
		default:
			logger.InfoContext(ctx, "an error occured while performing query", "error", err)
			return nil, err
		}
	}

	return &user, nil
}

func (u UserModel) List(ctx context.Context) ([]*User, error) {
	logger := loggingutils.LoggerFromContext(ctx)

	stmt := `
SELECT id, username, created_at, last_updated
FROM public.users
ORDER BY name DESC;
	`

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", stmt),
		),
	)

	ctx, cancel := context.WithTimeout(ctx, *u.Timeout)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, stmt)
	if err != nil {
		logger.ErrorContext(ctx, "error executing query", "error", err)
		return nil, err
	}
	defer rows.Close()

	var ul []*User

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.CreatedAt,
			&user.LastUpdated,
		)
		if err != nil {
			logger.ErrorContext(ctx, "error scanning row", "error", err)
			return nil, err
		}

		ul = append(ul, &user)
	}

	if err = rows.Err(); err != nil {
		logger.ErrorContext(ctx, "error with row iteration", "error", err)
		return nil, err
	}

	logger.InfoContext(ctx, "successfully fetched all users")

	return ul, nil
}

func (u UserModel) Insert(ctx context.Context, us *User) (*User, error) {
	logger := loggingutils.LoggerFromContext(ctx)

	stmt := `
INSERT INTO public.users (
    id, username, password, created_at, last_updated
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, username, password, created_at, last_updated;
    `

	ctx, cancel := context.WithTimeout(ctx, *u.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", stmt),
			"user", us,
		),
	)

	var result User

	logger.InfoContext(ctx, "performing query")
	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		us.ID,
		us.Username,
		us.Password,
		us.CreatedAt,
		us.LastUpdated,
	).Scan(
		&result.ID,
		&result.Username,
		&result.Password,
		&result.CreatedAt,
		&result.LastUpdated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.InfoContext(ctx, "no rows found")
			return nil, ErrRecordNotFound
		default:
			logger.InfoContext(ctx, "an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	return &result, nil
}
