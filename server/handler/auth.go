package handler

import (
	"log"
	"net/http"

	"github.com/zvxte/kera/store"
)

type AuthHandler struct {
	userStore    *store.UserStore
	sessionStore *store.SessionStore
	logger       *log.Logger
}

func NewAuthHandler(
	userStore *store.UserStore,
	sessionStore *store.SessionStore,
	logger *log.Logger,
) *AuthHandler {
	return &AuthHandler{userStore: userStore, sessionStore: sessionStore, logger: logger}
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) (int, error) {
	return 200, nil
}
