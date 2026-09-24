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

// SetPIN назначает карте новый PIN.
//
// Значение не принимается: PIN выбирает процессинг и наружу не отдаёт, до
// держателя карты он доходит SMS-оповещением. Присланный pin отклоняем, а не
// игнорируем - иначе клиент будет считать, что установлен его PIN.
//
// TODO: ручной ввод PIN вернём. Тогда при непустом pin вместо отказа снова
// вызывается xmiss/setPIN - реализация была в service.SetPinG2b, её можно
// достать из истории (удалена в 3083cae). Криптография для этого на месте:
// buildPinRequest живёт в pinBlock.go и используется verifyPIN.
func SetPIN(c *gin.Context) {
	var req PinChangeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Error binding PinChageReq: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Error binding PinChageReq"})
		return
	}

	if req.PIN != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pin must be empty: PIN задаётся процессингом, ручной ввод отключён"})
		return
	}
	if req.PAN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pan is required"})
		return
	}
	if req.ExpiryDate == "" {
		var err error
		req.ExpiryDate, err = service.GetExpDateByPan(req.PAN)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "expiryDate is required: " + err.Error()})
			return
		}
	}

	if err := service.GeneratePIN(req.PAN, req.ExpiryDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"generated": true})
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
