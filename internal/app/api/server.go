package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/PritOriginal/cryptolabs-back/internal/config"
	"github.com/PritOriginal/cryptolabs-back/internal/handler"
	slogger "github.com/PritOriginal/problem-map-server/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type App struct {
	server *http.Server
	log    *slog.Logger
	router *chi.Mux
	port   int
}

func New(log *slog.Logger, cfg *config.Config) *App {
	router := handler.GetRouter(log)

	server := &http.Server{
		Addr:         cfg.REST.Host + ":" + strconv.Itoa(cfg.REST.Port),
		Handler:      router,
		ReadTimeout:  cfg.REST.Timeout.Read,
		WriteTimeout: cfg.REST.Timeout.Write,
		IdleTimeout:  cfg.REST.Timeout.Idle,
	}

	return &App{
		server: server,
		log:    log,
		router: router,
		port:   cfg.REST.Port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "rest.Run"

	a.log.Info("server started", slog.String("address", ":"+strconv.Itoa(a.port)))
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.log.Error("failed to start server")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "rest.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping REST server", slog.Int("port", a.port))

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.log.Error("an error occurred while stopping the server", slogger.Err(err))
	}
}
