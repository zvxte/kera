package handler

type handlerError struct {
	Message string `json:"message"`
}

func newHandlerError(message string) handlerError {
	return handlerError{Message: message}
}

func (handlerError handlerError) Error() string {
	return handlerError.Message
}
