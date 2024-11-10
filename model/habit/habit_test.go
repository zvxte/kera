package habit

import (
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
			_, err := New(test.title, test.description, test.weekDays...)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"New(%q, %q, %v), error=%v, shouldErr=%v",
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
