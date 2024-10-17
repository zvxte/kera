package store

import (
	"context"
	"database/sql"

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
