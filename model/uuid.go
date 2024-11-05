package model

import (
	"fmt"

	"github.com/google/uuid"
)

// UUID represents an unique identifier used for model IDs.
type UUID [16]byte

// NewUUIDv7 returns a new version 7 UUID.
// It fails if the system's source of randomness is unavailable.
func NewUUIDv7() (UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return UUID{}, fmt.Errorf("failed to generate UUIDv7: %w", err)
	}
	return UUID(id), nil
}

// ParseUUID returns UUID parsed from given string.
// It fails if the given string is an invalid UUID.
func ParseUUID(value string) (UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return UUID{}, fmt.Errorf("failed to parse UUID: %w", err)
	}
	return UUID(id), nil
}

// String returns a string representation of UUID.
// "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" or "" if UUID is invalid.
func (id UUID) String() string {
	return uuid.UUID(id).String()
}
