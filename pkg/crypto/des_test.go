package crypto

import (
	"crypto/des"
	"testing"
)

func tripleDESDecryptECB(key, ciphertext []byte) ([]byte, error) {
	k24, err := resolveTripleDESKey(key)
	if err != nil {
		return nil, err
	}
	cipher, err := des.NewTripleDESCipher(k24)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(ciphertext))
	for off := 0; off < len(ciphertext); off += 8 {
		cipher.Decrypt(out[off:off+8], ciphertext[off:off+8])
	}
	return out, nil
}

func TestTripleDESEncrypter_EncryptECB_RoundTrip(t *testing.T) {
	e := NewTripleDESEncrypter()
	key := make([]byte, 16)
	for i := range key {
		key[i] = byte(i + 1)
		key[i] = ApplyOddParity(key[i])
	}
	pt := []byte{0x04, 0x12, 0x6D, 0xC6, 0x66, 0x66, 0x66, 0x66}
	ct, err := e.EncryptECB(key, pt)
	if err != nil {
		t.Fatal(err)
	}
	dec, err := tripleDESDecryptECB(key, ct)
	if err != nil {
		t.Fatal(err)
	}
	for i := range pt {
		if dec[i] != pt[i] {
			t.Fatalf("byte %d: got %02x want %02x", i, dec[i], pt[i])
		}
	}
}

func TestTripleDESEncrypter_InvalidKey(t *testing.T) {
	e := NewTripleDESEncrypter()
	if _, err := e.EncryptECB([]byte{1, 2, 3}, []byte{1, 2, 3, 4, 5, 6, 7, 8}); err == nil {
		t.Fatal("expected error")
	}
}

func TestTripleDESEncrypter_EncryptECB_32ByteKey_RoundTrip(t *testing.T) {
	e := NewTripleDESEncrypter()
	key := make([]byte, 32)
	for i := range key {
		key[i] = ApplyOddParity(byte(i + 1))
	}
	pt := []byte{0x04, 0x12, 0x6D, 0xC6, 0x66, 0x66, 0x66, 0x66}
	ct, err := e.EncryptECB(key, pt)
	if err != nil {
		t.Fatal(err)
	}
	dec, err := tripleDESDecryptECB(key, ct)
	if err != nil {
		t.Fatal(err)
	}
	for i := range pt {
		if dec[i] != pt[i] {
			t.Fatalf("byte %d: got %02x want %02x", i, dec[i], pt[i])
		}
	}
}
