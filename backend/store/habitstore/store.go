package habitstore

import (
	"context"

	"github.com/zvxte/kera/model/date"
	"github.com/zvxte/kera/model/habit"
	"github.com/zvxte/kera/model/uuid"
)

type Store interface {
	// Create inserts a new habit into the store.
	// It returns an error if there is a connection issue.
	Create(ctx context.Context, habit *habit.Habit, userID uuid.UUID) error

	// GetAll returns a session slice from the store or a nil slice.
	// It fails if there is a connection issue.
	GetAll(ctx context.Context, userID uuid.UUID) ([]*habit.Habit, error)

	// Update updates a habit in the store.
	// It fails if there is a connection issue.
	// It returns [store.ErrInvalidColumn] or [store.ErrInvalidColumnValue]
	// if unsupported column or invalid column value is provided.
	// Supported columns: [habitstore.TitleColumn], [habitstore.DescriptionColumn].
	Update(
		ctx context.Context, id uuid.UUID, col Column, value any, userID uuid.UUID,
	) error

	// Delete deletes a habit from the store.
	// It fails if there is a connection issue.
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// End ends a habit in the store.
	// It sets the status to [habit.Ended] and end date to now
	// if it's not already ended.
	// It fails if there is a connection issue.
	End(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// UpdateHistory updates a habit's history.
	// It sets the provided date to [habit.DayDone],
	// or unsets if it was already set.
	// It fails if there is a connection issue.
	UpdateHistory(
		ctx context.Context, id uuid.UUID, historyDate date.Date, userID uuid.UUID,
	) error

	// GetMonthHistory returns a month of the habit's history from the provided date.
	// It fails if there is a connection issue.
	GetMonthHistory(
		ctx context.Context, id uuid.UUID, historyDate date.Date, userID uuid.UUID,
	) (habit.History, error)
}

// Column represents a store column.
type Column uint8

const (
	TitleColumn Column = iota
	DescriptionColumn
)

func (c Column) String() string {
	switch c {
	case TitleColumn:
		return "title"
	case DescriptionColumn:
		return "description"
	default:
		return ""
	}
}
