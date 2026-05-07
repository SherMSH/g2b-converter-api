package getcardinfo

import (
	service "converterapi/internal/service/G2B"
	"fmt"
)

func Svc(sb *Body) (resp *Response, err error) {
	cardInfo, err := service.GetCardInfo(sb.SoapRq.Req.PAN, sb.SoapRq.Req.ExpirationDate)
	if err != nil {
		return nil, err
	}
	resp.Product = sb.SoapRq.Req.Product
	resp.Response = "1"
	resp.Ver = sb.SoapRq.Req.Ver

	resp.Accounts = Accounts{
		Row: AccountRow{
			AcctNo:        string(cardInfo.CardAccounts[0].AcctNo),
			Status:        "3",
			LedgerBalance: "0",
			AvailBalance:  "0",
			Currency:      "",
			Type:          "1",
			AccountStatus: "3",
		},
	}
	resp.Acct2CardAttachType = "0"
	resp.CNSDisabled = "1"
	resp.CardAllowedEMVScript = ""
	resp.CardProfiles = CardProfiles{
		Row: CardProfileRow{
			Id:    fmt.Sprint(cardInfo.Lkey.LkeyId),
			Title: cardInfo.Title,
		},
	}
	resp.ContactlessStatus = "0"
	resp.ECNeedCAPAuth = "0"
	resp.ECNeedDynPwdAuth = "0"
	resp.ECNeedStaticAuth = "0"
	resp.ECNeedTokenAuth = "0"
	resp.ECStatus = "0"
	resp.ECUseDecoupledAuth = "0"
	resp.EMVOptionsCheckDisabled = "0"
	resp.ExpDate = cardInfo.ExpDate
	resp.FoundMBR = "0"
	resp.FoundPAN = cardInfo.Lkey.MaskedPan
	resp.IB_Registered = "0"
	resp.InstName = "ARVD"
	resp.IssueTechnology = "0"
	resp.LastATMUsed = ""
	resp.LastChangeStatusTime = ""
	resp.LastPOSUsed = ""
	resp.LastPVVChangeTime = ""
	resp.LastRefreshTime = ""

	lcar := len(cardInfo.CardAuthRestrictions)
	if lcar > 0 {
		resp.LastTranId = fmt.Sprint(cardInfo.CardAuthRestrictions[lcar-1].TlId)
		resp.LastTranTime = cardInfo.CardAuthRestrictions[lcar-1].When_created
	}
	resp.MaskBalances = "0"
	resp.MaskPVV = "0"
	resp.NameOnCard = cardInfo.EmbossName
	resp.PINVerifyType = ""
	resp.PVV = cardInfo.PVV
	resp.PasswordFlag = ""

	resp.PersonConfidential = PersonConfidential{
		Row: ConfidentialRow{
			What:         "",
			Value:        "",
			IsAllowedCST: "0",
			IsAllowedADS: "0",
			IsAllowedTB:  "0",
		},
	}
	resp.PersonExtId = cardInfo.CustomerCode
	resp.PersonFIO = cardInfo.LastName + " " + cardInfo.FirstName
	resp.PersonId = cardInfo.CustomerCode
	resp.PersonVIP = "0"
	resp.RequiredPasswordVersion = "1"
	resp.RiskControlDisabled = "0"
	resp.RiskLevel = "1"
	resp.Status = "1"
	resp.TmpECStatus = "-1"
	resp.Type = "1"

	return resp, nil
}
