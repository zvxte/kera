package main

import (
	"encoding/json"
	"net/http"
)

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status int `json:"status"`
	}{200})
}
