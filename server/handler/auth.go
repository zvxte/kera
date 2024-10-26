package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/zvxte/kera/model"
	"github.com/zvxte/kera/store"
)

var (
	ErrUsernameAlreadyTaken = errors.New("username is already taken")
)

type userIn struct {
	Username      string `json:"username"`
	PlainPassword string `json:"plain_password"`
}

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
	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, ErrMethodNotAllowed
	}
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusUnsupportedMediaType, ErrUnsupportedMediaType
	}

	var in userIn
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		return http.StatusBadRequest, ErrBadRequest
	}

	user, err := model.NewUser(in.Username, in.PlainPassword)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	isTaken, err := h.userStore.IsTaken(ctx, user.Username)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}
	if isTaken {
		return http.StatusConflict, ErrUsernameAlreadyTaken
	}

	err = h.userStore.Create(ctx, user)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}

	return http.StatusCreated, nil
}
