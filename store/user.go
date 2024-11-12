package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zvxte/kera/model/date"
	"github.com/zvxte/kera/model/user"
	"github.com/zvxte/kera/model/uuid"
)

type UserStore interface {
	Create(ctx context.Context, user *user.User) error
	IsTaken(ctx context.Context, username string) (bool, error)
	UpdateDisplayName(ctx context.Context, id uuid.UUID, displayName string) error
	UpdateHashedPassword(ctx context.Context, id uuid.UUID, hashedPassword string) error
	GetByUsername(ctx context.Context, username string) (*user.User, error)
	Get(ctx context.Context, userID uuid.UUID) (*user.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}

type SqlUserStore struct {
	db *sql.DB
}

func NewSqlUserStore(db *sql.DB) (*SqlUserStore, error) {
	if db == nil {
		return nil, NilDBPointerError
	}
	return &SqlUserStore{db}, nil
}

func (s SqlUserStore) Create(ctx context.Context, user *user.User) error {
	query := `
	INSERT INTO users(id, username, username_lower,
					  display_name, hashed_password, creation_date)
	VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := s.db.ExecContext(
		ctx, query,
		user.ID, user.Username, strings.ToLower(user.Username), user.DisplayName,
		user.HashedPassword, time.Time(user.CreationDate),
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s SqlUserStore) IsTaken(ctx context.Context, username string) (bool, error) {
	query := `
	SELECT 1
	FROM users
	WHERE username_lower = $1;
	`
	row := s.db.QueryRowContext(
		ctx, query, strings.ToLower(username),
	)

	var isTaken uint8
	err := row.Scan(&isTaken)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if isTaken == 1 {
		return true, nil
	}

	return false, fmt.Errorf("failed to query: %w", err)
}

func (s SqlUserStore) UpdateDisplayName(
	ctx context.Context, id uuid.UUID, displayName string,
) error {
	query := `
	UPDATE users
	SET display_name = $1
	WHERE id = $2;
	`
	_, err := s.db.ExecContext(
		ctx, query,
		displayName, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update display name: %w", err)
	}

	return nil
}

func (s SqlUserStore) UpdateHashedPassword(
	ctx context.Context, id uuid.UUID, hashedPassword string,
) error {
	query := `
	UPDATE users
	SET hashed_password = $1
	WHERE id = $2;
	`
	_, err := s.db.ExecContext(
		ctx, query,
		hashedPassword, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update hashed password: %w", err)
	}

	return nil
}

func (s SqlUserStore) GetByUsername(
	ctx context.Context, username string,
) (*user.User, error) {
	query := `
	SELECT id, username, display_name, hashed_password, creation_date
	FROM users
	WHERE username_lower = $1;
	`
	row := s.db.QueryRowContext(
		ctx, query, strings.ToLower(username),
	)

	var rawUserID, dbUsername, displayName, hashedPassword string
	var creationDate time.Time
	err := row.Scan(
		&rawUserID, &dbUsername, &displayName,
		&hashedPassword, &creationDate,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	id, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user, err := user.Load(
		id, dbUsername, displayName,
		hashedPassword, date.Load(creationDate),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s SqlUserStore) Get(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	query := `
	SELECT username, display_name, hashed_password, creation_date
	FROM users
	WHERE id = $1;
	`
	row := s.db.QueryRowContext(ctx, query, userID[:])

	var dbUsername, displayName, hashedPassword string
	var creationDate time.Time
	err := row.Scan(
		&dbUsername, &displayName, &hashedPassword, &creationDate,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user, err := user.Load(
		userID, dbUsername, displayName,
		hashedPassword, date.Load(creationDate),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s SqlUserStore) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `
	DELETE FROM users
	WHERE id = $1;
	`
	_, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
