package handler

import "errors"

var (
	ErrInternalServer       = errors.New("internal server error")
	ErrMethodNotAllowed     = errors.New("method not allowed")
	ErrBadRequest           = errors.New("bad request")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	ErrUnauthorized         = errors.New("unauthorized")
)

var (
	ErrUsernameAlreadyTaken = errors.New("username is already taken")
	ErrInvalidCredentials   = errors.New("invalid credentials")
)

type handlerError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func newHandlerError(statusCode int, message string) handlerError {
	return handlerError{StatusCode: statusCode, Message: message}
}

func (handlerError handlerError) Error() string {
	return handlerError.Message
}
