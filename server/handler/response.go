package handler

import (
	"encoding/json"
	"net/http"
)

func jsonResponse[T any](w http.ResponseWriter, statusCode int, body T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(body)
}
