package userstore

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zvxte/kera/model/date"
	"github.com/zvxte/kera/model/user"
	"github.com/zvxte/kera/model/uuid"
	"github.com/zvxte/kera/store"
)

// Sql represents an relational database implementation
// of the [userstore.Store] interface.
// It uses an [*sql.DB] pool to interact with the database.
type Sql struct {
	db *sql.DB
}

func NewSql(db *sql.DB) (Sql, error) {
	if db == nil {
		return Sql{}, store.ErrNilDB
	}

	return Sql{db}, nil
}

func (s Sql) Create(ctx context.Context, user *user.User) error {
	const query = `
	INSERT INTO users(
		id, username, username_lower, display_name, hashed_password, creation_date
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (username_lower) DO NOTHING
	RETURNING 1;
	`

	var result uint8

	row := s.db.QueryRowContext(
		ctx, query,
		user.ID, user.Username, strings.ToLower(user.Username),
		user.DisplayName, user.HashedPassword, time.Time(user.CreationDate),
	)
	err := row.Scan(&result)
	if err == sql.ErrNoRows {
		return ErrUsernameAlreadyTaken
	}
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s Sql) Get(
	ctx context.Context, col Column, value any,
) (*user.User, error) {
	const (
		idQuery = `
		SELECT id, username, display_name, hashed_password, creation_date
		FROM users
		WHERE id = $1;
		`
		usernameQuery = `
		SELECT id, username, display_name, hashed_password, creation_date
		FROM users
		WHERE username_lower = $1;
		`
	)

	var query string
	switch col {
	case IDColumn:
		if _, ok := value.(uuid.UUID); !ok {
			return nil, store.ErrInvalidColumnValue
		}
		query = idQuery
	case UsernameColumn:
		value, ok := value.(string)
		if !ok {
			return nil, store.ErrInvalidColumnValue
		}
		value = strings.ToLower(value)
		query = usernameQuery
	default:
		return nil, store.ErrInvalidColumn
	}

	var rawUserID, username, displayName, hashedPassword string
	var creationDate time.Time

	row := s.db.QueryRowContext(ctx, query, value)
	err := row.Scan(
		&rawUserID, &username, &displayName, &hashedPassword, &creationDate,
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
		id, username, displayName,
		hashedPassword, date.Load(creationDate),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s Sql) Update(
	ctx context.Context, id uuid.UUID, col Column, value any,
) error {
	const (
		displayNameQuery = `
		UPDATE users SET display_name = $1 WHERE id = $2;
		`
		hashedPasswordQuery = `
		UPDATE users SET hashed_password = $1 WHERE id = $2;
		`
	)

	var query string
	switch col {
	case DisplayNameColumn:
		if _, ok := value.(string); !ok {
			return store.ErrInvalidColumnValue
		}
		query = displayNameQuery
	case HashedPasswordColumn:
		if _, ok := value.(string); !ok {
			return store.ErrInvalidColumnValue
		}
		query = hashedPasswordQuery
	default:
		return store.ErrInvalidColumn
	}

	_, err := s.db.ExecContext(ctx, query, value, id)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", col, err)
	}

	return nil
}

func (s Sql) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
	DELETE FROM users
	WHERE id = $1;
	`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
