package server

import (
	"Sprint-13-14/pkg/api"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	logger *log.Logger
	server *http.Server
}

func NewServer(logger *log.Logger) *Server {
	webDir := "./web"

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	addr := ":" + port

	router := http.NewServeMux()
	router.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init(router)

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		logger: logger,
		server: server,
	}
}

func (s *Server) Start() {
	if err := s.server.ListenAndServe(); err != nil {
		s.logger.Fatal("Ошибка запуска сервера:", err)
		return
	}
}
