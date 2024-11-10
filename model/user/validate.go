package user

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	usernameMinChars   = 4
	usernameMaxChars   = 16
	usernameCharset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
	usernameCharsetLen = len(usernameCharset)

	displayNameMinChars = 4
	displayNameMaxChars = 16

	plainPasswordMinChars = 8
	plainPasswordMaxChars = 128
)

var usernameCharsetSet = func() map[rune]bool {
	s := make(map[rune]bool, usernameCharsetLen)
	for _, r := range usernameCharset {
		s[r] = true
	}
	return s
}()

// ValidateUsername fails if the provided username
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidateUsername(username string) error {
	length := len(username)
	if length < usernameMinChars {
		return ErrUsernameTooShort
	}
	if length > usernameMaxChars {
		return ErrUsernameTooLong
	}

	for _, r := range username {
		if !usernameCharsetSet[r] {
			return ErrUsernameInvalid
		}
	}

	underscoreCount := 0
	for i := 0; i < length; i++ {
		if username[i] == '_' {
			underscoreCount++
		}
	}

	if (length - underscoreCount) < usernameMinChars {
		return ErrUsernameInvalid
	}

	return nil
}

// ValidateDisplayName fails if the provided display name
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidateDisplayName(displayName string) error {
	// Prevents from counting runes on a large string
	if len(displayName) > displayNameMaxChars*4 {
		return ErrDisplayNameTooLong
	}

	length := utf8.RuneCountInString(displayName)
	if length < displayNameMinChars {
		return ErrDisplayNameTooShort
	}
	if length > displayNameMaxChars {
		return ErrDisplayNameTooLong
	}

	if !utf8.ValidString(displayName) {
		return ErrDisplayNameInvalid
	}

	for _, c := range displayName {
		if unicode.IsControl(c) || (unicode.IsSpace(c) && c != ' ') {
			return ErrDisplayNameInvalid
		}
	}

	spaceCount := 0
	for _, r := range displayName {
		if r == ' ' {
			spaceCount++
		}
	}
	if (length - spaceCount) < displayNameMinChars {
		return ErrDisplayNameTooShort
	}

	if strings.HasPrefix(displayName, " ") ||
		strings.HasSuffix(displayName, " ") {
		return ErrDisplayNameInvalid
	}

	return nil
}

// ValidatePlainPassword fails if the provided plain password
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidatePlainPassword(plainPassword string) error {
	// Prevents from counting runes on a large string
	if len(plainPassword) > plainPasswordMaxChars*4 {
		return ErrPasswordTooLong
	}

	length := utf8.RuneCountInString(plainPassword)
	if length < plainPasswordMinChars {
		return ErrPasswordTooShort
	}
	if length > plainPasswordMaxChars {
		return ErrPasswordTooLong
	}

	return nil
}

// ValidateLocationName fails if the provided location name
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidateLocationName(locationName string) error {
	_, err := time.LoadLocation(locationName)
	if err != nil {
		return ErrLocationInvalid
	}
	return nil
}
