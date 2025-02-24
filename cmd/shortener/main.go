package main

import (
	"log"
	"net/http"

	"github.com/msorokin-hash/urlshortener/internal/app"
)

func main() {
	port := ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.AddURLHandler)
	mux.HandleFunc("/{hash}", app.GetURLHandler)

	log.Fatal(http.ListenAndServe(port, mux))
}
