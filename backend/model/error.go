package model

import "errors"

// ErrUnexpected represents an unexpected error that occurred.
// It's used as an error that can be safe for a client-side message.
var ErrUnexpected = errors.New(
	"an unexpected error occurred, please try again later",
)
