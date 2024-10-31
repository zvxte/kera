package store

import (
	"context"
	"database/sql"
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
