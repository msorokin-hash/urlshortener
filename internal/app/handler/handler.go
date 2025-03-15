package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/msorokin-hash/urlshortener/internal/app/config"
	"github.com/msorokin-hash/urlshortener/internal/app/entity"
	"github.com/msorokin-hash/urlshortener/internal/app/service"
	"github.com/msorokin-hash/urlshortener/internal/app/util"
)

type Handler struct {
	Config  *config.Config
	Storage service.Storage
}

func NewHandler(config *config.Config, storage service.Storage) *Handler {
	return &Handler{Config: config, Storage: storage}
}

func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	res, err := h.Storage.Get(parts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) AddURLShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusBadRequest)
		return
	}

	var req entity.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashed := util.HashStringData(req.URL)
	_ = h.Storage.Insert(hashed, req.URL)
	resp := entity.Response{Result: h.Config.BaseShortURL + "/" + hashed}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
	_ = h.Storage.Insert(hashed, string(body))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.Config.BaseShortURL + "/" + hashed))
}

func (h *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.Storage.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
