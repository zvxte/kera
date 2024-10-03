package main

import (
	"net/http"
)

type Router struct {
	Mux *http.ServeMux
}

func NewRouter() Router {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthcheck", makeHandlerFunc(healthCheck))

	router := Router{mux}
	return router
}
