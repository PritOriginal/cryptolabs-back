package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PritOriginal/cryptolabs-back/internal/app/api"
	"github.com/PritOriginal/cryptolabs-back/internal/config"
	slogger "github.com/PritOriginal/problem-map-server/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	logger, err := slogger.SetupLogger(cfg.Env)
	if err != nil {
		log.Fatalf("error init logger: %v", err)
	}

	app := api.New(logger, cfg)

	go func() {
		app.MustRun()
	}()

	// Graceful shutdown

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done

	app.Stop()

	logger.Info("server stopped")
}
