package crypto

import (
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
)

func Generate3DESKey() ([]byte, error) {
	// Генерируем случайный 16-байтовый ключ
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	// Настраиваем каждый байт для обеспечения нечетной четности
	for i := range key {
		keyByte := key[i]
		// Подсчитываем количество установленных бит
		if countBits(keyByte)%2 == 0 {
			// Если четность четная, инвертируем младший бит
			keyByte ^= 0x01
		}
		key[i] = keyByte
	}

	return key, nil
}

func countBits(b byte) int {
	count := 0
	for b > 0 {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

func ReadPublicKey(filename string) (*rsa.PublicKey, error) {
	// data, err := os.ReadFile(filename)
	// if err != nil {
	// 	return nil, err
	// }
	// if key, err := x509.ParsePKIXPublicKey(data); err == nil {
	// 	pub, ok := key.(*rsa.PublicKey)
	// 	if !ok {
	// 		return nil, fmt.Errorf("%s: not an RSA public key", filename)
	// 	}
	// 	return pub, nil
	// }
	// return x509.ParsePKCS1PublicKey(data)

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// Парсим DER
	key, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		// Пробуем парсить как PKCS1
		block, _ := pem.Decode(data)
		if block != nil {
			key, err = x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("parse PEM: %w", err)
			}
		} else {
			return nil, fmt.Errorf("parse DER: %w", err)
		}
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaKey, nil
}

// EncryptWithRSA шифрует данные RSA публичным ключом (PKCS1 v1.5)
func EncryptWithRSA(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	// PKCS1 v1.5 padding
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
}

// EncryptWith3DES шифрует данные 3DES ключом в режиме ECB
func EncryptWith3DES(key, data []byte) ([]byte, error) {
	// Создаем 3DES шифр
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	// Проверяем, что данные кратны блоку (8 байт)
	if len(data)%des.BlockSize != 0 {
		return nil, fmt.Errorf("data length %d is not a multiple of block size %d", len(data), des.BlockSize)
	}

	// ECB режим - просто шифруем каждый блок
	encrypted := make([]byte, len(data))
	for i := 0; i < len(data); i += des.BlockSize {
		block.Encrypt(encrypted[i:i+des.BlockSize], data[i:i+des.BlockSize])
	}

	return encrypted, nil
}

func HexUpper(data []byte) string {
	return hex.EncodeToString(data)
}
