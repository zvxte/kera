package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/zvxte/kera/hash/sha256"
	"github.com/zvxte/kera/model"
	"github.com/zvxte/kera/store"
)

type userOut struct {
	Username     string    `json:"username"`
	DisplayName  string    `json:"display_name"`
	Location     string    `json:"location"`
	CreationDate time.Time `json:"creation_date"`
}

func NewMeMux(
	userStore store.UserStore, sessionStore store.SessionStore, logger *log.Logger,
) *http.ServeMux {
	h := &meHandler{
		userStore:    userStore,
		sessionStore: sessionStore,
		logger:       logger,
	}

	m := http.NewServeMux()
	m.HandleFunc("GET /{$}", MakeHandlerFunc(h.Get))
	m.HandleFunc("POST /logout", MakeHandlerFunc(h.Logout))
	return m
}

type meHandler struct {
	userStore    store.UserStore
	sessionStore store.SessionStore
	logger       *log.Logger
}

func (h *meHandler) Get(w http.ResponseWriter, r *http.Request) (int, error) {
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

	err = jsonResponse(w, http.StatusOK, userOut{
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

func (h *meHandler) Logout(w http.ResponseWriter, r *http.Request) (int, error) {
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
