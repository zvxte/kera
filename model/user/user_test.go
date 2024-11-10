package user

import (
	"testing"
	"time"

	"github.com/zvxte/kera/model/date"
	"github.com/zvxte/kera/model/uuid"
)

func TestNew(t *testing.T) {
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
			"Invalid: username",
			"aaa",
			"password",
			true,
		},
		{
			"Invalid: password",
			"username",
			"aaa",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.username, test.plainPassword)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"New(%q, %q), error=%v, shouldErr=%v",
					test.username, test.plainPassword, err, test.shouldErr,
				)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		id             uuid.UUID
		username       string
		displayName    string
		hashedPassword string
		location       *time.Location
		creationDate   date.Date
		shouldErr      bool
	}{
		{
			"Valid",
			uuid.UUID{},
			"username",
			"display name",
			"hashed password",
			time.UTC,
			date.Now(),
			false,
		},
		{
			"Invalid: username",
			uuid.UUID{},
			"aaa",
			"display name",
			"hashed password",
			time.UTC,
			date.Now(),
			true,
		},
		{
			"Invalid: display name",
			uuid.UUID{},
			"username",
			"  display name  ",
			"hashed password",
			time.UTC,
			date.Now(),
			true,
		},
		{
			"Invalid: location",
			uuid.UUID{},
			"username",
			"display name",
			"hashed password",
			nil,
			date.Now(),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Load(
				test.id, test.username, test.displayName,
				test.hashedPassword, test.location, test.creationDate,
			)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"Load(%q, %q, %q, %q, %q, %q), error=%v, shouldErr=%v",
					test.id, test.username, test.displayName,
					test.hashedPassword, test.location, test.creationDate,
					err, test.shouldErr,
				)
			}
		})
	}
}
