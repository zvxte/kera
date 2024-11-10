package session

import (
	"time"

	"github.com/zvxte/kera/hash/sha256"
	"github.com/zvxte/kera/model/date"
)

const expirationDuration = 24 * time.Hour * 30

// HashedID represents a hashed session ID.
type HashedID [32]byte

// Session represents a user's session in the application.
type Session struct {
	HashedID       HashedID
	CreationDate   date.Date
	ExpirationDate date.Date
}

// New returns a new *Session.
// The sessionID is hashed using sha256.
// The CreationDate field is set to the current Date value.
// The ExpirationDate field is set to the CreationDate + sessionExpirationDuration constant.
func New(sessionID string) *Session {
	hashedID := sha256.Hash(sessionID)

	creationDate := date.Now()
	expirationDate := creationDate.Add(expirationDuration)

	return &Session{
		HashedID:       hashedID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
	}
}

// Load returns a *Session from provided parameters.
func Load(
	hashedID HashedID, creationDate, expirationDate date.Date,
) *Session {
	return &Session{
		HashedID:       hashedID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
	}
}
