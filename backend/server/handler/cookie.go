package handler

import (
	"net/http"
	"time"
)

const sessionIDCookieName = "session_id"

var sessionIDUnsetCookie = &http.Cookie{
	Name:     sessionIDCookieName,
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

func setSessionIDCookie(w http.ResponseWriter, sessionID string, expirationDate time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionIDCookieName,
		Value:    sessionID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  expirationDate,
		SameSite: http.SameSiteStrictMode,
	})
}
