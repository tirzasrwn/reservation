package main

import (
	"net/http"
	"testing"

	"github.com/tirzasrwn/reservation/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)
	switch v := mux.(type) {
	case http.Handler:
		// Do nothing, test passed!
	default:
		t.Errorf("Type is not http.Handler, but type is %T.", v)
	}
}
