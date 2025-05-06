package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/PritOriginal/cryptolabs-back/internal/handler"
	slogger "github.com/PritOriginal/problem-map-server/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func main() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	logger := slogger.SetupLogger("dev", f)

	r := handler.GetRouter(logger)

	server := New(logger, r)
	server.Start()
}

type Server struct {
	log    *slog.Logger
	router *chi.Mux
}

func New(log *slog.Logger, router *chi.Mux) *Server {
	return &Server{log: log, router: router}
}

func (s *Server) Start() {
	go func() {
		if err := http.ListenAndServe(":3333", s.router); err != nil {
			s.log.Error("failed to start server")
		}
	}()
	s.log.Info("server started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	s.log.Info("server stopped")
}
