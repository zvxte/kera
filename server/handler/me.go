package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/zvxte/kera/hash/argon2id"
	"github.com/zvxte/kera/hash/sha256"
	"github.com/zvxte/kera/model/session"
	"github.com/zvxte/kera/model/user"
	"github.com/zvxte/kera/model/uuid"
	"github.com/zvxte/kera/store/sessionstore"
	"github.com/zvxte/kera/store/userstore"
)

func NewMeMux(
	userStore userstore.Store, sessionStore sessionstore.Store, logger *log.Logger,
) *http.ServeMux {
	h := &meHandler{
		userStore:    userStore,
		sessionStore: sessionStore,
		logger:       logger,
	}

	m := http.NewServeMux()
	m.HandleFunc("GET /{$}", makeHandlerFunc(h.get))
	m.HandleFunc("DELETE /{$}", makeHandlerFunc(h.delete))
	m.HandleFunc("PATCH /display-name", makeHandlerFunc(h.patchDisplayName))
	m.HandleFunc("PATCH /password", makeHandlerFunc(h.patchPassword))
	m.HandleFunc("POST /logout", makeHandlerFunc(h.logout))
	m.HandleFunc("GET /sessions", makeHandlerFunc(h.getSessionsCount))
	m.HandleFunc("DELETE /sessions", makeHandlerFunc(h.deleteSessions))
	return m
}

type meHandler struct {
	userStore    userstore.Store
	sessionStore sessionstore.Store
	logger       *log.Logger
}

func (h *meHandler) get(w http.ResponseWriter, r *http.Request) response {
	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.userStore.Get(ctx, userstore.IDColumn, userID)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	if user == nil {
		unsetSessionIDCookie(w)
		return unauthorizedResponse
	}

	return newJsonResponse(
		http.StatusOK,
		struct {
			Username     string    `json:"username"`
			DisplayName  string    `json:"display_name"`
			CreationDate time.Time `json:"creation_date"`
		}{
			Username:     user.Username,
			DisplayName:  user.DisplayName,
			CreationDate: time.Time(user.CreationDate),
		},
	)
}

func (h *meHandler) delete(w http.ResponseWriter, r *http.Request) response {
	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.userStore.Delete(ctx, userID)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	unsetSessionIDCookie(w)
	return noContentResponse{}
}

func (h *meHandler) patchDisplayName(w http.ResponseWriter, r *http.Request) response {
	if r.Header.Get("Content-Type") != "application/json" {
		return unsupportedMediaTypeResponse
	}

	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	var in struct {
		DisplayName string `json:"display_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return badRequestResponse
	}

	err := user.ValidateDisplayName(in.DisplayName)
	if err != nil {
		return newJsonResponse(
			http.StatusBadRequest,
			newHandlerError(http.StatusBadRequest, err.Error()),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.userStore.Update(
		ctx, userID, userstore.DisplayNameColumn, in.DisplayName,
	)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	return noContentResponse{}
}

func (h *meHandler) patchPassword(w http.ResponseWriter, r *http.Request) response {
	if r.Header.Get("Content-Type") != "application/json" {
		return unsupportedMediaTypeResponse
	}

	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	var in struct {
		PlainPassword    string `json:"password"`
		NewPlainPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return badRequestResponse
	}

	err := user.ValidatePlainPassword(in.NewPlainPassword)
	if err != nil {
		return newJsonResponse(
			http.StatusBadRequest,
			newHandlerError(http.StatusBadRequest, err.Error()),
		)
	}

	err = user.ValidatePlainPassword(in.PlainPassword)
	if err != nil {
		return invalidCredentialsResponse
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.userStore.Get(ctx, userstore.IDColumn, userID)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	isValid, err := argon2id.VerifyHash(
		in.PlainPassword, user.HashedPassword,
	)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}
	if !isValid {
		return invalidCredentialsResponse
	}

	newHashedPassword, err := argon2id.Hash(
		in.NewPlainPassword, argon2id.DefaultParams,
	)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	err = h.userStore.Update(
		ctx, user.ID, userstore.HashedPasswordColumn, newHashedPassword,
	)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	return noContentResponse{}
}

func (h *meHandler) logout(w http.ResponseWriter, r *http.Request) response {
	go func() {
		sessionID := r.Header.Get(sessionIDHeaderName)
		if sessionID == "" {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		hashedSessionID := session.HashedID(sha256.Hash(sessionID))

		err := h.sessionStore.Delete(
			ctx, sessionstore.HashedIDColumn, hashedSessionID,
		)
		if err != nil {
			h.logger.Println(err)
		}
	}()

	unsetSessionIDCookie(w)
	return noContentResponse{}
}

func (h *meHandler) getSessionsCount(w http.ResponseWriter, r *http.Request) response {
	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := h.sessionStore.Count(ctx, userID)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	return newJsonResponse(
		http.StatusOK,
		struct {
			Count uint `json:"count"`
		}{Count: count},
	)
}

func (h *meHandler) deleteSessions(w http.ResponseWriter, r *http.Request) response {
	userID, ok := r.Context().Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return internalServerErrorResponse
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.sessionStore.Delete(ctx, sessionstore.UserIDColumn, userID)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	unsetSessionIDCookie(w)
	return noContentResponse{}
}
