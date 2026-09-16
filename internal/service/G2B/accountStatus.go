package service

import (
	"bytes"
	"converterapi/internal/config"
	d8corp "converterapi/internal/models/D8CORP"
	d8procweb "converterapi/pkg/d8-proc-web"
	"converterapi/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
)

// SetAccountStatusG2b меняет статус счёта.
//
// Отдельного метода для смены статуса счёта в xapi процессинга нет, а запись
// MDI типа ACCOUNT с действием UPDATE не находит счёт, если он заведён на
// филиал, а его клиент - на головную организацию. Поэтому статус меняется
// через updateAccount: процессинг требует прислать запись целиком, так что
// недостающие поля берём из текущего состояния счёта.
func SetAccountStatusG2b(accnum, statCode string) (err error) {
	acct, err := GetAcctInfoG2b(accnum)
	if err != nil {
		return err
	}
	if acct.Statcode == statCode {
		logger.Infof("[SERVICE] D8 G2b updateAccount: счёт %v уже в статусе %v", accnum, statCode)
		return nil
	}

	reqJSON, err := json.Marshal(d8procweb.RequestBody{
		Data: d8procweb.AccountUpdate{
			ID:         acct.ID,
			Recver:     acct.Recver,
			CompanyID:  acct.CompanyID,
			CustomerID: acct.CustomerID,
			Typecode:   acct.Typecode,
			Statcode:   statCode,
		},
	})
	if err != nil {
		return fmt.Errorf("updateAccount request marshaling error: %v", err)
	}

	// GetAcctInfoG2b закрывает свою сессию, поэтому логинимся заново
	signin, err := d8procweb.Signin()
	if err != nil {
		return fmt.Errorf("d8procweb signin error %v", err)
	}
	if signin.StatusCode != 200 {
		return fmt.Errorf("d8procweb signin status: %v", signin.StatusCode)
	}
	defer d8procweb.Signout()

	resp, err := d8procweb.Client.Post(config.Config.Processing.Address+"/api/miss/v1/updateAccount",
		"application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		return fmt.Errorf("updateAccount request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("updateAccount response reading error: %v", err)
	}
	logger.Infof("[SERVICE] D8 G2b updateAccount resp status: %v, body: %v (req %v)", resp.StatusCode, string(body), string(reqJSON))

	updResp := d8corp.CommonResp{}
	if err = json.Unmarshal(body, &updResp); err != nil {
		return fmt.Errorf("updateAccount resp unmarshaling error: %v", err)
	}
	if updResp.Status.Code != "0" || resp.StatusCode != 200 {
		return fmt.Errorf("%s - %s", updResp.Status.RspCode, updResp.Status.Message)
	}
	return nil
}
