package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/PritOriginal/cryptolabs-back/internal/handler"
	slogger "github.com/PritOriginal/problem-map-server/pkg/logger"
)

func main() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	logger := slogger.SetupLogger("dev", f)

	r := handler.GetRouter(logger)
	go func() {
		if err := http.ListenAndServe(":3333", r); err != nil {
			logger.Error("failed to start server")
		}
	}()
	logger.Info("server started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	logger.Info("server stopped")
}
