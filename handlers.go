package main

import (
	"encoding/json"
	"errors"
	"log"
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

func jsonResponse[T any](w http.ResponseWriter, status int, value T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(value)
}

type handlerFuncWithError func(w http.ResponseWriter, r *http.Request) (int, error)

func makeHandlerFunc(f handlerFuncWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if status, err := f(w, r); err != nil {
			err = jsonResponse(w, status, newHandlerError(err.Error()))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, errors.New("method not allowed")
	}
	w.WriteHeader(http.StatusOK)
	return http.StatusOK, nil
}
