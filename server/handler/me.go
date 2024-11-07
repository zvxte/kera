package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/zvxte/kera/hash/argon2id"
	"github.com/zvxte/kera/hash/sha256"
	"github.com/zvxte/kera/model"
	"github.com/zvxte/kera/store"
)

func NewMeMux(
	userStore store.UserStore, sessionStore store.SessionStore, logger *log.Logger,
) *http.ServeMux {
	h := &meHandler{
		userStore:    userStore,
		sessionStore: sessionStore,
		logger:       logger,
	}

	m := http.NewServeMux()
	m.HandleFunc("GET /{$}", MakeHandlerFunc(h.get))
	m.HandleFunc("PATCH /display-name", MakeHandlerFunc(h.patchDisplayName))
	m.HandleFunc("PATCH /location", MakeHandlerFunc(h.patchLocation))
	m.HandleFunc("PATCH /password", MakeHandlerFunc(h.patchPassword))
	m.HandleFunc("POST /logout", MakeHandlerFunc(h.logout))
	return m
}

type meHandler struct {
	userStore    store.UserStore
	sessionStore store.SessionStore
	logger       *log.Logger
}

func (h *meHandler) get(w http.ResponseWriter, r *http.Request) (int, error) {
	userID, ok := r.Context().Value(UserIDContextKey).(model.UUID)
	if !ok {
		return http.StatusInternalServerError, ErrInternalServer
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.userStore.Get(ctx, userID)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}

	if user == nil {
		unsetSessionIDCookie(w)
		return http.StatusUnauthorized, ErrUnauthorized
	}

	err = jsonResponse(w, http.StatusOK, struct {
		Username     string    `json:"username"`
		DisplayName  string    `json:"display_name"`
		Location     string    `json:"location"`
		CreationDate time.Time `json:"creation_date"`
	}{
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		Location:     user.Location.String(),
		CreationDate: user.CreationDate,
	})
	if err != nil {
		return http.StatusInternalServerError, ErrInternalServer
	}

	return http.StatusOK, nil
}

func (h *meHandler) patchDisplayName(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusUnsupportedMediaType, ErrUnsupportedMediaType
	}

	userID, ok := r.Context().Value(UserIDContextKey).(model.UUID)
	if !ok {
		return http.StatusInternalServerError, ErrInternalServer
	}

	var in struct {
		DisplayName string `json:"display_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return http.StatusBadRequest, ErrBadRequest
	}

	err := model.ValidateDisplayName(in.DisplayName)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.userStore.UpdateDisplayName(ctx, userID, in.DisplayName)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}

	return http.StatusNoContent, nil
}

func (h *meHandler) patchLocation(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusUnsupportedMediaType, ErrUnsupportedMediaType
	}

	userID, ok := r.Context().Value(UserIDContextKey).(model.UUID)
	if !ok {
		return http.StatusInternalServerError, ErrInternalServer
	}

	var in struct {
		Location string `json:"location"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return http.StatusBadRequest, ErrBadRequest
	}

	err := model.ValidateLocationName(in.Location)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}

	location, _ := time.LoadLocation(in.Location)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.userStore.UpdateLocation(ctx, userID, location)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}

	return http.StatusNoContent, nil
}

func (h *meHandler) patchPassword(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusUnsupportedMediaType, ErrUnsupportedMediaType
	}

	userID, ok := r.Context().Value(UserIDContextKey).(model.UUID)
	if !ok {
		return http.StatusInternalServerError, ErrInternalServer
	}

	var in struct {
		PlainPassword    string `json:"password"`
		NewPlainPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return http.StatusBadRequest, ErrBadRequest
	}

	err := model.ValidatePlainPassword(in.NewPlainPassword)
	if err != nil {
		return http.StatusUnprocessableEntity, err
	}

	err = model.ValidatePlainPassword(in.PlainPassword)
	if err != nil {
		return http.StatusUnprocessableEntity, ErrInvalidCredentials
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.userStore.Get(ctx, userID)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}

	isValid, err := argon2id.VerifyHash(in.PlainPassword, user.HashedPassword)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}
	if !isValid {
		return http.StatusUnprocessableEntity, ErrInvalidCredentials
	}

	newHashedPassword, err := argon2id.Hash(in.NewPlainPassword, argon2id.DefaultParams)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}

	err = h.userStore.UpdateHashedPassword(ctx, user.ID, newHashedPassword)
	if err != nil {
		h.logger.Println(err)
		return http.StatusInternalServerError, ErrInternalServer
	}

	return http.StatusNoContent, nil
}

func (h *meHandler) logout(w http.ResponseWriter, r *http.Request) (int, error) {
	unsetSessionIDCookie(w)
	w.WriteHeader(http.StatusNoContent)

	go func() {
		sessionID := r.Header.Get(SessionIDHeaderName)
		if sessionID == "" {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		hashedSessionID := model.HashedSessionID(sha256.Hash(sessionID))

		err := h.sessionStore.Delete(ctx, hashedSessionID)
		if err != nil {
			h.logger.Println(err)
		}
	}()

	return http.StatusNoContent, nil
}

var sessionIDUnsetCookie = &http.Cookie{
	Name:     "session_id",
	Value:    "",
	Path:     "/",
	Secure:   true,
	HttpOnly: true,
	Expires:  time.Time{},
	MaxAge:   -1,
	SameSite: http.SameSiteStrictMode,
}

func unsetSessionIDCookie(w http.ResponseWriter) {
	http.SetCookie(w, sessionIDUnsetCookie)
}
