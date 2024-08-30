package main

import (
	"bytes"
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

// example of composite types a template might receive page
// and parse page.nav.bar, page.head, page.footer, etc
type page struct {
	nav    navbar
	head   head
	footer footer
}
type navbar struct{}
type head struct{}
type footer struct{}

func newPage() (page, error) {
	// return page{
	// 	nav:    navbar{},
	// 	head:   head{
	// 		Title: "nailivic",

	// 	},
	// 	footer: footer{
	// 		Year: 2021,
	// 	},
	// }, nil
	return page{}, fmt.Errorf("not implemented")
}

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

	// ROUTES
	// full pages
	http.HandleFunc("/", logMW(serveRoot))
	http.HandleFunc("/login", logMW(serveLogin))
	// http.HandleFunc("/about", logMW(serveAbout)) // TODO (maybe?)
	// htmx components
	http.HandleFunc("/htmx/{component}", logMW(serveHtmx))
	// static files
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

type index struct {
	Name       string
	Title      string
	Stylesheet string
}

func serveLogin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusTeapot)
}

// serveRoot is the base handler for the root (bare) path ("/")
func serveRoot(w http.ResponseWriter, r *http.Request) {
	// order matters (parent -> child)
	templates := []string{
		"static/html/index.html",
		"static/html/head.html",
		"static/html/footer.html",
	}
	data := index{
		Name:       "Nailivic Studios!!",
		Title:      "nailivic",
		Stylesheet: "static/css/style.css",
	}
	w.Header().Set("X-Custom-Header", "special :)")
	err := writeTemplate(w, templates, data)
	if err != nil {
		log.Error("failed to write template",
			"error", err,
			"templates", templates,
		)
	}
	log.Debug("root served", "templates", templates)
}

// serveHtmx dynamically serves htmx components based on the path
func serveHtmx(w http.ResponseWriter, r *http.Request) {
	// get the component name from the path
	componentName := r.PathValue("component")
	log.Info("htmx component requested", "name", componentName)

	// serve the appropriate htmx component based on name from path
	var err error
	w.Header().Set("X-htmx-component-name", componentName)
	switch componentName {
	case "special":
		err = writeTemplate(w, []string{"static/html/special.html"}, nil)
	default:
		http.Error(w, "missing or invalid htmx component name", http.StatusBadRequest)
	}
	log.Warn("wha happen")
	if err != nil {
		log.Error("failed to write htmx component",
			"error", err,
			"component", componentName,
		)
	}
}

func writeTemplate(w http.ResponseWriter, templatePaths []string, data interface{}) error {
	tmpl, err := template.ParseFS(content, templatePaths...)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	// write to buffer first to allow inspection
	// because if a child template is called before a parent template,
	// the output will be empty
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	if buf.Len() == 0 {
		return fmt.Errorf("template output is empty")
	}
	b, err := buf.WriteTo(w)
	if err != nil {
		return fmt.Errorf("failed to write template to response: %w", err)
	}
	log.Debug("template executed",
		"bytes_read", buf.Len(),
		"bytes_written", b,
	)
	return nil
}
