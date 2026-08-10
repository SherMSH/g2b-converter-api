package service

import (
	"converterapi/internal/config"
	"converterapi/internal/utils"
	"os"
	"testing"
)

// Живые проверки против стенда процессинга.
//
// По умолчанию пропускаются: они требуют доступа к D8 и обращаются к реальным
// данным. Запуск:
//
//	D8_LIVE=1 go test ./internal/service/G2B/ -run TestLive -v
func liveSetup(t *testing.T) {
	t.Helper()
	if os.Getenv("D8_LIVE") == "" {
		t.Skip("живой тест: нужен D8_LIVE=1 и доступ к стенду")
	}
	if err := config.ReadFileConfigs("../../config/config.json"); err != nil {
		t.Fatalf("конфиг не прочитан: %v", err)
	}
	utils.Init()
}

// Расчёт комиссии не двигает средства, поэтому его можно гонять свободно.
func TestLiveCalculateAcqCommission(t *testing.T) {
	liveSetup(t)

	ecTxRefNo, err := InitiateTransaction()
	if err != nil {
		t.Fatalf("InitiateTransaction: %v", err)
	}
	t.Logf("ecTxRefno = %s", *ecTxRefNo)

	in := trnInput{
		txnType: utils.Sales,
		pan:     "5058270530003879",
	}

	commission, err := CalculateAcqCommission(in, *ecTxRefNo)
	if err != nil {
		t.Fatalf("CalculateAcqCommission: %v", err)
	}

	t.Logf("итого комиссия: %.2f %s", commission.TotalCommAmount, commission.CommCurrency)
	for _, item := range commission.AcqCommissionItems {
		t.Logf("  %s: %.2f с суммы %.2f", item.DefCommCode, item.AmtComm, item.AmtSource)
	}
}

// Расчёт комиссии для перевода: показывает, доходит ли до процессинга сам тип
// операции TRANSF_C2C, если авторизация по нему отвечает System malfunction.
func TestLiveCalculateAcqCommissionTransfer(t *testing.T) {
	liveSetup(t)

	ecTxRefNo, err := InitiateTransaction()
	if err != nil {
		t.Fatalf("InitiateTransaction: %v", err)
	}

	in := trnInput{
		txnType:      utils.C2C,
		pan:          "5058270530003879",
		recipientPan: "5058270530000073",
	}

	commission, err := CalculateAcqCommission(in, *ecTxRefNo)
	if err != nil {
		// На тестовом стенде переводы не настроены: authorizeTransaction по ним
		// отвечает System malfunction, расчёт комиссии - status code 1.
		t.Skipf("переводы не поддерживаются стендом: %v", err)
	}
	t.Logf("комиссия перевода: %.2f %s, составляющих: %d",
		commission.TotalCommAmount, commission.CommCurrency, len(commission.AcqCommissionItems))
}
