package userstore

import (
	"context"
	"errors"

	"github.com/zvxte/kera/model/user"
	"github.com/zvxte/kera/model/uuid"
)

var ErrUsernameAlreadyTaken = errors.New("username is already taken")

type Store interface {
	// Create inserts a new user into the store.
	// It returns an error if there is a connection issue,
	// or [userstore.ErrUsernameAlreadyTaken] on a username conflict.
	Create(ctx context.Context, user *user.User) error

	// Get returns a user from the store or nil.
	// It fails if there is a connection issue.
	// It returns [store.ErrInvalidColumn] or [store.ErrInvalidColumnValue]
	// if unsupported column or invalid column value is provided.
	// Supported columns: [userstore.IDColumn], [userstore.UsernameColumn].
	Get(ctx context.Context, col Column, value any) (*user.User, error)

	// Update updates a user in the store.
	// It fails if there is a connection issue.
	// It returns [store.ErrInvalidColumn] or [store.ErrInvalidColumnValue]
	// if unsupported column or invalid column value is provided.
	// Supported columns: [userstore.DisplayNameColumn], [userstore.HashedPasswordColumn].
	Update(ctx context.Context, id uuid.UUID, col Column, value any) error

	// Delete deletes a user from the store.
	// It fails if there is a connection issue.
	Delete(ctx context.Context, id uuid.UUID) error
}

// Column represents a store column.
type Column uint8

const (
	IDColumn Column = iota
	UsernameColumn
	DisplayNameColumn
	HashedPasswordColumn
)

func (c Column) String() string {
	switch c {
	case IDColumn:
		return "id"
	case UsernameColumn:
		return "username"
	case DisplayNameColumn:
		return "display_name"
	case HashedPasswordColumn:
		return "hashed_password"
	default:
		return ""
	}
}
