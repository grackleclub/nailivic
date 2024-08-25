package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"text/template"
)

const embedDir = "static"

//go:embed static
var content embed.FS // object representing the embedded directory

var log *slog.Logger

func init() {
	// setup logger
	opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stderr, &opts)
	log = slog.New(handler)
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

	// log found files in embed dir
	// for _, f := range fs {
	// 	// read that file
	// 	path := path.Join(embedDir, f.Name())
	// 	fileContents, err := content.ReadFile(path)
	// 	if err != nil {
	// 		log.Error("failed to read file",
	// 			"error", err,
	// 			"file", f.Name(),
	// 		)
	// 		continue
	// 	}
	// 	log.Debug("file found",
	// 		"file", f.Name(),
	// 		"contents", string(fileContents),
	// 	)
	// }

	// routes
	http.HandleFunc("/", logMiddleware(serveRoot))
	http.Handle("/static/",
		http.FileServer(http.FS(content)),
	)
	// templates test
	http.HandleFunc("/templates", logMiddleware(serveTemplates))

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

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("request received",
			"method", r.Method,
			"remote", r.RemoteAddr,
			"path", r.URL.Path,
			"refer", r.Referer(),
			"user-agent", r.UserAgent(),
			"opaque", r.URL.Opaque,
		)
		next(w, r)
	}
}

func serveTemplates(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFS(content,
		"static/templates/index.html",
		"static/templates/header.html",
	)
	log.Info("templates paresed")
	if err != nil {
		log.Error("failed to parse template",
			"error", err,
		)
		http.Error(w, "failed to parse template", http.StatusInternalServerError)
	}
	data := struct {
		Name  string
		Title string
	}{
		Name:  "Nailivic Studios!!",
		Title: "-nailivic-",
	}
	log.Info("data parsed", "data", data)
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Error("failed to execute template",
			"error", err,
		)
		http.Error(w, "failed to execute template", http.StatusInternalServerError)
	}
}

// serveRoot is the base handler for the root (bare) path ("/")
func serveRoot(w http.ResponseWriter, r *http.Request) {
	log.Debug("serving response",
		"request", r.URL.Path,
		"method", r.Method,
		"remote", r.RemoteAddr,
	)

	// get /static/html/index.html
	fileName := "index.html"
	path := path.Join("static", "html", fileName)
	file, err := content.ReadFile(path)
	if err != nil {
		log.Error("failed to read file",
			"error", err,
			"file", path,
		)
		msg := fmt.Sprintf("failed to read file %s", fileName)
		http.Error(w, msg, http.StatusInternalServerError)
	}
	w.Header().Set("Special-Status", "super special")

	size, err := w.Write(file)
	if err != nil {
		log.Error("failed to write response",
			"error", err,
			"file", path,
			"bytes", size,
		)
	}
	log.Warn("response written", "bytes", size)
}
