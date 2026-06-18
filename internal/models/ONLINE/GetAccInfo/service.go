package getaccinfo

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"fmt"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	if len(sb.SoapRq.Req.Account) < 16 {
		return nil, fmt.Errorf("wrong mandatory field `fimi1:Account`")
	}
	foundAcc, err := service.GetAcctInfoG2b(sb.SoapRq.Req.Account)
	if err != nil {
		return nil, err
	}
	cards, err := service.GetCardsListG2b(foundAcc.Custcode, foundAcc.Currcode)
	if err != nil {
		return nil, err
	}
	var cardRows []CardRow
	for _, v := range cards {
		var cardPan string
		switch v.PAN {
		case "":
			cardPan = v.LkeyDisplay
		default:
			cardPan = v.PAN
		}
		cardRows = append(cardRows, CardRow{
			PAN:    cardPan,
			MBR:    "0",
			Status: utils.CardStatuses[v.StatCode],
			Type:   utils.CardTypes[v.ProductType],
		},
		)
	}
	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	soapResp.Body.GetAcctInfoRp.Response = Response{
		Echo:         sb.SoapRq.Req.Echo,
		Product:      sb.SoapRq.Req.Product,
		ResponseAttr: "1",
		TranId:       utils.GenerateTimestampID(),
		Ver:          "1.0",

		Avail:                 fmt.Sprintf("%.2f", foundAcc.AvlBal),
		Bonus:                 fmt.Sprintf("%.2f", foundAcc.Balincr),
		Cards:                 Rows{Rows: cardRows},
		CreditHold:            "0",
		Currency:              utils.Currencies[foundAcc.Currcode],
		DebitHold:             "0",
		DropTmpOverOnRefresh:  "0",
		ExtendedAccountNumber: foundAcc.Accnum,
		FoundAccount:          foundAcc.Accnum,
		LastDepAmount:         fmt.Sprintf("%.2f", foundAcc.AvlbalUnset),
		LastDepTime:           utils.ConvertDate(foundAcc.OpenDate),
		LastRefreshTime:       utils.ConvertDate(foundAcc.LastUsage),
		LastTranId:            "",
		LastWdlAmount:         fmt.Sprintf("%.2f", foundAcc.BlkamtUnset),
		LastWdlTime:           utils.ConvertDate(foundAcc.LastUsage),
		Ledger:                fmt.Sprintf("%.2f", foundAcc.AvlBal+foundAcc.BlkAmt),
		MaskBalances:          "0",
		PermissibleExcessType: "-1",
		PersonExtId:           "",
		PersonFIO:             foundAcc.Name,
		PersonId:              foundAcc.Custcode,
		Remain:                fmt.Sprintf("%.2f", foundAcc.AvlBal),
		Status:                utils.AccountStatuses[foundAcc.Statcode],
		TmpOverdraft:          fmt.Sprintf("%.2f", foundAcc.Balincr),
		Type:                  utils.AccountTypes[foundAcc.Typecode],
	}

	return soapResp, nil
}
