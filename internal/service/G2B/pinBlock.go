package service

import (
	d8corp "converterapi/internal/models/D8CORP"
	"converterapi/pkg/crypto"
	"converterapi/pkg/logger"
	"fmt"
)

// transportKeyPath - публичный ключ процессинга для шифрования одноразового
// DES-ключа. Путь относительный, поэтому сервис обязан запускаться из корня.
var transportKeyPath = "internal/app/files/transport_setpin.der"

// buildPinRequest собирает запрос с зашифрованным PIN-блоком.
//
// Схема по спецификации (7.8, 7.22): одноразовый 3DES-ключ шифруется публичным
// RSA-ключом процессинга, PIN-блок формата ANSI X9.8 шифруется этим же
// 3DES-ключом. Именно им процессинг и расшифрует блок, поэтому лишних слоёв
// шифрования быть не должно.
func buildPinRequest(pan, pin, expDate string) (d8corp.SetPinReq, error) {
	key3DES, err := crypto.Generate3DESKey()
	if err != nil {
		logger.Errorf("generate 3DES key error: %v", err)
		return d8corp.SetPinReq{}, fmt.Errorf("generate 3DES key error: %v", err)
	}

	publicKey, err := crypto.ReadPublicKey(transportKeyPath)
	if err != nil {
		logger.Errorf("read public key error: %v", err)
		return d8corp.SetPinReq{}, fmt.Errorf("read public key error: %v", err)
	}

	pinKeyUnderRSA, err := crypto.EncryptWithRSA(publicKey, key3DES)
	if err != nil {
		logger.Errorf("encrypt key with RSA: %v", err)
		return d8corp.SetPinReq{}, fmt.Errorf("encrypt key with RSA: %v", err)
	}

	clear, err := crypto.Format0(pin, pan)
	if err != nil {
		logger.Errorf("pin block format0: %v", err)
		return d8corp.SetPinReq{}, fmt.Errorf("pin block format0: %v", err)
	}

	pinBlock, err := crypto.EncryptWith3DES(key3DES, clear)
	if err != nil {
		logger.Errorf("encrypt pin block: %v", err)
		return d8corp.SetPinReq{}, fmt.Errorf("encrypt pin block: %v", err)
	}

	return d8corp.SetPinReq{
		CardKey: d8corp.CardKey{
			Pan:        pan,
			ExpiryDate: expDate,
		},
		PinKeyUnderRSA: crypto.HexUpper(pinKeyUnderRSA),
		PinBlock:       crypto.HexUpper(pinBlock),
		PinBlockType:   0, // 0 - ANSI X9.8
	}, nil
}
