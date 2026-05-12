package router

import (
	"converterapi/internal/handlers"
	"converterapi/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Init(h *handlers.Handler) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORS())
	router.Use(middlewares.Prometheus())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.StaticFile("/favicon.ico", "./internal/app/files/favicon.ico")
	router.MaxMultipartMemory = 10 << 20 // 10 MiB

	soap := router.Group("/g2b")
	soap.Use(middlewares.SOAPLogger())
	soap.Use(middlewares.CheckApiKey())
	{
		soap.GET("/ping", ping)
		soap.POST("/d8convert", handlers.D8Converter)
		soap.PUT("/convFile", handlers.PutConvFile)
		soap.GET("/convFile/:filename", handlers.GetConvFile)

		soap.POST("/PinChange", handlers.PinChange)
	}
	return router
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Pong!"})
}
