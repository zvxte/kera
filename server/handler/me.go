package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/zvxte/kera/model"
	"github.com/zvxte/kera/store"
)

type userOut struct {
	Username     string    `json:"username"`
	DisplayName  string    `json:"display_name"`
	Location     string    `json:"location"`
	CreationDate time.Time `json:"creation_date"`
}

func NewMeMux(userStore store.UserStore, logger *log.Logger) *http.ServeMux {
	h := NewMeHandler(userStore, logger)

	m := http.NewServeMux()
	m.HandleFunc("GET /{$}", MakeHandlerFunc(h.Get))
	return m
}

type MeHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewMeHandler(userStore store.UserStore, logger *log.Logger) *MeHandler {
	return &MeHandler{userStore: userStore, logger: logger}
}

func (h *MeHandler) Get(w http.ResponseWriter, r *http.Request) (int, error) {
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
		// TODO: unset the session cookie
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
