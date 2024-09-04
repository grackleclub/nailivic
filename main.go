package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"

	cookie "github.com/ddbgio/cookie/v2"
	logger "github.com/ddbgio/log"
)

const embedDir = "static"

//go:embed static
var content embed.FS // object representing the embedded directory

var log *slog.Logger

func init() {
	// setup logger
	opts := slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	var err error
	log, err = logger.New(opts)
	log = log.With("service", "nailivic")
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	// setup cookies
	// TODO this must be changed if more than one server is ever active
	cookieGlobalSecret, err = cookie.NewCookieSecret()
	if err != nil {
		log.Error("failed to generate secret",
			"error", err,
		)
		panic(err)
	}
	log.Warn("secret set", "secret", string(cookieGlobalSecret))

	// read embed dir
	_, err = content.ReadDir(embedDir)
	if err != nil {
		log.Error("failed to read directory",
			"error", err,
		)
		panic(err)
	}

	// setup database
	// TODO extract ddbgio/db into it's own package,
	// and use that for the backend

	// ROUTES
	// full pages
	http.HandleFunc("/secret", logMW(authMW(serveSecret)))
	http.HandleFunc("/crazy", logMW(serveCrazy))
	http.HandleFunc("/", logMW(serveRoot))
	// http.HandleFunc("/about", logMW(serveAbout)) // TODO (maybe?)
	// htmx components
	http.HandleFunc("/htmx/{component}", logMW(serveHtmx))
	// static files
	http.Handle("/static/",
		http.FileServer(http.FS(content)),
	)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// this doesn't work
		log.Debug("favicon requested")
		http.ServeFile(w, r, "static/img/favicon.ico")
	})

	// start server
	port := 8008
	log.Info("starting server",
		"port", port,
	)
	log.Info("server is running",
		"port", port,
	)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		slog.Error("failed to start server",
			"error", err,
			"port", port,
		)
	}
	log.Info("server stopped", "port", port)
}
