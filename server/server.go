package server

import "net/http"

type Server struct {
	authMux     *http.ServeMux
	usersMux    *http.ServeMux
	sessionsMux *http.ServeMux
	habitsMux   *http.ServeMux
}
