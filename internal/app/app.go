package app

import (
	"context"
	"converterapi/internal/config"
	"converterapi/internal/handlers"
	"converterapi/internal/repository"
	"converterapi/internal/router"
	"converterapi/internal/service"
	"converterapi/pkg/prometheus"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	repo    *repository.Repository
	service *service.Service
	handler *handlers.Handler
	gengine *gin.Engine
	server  *http.Server
}

func New() *App {
	app := &App{}

	app.repo = repository.New()
	app.service = service.New(&config.Config, app.repo)
	app.handler = handlers.New(app.service)
	app.gengine = router.Init(app.handler)
	app.server = &http.Server{
		Addr:    net.JoinHostPort(config.Config.App.Server.Host, config.Config.App.Server.Port),
		Handler: app.gengine,
	}
	prometheus.Init()
	return app
}

func (a *App) Run() error {
	go func() {
		for {
			prometheus.UpdateSystemMetrics()
			time.Sleep(1 * time.Minute)
		}
	}()
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(c context.Context) error {
	return a.server.Shutdown(c)
}
