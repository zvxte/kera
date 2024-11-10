package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zvxte/kera/model/date"
	"github.com/zvxte/kera/model/session"
	"github.com/zvxte/kera/model/uuid"
)

type SessionStore interface {
	Create(ctx context.Context, session *session.Session, userID uuid.UUID) error
	Get(ctx context.Context, hashedID session.HashedID) (
		*session.Session, uuid.UUID, error,
	)
	Delete(ctx context.Context, hashedID session.HashedID) error
	DeleteAll(ctx context.Context, userID uuid.UUID) error
	Count(ctx context.Context, userID uuid.UUID) (uint, error)
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
	ctx context.Context, session *session.Session, userID uuid.UUID,
) error {
	query := `
	INSERT INTO sessions(id, user_id, creation_date, expiration_date)
	VALUES ($1, $2, $3, $4);
	`
	_, err := s.db.ExecContext(
		ctx, query,
		session.HashedID[:], userID, time.Time(session.CreationDate),
		time.Time(session.ExpirationDate),
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (s SqlSessionStore) Get(
	ctx context.Context, hashedID session.HashedID,
) (*session.Session, uuid.UUID, error) {
	query := `
	SELECT user_id, creation_date, expiration_date
	FROM sessions
	WHERE id = $1;
	`
	row := s.db.QueryRowContext(ctx, query, hashedID[:])

	var rawUserID string
	var creation_date, expiration_date time.Time
	err := row.Scan(&rawUserID, &creation_date, &expiration_date)

	if err == sql.ErrNoRows {
		return nil, uuid.UUID{}, nil
	}

	if err != nil {
		return nil, uuid.UUID{}, fmt.Errorf("failed to get session: %w", err)
	}

	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, uuid.UUID{}, fmt.Errorf("failed to get session: %w", err)
	}

	session := session.Load(
		hashedID, date.Load(creation_date),
		date.Load(expiration_date),
	)

	return session, userID, nil
}

func (s SqlSessionStore) Delete(ctx context.Context, hashedID session.HashedID) error {
	query := `
	DELETE FROM sessions
	WHERE id = $1;
	`
	_, err := s.db.ExecContext(ctx, query, hashedID[:])
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s SqlSessionStore) DeleteAll(
	ctx context.Context, userID uuid.UUID,
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

func (s SqlSessionStore) Count(ctx context.Context, userID uuid.UUID) (uint, error) {
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
