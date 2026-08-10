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

// VerifyPIN проверяет PIN по карте.
//
// Неверный PIN - это не ошибка сервиса, поэтому отвечаем 200 с признаком
// verified: клиенту нужно отличать «PIN не подошёл» от «проверить не удалось».
// Учтите, что неудачные попытки увеличивают счётчик неверных вводов в
// процессинге и в итоге блокируют карту.
func VerifyPIN(c *gin.Context) {
	var req PinChangeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Error binding PinChageReq: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Error binding PinChageReq"})
		return
	}

	if err := service.VerifyPinG2b(req.PAN, req.PIN, req.ExpiryDate); err != nil {
		c.JSON(http.StatusOK, gin.H{"verified": false, "reason": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"verified": true})
}
