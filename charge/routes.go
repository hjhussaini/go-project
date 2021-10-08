package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/api/charges", app.giftCharge)
	mux.Get("/api/charges/{code}/users", app.getUsersByCharge)

	return mux
}
