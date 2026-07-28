package service

import (
	"converterapi/internal/models"
	d8procweb "converterapi/pkg/d8-proc-web"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
)

const (
	ArvCompanyId int = 41
)

func AddCompaniesG2b(input models.MDIface) (companies []d8procweb.CompaniesData, err error) {
	var (
		path string = "/api/miss/v1/addCardProduct"
	)
	for i, v := range input.GetRecords() {
		companyData := d8procweb.CompaniesData{
			RegNumber: v.ExtID,
			Name:      v.Name,
			ParentId:  ArvCompanyId,
		}
		err = nil

		raw, err := json.Marshal(companyData)
		if err != nil {
			logger.Errorf("[SERVICE] AddCompaniesG2b json marshal error: %v; record %d skipped", err, i)
			continue
		}
		resp, err := d8procweb.PutRequest(path, raw)
		if err != nil {
			logger.Errorf("[SERVICE] AddCompaniesG2b error: %v; record %d skipped", err, i)
			continue
		}

		newData := new(d8procweb.CompaniesData)
		err = json.Unmarshal(resp, newData)
		if err != nil {
			logger.Errorf("[SERVICE] AddCompaniesG2b json unmarshal error: %v; record %d skipped", err, i)
			continue
		}
		if newData.ID == 0 || len(newData.RegNumber) == 0 {
			logger.Errorf("[SERVICE] AddCompaniesG2b record %d skipped (D8 sent empty data)", i)
			continue
		}
		companyData.ID = newData.ID
		companyData.RegNumber = newData.RegNumber
		companies = append(companies, companyData)
	}

	if len(companies) == 0 {
		return nil, fmt.Errorf("[SERVICE] AddCompaniesG2b error occured (see logs)")
	}
	return companies, nil
}
