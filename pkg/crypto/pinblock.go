package crypto

import (
	"errors"
	"fmt"
	"unicode"
)

// PinBlockBuilder builds ISO 9564-1 Format 0 (ANSI X9.8) PIN blocks.
type PinBlockBuilder struct{}

// NewPinBlockBuilder returns a PinBlockBuilder.
func NewPinBlockBuilder() *PinBlockBuilder {
	return &PinBlockBuilder{}
}

// BuildFormat0 returns the 8-byte cleartext PIN block before 3DES.
// Algorithm: XOR of PIN field P1 and PAN field P2 (both 8 bytes).
func (b *PinBlockBuilder) BuildFormat0(pan, pin string) ([]byte, error) {
	if err := validatePIN(pin); err != nil {
		return nil, err
	}
	panDigits, err := digitsOnly(pan)
	if err != nil {
		return nil, err
	}
	if len(panDigits) < 13 || len(panDigits) > 19 {
		return nil, fmt.Errorf("pan: expected 13–19 digits, got %d", len(panDigits))
	}

	p1, err := buildP1(pin)
	if err != nil {
		return nil, err
	}
	p2, err := buildP2(panDigits)
	if err != nil {
		return nil, err
	}

	out := make([]byte, 8)
	for i := range out {
		out[i] = p1[i] ^ p2[i]
	}
	return out, nil
}

func validatePIN(pin string) error {
	n := len(pin)
	if n < 4 || n > 12 {
		return fmt.Errorf("pin: length must be 4–12, got %d", n)
	}
	for _, r := range pin {
		if r < '0' || r > '9' {
			return errors.New("pin: must contain only digits")
		}
	}
	return nil
}

func digitsOnly(s string) (string, error) {
	var b []byte
	for _, r := range s {
		if unicode.IsDigit(r) {
			b = append(b, byte(r))
		} else {
			return "", errors.New("pan: must contain only digits")
		}
	}
	return string(b), nil
}

// buildP1: nibble0=0 (format 0), nibble1=PIN length, then PIN in BCD padded with F; remaining bytes 0xFF.
func buildP1(pin string) ([]byte, error) {
	pinLen := len(pin)
	block := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	block[0] = byte((0 << 4) | pinLen)

	nibbles := make([]byte, 0, 24)
	for i := 0; i < pinLen; i++ {
		nibbles = append(nibbles, pin[i]-'0')
	}
	if len(nibbles)%2 == 1 {
		nibbles = append(nibbles, 0x0F)
	}
	idx := 1
	for i := 0; i < len(nibbles); i += 2 {
		if idx >= 8 {
			return nil, errors.New("pin: internal encoding overflow")
		}
		block[idx] = (nibbles[i] << 4) | nibbles[i+1]
		idx++
	}
	return block, nil
}

// buildP2: bytes 0-1 are 0x00; bytes 2-7 are 12 PAN digits (excluding check digit) as BCD.
func buildP2(panDigits string) ([]byte, error) {
	// Exclude check digit (last digit); take 12 rightmost of the remainder.
	body := panDigits[:len(panDigits)-1]
	if len(body) < 12 {
		return nil, fmt.Errorf("pan: need at least 12 digits excluding check digit, got %d", len(body))
	}
	start := len(body) - 12
	d12 := body[start:]

	block := []byte{0x00, 0x00, 0, 0, 0, 0, 0, 0}
	for i := 0; i < 6; i++ {
		high := d12[i*2] - '0'
		low := d12[i*2+1] - '0'
		block[2+i] = (high << 4) | low
	}
	return block, nil
}
