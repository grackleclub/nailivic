package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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

// serveRoot is the base handler for the root (bare) path ("/")
func serveRoot(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(content,
		"static/html/index.html",
		"static/html/head.html",
		"static/html/footer.html",
	)
	log.Debug("templates paresed")
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
	log.Debug("data parsed", "data", data)
	w.Header().Set("Special-Status", "super special")

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Error("failed to execute template",
			"error", err,
		)
		http.Error(w, "failed to execute template", http.StatusInternalServerError)
	}
}
