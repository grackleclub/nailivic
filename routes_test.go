package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	name       string
	method     string
	route      string
	handler    http.HandlerFunc
	statusCode int // expected status code
	// resBodyMustHave []string // expected response body
}

// testing new routes and handlers requires only setting more test cases
var cases = []test{
	{"root", "GET", "/", serveRoot, http.StatusOK},
	// {"root", "GET", "/dash", serveDash, http.StatusOK},
	// add more test cases here as application grows
}

func TestServeRoutes(t *testing.T) {
	strictTemplateChecking = true
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("testing %s: %s '%s'", test.name, test.method, test.route)
			req, err := http.NewRequest(test.method, test.route, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(test.handler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.statusCode, rr.Code)

			body := rr.Body.String()
			assert.NotEqual(t, 0, len(body))
			t.Logf("passed %s: %d\n%v", test.name, test.statusCode, body)
		})
	}
}
