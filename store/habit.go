package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zvxte/kera/model/date"
	"github.com/zvxte/kera/model/habit"
	"github.com/zvxte/kera/model/uuid"
)

// type HistoryParams struct {
// 	date time.Time
// }

type HabitStore interface {
	Create(ctx context.Context, habit *habit.Habit, userID uuid.UUID) error
	GetAll(ctx context.Context, userID uuid.UUID) ([]*habit.Habit, error)
	UpdateTitle(
		ctx context.Context, id uuid.UUID, title string, userID uuid.UUID,
	) error
	UpdateDescription(
		ctx context.Context, id uuid.UUID, description string, userID uuid.UUID,
	) error
	End(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	// GetHistory(
	// 	ctx context.Context,
	// 	habitID model.UUID,
	// 	params HistoryParams,
	// ) (*model.HabitHistory, error)
	// UpdateHistory(ctx context.Context, habitID model.UUID, date time.Time) error
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
	ctx context.Context, habit *habit.Habit, userID uuid.UUID,
) error {
	query := `
	INSERT INTO habits(id, user_id, status, title, description,
					   tracked_week_days, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err := s.db.ExecContext(
		ctx, query,
		habit.ID, userID, habit.Status, habit.Title, habit.Description,
		habit.TrackedWeekDays, time.Time(habit.StartDate), time.Time(habit.EndDate),
	)
	if err != nil {
		return fmt.Errorf("failed to create habit: %w", err)
	}

	return nil
}

func (s SqlHabitStore) GetAll(
	ctx context.Context, userID uuid.UUID,
) ([]*habit.Habit, error) {
	query := `
	SELECT id, status, title, description,
		   tracked_week_days, start_date, end_date
	FROM habits
	WHERE user_id = $1;
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all habits: %w", err)
	}
	defer rows.Close()

	var rawID, title, description string
	var status habit.Status
	var trackedWeekDays habit.TrackedWeekDays
	var startDate, endDate time.Time

	var habits []*habit.Habit

	for rows.Next() {
		err = rows.Scan(
			&rawID, &status, &title, &description,
			&trackedWeekDays, &startDate, &endDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get all habits: %w", err)
		}

		id, err := uuid.Parse(rawID)
		if err != nil {
			return nil, fmt.Errorf("failed to get all habits: %w", err)
		}

		habit, err := habit.Load(
			id, status, title, description, trackedWeekDays,
			date.Load(startDate), date.Load(endDate),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get all habits: %w", err)
		}

		habits = append(habits, habit)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to get all habits: %w", err)
	}

	return habits, nil
}

func (s SqlHabitStore) UpdateTitle(
	ctx context.Context, id uuid.UUID, title string, userID uuid.UUID,
) error {
	query := `
	UPDATE habits
	SET title = $1
	WHERE id = $2 AND user_id = $3;
	`
	_, err := s.db.ExecContext(
		ctx, query,
		title, id, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update title: %w", err)
	}

	return nil
}

func (s SqlHabitStore) UpdateDescription(
	ctx context.Context, id uuid.UUID, description string, userID uuid.UUID,
) error {
	query := `
	UPDATE habits
	SET description = $1
	WHERE id = $2 AND user_id = $3;
	`
	_, err := s.db.ExecContext(
		ctx, query,
		description, id, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update description: %w", err)
	}

	return nil
}

func (s SqlHabitStore) End(
	ctx context.Context, id uuid.UUID, userID uuid.UUID,
) error {
	query := `
	UPDATE habits
	SET status = $1, end_date = $2
	WHERE status = $3 AND id = $4 AND user_id = $5;
	`
	_, err := s.db.ExecContext(
		ctx, query,
		habit.Ended, time.Time(date.Now()), habit.Active, id, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to end habit: %w", err)
	}

	return nil
}

func (s SqlHabitStore) Delete(
	ctx context.Context, id uuid.UUID, userID uuid.UUID,
) error {
	query := `
	DELETE FROM habits
	WHERE id = $1 AND user_id = $2;
	`
	_, err := s.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	return nil
}
