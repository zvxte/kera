package user

import "errors"

var (
	ErrUsernameTooShort    = errors.New("username is too short")
	ErrUsernameTooLong     = errors.New("username is too long")
	ErrUsernameInvalid     = errors.New("username is invalid")
	ErrDisplayNameTooShort = errors.New("display name is too short")
	ErrDisplayNameTooLong  = errors.New("display name is too long")
	ErrDisplayNameInvalid  = errors.New("display name is invalid")
	ErrPasswordTooShort    = errors.New("password is too short")
	ErrPasswordTooLong     = errors.New("password is too long")
	ErrLocationInvalid     = errors.New("location is invalid")
)
