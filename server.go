package main

import "net/http"

type Server struct {
	Address string
	Router  Router
}

func NewServer(address string, router Router) Server {
	return Server{address, router}
}

func (s Server) Run() error {
	return http.ListenAndServe(s.Address, s.Router.Mux)
}
