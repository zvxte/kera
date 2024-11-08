package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zvxte/kera/model"
)

type SessionStore interface {
	Create(ctx context.Context, session *model.Session, userID model.UUID) error
	Get(ctx context.Context, hashedSessionID model.HashedSessionID) (
		*model.Session, model.UUID, error,
	)
	Delete(ctx context.Context, hashedSessionID model.HashedSessionID) error
	DeleteAll(ctx context.Context, userID model.UUID) error
	Count(ctx context.Context, userID model.UUID) (uint, error)
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

func (s SqlSessionStore) Create(
	ctx context.Context, session *model.Session, userID model.UUID,
) error {
	query := `
	INSERT INTO sessions(id, user_id, creation_date, expiration_date)
	VALUES ($1, $2, $3, $4);
	`
	_, err := s.db.ExecContext(
		ctx, query,
		session.HashedID[:], userID, session.CreationDate, session.ExpirationDate,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (s SqlSessionStore) Get(
	ctx context.Context, hashedSessionID model.HashedSessionID,
) (*model.Session, model.UUID, error) {
	query := `
	SELECT user_id, creation_date, expiration_date
	FROM sessions
	WHERE id = $1;
	`
	row := s.db.QueryRowContext(ctx, query, hashedSessionID[:])

	var rawUserID string
	var creation_date, expiration_date time.Time
	err := row.Scan(&rawUserID, &creation_date, &expiration_date)

	if err == sql.ErrNoRows {
		return nil, model.UUID{}, nil
	}

	if err != nil {
		return nil, model.UUID{}, fmt.Errorf("failed to get session: %w", err)
	}

	userID, err := model.ParseUUID(rawUserID)
	if err != nil {
		return nil, model.UUID{}, fmt.Errorf("failed to get session: %w", err)
	}

	session := model.LoadSession(
		hashedSessionID, creation_date, expiration_date,
	)

	return session, userID, nil
}

func (s SqlSessionStore) Delete(
	ctx context.Context, hashedSessionID model.HashedSessionID,
) error {
	query := `
	DELETE FROM sessions
	WHERE id = $1;
	`
	_, err := s.db.ExecContext(ctx, query, hashedSessionID[:])
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s SqlSessionStore) DeleteAll(
	ctx context.Context, userID model.UUID,
) error {
	query := `
	DELETE FROM sessions
	WHERE user_id = $1;
	`
	_, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete all sessions: %w", err)
	}

	return nil
}

func (s SqlSessionStore) Count(
	ctx context.Context, userID model.UUID,
) (uint, error) {
	query := `
	SELECT COUNT(id)
	FROM sessions
	WHERE user_id = $1;
	`
	row := s.db.QueryRowContext(ctx, query, userID)

	var count uint
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	return count, nil
}
