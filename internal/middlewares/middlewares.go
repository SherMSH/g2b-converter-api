package middlewares

import (
	"bytes"
	"converterapi/internal/config"
	d8procweb "converterapi/pkg/d8-proc-web"
	"converterapi/pkg/logger"
	"converterapi/pkg/prometheus"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		var key = config.Config.App.ApiKey

		reqToken := c.GetHeader("Authorization")
		if reqToken == "" {
			reqToken = c.GetHeader("X-API-Key")
		}
		if reqToken == "" {
			logger.Warnf("empty API key")
			c.XML(http.StatusUnauthorized, errors.New("secret API key is needed"))
			c.Abort()
			return
		}
		reqToken = strings.TrimPrefix(reqToken, "Bearer ")
		splitToken := strings.TrimPrefix(key, "Bearer ")

		if len(key) > 5 {
			logger.Infof("--------interchange API key-------- %v***", key[:5])
		}

		if reqToken != splitToken {
			logger.Warnf("wrong API key!")
			c.XML(http.StatusUnauthorized, "secret key is not valid")
			c.Abort()
			return
		}

	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime).Seconds()
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		status := strconv.Itoa(c.Writer.Status())
		if status >= "200" && status < "500" {
			prometheus.ServiceStatus.Set(1)
		} else {
			prometheus.ServiceStatus.Set(0)
		}

		prometheus.HttpRequestsTotal.WithLabelValues(c.Request.Method, endpoint, status).Inc()
		prometheus.RequestDuration.WithLabelValues(c.Request.Method, endpoint).Observe(duration)
		if status >= "400" {
			prometheus.ErrorsTotal.WithLabelValues(c.Request.Method, endpoint, status).Inc()
		}
		prometheus.RequestProcessingTime.WithLabelValues(c.Request.Method, endpoint, status).Observe(duration)
	}
}

func SOAPLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Читаем тело запроса для логирования
		body, _ := io.ReadAll(c.Request.Body)
		// Восстанавливаем тело для дальнейшего использования
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		logger.Infof("Incoming request: %s", string(body))
		c.Next()
	}
}

func D8ProcWebAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := d8procweb.Signin()
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		if resp.StatusCode != 200 {
			c.AbortWithError(resp.StatusCode, fmt.Errorf("status %v", resp.Status))
			return
		}
		c.Next()
		resp, err = d8procweb.Signout()
	}
}
