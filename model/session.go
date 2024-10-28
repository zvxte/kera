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

type HashedSessionID [32]byte

type Session struct {
	HashedID       HashedSessionID
	CreationDate   time.Time
	ExpirationDate time.Time
}

func NewSession(sessionID string) *Session {
	hashedID := sha256.Hash(sessionID)

	creationDate := dateNow()

	expirationDate := creationDate.Add(sessionExpirationDuration)

	return &Session{
		HashedID:       hashedID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
	}
}

func LoadSession(hashedID HashedSessionID, creationDate, expirationDate time.Time) *Session {
	return &Session{
		HashedID:       hashedID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
	}
}

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
