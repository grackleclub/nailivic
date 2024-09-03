package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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

// TODO set this somewhere
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
		// read the expected cookie from the 'backend'
		// TODO make this real
		for _, s := range secrets {
			if s.Username == cookie {
				// fetch the user's session cookie from backend storage
				// ensuring they match
				next(w, r)
				return
			}
		}
		http.Error(w, "invalid session cookie", http.StatusBadRequest)
	}
}

var secrets []secret

type secret struct {
	Username      string
	SessionSecret string
}

func getSecret(username string) (string, error) {
	for _, s := range secrets {
		if s.Username == username {
			return s.SessionSecret, nil
		}
	}
	return "", fmt.Errorf("no secret found for user %s", username)
}

func setSecret(username, secret secret) {
	secrets = append(secrets)
}

func newSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
