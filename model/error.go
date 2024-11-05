package model

import "errors"

// errInternalServer represents an internal server error.
// It's used as a safe error whenever unexpected failure occured.
var errInternalServer = errors.New("internal server error")
