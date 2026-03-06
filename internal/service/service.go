package service

import (
	"converterapi/internal/models"
	acctcredit "converterapi/internal/models/AcctCredit"
	acctdebit "converterapi/internal/models/AcctDebit"
	addcmsabonent "converterapi/internal/models/AddCmsAbonent"
	addpersoncmsabonent "converterapi/internal/models/AddPersonCMSAbonent"
	changecmsabonent "converterapi/internal/models/ChangeCMSAbonent"
	getaccinfo "converterapi/internal/models/GetAccInfo"
	getacctstatement "converterapi/internal/models/GetAcctStatement"
	getcvv "converterapi/internal/models/GetCVV"
	getcardinfo "converterapi/internal/models/GetCardInfo"
	getcardstatement "converterapi/internal/models/GetCardStatement"
	getpersoninfo "converterapi/internal/models/GetPersonInfo"
	gettransinfo "converterapi/internal/models/GetTransInfo"
	initsession "converterapi/internal/models/InitSession"
	relinkpreissuedcards "converterapi/internal/models/RelinkPreIssuedCards"
	removecmsabonent "converterapi/internal/models/RemoveCMSAbonent"
	removepersoncmsabonent "converterapi/internal/models/RemovePersonCMSAbonent"
	resetbadpintries "converterapi/internal/models/ResetBadPINTries"
	setcardstatus "converterapi/internal/models/SetCardStatus"
	updateperson "converterapi/internal/models/UpdatePerson"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func D8Converter(c *gin.Context) {
	// 1. Проверяем Content-Type
	if c.ContentType() != "application/soap+xml" && c.ContentType() != "text/xml" && c.ContentType() != "application/xml" {
		c.String(http.StatusUnsupportedMediaType, "Content-Type must be application/soap+xml or text/xml")
		return
	}

	// 2. Читаем тело запроса
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to read request body: %v", err)
		return
	}
	defer c.Request.Body.Close()

	// 3. Парсим SOAP-запрос
	var envelope models.SoapEnvelope
	err = xml.Unmarshal(body, &envelope)
	if err != nil {
		// В случае ошибки парсинга возвращаем SOAP-ошибку
		sendSoapFault(c, 500, "Client", "Invalid XML format: "+err.Error())
		return
	}

	// 4. Извлекаем данные запроса
	bodyStr := string(envelope.Body.XMLData)

	// 5. Парсинг-логика
	var reqBody any
	rqType := utils.GetRqType(bodyStr)

	switch rqType {
	case utils.Unknown:
		logger.Error("Unknown XML body!")
		sendSoapFault(c, 400, "Client", "Unknown XML body")
		return
	case utils.AcctCreditRq:
		var unmBody acctcredit.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("acctcredit.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.AcctDebitRq:
		var unmBody acctdebit.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("acctdebit.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.AddCMSAbonentRq:
		var unmBody addcmsabonent.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("addcmsabonent.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.AddPersonCMSAbonentRq:
		var unmBody addpersoncmsabonent.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("addpersoncmsabonentrq.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.ChangeCMSAbonentRq:
		var unmBody changecmsabonent.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("changecmsabonentrq.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.GetAcctInfoRq:
		var unmBody getaccinfo.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("getaccinfo.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.GetAcctStatementRq:
		var unmBody getacctstatement.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("getacctstatement.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.GetCardInfoRq:
		var unmBody getcardinfo.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("getcardinforq.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.GetCardStatementRq:
		var unmBody getcardstatement.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("getcardstatement.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody

	case utils.GetCVVRq:
		var unmBody getcvv.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("getcvv.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.GetPersonInfoRq:
		var unmBody getpersoninfo.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("getpersoninfo.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.GetTransInfoRq:
		var unmBody gettransinfo.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("gettransinfo.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.InitSessionRq:
		var unmBody initsession.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("initsession.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.UpdateCard2AcctLinkRq:
		var unmBody relinkpreissuedcards.SoapEnvelope
		err = xml.Unmarshal(body, &unmBody)
		if err != nil {
			logger.Error("relinkpreissuedcards.SoapEnvelope unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody.Body
	case utils.RemoveCMSAbonentRq:
		var unmBody removecmsabonent.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("removecmsabonent.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.RemovePersonCMSAbonentRq:
		var unmBody removepersoncmsabonent.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("removepersoncmsabonent.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.ResetBadPINTriesRq:
		var unmBody resetbadpintries.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("resetbadpintries.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.SetCardStatusRq:
		var unmBody setcardstatus.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("setcardstatus.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	case utils.UpdatePersonRq:
		var unmBody updateperson.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Error("updateperson.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		reqBody = unmBody
	default:
		logger.Error("Unknown XML body")
		sendSoapFault(c, 400, "Client", "Unknown XML body")
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, reqBody)
}

// Вспомогательная функция для отправки ответа в формате json
func sendJsonResponse(c *gin.Context, req interface{}) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, req)
}

// Вспомогательная функция для отправки SOAP-ошибки (Fault)
func sendSoapFault(c *gin.Context, status int, faultCode, faultString string) {
	fault := struct {
		XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Fault"`
		Code    struct {
			Value string `xml:"Value"`
		} `xml:"Code"`
		Reason struct {
			Text string `xml:"Text"`
		} `xml:"Reason"`
	}{
		Code: struct {
			Value string `xml:"Value"`
		}{
			Value: "soap:" + faultCode,
		},
		Reason: struct {
			Text string `xml:"Text"`
		}{
			Text: faultString,
		},
	}
	// Переопределяем Body для ошибки (в реальном коде нужно создать отдельную структуру)
	output, _ := xml.MarshalIndent(fault, "", "  ")

	c.Header("Content-Type", "application/soap+xml; charset=utf-8")
	if status == 0 {
		status = 500
	}
	c.Status(status)
	c.String(status, xml.Header+string(output))
}
