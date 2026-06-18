package getcardinfo

import (
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
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
	resp.Echo = sb.SoapRq.Req.Echo
	resp.Product = sb.SoapRq.Req.Product
	resp.ResponseAttr = "1"
	resp.Ver = sb.SoapRq.Req.Ver
	resp.TranId = utils.GenerateTimestampID()

	var (
		accs []AccountRow
		nums string
	)

	if cardInfo.CardAccounts != nil {
		accrow := AccountRow{
			AcctNo:        cardInfo.CardAccounts[0].AccountNumber,
			Status:        utils.AccountStatuses[cardInfo.CardAccounts[0].StatCode],
			LedgerBalance: fmt.Sprintf("%.2f", cardInfo.CardAccounts[0].AvlBal+cardInfo.CardAccounts[0].BlkAmt),
			AvailBalance:  fmt.Sprintf("%.2f", cardInfo.CardAccounts[0].AvlBal),
			Currency:      utils.Currencies[cardInfo.CardAccounts[0].Currency],
			Type:          utils.AccountTypes[cardInfo.CardAccounts[0].TypeCode],
			AccountStatus: utils.AccountStatuses[cardInfo.CardAccounts[0].StatCode],
		}
		accs = append(accs, accrow)
	}

	resp.Accounts = Accounts{
		Row: accs,
	}
	resp.Acct2CardAttachType = "1"
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
	resp.ECStatus = "-1"
	resp.ECUseCardSettingsAuth = "0"
	resp.ECUseDecoupledAuth = "0"
	resp.EMVOptionsCheckDisabled = "0"
	resp.ExpDate = utils.ConvertDate(cardInfo.CardBasicInfo.ExpiryDate)
	resp.FoundMBR = "0"
	resp.FoundPAN = cardInfo.CardBasicInfo.Lkey.Pan
	resp.IB_Registered = "0"
	resp.InstName = "ARVD"
	resp.IssueTechnology = "1"
	if len(cardInfo.CardTransactions) != 0 {
		resp.LastATMUsed = utils.ConvertD8Tmstmp(cardInfo.CardTransactions[0].Tstamp_insert)
	}
	resp.LastChangeStatusTime = utils.ConvertD8Tmstmp(cardInfo.CardBasicInfo.StatChangeTime)
	if len(cardInfo.CardTransactions) != 0 {
		resp.LastPOSUsed = utils.ConvertD8Tmstmp(cardInfo.CardTransactions[0].Tstamp_insert)
	}
	resp.LastPVVChangeTime = ""
	resp.LastRefreshTime = ""

	lcar := len(cardInfo.CardTransactions)
	if lcar > 0 {
		resp.LastTranId = fmt.Sprint(cardInfo.CardTransactions[lcar-1].TlId)
		resp.LastTranTime = utils.ConvertD8Tmstmp(cardInfo.CardTransactions[lcar-1].When_created)
	}
	resp.MaskBalances = "0"
	resp.MaskPVV = "0"
	resp.NameOnCard = cardInfo.CardBasicInfo.EmbossName
	resp.PINVerifyType = ""
	resp.PVV = cardInfo.CardBasicInfo.Pvv
	resp.PasswordFlag = "0"
	resp.UseUdCVV2 = fmt.Sprintf("%d", cardInfo.CardBasicInfo.Cvv2Type)

	if cardInfo.CardNotifications != nil {
		nums = cardInfo.CardNotifications[0].NotificationTarget
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
	resp.PersonFIO = fmt.Sprintf("%v %v", nums, cardInfo.CardBasicInfo.EmbossName)
	resp.PersonId = cardInfo.CardBasicInfo.CustomerCode
	resp.PersonVIP = "0"
	resp.RequiredPasswordVersion = "1"
	resp.RiskControlDisabled = "0"
	resp.RiskLevel = "1"
	resp.Status = "1" //cardInfo.CardBasicInfo.StatCode
	resp.TmpECStatus = "-1"
	resp.Type = utils.CardTypes[cardInfo.CardBasicInfo.ProductType]
	soapResp.Body = RespBody{
		GetCardInfoRp: GetCardInfoRp{
			Response: resp,
		},
	}
	return soapResp, nil
}
