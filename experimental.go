package main

// TODO this package doesn't work yet, but outlines the needed components

// import (
// 	"fmt"
// 	"net/http"
// )

// // make a func to setup routes
// type routes map[string]routeConfig

// func New([]routeConfig) routes {
// 	routesMap := make(routes)
// 	for _, route := range routesMap {
// 		routesMap[route.Path] = route
// 	}
// 	return routesMap
// }

// // usage
// // registerRoutes(mux, mainConfigs)
// // routesMap := New(mainConfigs)

// var mainConfigs = []routeConfig{
// 	routeConfig{
// 		Path:    "/wild",
// 		Handler: logMW(serveWild),
// 		Templates: []string{
// 			"static/html/wild.html",
// 			"static/html/head.html",
// 			"static/html/footer.html",
// 		},
// 		ContentFunc: wildData,
// 	},
// }

// type appRouteConfigs []routeConfig

// type routeConfig struct {
// 	Path        string
// 	Handler     http.HandlerFunc
// 	Templates   []string
// 	ContentFunc func() (interface{}, error)
// }

// func wildData() (interface{}, error) {
// 	return getNewPage("wild")
// }

// func (c *routeConfig) do(w http.ResponseWriter) error {
// 	content, err := c.ContentFunc()
// 	if err != nil {
// 		return fmt.Errorf("unable to fetch content: %w", err)
// 	}
// 	err = writeTemplate(w, c.Templates, content)
// 	if err != nil {
// 		return fmt.Errorf("unable to write template: %w", err)
// 	}
// 	return nil
// }

// // registerRoutes registers routes based on the provided configuration
// func registerRoutes(mux *http.ServeMux, routes []routeConfig) {
// 	for _, route := range routes {
// 		mux.HandleFunc(route.Path, route.Handler)
// 	}
// }

// // serveCrazy is a simple example of a page handler using better structs
// func serveWild(w http.ResponseWriter, r *http.Request) {
// 	mapRoutes := New(mainConfigs)

// 	log.Debug("serving wild")
// 	path := r.URL.Path
// 	mapRoutes[path].do(w)
// }
