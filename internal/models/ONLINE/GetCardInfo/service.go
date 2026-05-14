package getcardinfo

import (
	service "converterapi/internal/service/G2B"
	"fmt"
)

func Svc(sb *Body) (soapResp *Envelope, err error) {
	cardInfo, err := service.GetCardInfo(sb.SoapRq.Req.PAN, sb.SoapRq.Req.ExpirationDate)
	if err != nil {
		return nil, err
	}

	soapResp = new(Envelope)
	soapResp.XmlnsM0 = "http://schemas.compassplus.com/two/1.0/fimi_types.xsd"
	soapResp.XmlnsM1 = "http://schemas.compassplus.com/two/1.0/fimi.xsd"
	soapResp.XmlnsS = "http://www.w3.org/2003/05/soap-envelope"

	resp := Response{}
	resp.Product = sb.SoapRq.Req.Product
	resp.ResponseAttr = "1"
	resp.Ver = sb.SoapRq.Req.Ver

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

	resp.Accounts = Accounts{
		Row: accs,
	}
	resp.Acct2CardAttachType = "0"
	resp.CNSDisabled = "1"
	resp.CardAllowedEMVScript = ""
	resp.CardProfiles = CardProfiles{
		Row: CardProfileRow{
			Id:    fmt.Sprint(cardInfo.CardBasicInfo.Lkey.LkeyId),
			Title: cardInfo.CardBasicInfo.Title,
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
	resp.ExpDate = cardInfo.CardBasicInfo.ExpiryDate
	resp.FoundMBR = "0"
	resp.FoundPAN = cardInfo.CardBasicInfo.Lkey.MaskedPan
	resp.IB_Registered = "0"
	resp.InstName = "ARVD"
	resp.IssueTechnology = "0"
	resp.LastATMUsed = ""
	resp.LastChangeStatusTime = cardInfo.CardBasicInfo.StatChangeTime
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
	resp.NameOnCard = cardInfo.CardBasicInfo.EmbossName
	resp.PINVerifyType = ""
	resp.PVV = cardInfo.CardBasicInfo.Pvv
	resp.PasswordFlag = ""
	resp.UseUdCVV2 = fmt.Sprintf("%d", cardInfo.CardBasicInfo.Cvv2Type)

	if cardInfo.CardNotifications != nil {
		resp.PersonConfidential = PersonConfidential{
			Row: ConfidentialRow{
				What:         "phone",
				Value:        cardInfo.CardNotifications[0].NotificationTarget,
				IsAllowedCST: "0",
				IsAllowedADS: "0",
				IsAllowedTB:  "0",
			},
		}
	}
	resp.PersonExtId = cardInfo.CardBasicInfo.CustomerCode
	resp.PersonFIO = cardInfo.CardBasicInfo.LastName + " " + cardInfo.CardBasicInfo.FirstName
	resp.PersonId = cardInfo.CardBasicInfo.CustomerCode
	resp.PersonVIP = "0"
	resp.RequiredPasswordVersion = "1"
	resp.RiskControlDisabled = "0"
	resp.RiskLevel = "1"
	resp.Status = cardInfo.CardBasicInfo.StatCode
	resp.TmpECStatus = "-1"
	resp.Type = fmt.Sprintf("%d", cardInfo.CardBasicInfo.ProductType)

	soapResp.Body = RespBody{
		GetCardInfoRp: GetCardInfoRp{
			Response: resp,
		},
	}
	return soapResp, nil
}
