package store

import "errors"

var NilDBPointerError = errors.New("function called with nil pointer to sql.DB")
