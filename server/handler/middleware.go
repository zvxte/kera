package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/zvxte/kera/hash/sha256"
	"github.com/zvxte/kera/model/date"
	"github.com/zvxte/kera/model/session"
	"github.com/zvxte/kera/store"
)

func SessionMiddleware(next http.Handler, store store.SessionStore) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) response {
		sessionID := r.Header.Get(sessionIDHeaderName)
		if sessionID == "" {
			return unauthorizedResponse
		}

		if !session.ValidateID(sessionID) {
			unsetSessionIDCookie(w)
			return unauthorizedResponse
		}

		hashedSessionID := session.HashedID(sha256.Hash(sessionID))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		session, userID, err := store.Get(ctx, hashedSessionID)
		if err != nil {
			return internalServerErrorResponse
		}

		if session == nil || session.ExpirationDate.Before(date.Now()) {
			unsetSessionIDCookie(w)
			return unauthorizedResponse
		}

		ctx = context.WithValue(r.Context(), userIDContextKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		return nil
	}

	return makeHandlerFunc(f)
}
