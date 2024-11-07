package model

import (
	"errors"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/zvxte/kera/hash/argon2id"
)

const (
	usernameMinChars      = 4
	usernameMaxChars      = 16
	displayNameMinChars   = 4
	displayNameMaxChars   = 16
	plainPasswordMinChars = 8
	plainPasswordMaxChars = 128
)

var (
	errUsernameTooShort    = errors.New("username is too short")
	errUsernameTooLong     = errors.New("username is too long")
	errUsernameInvalid     = errors.New("username is invalid")
	errDisplayNameTooShort = errors.New("display name is too short")
	errDisplayNameTooLong  = errors.New("display name is too long")
	errDisplayNameInvalid  = errors.New("display name is invalid")
	errPasswordTooShort    = errors.New("password is too short")
	errPasswordTooLong     = errors.New("password is too long")
	errLocationInvalid     = errors.New("location is invalid")
)

// User represents an application user.
type User struct {
	ID             UUID
	Username       string
	DisplayName    string
	HashedPassword string
	Location       *time.Location
	CreationDate   time.Time
}

// NewUser returns a new *User.
// It fails if the provided parameters do not meet the application requirements.
// The returned error is safe for client-side message.
// The plain password is hashed using Argon2ID.
// The Username and DisplayName fields are set to the given username,
// Location is set to UTC, and CreationDate is set to the current date in UTC.
func NewUser(username, plainPassword string) (*User, error) {
	if err := ValidateUsername(username); err != nil {
		return nil, err
	}

	if err := ValidatePlainPassword(plainPassword); err != nil {
		return nil, err
	}

	id, err := NewUUIDv7()
	if err != nil {
		return nil, errInternalServer
	}

	displayName := username

	hashedPassword, err := argon2id.Hash(plainPassword, argon2id.DefaultParams)
	if err != nil {
		return nil, errInternalServer
	}

	location := time.UTC

	creationDate := DateNow()

	return &User{
		ID:             id,
		Username:       username,
		DisplayName:    displayName,
		HashedPassword: hashedPassword,
		Location:       location,
		CreationDate:   creationDate,
	}, nil
}

// LoadUser returns a *User.
// It fails if the provided parameters do not meet the application requirements.
// The returned error is safe for client-side message.
// It returns errInternalServer if provided location is nil.
func LoadUser(
	id UUID,
	username, displayName, hashedPassword string,
	location *time.Location,
	creationDate time.Time,
) (*User, error) {
	if err := ValidateUsername(username); err != nil {
		return nil, err
	}

	if err := ValidateDisplayName(displayName); err != nil {
		return nil, err
	}

	if location == nil {
		return nil, errInternalServer
	}

	return &User{
		ID:             id,
		Username:       username,
		DisplayName:    displayName,
		HashedPassword: hashedPassword,
		Location:       location,
		CreationDate:   creationDate,
	}, nil
}

// ValidateUsername fails if the provided username
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidateUsername(username string) error {
	length := len(username)
	if length < usernameMinChars {
		return errUsernameTooShort
	}
	if length > usernameMaxChars {
		return errUsernameTooLong
	}

	for _, c := range username {
		if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '_' {
			return errUsernameInvalid
		}
	}

	if len(strings.ReplaceAll(username, "_", "")) < usernameMinChars {
		return errUsernameInvalid
	}

	return nil
}

// ValidateDisplayName fails if the provided display name
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidateDisplayName(displayName string) error {
	// Prevents from counting runes on a large string
	if len(displayName) > displayNameMaxChars*4 {
		return errDisplayNameTooLong
	}

	length := utf8.RuneCountInString(displayName)
	if length < displayNameMinChars {
		return errDisplayNameTooShort
	}
	if length > displayNameMaxChars {
		return errDisplayNameTooLong
	}

	for _, c := range displayName {
		if unicode.IsControl(c) || (unicode.IsSpace(c) && c != ' ') {
			return errDisplayNameInvalid
		}
	}

	if utf8.RuneCountInString(strings.ReplaceAll(displayName, " ", "")) < usernameMinChars {
		return errDisplayNameTooShort
	}

	if strings.HasPrefix(displayName, " ") || strings.HasSuffix(displayName, " ") {
		return errDisplayNameInvalid
	}

	return nil
}

// ValidatePlainPassword fails if the provided plain password
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidatePlainPassword(plainPassword string) error {
	// Prevents from counting runes on a large string
	if len(plainPassword) > plainPasswordMaxChars*4 {
		return errPasswordTooLong
	}

	length := utf8.RuneCountInString(plainPassword)
	if length < plainPasswordMinChars {
		return errPasswordTooShort
	}
	if length > plainPasswordMaxChars {
		return errPasswordTooLong
	}

	return nil
}

// ValidateLocationName fails if the provided location name
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidateLocationName(locationName string) error {
	_, err := time.LoadLocation(locationName)
	if err != nil {
		return errLocationInvalid
	}
	return nil
}
