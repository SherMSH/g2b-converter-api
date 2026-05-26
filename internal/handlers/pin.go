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

func PinChange(c *gin.Context) {

	var req PinChangeReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
		logger.Errorf("Error binding PinChageReq: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error binding PinChageReq"})
		return
	}

	if err := service.SetPinG2b(req.PAN, req.PIN, req.ExpiryDate); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

}
