package date

import "errors"

var (
	ErrInvalidYear  = errors.New("invalid year")
	ErrInvalidMonth = errors.New("invalid month")
)
