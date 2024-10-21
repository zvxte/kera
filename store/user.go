package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zvxte/kera/model"
)

type UserStore interface {
	CreateUser(ctx context.Context, user model.User) error
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

func (sus SqlUserStore) CreateUser(ctx context.Context, user model.User) error {
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
