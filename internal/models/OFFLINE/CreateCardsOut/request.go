package createcardsout

import (
	models "converterapi/internal/models/OFFLINE"
	service "converterapi/internal/service/G2B"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
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

func (r Root) Call() (respContent []byte, err error) {
	mdiData, err := service.AddCardsG2b(r)
	if err != nil {
		return []byte(err.Error()), err
	}

	if mdiData.Header.CActionCode != "0" {
		err = fmt.Errorf("%s - %s", mdiData.Header.CRspCode, mdiData.Header.IRejMsg)
		return []byte(err.Error()), err
	}

	// В пакете вместе с картами идут контракты оповещений, поэтому детали
	// ответа больше не совпадают с записями файла один к одному: берём только
	// карты, а отказы по оповещениям отмечаем в логе - выпуск они не отменяют.
	card := 0
	for _, v := range mdiData.Details {
		if v.ISS_RECTYPE != "CARD" {
			if v.C_ACTIONCODE != "0" {
				logger.Errorf("[CreateCardsOut] запись %s #%d отклонена: %s - %s",
					v.ISS_RECTYPE, v.ISS_RECNUM, v.C_RSPCODE, v.I_REJMSG)
			}
			continue
		}
		if card >= len(r.Records) {
			break
		}
		if v.C_ACTIONCODE != "0" {
			logger.Errorf("[CreateCardsOut] карта %s не выпущена: %s - %s",
				r.Records[card].ExternalID, v.C_RSPCODE, v.I_REJMSG)
			card++
			continue
		}
		pck := models.Pack{
			CustomerId:   r.Records[card].PCode,
			CustomerCode: r.Records[card].ExtID,
			AccNum:       r.Records[card].Account,
			CurrencyCode: r.Records[card].CurrencyNo,
			LkeyAlias:    r.Records[card].ExternalID,
			CardPan:      v.KL_LKEY_CLR,
		}
		respContent = append(respContent, pck.GetData()...)
		card++
	}
	return respContent, nil
}
