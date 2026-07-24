package relinkpreissuedcardsout

import (
	models "converterapi/internal/models/OFFLINE"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"encoding/xml"
	"fmt"
)

type Root struct {
	XMLName xml.Name         `xml:"ROOT"`
	Records []models.MRecord `xml:"RECORD"`
}

func (r Root) GetReqType() string {
	return string(utils.RelinkPreIssuedCardsOut)
}

func (r Root) GetRecords() []models.MRecord {
	return r.Records
}
func (r Root) GetRecordsCount() int {
	return len(r.Records)
}

func (r Root) Call() (respContent []byte, err error) {
	mdiData, err := service.RelinkPreissiedCardG2b(r)
	if err != nil {
		return nil, err
	}
	if mdiData.Header.CActionCode != "0" {
		return nil, fmt.Errorf("%s - %s", mdiData.Header.CActionCode, mdiData.Header.IRejMsg)
	}
	for i := range r.Records {
		if len(mdiData.Details) > 0 {
			if mdiData.Details[i].C_ACTIONCODE != "0" {
				break
			}
		}
		pck := models.Pack{
			CustomerId:   r.Records[i].PCode,
			CustomerCode: r.Records[i].ExtID,
			AccNum:       r.Records[i].Account,
			CurrencyCode: "972",
			LkeyAlias:    r.Records[i].ExternalID,
			CardPan:      r.Records[i].PAN,
		}
		respContent = append(respContent, pck.GetData()...)
	}
	return respContent, nil
}
