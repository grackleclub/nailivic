package main

import (
	"net/http"

	cookie "github.com/ddbgio/cookie"
)

// logMW is a middleware that logs all incoming requests
func logMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("request received",
			"method", r.Method,
			"remote", r.RemoteAddr,
			"path", r.URL.Path,
			"refer", r.Referer(),
			"user-agent", r.UserAgent(),
			"opaque", r.URL.Opaque,
			"bytes", r.ContentLength,
		)
		next(w, r)
	}
}

const cookieName string = "session"

// authMW is a middleware that checks cookie for authentication
func authMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := cookie.ReadEncrypted(r, cookieName, cookieSecret)
		if err != nil {
			log.Error("failed to read cookie",
				"error", err,
			)
			http.Error(w, "invalid session cookie", http.StatusBadRequest)
		}
		log.DebugContext(r.Context(), "encrypted cookie read", "value", cookie)
		if err != nil {
			log.InfoContext(r.Context(), "cookie validation failed",
				"error", err,
				"host", r.Host,
				"remote", r.RemoteAddr,
				"path", r.URL.Path,
			)
			http.Error(w, "invalid session cookie", http.StatusBadRequest)
		}

		next(w, r)
	}
}
