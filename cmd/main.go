package main

import (
	"converterapi/internal/app"
	"converterapi/internal/config"
	"converterapi/pkg/logger"
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

	if err := app.Run(&config.Config); err != nil {
		logger.Fatal("Application run failed")
	}
}

func beforeQuit() {
	logger.Info("[MAIN] Work has stopped!")
}
