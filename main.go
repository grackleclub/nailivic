package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"text/template"
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
	log, err = newLog(opts)
	log = log.With("service", "nailivic")
	if err != nil {
		panic(err)
	}
}

func main() {
	// read embed dir
	_, err := content.ReadDir(embedDir)
	if err != nil {
		log.Error("failed to read directory",
			"error", err,
		)
		panic(err)
	}

	// routes
	http.HandleFunc("/", logMdl(serveRoot))
	http.Handle("/static/",
		http.FileServer(http.FS(content)),
	)

	// start server
	port := 8008
	log.Info("starting server",
		"port", port,
	)
	fmt.Println("Server is running on port", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		slog.Error("failed to start server",
			"error", err,
			"port", port,
		)
	}
	log.Info("server stopped", "port", port)
}

// log is a middleware that logs all incoming requests
func logMdl(next http.HandlerFunc) http.HandlerFunc {
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

// serveRoot is the base handler for the root (bare) path ("/")
func serveRoot(w http.ResponseWriter, r *http.Request) {
	// order matters (parent -> child)
	tmpl, err := template.ParseFS(content,
		"static/html/index.html",
		"static/html/head.html",
		"static/html/footer.html",
	)
	if err != nil {
		log.Error("failed to parse template",
			"error", err,
		)
		http.Error(w, "parse error", http.StatusInternalServerError)
	}
	// TODO should this be structured and abstracted differently?
	data := struct {
		Name       string
		Title      string
		Stylesheet string
	}{
		Name:       "Nailivic Studios!!",
		Title:      "nailivic",
		Stylesheet: "static/css/style.css",
	}
	log.Debug("templates and data parsed", "data", data)
	w.Header().Set("X-Custom-Header", "special :)")

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Error("failed to execute template",
			"error", err,
		)
		http.Error(w, "serve error", http.StatusInternalServerError)
		return
	}
}
