package handlers

import (
	"converterapi/internal/models"
	acctcredit "converterapi/internal/models/ONLINE/AcctCredit"
	acctdebit "converterapi/internal/models/ONLINE/AcctDebit"
	addcmsabonent "converterapi/internal/models/ONLINE/AddCmsAbonent"
	addpersoncmsabonent "converterapi/internal/models/ONLINE/AddPersonCMSAbonent"
	changecmsabonent "converterapi/internal/models/ONLINE/ChangeCMSAbonent"
	getaccinfo "converterapi/internal/models/ONLINE/GetAccInfo"
	getacctstatement "converterapi/internal/models/ONLINE/GetAcctStatement"
	getcvv "converterapi/internal/models/ONLINE/GetCVV"
	getcardinfo "converterapi/internal/models/ONLINE/GetCardInfo"
	getcardstatement "converterapi/internal/models/ONLINE/GetCardStatement"
	getpersoninfo "converterapi/internal/models/ONLINE/GetPersonInfo"
	gettransinfo "converterapi/internal/models/ONLINE/GetTransInfo"
	initsession "converterapi/internal/models/ONLINE/InitSession"
	posrequestrq "converterapi/internal/models/ONLINE/POSRequest"
	relinkpreissuedcards "converterapi/internal/models/ONLINE/RelinkPreIssuedCards"
	removecmsabonent "converterapi/internal/models/ONLINE/RemoveCMSAbonent"
	removepersoncmsabonent "converterapi/internal/models/ONLINE/RemovePersonCMSAbonent"
	resetbadpintries "converterapi/internal/models/ONLINE/ResetBadPINTries"
	setcardstatus "converterapi/internal/models/ONLINE/SetCardStatus"
	updateperson "converterapi/internal/models/ONLINE/UpdatePerson"
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
	var resp any
	rqType := utils.GetRqType(bodyStr)

	switch rqType {
	case utils.Unknown:
		logger.Errorf("Unknown XML body!")
		sendSoapFault(c, 400, "Client", "Unknown XML body")
		return
	case utils.AcctCreditRq:
		var unmBody acctcredit.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("acctcredit.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.AcctDebitRq:
		var unmBody acctdebit.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("acctdebit.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.AddCMSAbonentRq:
		var unmBody addcmsabonent.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("addcmsabonent.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		err = unmBody.Call()
		if err != nil {
			logger.Errorf("addcmsabonent g2b svc call err: %s", err.Error())
			sendSoapFault(c, 500, err.Error(), "Service error")
			return
		}
		resp = unmBody
	case utils.AddPersonCMSAbonentRq:
		var unmBody addpersoncmsabonent.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("addpersoncmsabonentrq.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.ChangeCMSAbonentRq:
		var unmBody changecmsabonent.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("changecmsabonentrq.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.GetAcctInfoRq:
		var unmBody getaccinfo.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("getaccinfo.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.GetAcctStatementRq:
		var unmBody getacctstatement.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("getacctstatement.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.GetCardInfoRq:
		var unmBody getcardinfo.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("getcardinforq.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.GetCardStatementRq:
		var unmBody getcardstatement.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("getcardstatement.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody

	case utils.GetCVVRq:
		var unmBody getcvv.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("getcvv.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.GetPersonInfoRq:
		var unmBody getpersoninfo.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("getpersoninfo.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.GetTransInfoRq:
		var unmBody gettransinfo.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("gettransinfo.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.InitSessionRq:
		var unmBody initsession.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("initsession.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.POSRequestRq:
		var unmBody posrequestrq.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("posrequestrq.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		err := unmBody.Call()
		if err != nil {
			sendSoapFault(c, 500, "Client", err.Error())
			return
		}
		switch unmBody.SoapRq.Req.TranCode {
		case posrequestrq.Credit:
			resp = posrequestrq.Envelope{
				XmlnsS:  "http://www.w3.org/2003/05/soap-envelope",
				XmlnsM1: "http://schemas.compassplus.com/two/1.0/fimi.xsd",
				XmlnsM0: "http://schemas.compassplus.com/two/1.0/fimi_types.xsd",
				Body: posrequestrq.RespBody{
					POSRequestRp: posrequestrq.POSRequestRp{
						Response: posrequestrq.Response{
							Product:      unmBody.SoapRq.Req.Product,
							ResponseAttr: "1",
							TranId:       unmBody.SoapRq.Req.ThisTranId,
							Ver:          "16.37",

							AccountCurrency:      unmBody.SoapRq.Req.Currency,
							ApprovalCode:         unmBody.SoapRq.ApprovalCode,
							AuthRespCode:         unmBody.SoapRq.Req.RespCode,
							AuthRespCodeCategory: "0",
							AvailBalance:         "",
							BalanceCurrency:      "",
							BonusDebt:            "",
							CVxOK:                "-1",
							Currency:             unmBody.SoapRq.Req.Currency,
							Fee:                  "",
							FromAcct:             unmBody.SoapRq.Req.FromAccount,
							IssuerFee:            "",
							LedgerBalance:        "",
							MaskBalances:         "",
							RelatedTran:          posrequestrq.RelatedTran{},
							ThisTranId:           unmBody.SoapRq.Req.ThisTranId,
							ToAcct:               unmBody.SoapRq.Req.ToAccount,
						},
					},
				},
			}
		case posrequestrq.Debit:
			resp = posrequestrq.Envelope{
				XmlnsS:  "http://www.w3.org/2003/05/soap-envelope",
				XmlnsM1: "http://schemas.compassplus.com/two/1.0/fimi.xsd",
				XmlnsM0: "http://schemas.compassplus.com/two/1.0/fimi_types.xsd",
				Body: posrequestrq.RespBody{
					POSRequestRp: posrequestrq.POSRequestRp{
						Response: posrequestrq.Response{
							Product:      unmBody.SoapRq.Req.Product,
							ResponseAttr: "1",
							TranId:       unmBody.SoapRq.Req.ThisTranId,
							Ver:          "16.37",

							AccountCurrency:      unmBody.SoapRq.Req.Currency,
							ApprovalCode:         unmBody.SoapRq.ApprovalCode,
							AuthRespCode:         unmBody.SoapRq.Req.RespCode,
							AuthRespCodeCategory: "0",
							AvailBalance:         "",
							BalanceCurrency:      "",
							BonusDebt:            "",
							CVxOK:                "-1",
							Currency:             unmBody.SoapRq.Req.Currency,
							Fee:                  "",
							FromAcct:             unmBody.SoapRq.Req.FromAccount,
							IssuerFee:            "",
							LedgerBalance:        "",
							MaskBalances:         "",
							RelatedTran:          posrequestrq.RelatedTran{},
							ThisTranId:           unmBody.SoapRq.Req.ThisTranId,
							ToAcct:               unmBody.SoapRq.Req.ToAccount,
						},
					},
				},
			}
		default:
			resp = posrequestrq.Envelope{
				XmlnsS:  "http://www.w3.org/2003/05/soap-envelope",
				XmlnsM1: "http://schemas.compassplus.com/two/1.0/fimi.xsd",
				XmlnsM0: "http://schemas.compassplus.com/two/1.0/fimi_types.xsd",
				Body: posrequestrq.RespBody{
					POSRequestRp: posrequestrq.POSRequestRp{
						Response: posrequestrq.Response{
							Product:      unmBody.SoapRq.Req.Product,
							ResponseAttr: "1",
							TranId:       unmBody.SoapRq.Req.ThisTranId,
							Ver:          "16.37",

							AccountCurrency:      unmBody.SoapRq.Req.Currency,
							ApprovalCode:         unmBody.SoapRq.ApprovalCode,
							AuthRespCode:         unmBody.SoapRq.Req.RespCode,
							AuthRespCodeCategory: "0",
							AvailBalance:         "",
							BalanceCurrency:      "",
							BonusDebt:            "",
							CVxOK:                "-1",
							Currency:             unmBody.SoapRq.Req.Currency,
							Fee:                  "",
							FromAcct:             unmBody.SoapRq.Req.FromAccount,
							IssuerFee:            "",
							LedgerBalance:        "",
							MaskBalances:         "",
							RelatedTran:          posrequestrq.RelatedTran{},
							ThisTranId:           unmBody.SoapRq.Req.ThisTranId,
							ToAcct:               unmBody.SoapRq.Req.ToAccount,
						},
					},
				},
			}
		}
	case utils.UpdateCard2AcctLinkRq:
		var unmBody relinkpreissuedcards.SoapEnvelope
		err = xml.Unmarshal(body, &unmBody)
		if err != nil {
			logger.Errorf("relinkpreissuedcards.SoapEnvelope unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody.Body
	case utils.RemoveCMSAbonentRq:
		var unmBody removecmsabonent.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("removecmsabonent.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.RemovePersonCMSAbonentRq:
		var unmBody removepersoncmsabonent.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("removepersoncmsabonent.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.ResetBadPINTriesRq:
		var unmBody resetbadpintries.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("resetbadpintries.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.SetCardStatusRq:
		var unmBody setcardstatus.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("setcardstatus.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	case utils.UpdatePersonRq:
		var unmBody updateperson.Body
		err = xml.Unmarshal(envelope.Body.XMLData, &unmBody.SoapRq)
		if err != nil {
			logger.Errorf("updateperson.Body unmarshal err: %v", err)
			sendSoapFault(c, 500, "Client", "Internal server error")
			return
		}
		resp = unmBody
	default:
		logger.Errorf("Unknown XML body")
		sendSoapFault(c, 400, "Client", "Unknown XML body")
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.XML(http.StatusOK, resp)
}

// Вспомогательная функция для отправки ответа в формате json
func sendJsonResponse(c *gin.Context, req interface{}) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, req)
}

// Вспомогательная функция для отправки SOAP-ошибки (Fault)
func sendSoapFault(c *gin.Context, status int, faultCode, faultString string) {
	fault := struct {
		XMLName xml.Name `xml:"Fault"`
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
	output, _ := xml.MarshalIndent(fault, "", "  ")

	c.Header("Content-Type", "application/soap+xml; charset=utf-8")
	if status == 0 {
		status = 500
	}
	c.Status(status)
	c.String(status, xml.Header+string(output))
}
