package main

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
)

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

func isValid(username, password string) bool {
	if username == "user" && password == "pass" {
		return true
	}
	return false
}
