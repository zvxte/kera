package model

import (
	"errors"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
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
)

type User struct {
	ID             UUID
	Username       string
	DisplayName    string
	HashedPassword string
	Location       *time.Location
	CreationDate   time.Time
}

func NewUser(username, plainPassword string) (User, error) {
	if err := validateUsername(username); err != nil {
		return User{}, err
	}

	if err := validatePlainPassword(plainPassword); err != nil {
		return User{}, err
	}

	id, err := NewUUIDv7()
	if err != nil {
		return User{}, errInternalServer
	}

	displayName := username

	hashedPassword := plainPassword

	location := time.UTC

	creationDate := time.Now().UTC()

	return User{
		ID:             id,
		Username:       username,
		DisplayName:    displayName,
		HashedPassword: hashedPassword,
		Location:       location,
		CreationDate:   creationDate,
	}, nil
}

func LoadUser(
	id UUID,
	username, displayName, hashedPassword string,
	location *time.Location,
	creationDate time.Time,
) (User, error) {
	if err := validateUsername(username); err != nil {
		return User{}, err
	}

	if err := validateDisplayName(displayName); err != nil {
		return User{}, err
	}

	if location == nil {
		return User{}, errInternalServer
	}

	return User{
		ID:             id,
		Username:       username,
		DisplayName:    displayName,
		HashedPassword: hashedPassword,
		Location:       location,
		CreationDate:   creationDate,
	}, nil
}

func validateUsername(username string) error {
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

func validateDisplayName(displayName string) error {
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

func validatePlainPassword(plainPassword string) error {
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
