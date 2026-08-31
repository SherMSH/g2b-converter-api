package service

import (
	"converterapi/internal/config"
	"testing"
)

func TestGetExpDateByPan(t *testing.T) {
	config.Config.App.DebugMode = true
	config.Config.Processing.Address = "http://d8-tprocweb1.humo.lab"

	pan := "5058270530000172"
	wantedExpDate := "3607"
	if config.Config.App.DebugMode {
		pan = "5058270530000016"
		wantedExpDate = "3004"
	}

	got, err := GetExpDateByPan(pan)
	if err != nil {
		t.Fatalf("Error testing GetCardExpDateByPan: %v", err)
	}

	if got != wantedExpDate {
		t.Errorf("GOT: %v WANT: %v", got, wantedExpDate)
	}
}
