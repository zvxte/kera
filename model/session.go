package model

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/zvxte/kera/hash/sha256"
)

const (
	sessionCharset            = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	sessionCharsetLength      = len(sessionCharset)
	sessionExpirationDuration = time.Hour * 24 * 30
)

// HashedSessionID represents a hashed session ID.
type HashedSessionID [32]byte

// Session represents a user's session in the application.
type Session struct {
	HashedID       HashedSessionID
	CreationDate   time.Time
	ExpirationDate time.Time
}

// NewSession returns a new *Session.
// The sessionID is hashed using sha256.
// The CreationDate field is set to the current date in UTC.
// The ExpirationDate field is set to the CreationDate + sessionExpirationDuration constant.
func NewSession(sessionID string) *Session {
	hashedID := sha256.Hash(sessionID)

	creationDate := DateNow()

	expirationDate := creationDate.Add(sessionExpirationDuration)

	return &Session{
		HashedID:       hashedID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
	}
}

// LoadSession returns a *Session from provided parameters.
func LoadSession(hashedID HashedSessionID, creationDate, expirationDate time.Time) *Session {
	return &Session{
		HashedID:       hashedID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
	}
}

// NewSessionID returns a new randomly generated session ID as a string.
// It fails if the system's source of randomness is unavailable.
func NewSessionID() (string, error) {
	id := make([]byte, 32)
	for i := 0; i < 32; i++ {
		randomNum, err := rand.Int(
			rand.Reader, big.NewInt(int64(sessionCharsetLength)),
		)
		if err != nil {
			return "", fmt.Errorf("failed to generate session ID: %w", err)
		}
		id[i] = sessionCharset[randomNum.Int64()]
	}
	return string(id), nil
}
