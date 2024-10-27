package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zvxte/kera/model"
)

type UserStore interface {
	Create(ctx context.Context, user *model.User) error
	IsTaken(ctx context.Context, username string) (bool, error)
	UpdateLocation(ctx context.Context, id model.UUID, location *time.Location) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
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

func (sus SqlUserStore) Create(ctx context.Context, user *model.User) error {
	query := `
	INSERT INTO users(id, username, username_lower, display_name, hashed_password, location, creation_date)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err := sus.db.ExecContext(
		ctx, query,
		user.ID, user.Username, strings.ToLower(user.Username), user.DisplayName,
		user.HashedPassword, user.Location.String(), user.CreationDate,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (sus SqlUserStore) IsTaken(ctx context.Context, username string) (bool, error) {
	query := `
	SELECT 1
	FROM users
	WHERE username_lower = $1;
	`
	row := sus.db.QueryRowContext(
		ctx, query, strings.ToLower(username),
	)

	var isTaken uint8
	err := row.Scan(&isTaken)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if isTaken == 1 {
		return true, nil
	}

	return false, fmt.Errorf("failed to query: %w", err)
}

func (sus SqlUserStore) UpdateLocation(ctx context.Context, id model.UUID, location *time.Location) error {
	query := `
	UPDATE users
	SET location = $1
	WHERE id = $2;
	`
	_, err := sus.db.ExecContext(
		ctx, query,
		id, location.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	return nil
}

func (sus SqlUserStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
	SELECT id, username, display_name, hashed_password, location, creation_date
	FROM users
	WHERE username_lower = $1;
	`
	row := sus.db.QueryRowContext(
		ctx, query, strings.ToLower(username),
	)

	var idDB, usernameDB, displayName, hashedPassword, locationName string
	var creationDate time.Time
	err := row.Scan(
		&idDB, &usernameDB, &displayName, &hashedPassword, &locationName, &creationDate,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	id, err := model.ParseUUID(idDB)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	location, err := time.LoadLocation(locationName)
	if err != nil {
		location, _ = time.LoadLocation("UTC")
		err = sus.UpdateLocation(ctx, id, location)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}

	user, err := model.LoadUser(
		id, usernameDB, displayName, hashedPassword, location, creationDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
