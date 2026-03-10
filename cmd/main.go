package main

import (
	"context"
	"converterapi/internal/app"
	"converterapi/internal/config"
	"converterapi/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	config.Setup("internal/config/config.json")
	logger.Init()
}

// @title CONVERTER API-MAIN
// @version 1.0
// @description CONVERTER API for partner xml <-> json
// @host 192.168.145.74
func main() {
	logger.Info("[MAIN] Work has started!")
	defer beforeQuit()
	app := app.New()

	go func() {
		if err := app.Run(); err != nil && err != http.ErrServerClosed {
			logger.Warn("Exite application")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	q := <-quit
	logger.Info("[SERVER] Shutdown signal received %v", q)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := app.Shutdown(ctx)
	if err != nil {
		logger.Error("Server Shutdown err: %v", err)
	}
	<-ctx.Done()
	defer cancel()
}

// TODO:
// docker

func beforeQuit() {
	logger.Info("[MAIN] Work has stopped!")
}
