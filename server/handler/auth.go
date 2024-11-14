package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/zvxte/kera/hash/argon2id"
	"github.com/zvxte/kera/model/session"
	"github.com/zvxte/kera/model/user"
	"github.com/zvxte/kera/store/sessionstore"
	"github.com/zvxte/kera/store/userstore"
)

type userIn struct {
	Username      string `json:"username"`
	PlainPassword string `json:"password"`
}

func NewAuthMux(
	userStore userstore.Store,
	sessionStore sessionstore.Store,
	logger *log.Logger,
) *http.ServeMux {
	h := &authHandler{
		userStore:    userStore,
		sessionStore: sessionStore,
		logger:       logger,
	}

	m := http.NewServeMux()
	m.HandleFunc("POST /login", makeHandlerFunc(h.Login))
	m.HandleFunc("POST /register", makeHandlerFunc(h.Register))
	return m
}

type authHandler struct {
	userStore    userstore.Store
	sessionStore sessionstore.Store
	logger       *log.Logger
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) response {
	if r.Header.Get("Content-Type") != "application/json" {
		return unsupportedMediaTypeResponse
	}

	var in userIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return badRequestResponse
	}

	if err := user.ValidateUsername(in.Username); err != nil {
		return invalidCredentialsResponse
	}

	if err := user.ValidatePlainPassword(in.PlainPassword); err != nil {
		return invalidCredentialsResponse
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.userStore.Get(
		ctx, userstore.UsernameColumn, in.Username,
	)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	if user == nil {
		return invalidCredentialsResponse
	}

	isValid, err := argon2id.VerifyHash(in.PlainPassword, user.HashedPassword)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}
	if !isValid {
		return invalidCredentialsResponse
	}

	sessionID, err := session.NewID()
	if err != nil {
		return internalServerErrorResponse
	}

	session := session.New(sessionID, user.ID)

	err = h.sessionStore.Create(ctx, session)
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	setSessionIDCookie(w, sessionID, time.Time(session.ExpirationDate))
	return noContentResponse{}
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) response {
	if r.Header.Get("Content-Type") != "application/json" {
		return unsupportedMediaTypeResponse
	}

	var in userIn
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		return badRequestResponse
	}

	user, err := user.New(in.Username, in.PlainPassword)
	if err != nil {
		return newJsonResponse(
			http.StatusBadRequest,
			newHandlerError(http.StatusBadRequest, err.Error()),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.userStore.Create(ctx, user)
	if err == userstore.ErrUsernameAlreadyTaken {
		return usernameAlreadyTakenResponse
	}
	if err != nil {
		h.logger.Println(err)
		return internalServerErrorResponse
	}

	return createdResponse{}
}
