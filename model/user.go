package model

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	usernameMinChars       = 4
	usernameMaxChars       = 16
	displayNameMinChars    = 4
	displayNameMaxChars    = 16
	plainPasswordMinChars  = 8
	plainPasswordMaxChars  = 128
	hashedPasswordMinChars = 8
	hashedPasswordMaxChars = 256
	timezoneNameMaxChars   = 64
)

type User struct {
	ID             UUID
	Username       string
	DisplayName    string
	HashedPassword string
	Location       time.Location
	CreationDate   time.Time
}

func isUsernameValid(username string) bool {
	length := len(username)
	if length < usernameMinChars || length > usernameMaxChars {
		return false
	}

	for _, c := range username {
		if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '_' {
			return false
		}
	}

	if len(strings.ReplaceAll(username, "_", "")) < usernameMinChars {
		return false
	}

	return true
}

func isDisplayNameValid(displayName string) bool {
	// Prevents from counting runes on a large string
	if len(displayName) > displayNameMaxChars*4 {
		return false
	}

	length := utf8.RuneCountInString(displayName)
	if length < displayNameMinChars || length > displayNameMaxChars {
		return false
	}

	for _, c := range displayName {
		if unicode.IsControl(c) || (unicode.IsSpace(c) && c != ' ') {
			return false
		}
	}

	if utf8.RuneCountInString(strings.ReplaceAll(displayName, " ", "")) < usernameMinChars {
		return false
	}

	if strings.HasPrefix(displayName, " ") || strings.HasSuffix(displayName, " ") {
		return false
	}

	return true
}

func IsPlainPasswordValid(plainPassword string) bool {
	// Prevents from counting runes on a large string
	if len(plainPassword) > plainPasswordMaxChars*4 {
		return false
	}

	length := utf8.RuneCountInString(plainPassword)
	if length < plainPasswordMinChars || length > plainPasswordMaxChars {
		return false
	}

	return true
}

func isHashedPasswordValid(hashedPassword string) bool {
	// Prevents from counting runes on a large string
	if len(hashedPassword) > hashedPasswordMaxChars*4 {
		return false
	}

	length := utf8.RuneCountInString(hashedPassword)
	if length < hashedPasswordMinChars || length > hashedPasswordMaxChars {
		return false
	}

	return true
}

func isLocationValid(location string) bool {
	if len(location) > timezoneNameMaxChars {
		return false
	}

	if _, err := time.LoadLocation(location); err != nil {
		return false
	}

	return true
}
