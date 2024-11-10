package session

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	idLen        = 32
	idCharset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	idCharsetLen = len(idCharset)
)

// NewID returns a new randomly generated session ID as a string.
// It fails if the system's source of randomness is unavailable.
func NewID() (string, error) {
	id := make([]byte, idLen)
	for i := 0; i < idLen; i++ {
		randomNum, err := rand.Int(
			rand.Reader, big.NewInt(int64(idCharsetLen)),
		)
		if err != nil {
			return "", fmt.Errorf("failed to generate session ID: %w", err)
		}
		id[i] = idCharset[randomNum.Int64()]
	}
	return string(id), nil
}
