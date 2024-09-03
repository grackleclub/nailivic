package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

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

const cookieName string = "nailivic-session"

var secretsBackend []sessionKey

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
		cookie, err := cookie.ReadEncrypted(r, cookieName, cookieSecret)
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
		log.DebugContext(r.Context(), "encrypted cookie read", "value", cookie)
		// read the expected cookie from the 'backend'
		// TODO make this real
		// fetch the user's session cookie from backend storage
		// ensuring they match
		for _, s := range secretsBackend {
			// TODO check user and key, not just key
			if s.SessionSecret == cookie {
				if s.isExpired() {
					log.InfoContext(r.Context(), "session key expired",
						"username", s.Username,
					)
					http.Error(w, "session key expired", http.StatusForbidden)
					return
				}
				log.Info("session authenticated", "username", s.Username)
				next(w, r)
				return
			}
		}
		http.Error(w, "invalid session cookie", http.StatusBadRequest)
	}
}
