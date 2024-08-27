package main

import (
	"fmt"
	"net/http"
)

type index struct {
	Name       string
	Title      string
	Stylesheet string
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

func serveLogin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusTeapot)
}

func serveInventory(w http.ResponseWriter, r *http.Request) {
	product := r.URL.Query().Get("product")
	a := r.URL.Query().Get("a")
	b := r.URL.Query().Get("b")
	c := r.URL.Query().Get("c")

	log.Info("inventory requested",
		"method", r.Method,
		"product", product,
		"a", a,
		"b", b,
		"c", c,
	)

	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("get inventory"))
	case http.MethodPost:
		w.Write([]byte("post inventory"))
	default:
		http.Error(w, "method not implemented", http.StatusBadRequest)
	}

}

// EXAMPLES

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
