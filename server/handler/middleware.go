package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/zvxte/kera/hash/sha256"
	"github.com/zvxte/kera/model"
	"github.com/zvxte/kera/store"
)

func SessionMiddleware(next http.Handler, store store.SessionStore) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) response {
		sessionID := r.Header.Get(SessionIDHeaderName)
		if sessionID == "" {
			return unauthorizedResponse
		}

		if !model.ValidateSessionID(sessionID) {
			unsetSessionIDCookie(w)
			return unauthorizedResponse
		}

		hashedSessionID := model.HashedSessionID(sha256.Hash(sessionID))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		session, userID, err := store.Get(ctx, hashedSessionID)
		if err != nil {
			return internalServerErrorResponse
		}

		if session == nil || session.ExpirationDate.Before(model.DateNow()) {
			unsetSessionIDCookie(w)
			return unauthorizedResponse
		}

		ctx = context.WithValue(r.Context(), UserIDContextKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		return nil
	}

	return MakeHandlerFunc(f)
}
