package handler

import (
	"log"
	"net/http"

	"github.com/zvxte/kera/store"
)

func NewAuthMux(
	userStore store.UserStore,
	sessionStore store.SessionStore,
	logger *log.Logger,
) *http.ServeMux {
	authHandler := NewAuthHandler(userStore, sessionStore, logger)

	authMux := http.NewServeMux()
	authMux.HandleFunc("/login", MakeHandlerFunc(authHandler.Login))
	authMux.HandleFunc("/register", MakeHandlerFunc(authHandler.Register))

	return authMux
}

type AuthHandler struct {
	userStore    store.UserStore
	sessionStore store.SessionStore
	logger       *log.Logger
}

func NewAuthHandler(
	userStore store.UserStore,
	sessionStore store.SessionStore,
	logger *log.Logger,
) AuthHandler {
	return AuthHandler{userStore: userStore, sessionStore: sessionStore, logger: logger}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}
