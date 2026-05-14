package crypto

import (
	"crypto/des"
	"fmt"
)

// TripleDESEncrypter encrypts 8-byte blocks with 3DES in ECB mode.
type TripleDESEncrypter struct{}

// NewTripleDESEncrypter returns a TripleDESEncrypter.
func NewTripleDESEncrypter() *TripleDESEncrypter {
	return &TripleDESEncrypter{}
}

// EncryptECB encrypts plaintext with 3DES in ECB mode.
// Supported key sizes: 16 bytes (two-key EDE, expanded to K1||K2||K1), 24 bytes (triple-length K1||K2||K3),
// 32 bytes (ZPK: first 24 bytes are used as triple-length 3DES material; remaining 8 are ignored for 3DES only).
// Plaintext length must be a multiple of 8.
func (e *TripleDESEncrypter) EncryptECB(key, plaintext []byte) ([]byte, error) {
	k24, err := resolveTripleDESKey(key)
	if err != nil {
		return nil, err
	}
	if len(plaintext) == 0 || len(plaintext)%8 != 0 {
		return nil, fmt.Errorf("3des: plaintext length must be non-zero multiple of 8, got %d", len(plaintext))
	}
	cipher, err := des.NewTripleDESCipher(k24)
	if err != nil {
		return nil, fmt.Errorf("3des: %w", err)
	}
	out := make([]byte, len(plaintext))
	for off := 0; off < len(plaintext); off += 8 {
		cipher.Encrypt(out[off:off+8], plaintext[off:off+8])
	}
	return out, nil
}

func resolveTripleDESKey(key []byte) ([]byte, error) {
	switch len(key) {
	case 16:
		return expandTwoKey168(key), nil
	case 24:
		k := make([]byte, 24)
		copy(k, key)
		return k, nil
	case 32:
		return key[:24], nil
	default:
		return nil, fmt.Errorf("3des: key must be 16, 24, or 32 bytes, got %d", len(key))
	}
}

// expandTwoKey168 turns a 16-byte 3DES two-key material into the 24-byte key expected by crypto/des.
func expandTwoKey168(key []byte) []byte {
	k := make([]byte, 24)
	copy(k[0:8], key[0:8])
	copy(k[8:16], key[8:16])
	copy(k[16:24], key[0:8])
	return k
}
