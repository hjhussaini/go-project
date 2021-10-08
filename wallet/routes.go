package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/api/wallets", app.createWallet)
	mux.Get("/api/wallets/{id}", app.getWalletByID)

	return mux
}
