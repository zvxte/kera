package habit

import (
	"math/bits"

	"github.com/zvxte/kera/model"
	"github.com/zvxte/kera/model/date"
	"github.com/zvxte/kera/model/uuid"
)

// Habit represents an application user's habit.
type Habit struct {
	ID              uuid.UUID
	Status          Status
	Title           string
	Description     string
	TrackedWeekDays TrackedWeekDays
	StartDate       date.Date
	EndDate         date.Date
}

// Status represents status of a habit.
type Status uint8

const (
	Active Status = iota
	Ended
)

// TrackedWeekDays represents days of the week that are tracked in a bitmap,
// where each bit (0 - untracked, 1 - tracked) represents a day
// starting from Monday as the first bit (LSB).
// The eighth bit should always be unset.
type TrackedWeekDays uint8

// WeekDays returns []WeekDay representing only the days that are tracked
// within the TrackedWeekDays bitmap.
// WeekDays always returns a non-nil and non-empty slice.
func (d TrackedWeekDays) WeekDays() []WeekDay {
	weekDays := make([]WeekDay, 0, bits.OnesCount8(uint8(d)))
	for i := 0; i <= int(Sunday); i++ {
		bit := (d >> i) & 1
		if bit == 1 {
			weekDays = append(weekDays, WeekDay(i))
		}
	}
	return weekDays
}

// WeekDay represents a day of the week.
type WeekDay uint8

const (
	Monday WeekDay = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// History represents a history of a habit.
type History []Day

// HabitDay represents a single day record in a history of a habit.
// It contains the status and the date of that record.
type Day struct {
	Status DayStatus
	Date   date.Date
}

// DayStatus represents a status of a single day record in a history of a habit.
type DayStatus uint8

const (
	DayUntracked DayStatus = iota
	DayDone
	DayMissed
	DayPending
)

// New returns a new *Habit.
// It fails if the provided parameters do not meet the application requirements.
// The returned error is safe for client-side message.
// The status field is set to Active value.
// The startDate field is set to the current Date value.
// The endDate field is set to the zero value of the date.Date type.
func New(
	title, description string,
	weekDays ...WeekDay,
) (*Habit, error) {
	if err := ValidateTitle(title); err != nil {
		return nil, err
	}

	if err := ValidateDescription(description); err != nil {
		return nil, err
	}

	trackedWeekDays, err := newTrackedWeekDays(weekDays...)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, model.ErrUnexpected
	}

	return &Habit{
		ID:              id,
		Status:          Active,
		Title:           title,
		Description:     description,
		TrackedWeekDays: trackedWeekDays,
		StartDate:       date.Now(),
		EndDate:         date.Date{},
	}, nil
}

// Load returns a *Habit.
// It fails if the provided parameters do not meet the application requirements.
// The returned error is safe for client-side message.
func Load(
	id uuid.UUID, status Status, title, description string,
	trackedWeekDays TrackedWeekDays, startDate, endDate date.Date,
) (*Habit, error) {
	if err := validateStatus(status); err != nil {
		return nil, err
	}

	if err := ValidateTitle(title); err != nil {
		return nil, err
	}

	if err := ValidateDescription(description); err != nil {
		return nil, err
	}

	if err := validateTrackedWeekDays(trackedWeekDays); err != nil {
		return nil, err
	}

	return &Habit{
		ID:              id,
		Status:          status,
		Title:           title,
		Description:     description,
		TrackedWeekDays: trackedWeekDays,
		StartDate:       startDate,
		EndDate:         endDate,
	}, nil
}

// newTrackedWeekDays returns new TrackedWeekDays value.
// It fails if the provided parameters are invalid days of the week
// or no parameters are provided.
func newTrackedWeekDays(weekDays ...WeekDay) (TrackedWeekDays, error) {
	var trackedWeekDays TrackedWeekDays
	for _, day := range weekDays {
		trackedWeekDays |= TrackedWeekDays(1 << day)
	}

	err := validateTrackedWeekDays(trackedWeekDays)
	if err != nil {
		return 0, err
	}

	return trackedWeekDays, nil
}
