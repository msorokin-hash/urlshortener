package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/msorokin-hash/urlshortener/internal/app/config"
	"github.com/msorokin-hash/urlshortener/internal/app/handler"
	"github.com/msorokin-hash/urlshortener/internal/app/logger"
)

type Server struct {
	Config  *config.Config
	Handler handler.Handler
}

func NewServer(c *config.Config, h handler.Handler) *Server {
	return &Server{Config: c, Handler: h}
}

func (s *Server) SetupRouter() *chi.Mux {

	_ = logger.Initialize(s.Config.LogLevel)

	r := chi.NewRouter()
	r.Use(logger.LoggingMiddleware)
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)

	r.Post("/", s.Handler.AddURLHandler)
	r.Get("/{hash}", s.Handler.GetURLHandler)

	return r
}

func (s *Server) StartServer() {
	r := s.SetupRouter()
	log.Fatal(http.ListenAndServe(s.Config.Address, r))
}
