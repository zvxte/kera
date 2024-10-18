package model

import (
	"strings"
	"time"
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

func IsUsernameValid(username string) bool {
	length := len(username)
	if length < usernameMinChars || length > usernameMaxChars {
		return false
	}

	for _, c := range username {
		if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && (c < '0' || c > '9') && c != '_' {
			return false
		}
	}

	if len(strings.Trim(username, "_")) == 0 {
		return false
	}

	return true
}
