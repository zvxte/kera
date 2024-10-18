package model

import (
	"fmt"

	"github.com/google/uuid"
)

type UUID [16]byte

func NewUUIDv7() (UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return UUID(id), fmt.Errorf("failed to generate UUID v7: %w", err)
	}
	return UUID(id), nil
}

func ParseUUID(value string) (UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return UUID(id), fmt.Errorf("failed to parse UUID: %w", err)
	}
	return UUID(id), nil
}

func (id UUID) String() string {
	return uuid.UUID(id).String()
}
