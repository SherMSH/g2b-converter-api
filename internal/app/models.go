package app

import (
	"converterapi/internal/config"
	"converterapi/internal/handler"
	"converterapi/internal/repository"
	"converterapi/internal/service"
	"converterapi/pkg/prometheus"
	"net"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	repo    *repository.Repository
	service *service.Service
	handler *handler.Handler
	gengine *gin.Engine
}

func New() *App {
	app := &App{}

	app.repo = repository.New()
	app.service = service.New(app.repo)
	app.handler = handler.New(app.service)
	app.gengine = handler.Init(app.handler)
	prometheus.Init()
	return app
}

func (a *App) Run(cfg *config.Configs) error {
	go func() {
		for {
			prometheus.UpdateSystemMetrics()
			time.Sleep(1 * time.Minute)
		}
	}()
	return a.gengine.Run(net.JoinHostPort(cfg.App.Server.Host, cfg.App.Server.Port))
}
