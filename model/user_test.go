package model

import (
	"strings"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		plainPassword string
		shouldErr     bool
	}{
		{
			"Valid",
			"username",
			"password",
			false,
		},
		{
			"Invalid username",
			"aaa",
			"password",
			true,
		},
		{
			"Invalid password",
			"username",
			"aaa",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewUser(test.username, test.plainPassword)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"NewUser(%q, %q), error=%v, shouldErr=%v",
					test.username, test.plainPassword, err, test.shouldErr,
				)
			}
		})
	}
}

func TestLoadUser(t *testing.T) {
	tests := []struct {
		name           string
		id             UUID
		username       string
		displayName    string
		hashedPassword string
		location       *time.Location
		creationDate   time.Time
		shouldErr      bool
	}{
		{
			"Valid",
			UUID{},
			"username",
			"display name",
			"hashed password",
			time.UTC,
			time.Now().UTC(),
			false,
		},
		{
			"Invalid username",
			UUID{},
			"aaa",
			"display name",
			"hashed password",
			time.UTC,
			time.Now().UTC(),
			true,
		},
		{
			"Invalid display name",
			UUID{},
			"username",
			"  display name  ",
			"hashed password",
			time.UTC,
			time.Now().UTC(),
			true,
		},
		{
			"Invalid location",
			UUID{},
			"username",
			"display name",
			"hashed password",
			nil,
			time.Now().UTC(),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := LoadUser(
				test.id, test.username, test.displayName,
				test.hashedPassword, test.location, test.creationDate,
			)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"LoadUser(%q, %q, %q, %q, %q, %q), error=%v, shouldErr=%v",
					test.id, test.username, test.displayName,
					test.hashedPassword, test.location, test.creationDate,
					err, test.shouldErr,
				)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		shouldErr bool
	}{
		{"Valid", "username", false},
		{"Valid", "user_name", false},
		{"Too short", "aaa", true},
		{"Too short", "_aaa_", true},
		{"Too long", "aaaaAAAAaaaaAAAAa", true},
		{"Invalid", "____", true},
		{"Invalid", "aaaa.", true},
		{"Invalid", "aa aa", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateUsername(test.username)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"ValidateUsername(%q), error=%v, shouldErr=%v",
					test.username, err, test.shouldErr,
				)
			}
		})
	}
}

func TestValidateDisplayName(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		shouldErr   bool
	}{
		{"Valid", "display name", false},
		{"Valid", "display name!", false},
		{"Valid", "display_name 💻", false},
		{"Too short", "aaa", true},
		{"Too short", "a a a", true},
		{"Too short", "    ", true},
		{"Too long", "aaaaAAAAaaaaAAAAa", true},
		{"Invalid", "display name \t", true},
		{"Invalid", " display name ", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateDisplayName(test.displayName)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"ValidateDisplayName(%q), error=%v, shouldErr=%v",
					test.displayName, err, test.shouldErr,
				)
			}
		})
	}
}

func TestValidatePlainPassword(t *testing.T) {
	tests := []struct {
		name          string
		plainPassword string
		shouldErr     bool
	}{
		{"Valid", "password", false},
		{"Valid", "aa aa aa", false},
		{"Valid", "display_name 💻 \t", false},
		{"Too short", "aaaaAAA", true},
		{"Too long", strings.Repeat("a", plainPasswordMaxChars+1), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidatePlainPassword(test.plainPassword)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"ValidatePlainPassword(%q), error=%v, shouldErr=%v",
					test.plainPassword, err, test.shouldErr,
				)
			}
		})
	}
}
