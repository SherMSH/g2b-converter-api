package reversetransaction

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"fmt"
	"strconv"
)

// defaultReasonCode - причина отмены по умолчанию.
// Партнёр может прислать свой код в поле ReasonCode (Appendix C спецификации D8).
const defaultReasonCode = 4000

// Svc отменяет ранее проведённую операцию.
//
// Партнёр присылает только Id операции, поэтому порядок такой:
//  1. getTransactionDetails по Id - берём ссылку, сумму и валюту оригинала;
//  2. initiateTransaction - реверсалу нужна собственная ссылка (8.6);
//  3. reverseTransaction.
//
// Сумма отмены берётся из оригинальной операции; если партнёр прислал Amount,
// используется он - это частичная отмена.
func Svc(sb *Body) (soapResp *Envelope, err error) {
	req := sb.SoapRq.Req

	if len(req.Id) == 0 {
		return nil, fmt.Errorf("wrong mandatory field `fimi1:Id`")
	}

	original, err := service.GetTransactionDetailsG2b(req.Id, "")
	if err != nil {
		logger.Errorf("[SERVICE] reverseTransaction: original tx %s not found: %v", req.Id, err)
		return nil, err
	}
	if original.Details.EcTxRefno == "" {
		return nil, fmt.Errorf("original transaction %s has no ecTxRefno", req.Id)
	}

	amount := original.Details.TxnAmount
	if len(req.Amount) != 0 {
		partial, errAmt := strconv.ParseFloat(req.Amount, 64)
		if errAmt != nil || partial <= 0 {
			return nil, fmt.Errorf("wrong field `fimi1:Amount`")
		}
		if partial > original.Details.TxnAmount {
			return nil, fmt.Errorf("reversal amount exceeds original transaction amount")
		}
		amount = partial
	}

	reasonCode := defaultReasonCode
	if len(req.ReasonCode) != 0 {
		if code, errCode := strconv.Atoi(req.ReasonCode); errCode == nil {
			reasonCode = code
		} else {
			logger.Warnf("[SERVICE] reverseTransaction: wrong ReasonCode %q, using default %d", req.ReasonCode, defaultReasonCode)
		}
	}

	ecTxRefNo, err := service.InitiateTransaction()
	if err != nil {
		logger.Errorf("[SERVICE] reverseTransaction: initiate err: %v", err)
		return nil, err
	}

	resp, err := service.ReverseTransaction(*ecTxRefNo, original.Details.EcTxRefno, amount, original.Details.TxnCurrency, reasonCode)
	if err != nil {
		logger.Errorf("[SERVICE] reverseTransaction err: %v", err)
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body = RespBody{
		ReverseTransactionRp: ReverseTransactionRp{
			Response: Response{
				Echo:         req.Echo,
				Product:      req.Product,
				ResponseAttr: "1",
				TranId:       utils.GenerateTimestampID(),
				Ver:          "1.0",

				ThisTranId:   req.Id,
				AuthRespCode: resp.Status.RspCode,
				AuthRespText: resp.Status.Message,
				ExtRespCode:  resp.Status.Code,
			},
		},
	}
	return soapResp, nil
}
