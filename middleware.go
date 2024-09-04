package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	cookie "github.com/ddbgio/cookie/v2"
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
		next.ServeHTTP(w, r)
	}
}

const cookieName string = "nailivic-session"

// var secretsBackend []sessionKey

var backend = make(map[int]sessionKey)

const sessionDefaultExpiry = 1 * time.Hour

type sessionKey struct {
	Username      string
	SessionSecret string
	Exipiry       time.Time
}

func (s sessionKey) isExpired() bool {
	return s.Exipiry.Before(time.Now())
}

func newSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// authMW is a middleware that checks cookie for authentication
func authMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, clientSessionToken, err := cookie.ReadEncrypted(r, cookieName, cookieSecret)
		if err != nil {
			log.InfoContext(r.Context(), "cookie validation failed",
				"error", err,
				"host", r.Host,
				"remote", r.RemoteAddr,
				"path", r.URL.Path,
			)
			http.Error(w, "invalid session cookie", http.StatusBadRequest)
			return
		}
		log.DebugContext(r.Context(), "encrypted cookie read", "value", clientSessionToken)
		// read the expected cookie from the 'backend'
		// TODO make this real
		// fetch the user's session cookie from backend storage
		// ensuring they match
		key, ok := backend[userID]
		if !ok {
			log.InfoContext(r.Context(), "session key not found",
				"userID", userID,
			)
			// TODO reconsider this
			http.Error(w, "invalid session cookie", http.StatusBadRequest)
			return
		}
		if key.isExpired() {
			log.InfoContext(r.Context(), "session key expired",
				"userID", userID,
			)
			http.Error(w, "session key expired", http.StatusForbidden)
			return
		}
		if key.SessionSecret != clientSessionToken {
			log.InfoContext(r.Context(), "session key mismatch",
				"userID", userID,
			)
			http.Error(w, "invalid session cookie", http.StatusForbidden)
			return
		}
		log.Info("session authenticated",
			"id", userID,
			"username", key.Username,
			"expiry", key.Exipiry,
			"token", clientSessionToken,
		)
		next.ServeHTTP(w, r)

		// http.Error(w, "invalid session cookie", http.StatusBadRequest)
	}
}
