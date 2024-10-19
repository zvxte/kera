package model

import (
	"errors"
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

var (
	errInvalidUsername       = errors.New("invalid username")
	errInvalidDisplayName    = errors.New("invalid display name")
	errInvalidHashedPassword = errors.New("invalid hashed password")
	errInvalidTimezoneName   = errors.New("invalid timezone name")
)

type User struct {
	ID             UUID
	Username       string
	DisplayName    string
	HashedPassword string
	Location       *time.Location
	CreationDate   time.Time
}

func NewUser(
	id UUID,
	username, displayName, hashedPassword, timezoneName string,
	creationDate time.Time,
) (User, error) {
	if !isUsernameValid(username) {
		return User{}, errInvalidUsername
	}
	if !isDisplayNameValid(displayName) {
		return User{}, errInvalidDisplayName
	}
	if !isHashedPasswordValid(hashedPassword) {
		return User{}, errInvalidHashedPassword
	}
	if !isTimezoneNameValid(timezoneName) {
		return User{}, errInvalidTimezoneName
	}
	location, _ := time.LoadLocation(timezoneName)
	return User{
		ID:             id,
		Username:       username,
		DisplayName:    displayName,
		HashedPassword: hashedPassword,
		Location:       location,
		CreationDate:   creationDate,
	}, nil
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

func isTimezoneNameValid(timezoneName string) bool {
	if len(timezoneName) > timezoneNameMaxChars {
		return false
	}

	if _, err := time.LoadLocation(timezoneName); err != nil {
		return false
	}

	return true
}
