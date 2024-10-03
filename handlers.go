package main

import (
	"encoding/json"
	"net/http"
)

type handlerError struct {
	Message string `json:"message"`
}

func (handlerError handlerError) Error() string {
	return handlerError.Message
}

func newHandlerError(message string) handlerError {
	return handlerError{Message: message}
}

type handlerFuncWithError func(w http.ResponseWriter, r *http.Request) (int, error)

func makeHandlerFunc(f handlerFuncWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if status, err := f(w, r); err != nil {
			jsonResponse(w, status, newHandlerError(err.Error()))
		}
	}
}

func jsonResponse[T any](w http.ResponseWriter, status int, value T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(value)
}

func healthCheck(w http.ResponseWriter, _ *http.Request) (int, error) {
	w.WriteHeader(http.StatusOK)
	return http.StatusOK, nil
}
