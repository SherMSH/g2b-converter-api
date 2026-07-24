package reissuecardsout

import (
	models "converterapi/internal/models/OFFLINE"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"encoding/xml"
	"fmt"
)

// Root - корневой элемент XML
type Root struct {
	XMLName xml.Name         `xml:"ROOT"`
	Records []models.MRecord `xml:"RECORD"`
}

func (r Root) GetReqType() string {
	return string(utils.ReissueCardsOut)
}

func (r Root) GetRecords() []models.MRecord {
	return r.Records
}
func (r Root) GetRecordsCount() int {
	return len(r.Records)
}

func (r Root) Call() (respContent []byte, err error) {
	mdiData, cardsData, err := service.ReissueCardG2b(r)
	if err != nil {
		return []byte(err.Error()), err
	}

	if mdiData.Header.CActionCode != "0" {
		err = fmt.Errorf("%s - %s", mdiData.Header.CRspCode, mdiData.Header.IRejMsg)
		return []byte(err.Error()), err
	}

	for i := range mdiData.Details {
		if mdiData.Details[i].C_ACTIONCODE != "0" {
			break
		}
		if cardsData[i].CardBasicInfo.Lkey.LkeyId == 0 {
			continue
		}
		accnum := ""
		if len(cardsData[i].CardAccounts) != 0 {
			accnum = cardsData[i].CardAccounts[0].AccountNumber
		}
		pck := models.Pack{
			CustomerId:   cardsData[i].CardBasicInfo.Lkey.LkeyAlias,
			CustomerCode: cardsData[i].CardBasicInfo.CustomerCode,
			AccNum:       accnum,
			CurrencyCode: cardsData[i].CardBasicInfo.Currcode,
			LkeyAlias:    r.Records[i].ExternalID,
			CardPan:      mdiData.Details[i].KL_LKEY_CLR,
		}
		respContent = append(respContent, pck.GetData()...)
	}
	return respContent, nil
}
