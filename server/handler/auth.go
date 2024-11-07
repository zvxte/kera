package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/zvxte/kera/hash/argon2id"
	"github.com/zvxte/kera/model"
	"github.com/zvxte/kera/store"
)

type userIn struct {
	Username      string `json:"username"`
	PlainPassword string `json:"password"`
}

func NewAuthMux(
	userStore store.UserStore,
	sessionStore store.SessionStore,
	logger *log.Logger,
) *http.ServeMux {
	h := &authHandler{
		userStore:    userStore,
		sessionStore: sessionStore,
		logger:       logger,
	}

	m := http.NewServeMux()
	m.HandleFunc("POST /login", MakeHandlerFunc(h.Login))
	m.HandleFunc("POST /register", MakeHandlerFunc(h.Register))

	return m
}

type authHandler struct {
	userStore    store.UserStore
	sessionStore store.SessionStore
	logger       *log.Logger
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusUnsupportedMediaType, ErrUnsupportedMediaType
	}

	var in userIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return http.StatusBadRequest, ErrBadRequest
	}

	if err := model.ValidateUsername(in.Username); err != nil {
		return http.StatusUnprocessableEntity, ErrInvalidCredentials
	}

	if err := model.ValidatePlainPassword(in.PlainPassword); err != nil {
		return http.StatusUnprocessableEntity, ErrInvalidCredentials
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.userStore.GetByUsername(ctx, in.Username)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}

	if user == nil {
		return http.StatusUnprocessableEntity, ErrInvalidCredentials
	}

	isValid, err := argon2id.VerifyHash(in.PlainPassword, user.HashedPassword)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}

	if !isValid {
		return http.StatusUnprocessableEntity, ErrInvalidCredentials
	}

	sessionID, err := model.NewSessionID()
	if err != nil {
		return http.StatusInternalServerError, ErrInternalServer
	}

	session := model.NewSession(sessionID)

	err = h.sessionStore.Create(ctx, session, user.ID)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  session.ExpirationDate,
		SameSite: http.SameSiteStrictMode,
	})
	return http.StatusNoContent, nil
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) (int, error) {
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
