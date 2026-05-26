package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestPinBlockBuilder_BuildFormat0_Golden(t *testing.T) {
	b := NewPinBlockBuilder()
	got, err := b.BuildFormat0("5375939999999993", "1234")
	if err != nil {
		t.Fatal(err)
	}
	want, _ := hex.DecodeString("04126DC666666666")
	if !bytes.Equal(got, want) {
		t.Fatalf("pin block mismatch\n got: %s\nwant: %s", hex.EncodeToString(got), hex.EncodeToString(want))
	}
}

func TestPinBlockBuilder_BuildFormat0_OddLengthPIN(t *testing.T) {
	b := NewPinBlockBuilder()
	got, err := b.BuildFormat0("5375939999999993", "12345")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 8 {
		t.Fatalf("len %d", len(got))
	}
}

func TestPinBlockBuilder_InvalidPIN(t *testing.T) {
	b := NewPinBlockBuilder()
	if _, err := b.BuildFormat0("5375939999999993", "123"); err == nil {
		t.Fatal("expected error for short PIN")
	}
	if _, err := b.BuildFormat0("5375939999999993", "12a4"); err == nil {
		t.Fatal("expected error for non-digit")
	}
}

func TestPinBlockBuilder_ShortPAN(t *testing.T) {
	b := NewPinBlockBuilder()
	if _, err := b.BuildFormat0("123456789012", "1234"); err == nil {
		t.Fatal("expected error for pan shorter than 13 digits")
	}
}
