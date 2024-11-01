package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zvxte/kera/model"
)

type HistoryParams struct {
	date time.Time
}

type HabitStore interface {
	Create(ctx context.Context, habit *model.Habit, userID model.UUID) error
	GetAllByUserID(ctx context.Context, userID model.UUID) ([]*model.Habit, error)
	GetHistory(
		ctx context.Context,
		habitID model.UUID,
		params HistoryParams,
	) (*model.HabitHistory, error)
	UpdateHistory(ctx context.Context, habitID model.UUID, date time.Time) error
}

type SqlHabitStore struct {
	db *sql.DB
}

func NewSqlHabitStore(db *sql.DB) (*SqlHabitStore, error) {
	if db == nil {
		return nil, NilDBPointerError
	}
	return &SqlHabitStore{db}, nil
}

func (s SqlHabitStore) Create(
	ctx context.Context, habit *model.Habit, userID model.UUID,
) error {
	query := `
	INSERT INTO habits(id, user_id, status, title, description,
					   tracked_week_days, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err := s.db.ExecContext(
		ctx, query,
		habit.ID, userID, habit.Status, habit.Title, habit.Description,
		habit.TrackedWeekDays, habit.StartDate, habit.EndDate,
	)
	if err != nil {
		return fmt.Errorf("failed to create habit: %w", err)
	}

	return nil
}
