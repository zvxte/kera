package model

import (
	"strings"
	"testing"
)

func TestNewHabit(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		weekDays    []WeekDay
		shouldErr   bool
	}{
		{
			"Valid",
			"Title",
			"Description",
			[]WeekDay{Monday, Tuesday, Sunday},
			false,
		},
		{
			"Valid: empty description",
			"Title",
			"",
			[]WeekDay{Friday},
			false,
		},
		{
			"Valid: utf-8",
			"Title 😀",
			"My description 😓",
			[]WeekDay{Wednesday},
			false,
		},
		{
			"Invalid: title",
			"Title\n",
			"",
			[]WeekDay{Sunday},
			true,
		},
		{
			"Invalid: description",
			"Title",
			"\r\n",
			[]WeekDay{Saturday},
			true,
		},
		{
			"Invalid: empty week days",
			"Title",
			"Description",
			[]WeekDay{},
			true,
		},
		{
			"Invalid: week days",
			"Title",
			"Description",
			[]WeekDay{WeekDay(1 << 7)},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewHabit(test.title, test.description, test.weekDays...)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"NewHabit(%q, %q, %v), error=%v, shouldErr=%v",
					test.title, test.description, test.weekDays, err, test.shouldErr,
				)
			}
		})
	}
}

func TestNewTrackedWeekDays(t *testing.T) {
	tests := []struct {
		name      string
		weekDays  []WeekDay
		shouldBe  TrackedWeekDays
		shouldErr bool
	}{
		{
			"Valid",
			[]WeekDay{Monday, Sunday},
			TrackedWeekDays((1 << Monday) | (1 << Sunday)),
			false,
		},
		{
			"Valid",
			[]WeekDay{Tuesday},
			TrackedWeekDays(1 << Tuesday),
			false,
		},
		{
			"Invalid: empty",
			[]WeekDay{},
			TrackedWeekDays(0),
			true,
		},
		{
			"Invalid: day",
			[]WeekDay{WeekDay(1 << 7)},
			TrackedWeekDays(0),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			trackedWeekDays, err := newTrackedWeekDays(test.weekDays...)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"newTrackedWeekDays(%v), error=%v, shouldErr=%v",
					test.weekDays, err, test.shouldErr,
				)
			}
			if trackedWeekDays != test.shouldBe {
				t.Errorf(
					"newTrackedWeekDays(%v), got=%v, expected=%v",
					test.weekDays, trackedWeekDays, test.shouldBe,
				)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		shouldErr bool
	}{
		{"Valid", "Title", false},
		{"Valid", "My title", false},
		{"Valid: short", "aa", false},
		{"Valid: long", strings.Repeat("a", titleMaxChars), false},
		{"Valid: utf-8", "Title 💻π", false},
		{"Invalid: escape character", "Title \t", true},
		{"Invalid: escape character", "Title \n", true},
		{"Invalid: short", "a", true},
		{"Invalid: long", strings.Repeat("a", titleMaxChars+1), true},
		{"Invalid: spaces around", " a ", true},
		{"Invalid: spaces around", " aa ", true},
		{"Invalid: only spaces", "   ", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateTitle(test.title)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"validateTitle(%q), error=%v, shouldErr=%v",
					test.title, err, test.shouldErr,
				)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		shouldErr   bool
	}{
		{"Valid", "Description", false},
		{"Valid", "My description", false},
		{"Valid: empty", "", false},
		{"Valid: short", "aa", false},
		{"Valid: long", strings.Repeat("a", descriptionMaxChars), false},
		{"Valid: utf-8", "💻π", false},
		{"Invalid: escape character", "Description \t", true},
		{"Invalid: escape character", "Description \n", true},
		{"Invalid: too long", strings.Repeat("a", descriptionMaxChars+1), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateDescription(test.description)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"validateDescription(%q), error=%v, shouldErr=%v",
					test.description, err, test.shouldErr,
				)
			}
		})
	}
}

func TestValidateTrackedWeekDays(t *testing.T) {
	tests := []struct {
		name            string
		trackedWeekDays TrackedWeekDays
		shouldErr       bool
	}{
		{
			"Valid: first day",
			TrackedWeekDays(0b_00000001),
			false,
		},
		{
			"Valid: last day",
			TrackedWeekDays(0b_01000000),
			false,
		},
		{
			"Valid: entire week",
			TrackedWeekDays(0b_01111111),
			false,
		},
		{
			"Invalid: day",
			TrackedWeekDays(0b_10000000),
			true,
		},
		{
			"Invalid: empty",
			TrackedWeekDays(0b_00000000),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateTrackedWeekDays(test.trackedWeekDays)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"validateTrackedWeekDays(%v), error=%v, shouldErr=%v",
					test.trackedWeekDays, err, test.shouldErr,
				)
			}
		})
	}
}

func TestValidateHabitStatus(t *testing.T) {
	tests := []struct {
		name      string
		status    HabitStatus
		shouldErr bool
	}{
		{
			"Valid: active",
			HabitActive,
			false,
		},
		{
			"Valid: ended",
			HabitEnded,
			false,
		},
		{
			"Invalid: out of range",
			HabitStatus(2),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateHabitStatus(test.status)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"validateHabitStatus(%q), error=%v, shouldErr=%v",
					test.status, err, test.shouldErr,
				)
			}
		})
	}
}
