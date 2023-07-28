package main

import (
	"testing"

	"github.com/AdiAkhileshSingh15/bookmyroom/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing
	default:
		t.Errorf("type is not *chi.Mux, but it %T", v)
	}

}
