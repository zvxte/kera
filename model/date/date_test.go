package date

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	warsawLoc, _ := time.LoadLocation("Europe/Warsaw")

	tests := []struct {
		name     string
		time     time.Time
		expected Date
	}{
		{
			"Valid: with time portion",
			time.Date(2024, 11, 10, 2, 6, 0, 0, time.UTC),
			Date(time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)),
		},
		{
			"Valid: not UTC location",
			time.Date(2024, 11, 10, 2, 6, 0, 0, warsawLoc),
			Date(time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)),
		},
		{
			"Valid: zero value",
			time.Time{},
			Date{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Load(test.time)
			if got != test.expected {
				t.Errorf(
					"Load(%q), got=%q, expected=%q",
					test.time, got, test.expected,
				)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		date     Date
		duration time.Duration
		expected Date
	}{
		{
			"Valid",
			Date(time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)),
			24 * time.Hour,
			Date(time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC)),
		},
		{
			"Valid: with time portion",
			Date(time.Date(2024, 11, 10, 2, 6, 0, 0, time.UTC)),
			24 * time.Hour,
			Date(time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.date.Add(test.duration)
			if got != test.expected {
				t.Errorf(
					"Date(%q).Add(%v), got=%q, expected=%q",
					test.date, test.duration, got, test.expected,
				)
			}
		})
	}
}

func TestBefore(t *testing.T) {
	tests := []struct {
		name     string
		date     Date
		other    Date
		expected bool
	}{
		{
			"Valid: before by one day",
			Date(time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)),
			Date(time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC)),
			true,
		},
		{
			"Valid: before by two days",
			Now(),
			Now().Add(48 * time.Hour),
			true,
		},
		{
			"Invalid: after three days",
			Date(time.Date(2024, 11, 10, 0, 0, 0, 0, time.UTC)),
			Date(time.Date(2024, 11, 7, 0, 0, 0, 0, time.UTC)),
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.date.Before(test.other)
			if got != test.expected {
				t.Errorf(
					"Date(%q).Before(%q), got=%v, expected=%v",
					test.date, test.other, got, test.expected,
				)
			}
		})
	}
}
