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
	cfg := config.InitConfig()

	dbConn, err := storage.InitDB()
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer dbConn.Close()

	app := &handler.App{
		Config: cfg,
		DB:     dbConn,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", app.AddURLHandler)
	r.Get("/{hash}", app.GetURLHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
