package sessionstore

import (
	"context"

	"github.com/zvxte/kera/model/session"
	"github.com/zvxte/kera/model/uuid"
)

type Store interface {
	// Create inserts a new session into the store.
	// It returns an error if there is a connection issue.
	Create(ctx context.Context, session *session.Session) error

	// Get returns a session from the store or nil.
	// It fails if there is a connection issue.
	// It returns [store.ErrInvalidColumn] or [store.ErrInvalidColumnValue]
	// if unsupported column or invalid column value is provided.
	// Supported columns: [sessionstore.HashedIDColumn].
	Get(ctx context.Context, col Column, value any) (*session.Session, error)

	// Delete deletes a session from the store.
	// It fails if there is a connection issue.
	// It returns [store.ErrInvalidColumn] or [store.ErrInvalidColumnValue]
	// if unsupported column or invalid column value is provided.
	// Supported columns: [sessionstore.HashedIDColumn], [sessionstore.UserIDColumn].
	Delete(ctx context.Context, col Column, value any) error

	// Count returns the number of sessions of the provided user.
	// It fails if there is a connection issue.
	Count(ctx context.Context, userID uuid.UUID) (uint, error)
}

// Column represents a store column.
type Column uint8

const (
	HashedIDColumn Column = iota
	UserIDColumn
)

func (c Column) String() string {
	switch c {
	case HashedIDColumn:
		return "id"
	case UserIDColumn:
		return "user_id"
	default:
		return ""
	}
}
