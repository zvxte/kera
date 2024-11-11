package user

import (
	"github.com/zvxte/kera/hash/argon2id"
	"github.com/zvxte/kera/model"
	"github.com/zvxte/kera/model/date"
	"github.com/zvxte/kera/model/uuid"
)

// User represents an application user.
type User struct {
	ID             uuid.UUID
	Username       string
	DisplayName    string
	HashedPassword string
	CreationDate   date.Date
}

// New returns a new *User.
// It fails if the provided parameters do not meet the application requirements.
// The returned error is safe for client-side message.
// The plain password is hashed using Argon2ID.
// The Username and DisplayName fields are set to the given username,
func New(username, plainPassword string) (*User, error) {
	if err := ValidateUsername(username); err != nil {
		return nil, err
	}

	if err := ValidatePlainPassword(plainPassword); err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, model.ErrUnexpected
	}

	hashedPassword, err := argon2id.Hash(plainPassword, argon2id.DefaultParams)
	if err != nil {
		return nil, model.ErrUnexpected
	}

	return &User{
		ID:             id,
		Username:       username,
		DisplayName:    username,
		HashedPassword: hashedPassword,
		CreationDate:   date.Now(),
	}, nil
}

// Load returns a *User.
// It fails if the provided parameters do not meet the application requirements.
// The returned error is safe for client-side message.
func Load(
	id uuid.UUID,
	username, displayName, hashedPassword string,
	creationDate date.Date,
) (*User, error) {
	if err := ValidateUsername(username); err != nil {
		return nil, err
	}

	if err := ValidateDisplayName(displayName); err != nil {
		return nil, err
	}

	return &User{
		ID:             id,
		Username:       username,
		DisplayName:    displayName,
		HashedPassword: hashedPassword,
		CreationDate:   creationDate,
	}, nil
}
