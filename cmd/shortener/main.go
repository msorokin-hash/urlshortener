package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/msorokin-hash/urlshortener/internal/app"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", app.AddURLHandler)
	r.Get("/{hash}", app.GetURLHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
