package handler

import (
	"converterapi/internal/middlewares"
	"converterapi/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Init(h *Handler) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORS())
	router.Use(middlewares.Prometheus())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.StaticFile("/favicon.ico", "./internal/app/files/favicon.ico")

	soap := router.Group("/g2b")
	soap.Use(middlewares.SOAPLogger())
	soap.POST("/d8convert", service.D8Converter)

	// json := router.Group("/json")
	// json.POST("/convert2xml")
	return router
}
