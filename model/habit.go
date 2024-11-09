package model

import (
	"errors"
	"math/bits"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	titleMinChars = 2
	titleMaxChars = 64

	descriptionMaxChars = 256

	trackedWeekDaysMin = 1
	trackedWeekDaysMax = (1 << 7) - 1

	weekDayMax = 6

	habitStatusMin = 0
	habitStatusMax = 1
)

var (
	errTitleTooShort          = errors.New("title is too short")
	errTitleTooLong           = errors.New("title is too long")
	errTitleInvalid           = errors.New("title is invalid")
	errDescriptionTooLong     = errors.New("description is too long")
	errDescriptionInvalid     = errors.New("description is invalid")
	errTrackedWeekDaysInvalid = errors.New(
		"tracked days of the week are invalid: unrecognized day specified",
	)
	errTrackedWeekDaysEmpty = errors.New(
		"at least one day of the week must be specified for tracking",
	)
)

// Habit represents an application user's habit.
type Habit struct {
	ID              UUID
	Status          HabitStatus
	Title           string
	Description     string
	TrackedWeekDays TrackedWeekDays
	StartDate       time.Time
	EndDate         time.Time
}

// HabitStatus represents status of a habit.
type HabitStatus uint8

const (
	HabitActive HabitStatus = iota
	HabitEnded
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
	for i := 0; i <= weekDayMax; i++ {
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

// HabitHistory represents a history of a habit.
type HabitHistory []HabitDay

// HabitDay represents a single day record in a history of a habit.
// It contains the status and the date of that record.
type HabitDay struct {
	Status HabitDayStatus
	Date   time.Time
}

// HabitDayStatus represents a status of a single day record in a history of a habit.
type HabitDayStatus uint8

const (
	DayUntracked HabitDayStatus = iota
	DayDone
	DayMissed
	DayPending
)

// NewHabit returns a new *Habit.
// It fails if the provided parameters do not meet the application requirements.
// The returned error is safe for client-side message.
// The status field is set to HabitActive.
// The startDate field is set to the current date in UTC.
// The endDate field is set to the zero value of the time.Time type.
func NewHabit(
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

	id, err := NewUUIDv7()
	if err != nil {
		return nil, errInternalServer
	}

	status := HabitActive

	startDate := DateNow()

	endDate := time.Time{}

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

// LoadHabit returns a *Habit.
// It fails if the provided parameters do not meet the application requirements.
// The returned error is safe for client-side message.
func LoadHabit(
	id UUID, status HabitStatus, title, description string,
	trackedWeekDays TrackedWeekDays, startDate, endDate time.Time,
) (*Habit, error) {
	if err := validateHabitStatus(status); err != nil {
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

// ValidateTitle fails if the provided title
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidateTitle(title string) error {
	// Prevents from counting runes on a large string
	if len(title) > titleMaxChars*4 {
		return errTitleTooLong
	}

	length := utf8.RuneCountInString(title)
	if length < titleMinChars {
		return errTitleTooShort
	}
	if length > titleMaxChars {
		return errTitleTooLong
	}

	for _, c := range title {
		if unicode.IsControl(c) || (unicode.IsSpace(c) && c != ' ') {
			return errTitleInvalid
		}
	}

	if utf8.RuneCountInString(strings.ReplaceAll(title, " ", "")) < titleMinChars {
		return errTitleTooShort
	}

	if strings.HasPrefix(title, " ") || strings.HasSuffix(title, " ") {
		return errTitleInvalid
	}

	return nil
}

// ValidateDescription fails if the provided description
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidateDescription(description string) error {
	// Prevents from counting runes on a large string
	if len(description) > descriptionMaxChars*4 {
		return errDescriptionTooLong
	}

	length := utf8.RuneCountInString(description)
	if length > descriptionMaxChars {
		return errDescriptionTooLong
	}

	for _, c := range description {
		if unicode.IsControl(c) || (unicode.IsSpace(c) && c != ' ') {
			return errDescriptionInvalid
		}
	}

	return nil
}

// validateTrackedWeekDays fails if the provided tracked week days
// do not meet the application requirements.
// The returned error is safe for client-side message.
func validateTrackedWeekDays(trackedWeekDays TrackedWeekDays) error {
	if trackedWeekDays < trackedWeekDaysMin {
		return errTrackedWeekDaysEmpty
	}
	if trackedWeekDays > trackedWeekDaysMax {
		return errTrackedWeekDaysInvalid
	}

	return nil
}

// validateHabitStatus fails if the provided status
// does not meet the application requirements.
// The returned error is safe for client-side message.
func validateHabitStatus(status HabitStatus) error {
	if status < habitStatusMin || status > habitStatusMax {
		return errInternalServer
	}
	return nil
}
