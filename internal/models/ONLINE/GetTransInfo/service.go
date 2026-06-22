package gettransinfo

import (
	service "converterapi/internal/service/G2B"
	"converterapi/pkg/logger"
	"fmt"
	"time"
)

func Svc(b *Body) (soapResp *Envelope, err error) {
	trn, err := service.GetTransactionDetailsG2b(b.SoapRq.Req.Id, b.SoapRq.Req.TranNumber)
	if err != nil {
		return nil, err
	}
	if trn.Details.Lkey.Pan != b.SoapRq.Req.PAN {
		return nil, fmt.Errorf("Transaction card PAN validation fail!")
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	tstamp, err := time.Parse("20060102T150405", trn.Details.DateLocal+"T"+trn.Details.TimeLocal)
	if err != nil {
		logger.Errorf("[SERVICE] gettransinfo time parsing err: %v", err)
		err = nil
	}
	var rspCode string = "0"
	if trn.Details.RspCode == "00" {
		rspCode = "1"
	}
	tranListArr := []TranListRow{
		{
			Id:                   fmt.Sprintf("%d", trn.Details.TlId),
			Type:                 fmt.Sprintf("%d", trn.Details.FnCode),
			Time:                 tstamp.Format("2006-01-02T15:04:05"),
			Phase:                "21",
			TermClass:            trn.Details.Termtype,
			TermName:             trn.Details.TermCode,
			TermDate:             tstamp.Format("2006-01-02") + "T00:00:00",
			TranCode:             fmt.Sprintf("%d", trn.Details.TxnCode),
			DraftCapture:         "1",
			FromAcct:             trn.Details.Lkey.MaskedPan,
			Amount:               fmt.Sprintf("%.2f", trn.Details.TxnAmount),
			Amount2:              "0",
			Fee:                  "0",
			IssuerFee:            "0",
			Currency:             trn.Details.TxnCurrency,
			PAN:                  trn.Details.Lkey.Pan,
			CardMember:           "0",
			RespCode:             rspCode,
			RetainCard:           "0",
			ApprovalCode:         trn.Details.Aprvlcode,
			LedgerBalance:        fmt.Sprintf("%.2f", trn.Details.Amtbillbal),
			AvailBalance:         fmt.Sprintf("%.2f", trn.Details.Amtbillbal),
			BalanceCurrencyAcct:  "0",
			CurrencyAcct:         trn.Details.Curbill,
			AmountAcct:           fmt.Sprintf("%.2f", trn.Details.Amtbill),
			ExchangeRateAcct:     fmt.Sprintf("%.2f", trn.Details.Ratebill),
			RevRequestId:         "",
			Error:                "0",
			OrigType:             "9",
			TermFIName:           trn.Details.CrdacptID,
			TermInstID:           trn.Details.TermCode,
			TermRetailerName:     trn.Details.CrdacptlocName,
			TermSIC:              fmt.Sprintf("%d", trn.Details.CrdacptBus),
			TermSICName:          trn.Details.TerminalPhysicalAddress.AdditionalInfo,
			TermCountry:          trn.Details.CrdacptlocCountry,
			TermCountryName:      "TAJIKISTAN",
			TermCity:             trn.Details.CrdacptlocCity,
			TermOwner:            trn.Details.CrdacptID,
			Track2:               trn.Details.Lkey.Pan + "=" + trn.Details.DateExp,
			AuthFIName:           trn.Details.Issrtcode,
			RevActualAmount:      fmt.Sprintf("%.2f", trn.Details.Amtbill),
			TranNumber:           trn.Details.EcTxRefno,
			POSCondition:         "0",
			POSEntryMode:         "000",
			FromAcctType:         "1",
			ToAcctType:           trn.Details.DestinationAccountType,
			CNSent:               "1",
			OverdraftLimit:       "0",
			TmpOverdraft:         "0",
			PrevTran:             fmt.Sprintf("%d", trn.Details.TlId),
			OrigTime:             tstamp.Format("2006-01-02T15:04:05"),
			DebitHold:            "0",
			CreditHold:           "0",
			Bonus:                "0",
			LedgerBalanceBefore:  fmt.Sprintf("%.2f", trn.Details.Amtbillbal),
			AvailBalanceBefore:   fmt.Sprintf("%.2f", trn.Details.Amtbillbal),
			DebitHoldBefore:      "0",
			CreditHoldBefore:     "0",
			BonusBefore:          "0",
			Reason:               fmt.Sprintf("%d", trn.Details.TxStatus),
			ICC_IssuerScript1:    "0",
			ICC_IssuerScript2:    "0",
			ICC_TranType:         fmt.Sprintf("%d", trn.Details.FnCode),
			Clerk:                trn.Details.TermCode,
			PINEntry:             "0",
			Host:                 "100",
			CNSId:                fmt.Sprintf("%d", trn.Details.Lkey.LkeyId),
			LedgerBalance2:       fmt.Sprintf("%.2f", trn.Details.Amtbillbal),
			AvailBalance2:        fmt.Sprintf("%.2f", trn.Details.Amtbillbal),
			LedgerBalance2Before: fmt.Sprintf("%.2f", trn.Details.Amtbillbal),
			AvailBalance2Before:  fmt.Sprintf("%.2f", trn.Details.Amtbillbal),
			ReversalAllowed:      "1",
			ReceiptFlag:          "0",
			IsContainCAVV:        "0",
			TermInstCountry:      trn.Details.TerminalPhysicalAddress.Country,
			TermInstCountryName:  trn.Details.TerminalPhysicalAddress.LocationName,
			TranTime:             tstamp.Format("2006-01-02T15:04:05"),
			RevActualAmountAcct:  fmt.Sprintf("%.2f", trn.Details.Amtbill),
		},
	}

	soapResp.Body.GetTransInfoRp.Response = Response{
		Echo:         b.SoapRq.Req.Echo,
		Product:      b.SoapRq.Req.Product,
		ResponseAttr: "1",
		Ver:          "1.0",

		MaskBalances: "0",
		TranList: TranList{
			Rows: tranListArr,
		},
	}
	return
}
