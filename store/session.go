package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zvxte/kera/model"
)

type SessionStore interface {
	Create(ctx context.Context, session *model.Session, userID model.UUID) error
}

type SqlSessionStore struct {
	db *sql.DB
}

func NewSqlSessionStore(db *sql.DB) (*SqlSessionStore, error) {
	if db == nil {
		return nil, NilDBPointerError
	}
	return &SqlSessionStore{db}, nil
}

func (sss SqlSessionStore) Create(
	ctx context.Context, session *model.Session, userID model.UUID,
) error {
	query := `
	INSERT INTO sessions(id, user_id, creation_date, expiration_date)
	VALUES ($1, $2, $3, $4);
	`
	_, err := sss.db.ExecContext(
		ctx, query,
		session.HashedID[:], userID, session.CreationDate, session.ExpirationDate,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}
