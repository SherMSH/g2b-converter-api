package handlers

import (
	service "converterapi/internal/service/G2B"
	"converterapi/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PinChangeReq struct {
	PAN        string `json:"pan"`
	ExpiryDate string `json:"expiryDate"`
	PIN        string `json:"pin"`
}

func SetPIN(c *gin.Context) {

	var req PinChangeReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Errorf("Error binding PinChageReq: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Error binding PinChageReq"})
		return
	}

	if err := service.SetPinG2b(req.PAN, req.PIN, req.ExpiryDate); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

}
