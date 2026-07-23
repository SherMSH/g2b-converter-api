package reissuecardsout

import (
	models "converterapi/internal/models/OFFLINE"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"encoding/xml"
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
	mdiData, err := service.ReissueCardG2b(r)
	if err != nil {
		return nil, err
	}

	for i, v := range mdiData.Details {
		if mdiData.Details[i].C_ACTIONCODE != "0" {
			break
		}
		pck := models.Pack{
			CustomerId:   r.Records[i].PCode,
			CustomerCode: r.Records[i].ExtID,
			AccNum:       r.Records[i].Account,
			CurrencyCode: r.Records[i].CurrencyNo,
			LkeyAlias:    r.Records[i].ExternalID,
			CardPan:      v.KL_LKEY_CLR,
		}
		respContent = append(respContent, pck.GetData()...)
	}
	return respContent, nil
}
