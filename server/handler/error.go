package handler

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
