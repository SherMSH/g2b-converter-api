package getcardinfo

import (
	service "converterapi/internal/service/G2B"
	"fmt"
)

func Svc(sb *Body) (rsp *Envelope, err error) {
	cardInfo, err := service.GetCardInfo(sb.SoapRq.Req.PAN, sb.SoapRq.Req.ExpirationDate)
	if err != nil {
		return nil, err
	}

	rsp = new(Envelope)
	rsp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	rsp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	rsp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	gcirp := GetCardInfoRp{}
	gcirp.Product = sb.SoapRq.Req.Product
	gcirp.Response = "1"
	gcirp.Ver = sb.SoapRq.Req.Ver

	var accs []AccountRow

	if cardInfo.CardAccounts != nil {
		accrow := AccountRow{
			AcctNo:        cardInfo.CardAccounts[0].AccountNumber,
			Status:        cardInfo.CardAccounts[0].StatCode,
			LedgerBalance: "",
			AvailBalance:  fmt.Sprintf("%.2f", cardInfo.CardAccounts[0].AvlBal),
			Currency:      cardInfo.CardAccounts[0].Currency,
			Type:          cardInfo.CardAccounts[0].TypeCode,
			AccountStatus: cardInfo.CardAccounts[0].StatCode,
		}
		accs = append(accs, accrow)
	}

	gcirp.Accounts = Accounts{
		Row: accs,
	}
	gcirp.Acct2CardAttachType = "0"
	gcirp.CNSDisabled = "1"
	gcirp.CardAllowedEMVScript = ""
	gcirp.CardProfiles = CardProfiles{
		Row: CardProfileRow{
			Id:    fmt.Sprint(cardInfo.CardBasicInfo.Lkey.LkeyId),
			Title: cardInfo.CardBasicInfo.Title,
		},
	}
	gcirp.ContactlessStatus = "0"
	gcirp.ECNeedCAPAuth = "0"
	gcirp.ECNeedDynPwdAuth = "0"
	gcirp.ECNeedStaticAuth = "0"
	gcirp.ECNeedTokenAuth = "0"
	gcirp.ECStatus = "0"
	gcirp.ECUseDecoupledAuth = "0"
	gcirp.EMVOptionsCheckDisabled = "0"
	gcirp.ExpDate = cardInfo.CardBasicInfo.ExpiryDate
	gcirp.FoundMBR = "0"
	gcirp.FoundPAN = cardInfo.CardBasicInfo.Lkey.MaskedPan
	gcirp.IB_Registered = "0"
	gcirp.InstName = "ARVD"
	gcirp.IssueTechnology = "0"
	gcirp.LastATMUsed = ""
	gcirp.LastChangeStatusTime = ""
	gcirp.LastPOSUsed = ""
	gcirp.LastPVVChangeTime = ""
	gcirp.LastRefreshTime = ""

	lcar := len(cardInfo.CardAuthRestrictions)
	if lcar > 0 {
		gcirp.LastTranId = fmt.Sprint(cardInfo.CardAuthRestrictions[lcar-1].TlId)
		gcirp.LastTranTime = cardInfo.CardAuthRestrictions[lcar-1].When_created
	}
	gcirp.MaskBalances = "0"
	gcirp.MaskPVV = "0"
	gcirp.NameOnCard = cardInfo.CardBasicInfo.EmbossName
	gcirp.PINVerifyType = ""
	gcirp.PVV = cardInfo.CardBasicInfo.Pvv
	gcirp.PasswordFlag = ""

	gcirp.PersonConfidential = PersonConfidential{
		Row: ConfidentialRow{
			What:         "",
			Value:        "",
			IsAllowedCST: "0",
			IsAllowedADS: "0",
			IsAllowedTB:  "0",
		},
	}
	gcirp.PersonExtId = cardInfo.CardBasicInfo.CustomerCode
	gcirp.PersonFIO = cardInfo.CardBasicInfo.LastName + " " + cardInfo.CardBasicInfo.FirstName
	gcirp.PersonId = cardInfo.CardBasicInfo.CustomerCode
	gcirp.PersonVIP = "0"
	gcirp.RequiredPasswordVersion = "1"
	gcirp.RiskControlDisabled = "0"
	gcirp.RiskLevel = "1"
	gcirp.Status = "1"
	gcirp.TmpECStatus = "-1"
	gcirp.Type = "1"

	rsp.Body = RespBody{
		GetCardInfoRp: gcirp,
	}
	return rsp, nil
}
