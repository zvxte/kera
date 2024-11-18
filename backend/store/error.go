package store

import "errors"

var (
	ErrNilDB              = errors.New("function called with nil *sql.DB")
	ErrInvalidColumn      = errors.New("column is invalid")
	ErrInvalidColumnValue = errors.New("column value is invalid")
)
