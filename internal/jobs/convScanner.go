package jobs

import (
	"converterapi/internal/config"
	createcardsout "converterapi/internal/models/OFFLINE/CreateCardsOut"
	createcustomerandaccount "converterapi/internal/models/OFFLINE/CreateCustomerAndAccount"
	createorganizations "converterapi/internal/models/OFFLINE/CreateOrganizations"
	createpreissuedcards "converterapi/internal/models/OFFLINE/CreatePreIssuedCards"
	createstatusactivationsout "converterapi/internal/models/OFFLINE/CreateStatusActivationsOut"
	reissuecardsout "converterapi/internal/models/OFFLINE/ReissueCardsOut"
	relinkpreissuedcardstatusactivationsout "converterapi/internal/models/OFFLINE/RelinkPreIssuedCardStatusActivationsOut"
	relinkpreissuedcardsout "converterapi/internal/models/OFFLINE/RelinkPreIssuedCardsOut"
	"converterapi/internal/service"
	"converterapi/internal/utils"
	"converterapi/pkg/logger"
	"converterapi/pkg/storage"
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

func ConvScanner() {
	logger.Infof("[JOBS] Converter scanner")

	for _, v := range utils.OfflineReqTypes {
		reqOf, err := unmarshalFromFile(v)
		if err != nil {
			if !os.IsNotExist(err) {
				logger.Errorf("Unmarshal from file %v Error: %v", v, err)
			}
			continue
		}
		logger.Infof("Converter Scans %v req", v)

		err = reqOf.Call()
		if err != nil {
			logger.Errorf("Converter Scanner service %v call error: %v", v, err)
			continue
		}
		sourcePath := config.Config.App.Storage.Basepath + config.Config.App.Storage.In + "/" + string(v)
		destPath := config.Config.App.Storage.Basepath + config.Config.App.Storage.Out + "/" + time.Now().Format("2006_01_02T15_04_05Z07_00") + string(v)
		storage.MoveFile(sourcePath, destPath)
	}
}

func unmarshalFromFile(ort utils.OfflineReqType) (service.G2bServiceIface, error) {
	source := config.Config.App.Storage.Basepath + config.Config.App.Storage.In + "/" + string(ort)
	data, err := storage.DownloadFile(source)
	if err != nil {
		return nil, err
	}

	switch ort {
	case utils.CreateCardsOut:
		var root createcardsout.Root
		err = xml.Unmarshal(data, &root)
		if err != nil {
			logger.Errorf("xml unmarshal from file err: %v", err)
			return nil, fmt.Errorf("ошибка парсинга %s: %w", ort, err)
		}
		return root, nil
	case utils.CreateCustomerAndAccount:
		var root createcustomerandaccount.Root
		err = xml.Unmarshal(data, &root)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга %s: %w", ort, err)
		}
		return root, nil
	case utils.CreateOrganizations:
		var root createorganizations.Root
		err = xml.Unmarshal(data, &root)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга %s: %w", ort, err)
		}
		return root, nil
	case utils.CreatePreIssuedCards:
		var root createpreissuedcards.Root
		err = xml.Unmarshal(data, &root)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга %s: %w", ort, err)
		}
		return root, nil
	case utils.CreateStatusActivationsOut:
		var root createstatusactivationsout.Root
		err = xml.Unmarshal(data, &root)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга %s: %w", ort, err)
		}
		return root, nil
	case utils.ReissueCardsOut:
		var root reissuecardsout.Root
		err = xml.Unmarshal(data, &root)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга %s: %w", ort, err)
		}
		return root, nil
	case utils.RelinkPreIssuedCardsOut:
		var root relinkpreissuedcardsout.Root
		err = xml.Unmarshal(data, &root)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга %s: %w", ort, err)
		}
		return root, nil
	case utils.RelinkPreIssuedCardStatusActivationsOut:
		var root relinkpreissuedcardstatusactivationsout.Root
		err = xml.Unmarshal(data, &root)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга %s: %w", ort, err)
		}
		return root, nil
	default:
		return nil, fmt.Errorf("Неизвестная ошибка code: 20.00")
	}
}
