package crypto

import (
	"crypto/aes"
	"crypto/des"
	"crypto/rand"
	"errors"
	"fmt"
)

const (
	ZPKBytes    = 16 // 32 hex digits (double-length 3DES ZPK)
	PINBlockLen = 8
)

// Format0 builds a clear ISO 9564-1 Format 0 PIN block (8 bytes).
func Format0(pin string) ([]byte, error) {
	if len(pin) < 4 || len(pin) > 12 {
		return nil, errors.New("PIN length must be 4..12 digits")
	}
	for _, c := range pin {
		if c < '0' || c > '9' {
			return nil, errors.New("PIN must contain only decimal digits")
		}
	}

	block := make([]byte, PINBlockLen)
	for i := range block {
		block[i] = 0xFF
	}
	block[0] = byte(len(pin) & 0x0F)

	for i, c := range pin {
		d := byte(c - '0')
		idx := 1 + i/2
		if i%2 == 0 {
			block[idx] = (d << 4) | (block[idx] & 0x0F)
		} else {
			block[idx] = (block[idx] & 0xF0) | d
		}
	}
	if len(pin)%2 == 1 {
		idx := 1 + len(pin)/2
		block[idx] = (block[idx] & 0xF0) | 0x0F
	}
	return block, nil
}

// GenerateZPK32 returns a random 16-byte ZPK with odd parity (32 hex digits).
func GenerateZPK32() ([]byte, error) {
	key := make([]byte, ZPKBytes)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	for i, b := range key {
		if bitsOnes(b)%2 == 0 {
			key[i] = b ^ 0x01
		}
	}
	return key, nil
}
func bitsOnes(b byte) int {
	n := 0
	for b != 0 {
		n += int(b & 1)
		b >>= 1
	}
	return n
}

// Encrypt3DES encrypts an 8-byte PIN block with ZPK using 3DES-ECB (2-key).
func Encrypt3DES(zpk, clear []byte) ([]byte, error) {
	if len(clear) != PINBlockLen {
		return nil, fmt.Errorf("clear PIN block must be %d bytes", PINBlockLen)
	}
	key24, err := expand2Key3DES(zpk)
	if err != nil {
		return nil, err
	}
	block, err := des.NewTripleDESCipher(key24)
	if err != nil {
		return nil, err
	}
	out := make([]byte, PINBlockLen)
	block.Encrypt(out, clear)
	return out, nil
}

func expand2Key3DES(key16 []byte) ([]byte, error) {
	if len(key16) != ZPKBytes {
		return nil, fmt.Errorf("ZPK must be %d bytes, got %d", ZPKBytes, len(key16))
	}
	out := make([]byte, 24)
	copy(out, key16)
	copy(out[16:], key16[:8])
	return out, nil
}

// GenerateZPK256 returns a random 32-byte key (AES-256).
func GenerateZPK256() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// EncryptAES encrypts an 8-byte PIN block with a 32-byte ZPK using AES-256-ECB.
func EncryptAES(zpk, clear []byte) ([]byte, error) {
	if len(zpk) != 32 {
		return nil, errors.New("AES ZPK must be 32 bytes")
	}
	if len(clear) != PINBlockLen {
		return nil, fmt.Errorf("clear PIN block must be %d bytes", PINBlockLen)
	}
	block, err := aes.NewCipher(zpk)
	if err != nil {
		return nil, err
	}
	out := make([]byte, PINBlockLen)
	block.Encrypt(out, clear)
	return out, nil
}
