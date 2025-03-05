package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/msorokin-hash/urlshortener/internal/app/config"
	"github.com/msorokin-hash/urlshortener/internal/app/storage"
	"github.com/msorokin-hash/urlshortener/internal/app/util"
)

type App struct {
	Config  *config.Config
	Storage *storage.Storage
}

func (app *App) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	res, err := app.Storage.GetURLByHash(parts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (app *App) AddURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashed := util.HashStringData(string(body))
	_ = app.Storage.CreateURL(hashed, string(body))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(app.Config.BaseShortURL + "/" + hashed))
}
