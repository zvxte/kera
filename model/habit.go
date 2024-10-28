package model

import (
	"errors"
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

type Habit struct {
	ID              UUID
	Status          HabitStatus
	Title           string
	Description     string
	TrackedWeekDays TrackedWeekDays
	History         HabitHistory
	StartDate       time.Time
	EndDate         time.Time
}

type HabitStatus uint8

const (
	HabitActive HabitStatus = iota
	HabitEnded
)

type TrackedWeekDays uint8

type WeekDay uint8

const (
	Monday WeekDay = 1 << iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type HabitHistory []HabitDay

type HabitDay struct {
	Status HabitDayStatus
	Date   time.Time
}

type HabitDayStatus uint8

const (
	DayUntracked HabitDayStatus = iota
	DayDone
	DayMissed
	DayPending
)

func NewHabit(
	title, description string,
	weekDays ...WeekDay,
) (*Habit, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	if err := validateDescription(description); err != nil {
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

	startDate := dateNow()

	endDate := time.Time{}

	history := HabitHistory{}

	return &Habit{
		ID:              id,
		Status:          status,
		Title:           title,
		Description:     description,
		TrackedWeekDays: trackedWeekDays,
		History:         history,
		StartDate:       startDate,
		EndDate:         endDate,
	}, nil
}

func newTrackedWeekDays(weekDays ...WeekDay) (TrackedWeekDays, error) {
	var trackedWeekDays TrackedWeekDays
	for _, day := range weekDays {
		trackedWeekDays |= TrackedWeekDays(day)
	}

	err := validateTrackedWeekDays(trackedWeekDays)
	if err != nil {
		return 0, err
	}

	return trackedWeekDays, nil
}

func validateTitle(title string) error {
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

func validateDescription(description string) error {
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

func validateTrackedWeekDays(trackedWeekDays TrackedWeekDays) error {
	if trackedWeekDays < trackedWeekDaysMin {
		return errTrackedWeekDaysEmpty
	}
	if trackedWeekDays > trackedWeekDaysMax {
		return errTrackedWeekDaysInvalid
	}

	return nil
}
