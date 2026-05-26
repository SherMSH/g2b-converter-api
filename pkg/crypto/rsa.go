package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

// RSAZPKEncrypter encrypts the ZPK with PKCS#1 v1.5 using a public key.
type RSAZPKEncrypter struct {
	pub *rsa.PublicKey
}

// NewRSAZPKEncrypter returns an RSAZPKEncrypter.
func NewRSAZPKEncrypter(pub *rsa.PublicKey) *RSAZPKEncrypter {
	return &RSAZPKEncrypter{pub: pub}
}

// Encrypt encrypts plaintext (ZPK, typically 32 bytes) with RSA PKCS1v15.
func (e *RSAZPKEncrypter) Encrypt(plaintext []byte) ([]byte, error) {
	if e.pub == nil {
		return nil, fmt.Errorf("rsa: nil public key")
	}
	if len(plaintext) == 0 {
		return nil, fmt.Errorf("rsa: empty plaintext")
	}
	ct, err := rsa.EncryptPKCS1v15(rand.Reader, e.pub, plaintext)
	if err != nil {
		return nil, fmt.Errorf("rsa: encrypt: %w", err)
	}
	return ct, nil
}
