package argon2id

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	Variant = "argon2id"
	Version = argon2.Version

	InputMinLength = 0
	InputMaxLength = (1 << 32) - 1

	MemoryMax      = (1 << 32) - 1
	IterationsMin  = 1
	IterationsMax  = (1 << 32) - 1
	ParallelismMin = 1
	ParallelismMax = (1 << 24) - 1

	KeyMinLength  = 4
	KeyMaxLength  = (1 << 32) - 1
	SaltMinLength = 8
	SaltMaxLength = (1 << 32) - 1
)

var MemoryMin = func(parallelism uint8) uint32 {
	return uint32(8 * parallelism)
}

var (
	ErrInvalidVariant     = errors.New("argon2id: variant is invalid")
	ErrInvalidVersion     = errors.New("argon2id: version is invalid")
	ErrInvalidInput       = errors.New("argon2id: input is invalid")
	ErrInvalidHashedInput = errors.New("argon2id: hashed input is invalid")
	ErrNilParamsPointer   = errors.New("argon2id: function called with nil Params pointer")
	ErrInvalidParams      = errors.New("argon2id: params are invalid")
)

var DefaultParams = &Params{
	memory:      19 * 1024,
	iterations:  2,
	parallelism: 1,
	keyLength:   32,
	saltLength:  16,
}

type Params struct {
	memory      uint32 // In KiB (mebibytes)
	iterations  uint32
	parallelism uint8
	keyLength   uint32 // In B (bytes)
	saltLength  uint32 // In B (bytes)
}

func NewParams(
	memory, iterations uint32,
	parallelism uint8,
	keyLength, saltLength uint32,
) (*Params, error) {
	if parallelism < ParallelismMin {
		return nil, ErrInvalidParams
	}

	if memory < MemoryMin(parallelism) || memory > MemoryMax {
		return nil, ErrInvalidParams
	}

	if iterations < IterationsMin || iterations > IterationsMax {
		return nil, ErrInvalidParams
	}

	if keyLength < KeyMinLength || keyLength > KeyMaxLength {
		return nil, ErrInvalidParams
	}

	if saltLength < SaltMinLength || saltLength > SaltMaxLength {
		return nil, ErrInvalidParams
	}

	return &Params{memory, iterations, parallelism, keyLength, saltLength}, nil
}

func Hash(input string, params *Params) (string, error) {
	inputLength := len(input)
	if inputLength < InputMinLength || inputLength > InputMaxLength {
		return "", ErrInvalidInput
	}

	if params == nil {
		return "", ErrNilParamsPointer
	}

	salt, err := generateSalt(params.saltLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	key := argon2.IDKey(
		[]byte(input),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)

	output := fmt.Sprintf(
		"$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		Variant, Version, params.memory, params.iterations, params.parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	)
	return output, nil
}

func VerifyHash(input, hashedInput string) (bool, error) {
	params, salt, otherKey, err := decodeHash(hashedInput)
	if err != nil {
		return false, err
	}

	key := argon2.IDKey(
		[]byte(input),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)

	if result := subtle.ConstantTimeEq(int32(len(key)), int32(len(otherKey))); result != 1 {
		return false, nil
	}

	if result := subtle.ConstantTimeCompare(key, otherKey); result != 1 {
		return false, nil
	}

	return true, nil
}

func decodeHash(hashedInput string) (*Params, []byte, []byte, error) {
	parts := strings.Split(hashedInput, "$")

	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHashedInput
	}

	if parts[1] != Variant {
		return nil, nil, nil, ErrInvalidVariant
	}

	var v uint
	_, err := fmt.Sscanf(parts[2], "v=%d", &v)
	if err != nil {
		return nil, nil, nil, ErrInvalidHashedInput
	}
	if v != Version {
		return nil, nil, nil, ErrInvalidVersion
	}

	var m, t uint32
	var p uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p)
	if err != nil {
		return nil, nil, nil, ErrInvalidHashedInput
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, ErrInvalidHashedInput
	}

	key, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, ErrInvalidHashedInput
	}

	params, err := NewParams(
		m, t, p,
		uint32(len(key)),
		uint32(len(salt)),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	return params, salt, key, nil
}

func generateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
