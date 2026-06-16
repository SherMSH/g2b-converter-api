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
	err = service.GetAcctInfoG2b()
	if err != nil {
		return nil, err
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

		Avail: "0",
		Bonus: "0",
		Cards: Rows{
			[]CardRow{
				{
					PAN:    "5058270530000016",
					MBR:    "0",
					Status: "11",
					Type:   "0",
				},
			},
		},
		CreditHold:            "0",
		Currency:              "972",
		DebitHold:             "0",
		DropTmpOverOnRefresh:  "0",
		ExtendedAccountNumber: sb.SoapRq.Req.Account,
		FoundAccount:          sb.SoapRq.Req.Account,
		LastDepAmount:         "0",
		LastDepTime:           "2026-04-21T15:45:34",
		LastRefreshTime:       "2026-04-21T15:45:34",
		LastTranId:            "",
		LastWdlAmount:         "0",
		LastWdlTime:           "2026-04-21T15:45:34",
		Ledger:                "0",
		MaskBalances:          "0",
		PermissibleExcessType: "-1",
		PersonExtId:           "120147",
		PersonFIO:             "006772212 Давронбек Болтабоев",
		PersonId:              "120147",
		Remain:                "0",
		Status:                "00",
		TmpOverdraft:          "0",
		Type:                  "00",
	}

	return soapResp, nil
}
