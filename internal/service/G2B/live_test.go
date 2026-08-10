package service

import (
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/pkg/crypto"
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
	// тест запускается из каталога пакета, путь к ключу считается от корня
	transportKeyPath = "../../app/files/transport_setpin.der"
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

// A/B-проверка сборки PIN-блока.
//
// В SetPinG2b блок шифруется дважды: сначала случайным ZPK, который никуда не
// передаётся, затем транспортным ключом. По спецификации (7.8) слой должен быть
// один - тем ключом, что уходит в pinKeyUnderRSA.
//
// Коды ответа verifyPIN различают ситуации: 1/17 - неверный PIN (значит блок
// расшифрован и разобран), 1/26 - invalid PIN block (значит не расшифрован).
func TestLivePinBlockVariants(t *testing.T) {
	liveSetup(t)

	const (
		pan     = "5058270530003879"
		expDate = "3004"
		pin     = "1234"
	)

	t.Run("один слой - как в спецификации", func(t *testing.T) {
		req, err := buildPinRequest(pan, pin, expDate)
		if err != nil {
			t.Fatalf("сборка запроса: %v", err)
		}
		t.Logf("результат: %v", verifyPinRequest(req))
	})

	t.Run("два слоя - как в текущем SetPinG2b", func(t *testing.T) {
		req, err := buildLegacyPinRequest(pan, pin, expDate)
		if err != nil {
			t.Fatalf("сборка запроса: %v", err)
		}
		t.Logf("результат: %v", verifyPinRequest(req))
	})
}

// buildLegacyPinRequest повторяет сборку PIN-блока из SetPinG2b: с лишним
// слоем шифрования случайным ZPK.
func buildLegacyPinRequest(pan, pin, expDate string) (d8corp.SetPinReq, error) {
	key3DES, err := crypto.Generate3DESKey()
	if err != nil {
		return d8corp.SetPinReq{}, err
	}
	publicKey, err := crypto.ReadPublicKey("../../app/files/transport_setpin.der")
	if err != nil {
		return d8corp.SetPinReq{}, err
	}
	pinKeyUnderRSA, err := crypto.EncryptWithRSA(publicKey, key3DES)
	if err != nil {
		return d8corp.SetPinReq{}, err
	}
	clear, err := crypto.Format0(pin, pan)
	if err != nil {
		return d8corp.SetPinReq{}, err
	}
	zpk, err := crypto.GenerateZPK32()
	if err != nil {
		return d8corp.SetPinReq{}, err
	}
	underZpk, err := crypto.Encrypt3DES(zpk, clear)
	if err != nil {
		return d8corp.SetPinReq{}, err
	}
	pinBlock, err := crypto.EncryptWith3DES(key3DES, underZpk)
	if err != nil {
		return d8corp.SetPinReq{}, err
	}
	return d8corp.SetPinReq{
		CardKey:        d8corp.CardKey{Pan: pan, ExpiryDate: expDate},
		PinKeyUnderRSA: crypto.HexUpper(pinKeyUnderRSA),
		PinBlock:       crypto.HexUpper(pinBlock),
		PinBlockType:   0,
	}, nil
}
