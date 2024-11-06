package model

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/zvxte/kera/hash/sha256"
)

const (
	sessionIDLength           = 32
	hashedSessionIDLength     = 32
	sessionIDCharset          = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	sessionIDCharsetLength    = len(sessionIDCharset)
	sessionExpirationDuration = time.Hour * 24 * 30
)

var sessionIDCharsetSet = func() map[rune]bool {
	s := make(map[rune]bool, sessionIDCharsetLength)
	for _, r := range sessionIDCharset {
		s[r] = true
	}
	return s
}()

// HashedSessionID represents a hashed session ID.
type HashedSessionID [hashedSessionIDLength]byte

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
	id := make([]byte, sessionIDLength)
	for i := 0; i < sessionIDLength; i++ {
		randomNum, err := rand.Int(
			rand.Reader, big.NewInt(int64(sessionIDCharsetLength)),
		)
		if err != nil {
			return "", fmt.Errorf("failed to generate session ID: %w", err)
		}
		id[i] = sessionIDCharset[randomNum.Int64()]
	}
	return string(id), nil
}

// ValidateSessionID returns true if the provided sessionID
// meets the application requirements, else false.
func ValidateSessionID(sessionID string) bool {
	if len(sessionID) != sessionIDLength {
		return false
	}
	for _, r := range sessionID {
		if !sessionIDCharsetSet[r] {
			return false
		}
	}

	return true
}
