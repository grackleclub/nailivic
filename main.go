package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
)

const embedDir = "temp"

//go:embed temp
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
	fs, err := content.ReadDir(embedDir)
	if err != nil {
		log.Error("failed to read directory",
			"error", err,
		)
		panic(err)
	}

	// log found files in embed dir
	for _, f := range fs {
		// read that file
		path := path.Join(embedDir, f.Name())
		fileContents, err := content.ReadFile(path)
		if err != nil {
			log.Error("failed to read file",
				"error", err,
				"file", f.Name(),
			)
			continue
		}
		log.Debug("file found",
			"file", f.Name(),
			"contents", string(fileContents),
		)
	}

	// start server
	port := 8008
	log.Info("starting server",
		"port", port,
	)
	fmt.Println("Server is running on port", port)

	// routes
	http.HandleFunc("/", serveRoot)
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.FS(content)),
		),
	)

	// start server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		slog.Error("failed to start server",
			"error", err,
			"port", port,
		)
	}
	log.Info("server stopped", "port", port)
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	log.Debug("serving response",
		"request", r.URL.Path,
		"method", r.Method,
		"remote", r.RemoteAddr,
	)

	// respond with "Hello, World!"
	fmt.Fprintf(w, "Hello, World!")
}
