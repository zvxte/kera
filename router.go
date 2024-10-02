package main

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /healthcheck", healthCheck)
	return router
}
