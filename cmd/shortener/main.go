package main

import (
	"log"
	"net/http"

	"github.com/msorokin-hash/urlshortener/internal/app"
)

func main() {
	port := ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.AddUrlHandler)
	mux.HandleFunc("/{hash}", app.GetUrlHandler)

	log.Fatal(http.ListenAndServe(port, mux))
}
