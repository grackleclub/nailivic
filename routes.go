package main

import (
	"net/http"
	"time"

	cookie "github.com/ddbgio/cookie/v2"
)

// serveRoot is the base handler for the root (bare) path ("/")
func serveRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// templates
		templates := []string{
			"static/html/login.html",
			"static/html/head.html",
			"static/html/footer.html",
		}
		// data
		data, err := getNewPage("index")
		if err != nil {
			log.Error("failed to get page data",
				"error", err,
			)
			http.Error(w, "failed to get page data", http.StatusInternalServerError)
			return
		}
		// other stufff
		w.Header().Set("I_am", "here")
		log.Info("serving login page",
			"templates", templates,
			"data", data,
		)
		// combine and write
		err = writeTemplate(w, templates, data)
		if err != nil {
			log.Error("failed to write template",
				"error", err,
				"templates", templates,
			)
			http.Error(w, "failed to write template", http.StatusInternalServerError)
		}
		return
	case http.MethodPost:
		// process a login request
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		valid := isValid(username, password)
		// TODO use bcrypt/auth package
		log.Info("login request",
			"username", username,
			"valid", valid,
		)
		if !valid {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		sessionToken, err := token(cookieTokenLen)
		if err != nil {
			log.Error("failed to generate secret",
				"error", err,
			)
			http.Error(w, "auth failure", http.StatusInternalServerError)
			return
		}

		// make a new secret for the user's session
		var sessionSecret = sessionKey{
			Username:      username,
			SessionSecret: sessionToken,
			Exipiry:       time.Now().Add(sessionDefaultExpiry),
		}

		// add that secret to the 'backend'
		// TODO make a real backend
		userID := 1234
		backend[userID] = sessionSecret

		// add that secret to a new cookie
		clientCookie := http.Cookie{
			Name:     cookieName,
			Value:    sessionToken,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   3600, // 1 hour
		}
		// encrypt the cookie
		err = cookie.WriteEncrypted(w, userID, clientCookie, cookieGlobalSecret)
		if err != nil {
			log.Error("failed to write cookie",
				"error", err,
			)
			http.Error(w, "failed to write cookie", http.StatusInternalServerError)
			return
		}
		serveDash(w, r)
		return
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// serveSecret is an authenticated endpoint example
func serveSecret(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("I am a teapot"))
}

// serveDash is unused?
func serveDash(w http.ResponseWriter, r *http.Request) {
	// order matters (parent -> child)
	templates := []string{
		"static/html/index.html",
		"static/html/head.html",
		"static/html/footer.html",
		// "static/html/login.html",
	}
	data, err := getNewPage("dash")
	if err != nil {
		log.Error("failed to get page data",
			"error", err,
		)
		http.Error(w, "failed to get page data", http.StatusInternalServerError)
		return
	}

	err = writeTemplate(w, templates, data)
	if err != nil {
		log.Error("failed to write template",
			"error", err,
			"templates", templates,
			"data", data,
		)
	}
	log.Debug("root served", "templates", templates)
}

// serveParts is a simple example of a page handler using better structs
func serveParts(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"static/html/parts.html",
		"static/html/head.html",
		"static/html/footer.html",
	}
	data, err := getNewPage("parts")
	if err != nil {
		log.Error("failed to get page data",
			"error", err,
		)
		http.Error(w, "failed to get page data", http.StatusInternalServerError)
		return
	}

	log.Debug("serving parts page",
		"templates", templates,
		"data", data,
	)

	err = writeTemplate(w, templates, data)
	if err != nil {
		log.Error("failed to write template",
			"error", err,
			"templates", templates,
			"data", data,
		)
		http.Error(w, "failed to write parts template", http.StatusInternalServerError)
		return
	}
	log.Debug("parts served", "templates", templates)
}

func serveInventory(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"static/html/inventory.html",
		"static/html/head.html",
		"static/html/footer.html",
	}
	data, err := getNewPage("inventory")
	if err != nil {
		log.Error("failed to get page data",
			"error", err,
		)
		http.Error(w, "failed to get page data", http.StatusInternalServerError)
		return
	}

	log.Debug("serving inventory page",
		"templates", templates,
		"data", data,
	)

	err = writeTemplate(w, templates, data)
	if err != nil {
		log.Error("failed to write template",
			"error", err,
			"templates", templates,
			"data", data,
		)
		http.Error(w, "failed to write inventory template", http.StatusInternalServerError)
		return
	}
	log.Debug("inventory served", "templates", templates)
}

// serveHtmx dynamically serves htmx components based on the path
func serveHtmx(w http.ResponseWriter, r *http.Request) {
	// get the component name from the path
	componentName := r.PathValue("component")
	log.Info("htmx component requested", "name", componentName)

	templates := []string{
		"static/html/head.html",
		"static/html/footer.html",
	}
	// serve the appropriate htmx component based on name from path
	var err error
	w.Header().Set("X-htmx-component-name", componentName)
	switch componentName {
	case "parts":
		pretemp := []string{"static/html/parts.html"}
		pretemp = append(pretemp, templates...)

		data, err := getNewPage(componentName)
		if err != nil {
			log.Error("failed to get page data",
				"error", err,
				"data", data,
			)
			http.Error(w, "failed to get page data", http.StatusInternalServerError)
			return
		}

		err = writeTemplate(w, pretemp, data)
		if err != nil {
			log.Error("failed to write htmx component",
				"error", err,
				"component", componentName,
			)
			http.Error(w, "failed to write htmx component", http.StatusInternalServerError)
		}
	default:
		http.Error(w, "missing or invalid htmx component name", http.StatusBadRequest)
	}
	if err != nil {
		log.Error("failed to write htmx component",
			"error", err,
			"component", componentName,
		)
	}
}
