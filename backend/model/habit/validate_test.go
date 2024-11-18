package habit

import (
	"strings"
	"testing"
)

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

func TestValidateStatus(t *testing.T) {
	tests := []struct {
		name      string
		status    Status
		shouldErr bool
	}{
		{
			"Valid: active",
			Active,
			false,
		},
		{
			"Valid: ended",
			Ended,
			false,
		},
		{
			"Invalid: out of range",
			Status(2),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateStatus(test.status)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"validateStatus(%q), error=%v, shouldErr=%v",
					test.status, err, test.shouldErr,
				)
			}
		})
	}
}
