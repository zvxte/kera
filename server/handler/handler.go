package handler

import "net/http"

const (
	SessionIDHeaderName = "session_id"
	UserIDContextKey    = "user_id"
)

type HandlerFuncWithError func(http.ResponseWriter, *http.Request) (int, error)

func MakeHandlerFunc(f HandlerFuncWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode, err := f(w, r)
		if err != nil {
			// Encoding HandlerError into JSON will never fail
			_ = jsonResponse(w, statusCode, newHandlerError(statusCode, err.Error()))
		}
	}
}
