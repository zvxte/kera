package handler

import (
	"net/http"
)

const (
	SessionIDHeaderName = "session_id"
	UserIDContextKey    = "user_id"
)

type HandlerFuncWithResponse func(http.ResponseWriter, *http.Request) response

func MakeHandlerFunc(f HandlerFuncWithResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := f(w, r)
		if response != nil {
			response.write(w)
		}
	}
}
