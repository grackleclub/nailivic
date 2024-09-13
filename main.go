package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"

	cookie "github.com/ddbgio/cookie/v2"
	"github.com/ddbgio/healthz"
	logger "github.com/ddbgio/log"
	"github.com/ddbgio/nailivic/db/sqlc"
	"github.com/ddbgio/postgres"
)

const embedDir = "static"

//go:embed static
var content embed.FS // object representing the embedded directory

var (
	log     *slog.Logger  // global logger
	queries *sqlc.Queries // global sqlc queries
)

func main() {
	// setup logger
	var logOpts slog.HandlerOptions
	_, ok := os.LookupEnv("DEBUG")
	if ok {
		logOpts = slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}
	} else {
		logOpts = slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		}
	}
	var err error
	log, err = logger.New(logOpts)
	log = log.With("service", "nailivic")
	if err != nil {
		panic(err)
	}
	log.Info("logger initialized", "debug", ok)

	// setup cookies
	// TODO this must be changed if more than one server is ever active
	cookieGlobalSecret, err = cookie.NewCookieSecret()
	if err != nil {
		log.Error("failed to generate secret",
			"error", err,
		)
		panic(err)
	}
	log.Info("parent cookie secret set")

	// read embed dir
	_, err = content.ReadDir(embedDir)
	if err != nil {
		log.Error("failed to read directory",
			"error", err,
		)
		panic(err)
	}
	log.Info("embed directory read", "dir", embedDir)

	// setup database
	ctx := context.Background()
	pgOpts := postgres.PostgresOpts{
		Host:     "localhost",
		User:     "postgres",
		Password: "mysecretpassword",
		Name:     "postgres",
		Sslmode:  "disable",
	}
	db, err := postgres.NewDB(ctx, pgOpts)
	if err != nil {
		log.Error("failed to create database",
			"error", err,
		)
		panic(err)
	}
	defer db.Close()
	// run migrations to prepare database
	migrationDir := path.Join("db", "migrations")
	migrations, err := postgres.Migrations(migrationDir, "up")
	if err != nil {
		log.Error("failed to read migrations",
			"error", err,
			"dir", migrationDir,
		)
		panic(err)
	}
	for _, m := range migrations {
		// stripped := strings.ReplaceAll(m.Content, "\n", " ")
		log.Debug("running migration",
			"migration", m.Filename,
			"direction", m.Direction,
			"content", m.Content,
		)
		result, err := db.Query(ctx, m.Content)
		if err != nil {
			log.Error("failed to run migration",
				"error", err,
				"migration", m.Filename,
				"direction", m.Direction,
			)
			panic(err)
		}
		log.Info("ran migration",
			"migration", m.Filename,
			"direction", m.Direction,
			"result", result,
		)
	}
	log.Info("db setup; migrations ran")

	// prepare sqlc queries
	pool, err := db.Pool(ctx)
	if err != nil {
		log.Error("failed to create connection pool",
			"error", err,
		)
		panic(err)
	}
	defer pool.Close()
	queries := sqlc.New(pool)

	user, err := queries.UserByID(ctx, 1)
	if err != nil {
		log.Error("failed to get user",
			"error", err,
		)
		panic(err)
	}
	log.Info("got user", "user", user)

	// ROUTES
	// full pages
	http.HandleFunc("/secret", logMW(authMW(serveSecret)))
	http.HandleFunc("/parts", logMW(serveParts))
	http.HandleFunc("/inventory", logMW(serveInventory))
	http.HandleFunc("/", logMW(serveRoot))
	http.HandleFunc("/healthz", logMW(healthz.Respond))
	// htmx components
	http.HandleFunc("/htmx/{component}", logMW(serveHtmx))
	// static files
	http.Handle("/static/",
		http.FileServer(http.FS(content)),
	)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
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
