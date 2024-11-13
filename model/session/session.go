package session

import (
	"time"

	"github.com/zvxte/kera/hash/sha256"
	"github.com/zvxte/kera/model/date"
	"github.com/zvxte/kera/model/uuid"
)

const (
	HashedIDLen        = 32
	expirationDuration = 24 * time.Hour * 30
)

// HashedID represents a hashed session ID.
type HashedID [HashedIDLen]byte

// Session represents a user's session in the application.
type Session struct {
	HashedID       HashedID
	UserID         uuid.UUID
	CreationDate   date.Date
	ExpirationDate date.Date
}

// New returns a new *Session.
// The sessionID is hashed using sha256.
// The CreationDate field is set to the current Date value.
// The ExpirationDate field is set to the CreationDate + sessionExpirationDuration constant.
func New(sessionID string, userID uuid.UUID) *Session {
	hashedID := sha256.Hash(sessionID)

	creationDate := date.Now()
	expirationDate := creationDate.Add(expirationDuration)

	return &Session{
		HashedID:       hashedID,
		UserID:         userID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
	}
}

// Load returns a *Session from provided parameters.
func Load(
	hashedID HashedID, userID uuid.UUID, creationDate, expirationDate date.Date,
) *Session {
	return &Session{
		HashedID:       hashedID,
		UserID:         userID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
	}
}
