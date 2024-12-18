package handler

import (
	"net/http"
)

const (
	sessionIDHeaderName = "session_id"
	userIDContextKey    = "user_id"
)

type handlerFuncWithResponse func(http.ResponseWriter, *http.Request) response

func makeHandlerFunc(f handlerFuncWithResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := f(w, r)
		if response != nil {
			response.write(w)
		}
	}
}
