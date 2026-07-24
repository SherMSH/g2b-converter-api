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
		return []byte(err.Error()), err
	}

	if mdiData.Header.CActionCode != "0" {
		err = fmt.Errorf("%s - %s", mdiData.Header.CRspCode, mdiData.Header.IRejMsg)
		return []byte(err.Error()), err
	}
	service.DeleteCardAcctLinkG2b(r)
	service.AddCardAcctLinkG2b(r)
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
			LkeyAlias:    r.Records[i].ExtID,
			CardPan:      r.Records[i].PAN,
		}
		respContent = append(respContent, pck.GetData()...)
	}
	return respContent, nil
}
