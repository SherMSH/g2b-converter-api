package createcardsout

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
	return string(utils.CreateCardsOut)
}

func (r Root) GetRecords() []models.MRecord {
	return r.Records
}
func (r Root) GetRecordsCount() int {
	return len(r.Records)
}

func (r Root) Call() ([]byte, error) {
	// err := service.CreateCustomer(r)
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = service.CreateAccountIfExist(r)
	// if err != nil {
	// 	return nil, err
	// }
	mdiData, err := service.AddCardG2b(r)
	if err != nil {
		return nil, err
	}
	if len(mdiData.Details) > r.GetRecordsCount() {
		return nil, fmt.Errorf("internal error: mdi response length mismatch")
	}
	for i, v := range mdiData.Details {
		r.Records[i].PAN = v.KL_LKEY_CLR
		r.Records[i].MBR = "0"
	}

	respContent, err := xml.Marshal(r)
	if err != nil {
		return nil, err
	}

	return respContent, nil
}
