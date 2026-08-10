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

// ReadPublicKey читает публичный RSA-ключ из файла.
//
// Поддерживаются оба представления, в которых процессинг отдаёт транспортный
// ключ: PKIX (SubjectPublicKeyInfo) и PKCS#1, как в сыром DER, так и в PEM.
// Транспортный ключ D8 приходит в PKCS#1, поэтому разбор только через PKIX
// на нём не работает.
func ReadPublicKey(filename string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// PEM - разбираем тело блока, иначе считаем содержимое сырым DER
	der := data
	if block, _ := pem.Decode(data); block != nil {
		der = block.Bytes
	}

	if key, err := x509.ParsePKIXPublicKey(der); err == nil {
		rsaKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("%s: not an RSA public key", filename)
		}
		return rsaKey, nil
	}

	rsaKey, err := x509.ParsePKCS1PublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse public key %s: %w", filename, err)
	}
	return rsaKey, nil
}

// EncryptWithRSA шифрует данные RSA публичным ключом (PKCS1 v1.5)
func EncryptWithRSA(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	// PKCS1 v1.5 padding
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
}

// EncryptWith3DES шифрует данные 3DES ключом в режиме ECB.
//
// Generate3DESKey отдаёт двойной ключ (16 байт), а Go принимает только тройной,
// поэтому ключ разворачивается в K1|K2|K1 - обычная схема 3DES EDE2.
func EncryptWith3DES(key, data []byte) ([]byte, error) {
	if len(key) == ZPKBytes {
		expanded, err := expand2Key3DES(key)
		if err != nil {
			return nil, err
		}
		key = expanded
	}

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
