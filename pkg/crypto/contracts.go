package crypto

import "converterapi/pkg/logger"

var RandomZPK []byte

// Format0PINBlockBuilder builds an ISO 9564 Format 0 cleartext PIN block.
type Format0PINBlockBuilder interface {
	BuildFormat0(pan, pin string) ([]byte, error)
}

// RandomZPKSource generates a 32-byte ZPK with odd parity per byte.
type RandomZPKSource interface {
	Generate() ([]byte, error)
}

// PINBlockEncrypter performs 3DES ECB on 8-byte PIN blocks.
type PINBlockEncrypter interface {
	EncryptECB(key, plaintext []byte) ([]byte, error)
}

// ZPKRSAEncrypter encrypts the ZPK with RSA PKCS1v15.
type ZPKRSAEncrypter interface {
	Encrypt(plaintext []byte) ([]byte, error)
}

func Init() {
	RandomZPK, err := NewZPKGenerator().Generate()
	if err != nil {
		logger.Errorf("Error generating zpk: %v", err.Error())
		return
	}
	RandomZPK = RandomZPK
}
