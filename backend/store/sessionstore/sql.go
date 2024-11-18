package sessionstore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zvxte/kera/model/date"
	"github.com/zvxte/kera/model/session"
	"github.com/zvxte/kera/model/uuid"
	"github.com/zvxte/kera/store"
)

// Sql represents an relational database implementation
// of the [sessionstore.Store] interface.
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

func (s Sql) Create(ctx context.Context, session *session.Session) error {
	const query = `
	INSERT INTO sessions(id, user_id, creation_date, expiration_date)
	VALUES ($1, $2, $3, $4);
	`

	_, err := s.db.ExecContext(
		ctx, query,
		session.HashedID[:], session.UserID, time.Time(session.CreationDate),
		time.Time(session.ExpirationDate),
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (s Sql) Get(
	ctx context.Context, col Column, value any,
) (*session.Session, error) {
	const hashedIDQuery = `
	SELECT id, user_id, creation_date, expiration_date
	FROM sessions
	WHERE id = $1;
	`

	var query string
	switch col {
	case HashedIDColumn:
		hashedID, ok := value.(session.HashedID)
		if !ok {
			return nil, store.ErrInvalidColumnValue
		}
		value = hashedID[:]
		query = hashedIDQuery
	default:
		return nil, store.ErrInvalidColumn
	}

	var rawHashedID, rawUserID string
	var creationDate, expirationDate time.Time

	row := s.db.QueryRowContext(ctx, query, value)
	err := row.Scan(
		&rawHashedID, &rawUserID, &creationDate, &expirationDate,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var hashedID session.HashedID
	n := copy(hashedID[:], []byte(rawHashedID))
	if n != session.HashedIDLen {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	session := session.Load(
		hashedID, userID, date.Load(creationDate),
		date.Load(expirationDate),
	)

	return session, nil
}

func (s Sql) Delete(ctx context.Context, col Column, value any) error {
	const (
		hashedIDQuery = `
		DELETE FROM sessions WHERE id = $1;
		`
		userIDQuery = `
		DELETE FROM sessions WHERE user_id = $1;
		`
	)

	var query string
	switch col {
	case HashedIDColumn:
		hashedID, ok := value.(session.HashedID)
		if !ok {
			return store.ErrInvalidColumnValue
		}
		value = hashedID[:]
		query = hashedIDQuery
	case UserIDColumn:
		if _, ok := value.(uuid.UUID); !ok {
			return store.ErrInvalidColumnValue
		}
		query = userIDQuery
	default:
		return store.ErrInvalidColumn
	}

	_, err := s.db.ExecContext(ctx, query, value)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s Sql) Count(ctx context.Context, userID uuid.UUID) (uint, error) {
	const query = `
	SELECT COUNT(id)
	FROM sessions
	WHERE user_id = $1;
	`

	var count uint

	row := s.db.QueryRowContext(ctx, query, userID)
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	return count, nil
}
