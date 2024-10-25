package sha256

import (
	"crypto/sha256"
)

func Hash(input string) [32]byte {
	return sha256.Sum256([]byte(input))
}

func VerifyHash(input string, hashedInput [32]byte) bool {
	return sha256.Sum256([]byte(input)) == hashedInput
}
