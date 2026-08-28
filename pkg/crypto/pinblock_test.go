package crypto

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestFormat0(t *testing.T) {
	pan := "5058270530000115"
	newPin := "7552"
	data, err := Format0(newPin, pan)
	if err != nil {
		t.Fatalf("Ошибка получения pinblock format0")
	}
	want := "0475D08FACFFFFEE"
	got := strings.ToUpper(hex.EncodeToString(data))
	if got != want {
		t.Errorf("GOT: %v; WANT: %s", got, want)
	}
}
