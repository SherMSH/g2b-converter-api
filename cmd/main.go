package main

import (
	"context"
	"converterapi/internal/app"
	"converterapi/internal/config"
	"converterapi/internal/jobs"
	"converterapi/internal/utils"
	"converterapi/pkg/crypto"
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
	utils.Init()
	crypto.Init()
}

// @title CONVERTER API-MAIN
// @version 1.0
// @description CONVERTER API for partner (from Compass) xml <-> json (to D8_G2B)
// @host 192.168.145.74
func main() {
	logger.Infof("[MAIN] Work has started!")
	defer beforeQuit()
	app := app.New()

	jobs.Start()
	go func() {
		if err := app.Run(); err != nil {
			if err == http.ErrServerClosed {
				logger.Infof("Exiting the application")
				return
			}
			logger.Fatalf("Application run err: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	q := <-quit
	logger.Infof("[SERVER] Shutdown signal received %v", q)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		logger.Errorf("Server Shutdown err: %v", err)
	}
	<-ctx.Done()
}

func beforeQuit() {
	logger.Infof("[MAIN] Work has stopped!")
}
