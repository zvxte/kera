package handler

import (
	"log"
	"net/http"

	"github.com/zvxte/kera/store"
)

func NewMeMux(userStore store.UserStore, logger *log.Logger) *http.ServeMux {
	_ = NewMeHandler(userStore, logger)

	meMux := http.NewServeMux()

	return meMux
}

type MeHandler struct {
	userStore    store.UserStore
	sessionStore store.SessionStore
	logger       *log.Logger
}

func NewMeHandler(userStore store.UserStore, logger *log.Logger) AuthHandler {
	return AuthHandler{userStore: userStore, logger: logger}
}
