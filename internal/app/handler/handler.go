package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/msorokin-hash/urlshortener/internal/app/config"
	"github.com/msorokin-hash/urlshortener/internal/app/util"
)

type Storage interface {
	Lookup(hash string) (string, error)
	Add(hash string, url string) error
}

type Handler struct {
	Config  *config.Config
	Storage Storage
}

func NewHandler(config *config.Config, storage Storage) *Handler {
	return &Handler{Config: config, Storage: storage}
}

func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	res, err := h.Storage.Lookup(parts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) AddURLHandler(w http.ResponseWriter, r *http.Request) {
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
	_ = h.Storage.Add(hashed, string(body))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.Config.BaseShortURL + "/" + hashed))
}
