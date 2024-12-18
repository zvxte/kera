package handler

import "errors"

var (
	ErrInternalServer       = errors.New("internal server error")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrBadRequest           = errors.New("bad request")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUsernameAlreadyTaken = errors.New("username is already taken")
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
