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
	"path/filepath"
)

func ConvScanner() {
	logger.Infof("[JOBS] Converter scanner")

	for _, v := range utils.OfflineReqTypes {
		reqOf, err := unmarshalFromFormat(v)
		if err != nil {
			if !os.IsNotExist(err) {
				logger.Errorf("Unmarshal from file %v Error: %v", v, err)
			}
			continue
		}
		logger.Infof("Converter Scans %v req", v)

		sourceMask := config.Config.App.Storage.Basepath + config.Config.App.Storage.In + "/" + string(v)
		matches, err := filepath.Glob(sourceMask)
		if err != nil {
			logger.Errorf("Error founding matches for %s: %v", sourceMask, err.Error())
			continue
		}
		if len(matches) == 0 {
			logger.Errorf("Empty matches for %s", sourceMask)
			continue
		}
		sourcePath := matches[0]
		base := filepath.Base(sourcePath)
		// prefix := string(v)[:len(v)-5]
		// suffix := ".xml"
		// timestamp := strings.TrimPrefix(base, prefix)
		// timestamp = strings.TrimSuffix(timestamp, suffix)

		destPath := config.Config.App.Storage.Basepath + config.Config.App.Storage.Out + "/" + base
		content, err := reqOf.Call()
		if err != nil {
			logger.Errorf("Converter Scanner service %v call error: %v", v, err)
			continue
		}
		err2 := storage.MoveFile(sourcePath, destPath, content)
		if err2 != nil {
			logger.Warnf("Error mv file %v: %v", sourcePath, err2)
		}
	}
}

func unmarshalFromFormat(ort utils.OfflineReqType) (service.G2bServiceIface, error) {
	source := config.Config.App.Storage.Basepath + config.Config.App.Storage.In + "/" + string(ort)
	matches, err := filepath.Glob(source)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, os.ErrNotExist
	}
	nameFile := matches[0]
	data, err := storage.LoadFile(nameFile)
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
