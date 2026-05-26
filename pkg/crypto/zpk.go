package crypto

import (
	"crypto/rand"
	"fmt"
)

// ZPKGenerator produces a 32-byte key with odd parity on each byte (full ZPK for RSA; first 24 bytes drive 3DES).
type ZPKGenerator struct{}

// NewZPKGenerator returns a ZPKGenerator.
func NewZPKGenerator() *ZPKGenerator {
	return &ZPKGenerator{}
}

// Generate returns 32 random bytes with odd parity applied to each byte.
func (g *ZPKGenerator) Generate() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("zpk: random: %w", err)
	}
	for i := range key {
		key[i] = ApplyOddParity(key[i])
	}
	return key, nil
}

// ApplyOddParity sets the least significant bit so the byte has an odd number of 1 bits.
// Rule from spec: count ones in the top 7 bits; if even, LSB=1; if odd, LSB=0.
func ApplyOddParity(b byte) byte {
	top7 := b >> 1
	var ones int
	for i := 0; i < 7; i++ {
		if (top7>>i)&1 == 1 {
			ones++
		}
	}
	lsb := byte(0)
	if ones%2 == 0 {
		lsb = 1
	}
	return (b & 0xFE) | lsb
}
