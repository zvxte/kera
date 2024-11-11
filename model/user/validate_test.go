package user

import (
	"strings"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		shouldErr bool
	}{
		{
			"Valid",
			"username",
			false,
		},
		{
			"Valid",
			"user_name",
			false,
		},
		{
			"Invalid: too short",
			"aaa",
			true,
		},
		{
			"Invalid: too short",
			"_aaa_",
			true,
		},
		{
			"Invalid: too long",
			"aaaaAAAAaaaaAAAAa",
			true,
		},
		{
			"Invalid: only underscores",
			"____",
			true,
		},
		{
			"Invalid: character",
			"aaaa.",
			true,
		},
		{
			"Invalid: space",
			"aa aa",
			true,
		},
		{
			"Invalid: byte sequence",
			"user\x80name",
			true,
		},
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
		{
			"Valid",
			"display name",
			false,
		},
		{
			"Valid",
			"display name!",
			false,
		},
		{
			"Valid",
			"display_name 💻",
			false,
		},
		{
			"Invalid: too short",
			"aaa",
			true,
		},
		{
			"Invalid: too short",
			"a a a",
			true,
		},
		{
			"Invalid: too long",
			"aaaaAAAAaaaaAAAAa",
			true,
		},
		{
			"Invalid: only spaces",
			"    ",
			true,
		},
		{
			"Invalid: escape character",
			"display name \t",
			true,
		},
		{
			"Invalid: byte sequence",
			"display \x80 name",
			true,
		},
		{
			"Invalid: spaces around",
			" display name ",
			true,
		},
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
		{
			"Valid",
			"password",
			false,
		},
		{
			"Valid",
			"aa aa aa",
			false,
		},
		{
			"Valid",
			"display_name 💻 \t",
			false,
		},
		{
			"Invalid: too short",
			"aaaaAAA",
			true,
		},
		{
			"Invalid: too long",
			strings.Repeat("a", plainPasswordMaxChars+1),
			true,
		},
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
