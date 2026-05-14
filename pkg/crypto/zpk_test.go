package crypto

import (
	"testing"
)

func TestApplyOddParity(t *testing.T) {
	// top7 = 0 -> even ones -> LSB 1
	if got := ApplyOddParity(0x00); got != 0x01 {
		t.Fatalf("0x00 -> 0x01, got %02x", got)
	}
	// top7 all ones (0xFE>>1 = 0x7F) -> 7 ones odd -> LSB 0
	if got := ApplyOddParity(0xFE); got&1 != 0 {
		t.Fatalf("expected even parity byte, got %02x", got)
	}
	for b := 0; b < 256; b++ {
		x := ApplyOddParity(byte(b))
		var bits int
		for i := 0; i < 8; i++ {
			if (x>>i)&1 == 1 {
				bits++
			}
		}
		if bits%2 == 0 {
			t.Fatalf("byte %02x adjusted to %02x does not have odd parity", b, x)
		}
	}
}

func TestZPKGenerator_Generate(t *testing.T) {
	g := NewZPKGenerator()
	key, err := g.Generate()
	if err != nil {
		t.Fatal(err)
	}
	if len(key) != 32 {
		t.Fatalf("len %d", len(key))
	}
}
