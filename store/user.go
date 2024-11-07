package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zvxte/kera/model"
)

type UserStore interface {
	Create(ctx context.Context, user *model.User) error
	IsTaken(ctx context.Context, username string) (bool, error)
	UpdateDisplayName(ctx context.Context, id model.UUID, displayName string) error
	UpdateHashedPassword(ctx context.Context, id model.UUID, hashedPassword string) error
	UpdateLocation(ctx context.Context, id model.UUID, location *time.Location) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Get(ctx context.Context, userID model.UUID) (*model.User, error)
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

func (s SqlUserStore) Create(ctx context.Context, user *model.User) error {
	query := `
	INSERT INTO users(id, username, username_lower, display_name, hashed_password, location, creation_date)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err := s.db.ExecContext(
		ctx, query,
		user.ID, user.Username, strings.ToLower(user.Username), user.DisplayName,
		user.HashedPassword, user.Location.String(), user.CreationDate,
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
	ctx context.Context, id model.UUID, displayName string,
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
	ctx context.Context, id model.UUID, hashedPassword string,
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

func (s SqlUserStore) UpdateLocation(
	ctx context.Context, id model.UUID, location *time.Location,
) error {
	query := `
	UPDATE users
	SET location = $1
	WHERE id = $2;
	`
	_, err := s.db.ExecContext(
		ctx, query,
		location.String(), id,
	)
	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	return nil
}

func (s SqlUserStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
	SELECT id, username, display_name, hashed_password, location, creation_date
	FROM users
	WHERE username_lower = $1;
	`
	row := s.db.QueryRowContext(
		ctx, query, strings.ToLower(username),
	)

	var rawUserID, dbUsername, displayName, hashedPassword, locationName string
	var creationDate time.Time
	err := row.Scan(
		&rawUserID, &dbUsername, &displayName, &hashedPassword, &locationName, &creationDate,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	id, err := model.ParseUUID(rawUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	location, err := time.LoadLocation(locationName)
	if err != nil {
		location = time.UTC
		err = s.UpdateLocation(ctx, id, location)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}

	user, err := model.LoadUser(
		id, dbUsername, displayName, hashedPassword, location, creationDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s SqlUserStore) Get(ctx context.Context, userID model.UUID) (*model.User, error) {
	query := `
	SELECT username, display_name, hashed_password, location, creation_date
	FROM users
	WHERE id = $1;
	`
	row := s.db.QueryRowContext(ctx, query, userID[:])

	var dbUsername, displayName, hashedPassword, locationName string
	var creationDate time.Time
	err := row.Scan(
		&dbUsername, &displayName, &hashedPassword, &locationName, &creationDate,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	location, err := time.LoadLocation(locationName)
	if err != nil {
		location = time.UTC
		err = s.UpdateLocation(ctx, userID, location)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}

	user, err := model.LoadUser(
		userID, dbUsername, displayName, hashedPassword, location, creationDate,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
