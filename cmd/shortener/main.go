package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/msorokin-hash/urlshortener/internal/app/config"
	"github.com/msorokin-hash/urlshortener/internal/app/handler"
	"github.com/msorokin-hash/urlshortener/internal/app/storage"
)

func main() {
	cfg := config.NewConfig()
	store := storage.NewStorage()

	app := &handler.App{
		Config:  cfg,
		Storage: store,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", app.AddURLHandler)
	r.Get("/{hash}", app.GetURLHandler)

	log.Fatal(http.ListenAndServe(app.Config.Address, r))
}
