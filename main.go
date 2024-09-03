package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"

	cookie "github.com/ddbgio/cookie"
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
	var err error
	// setup cookies
	// TODO this must be changed if more than one server is ever active
	cookieSecret, err := cookie.NewCookieSecret()
	if err != nil {
		log.Error("failed to generate secret",
			"error", err,
		)
		panic(err)
	}
	log.Warn("secret", "secret", string(cookieSecret))

	// read embed dir
	_, err = content.ReadDir(embedDir)
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

type index struct {
	Name        string
	Title       string
	Stylesheets []string // path to stylesheets (in order!)
}

func serveLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// serve login page
		templates := []string{
			"static/html/login.html",
			"static/html/head.html",
		}
		data := index{
			Name:  "Nailivic Studios Login",
			Title: "nailivic",
			Stylesheets: []string{
				"static/css/zero.css",
				"static/css/style.css",
			},
		}
		w.Header().Set("I_am", "here")
		err := writeTemplate(w, templates, data)
		if err != nil {
			log.Error("failed to write template",
				"error", err,
				"templates", templates,
			)
		}
		return
	case http.MethodPost:
		// process a login request
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		valid := isValid(username, password)
		log.Info("login request",
			"username", username,
			"valid", valid,
		)
		if !valid {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		serveRoot(w, r)
		return
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// serveRoot is the base handler for the root (bare) path ("/")
func serveRoot(w http.ResponseWriter, r *http.Request) {
	// order matters (parent -> child)
	templates := []string{
		"static/html/index.html",
		"static/html/head.html",
		"static/html/footer.html",
		// "static/html/login.html",
	}
	data := index{
		Name:  "Nailivic Studios!!",
		Title: "nailivic",
		Stylesheets: []string{
			"static/css/zero.css",
			"static/css/style.css",
		},
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
