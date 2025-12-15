package config

import (
	"testing"
)

// TestGenerateBoardPonOID verifies that dynamic OID generation matches expected values
// These test cases validate against the original hardcoded cfg.yaml values
func TestGenerateBoardPonOID(t *testing.T) {
	tests := []struct {
		name    string
		boardID int
		ponID   int
		want    *BoardPonConfig
		wantErr bool
	}{
		{
			name:    "Board 1 PON 1",
			boardID: 1,
			ponID:   1,
			want: &BoardPonConfig{
				OnuIDNameOID:              ".500.10.2.3.3.1.2.285278465",
				OnuTypeOID:                ".3.50.11.2.1.17.268501248",
				OnuSerialNumberOID:        ".500.10.2.3.3.1.18.285278465",
				OnuRxPowerOID:             ".500.20.2.2.2.1.10.285278465",
				OnuTxPowerOID:             ".3.50.12.1.1.14.268501248",
				OnuStatusOID:              ".500.10.2.3.8.1.4.285278465",
				OnuIPAddressOID:           ".3.50.16.1.1.10.268501248",
				OnuDescriptionOID:         ".500.10.2.3.3.1.3.285278465",
				OnuLastOnlineOID:          ".500.10.2.3.8.1.5.285278465",
				OnuLastOfflineOID:         ".500.10.2.3.8.1.6.285278465",
				OnuLastOfflineReasonOID:   ".500.10.2.3.8.1.7.285278465",
				OnuGponOpticalDistanceOID: ".500.10.2.3.10.1.2.285278465",
			},
			wantErr: false,
		},
		{
			name:    "Board 1 PON 2",
			boardID: 1,
			ponID:   2,
			want: &BoardPonConfig{
				OnuIDNameOID:              ".500.10.2.3.3.1.2.285278466",
				OnuTypeOID:                ".3.50.11.2.1.17.268501504",
				OnuSerialNumberOID:        ".500.10.2.3.3.1.18.285278466",
				OnuRxPowerOID:             ".500.20.2.2.2.1.10.285278466",
				OnuTxPowerOID:             ".3.50.12.1.1.14.268501504",
				OnuStatusOID:              ".500.10.2.3.8.1.4.285278466",
				OnuIPAddressOID:           ".3.50.16.1.1.10.268501504",
				OnuDescriptionOID:         ".500.10.2.3.3.1.3.285278466",
				OnuLastOnlineOID:          ".500.10.2.3.8.1.5.285278466",
				OnuLastOfflineOID:         ".500.10.2.3.8.1.6.285278466",
				OnuLastOfflineReasonOID:   ".500.10.2.3.8.1.7.285278466",
				OnuGponOpticalDistanceOID: ".500.10.2.3.10.1.2.285278466",
			},
			wantErr: false,
		},
		{
			name:    "Board 1 PON 16",
			boardID: 1,
			ponID:   16,
			want: &BoardPonConfig{
				OnuIDNameOID:              ".500.10.2.3.3.1.2.285278480",
				OnuTypeOID:                ".3.50.11.2.1.17.268505088",
				OnuSerialNumberOID:        ".500.10.2.3.3.1.18.285278480",
				OnuRxPowerOID:             ".500.20.2.2.2.1.10.285278480",
				OnuTxPowerOID:             ".3.50.12.1.1.14.268505088",
				OnuStatusOID:              ".500.10.2.3.8.1.4.285278480",
				OnuIPAddressOID:           ".3.50.16.1.1.10.268505088",
				OnuDescriptionOID:         ".500.10.2.3.3.1.3.285278480",
				OnuLastOnlineOID:          ".500.10.2.3.8.1.5.285278480",
				OnuLastOfflineOID:         ".500.10.2.3.8.1.6.285278480",
				OnuLastOfflineReasonOID:   ".500.10.2.3.8.1.7.285278480",
				OnuGponOpticalDistanceOID: ".500.10.2.3.10.1.2.285278480",
			},
			wantErr: false,
		},
		{
			name:    "Board 2 PON 1",
			boardID: 2,
			ponID:   1,
			want: &BoardPonConfig{
				OnuIDNameOID:              ".500.10.2.3.3.1.2.285278721",
				OnuTypeOID:                ".3.50.11.2.1.17.268566784",
				OnuSerialNumberOID:        ".500.10.2.3.3.1.18.285278721",
				OnuRxPowerOID:             ".500.20.2.2.2.1.10.285278721",
				OnuTxPowerOID:             ".3.50.12.1.1.14.268566784",
				OnuStatusOID:              ".500.10.2.3.8.1.4.285278721",
				OnuIPAddressOID:           ".3.50.16.1.1.10.268566784",
				OnuDescriptionOID:         ".500.10.2.3.3.1.3.285278721",
				OnuLastOnlineOID:          ".500.10.2.3.8.1.5.285278721",
				OnuLastOfflineOID:         ".500.10.2.3.8.1.6.285278721",
				OnuLastOfflineReasonOID:   ".500.10.2.3.8.1.7.285278721",
				OnuGponOpticalDistanceOID: ".500.10.2.3.10.1.2.285278721",
			},
			wantErr: false,
		},
		{
			name:    "Board 2 PON 16",
			boardID: 2,
			ponID:   16,
			want: &BoardPonConfig{
				OnuIDNameOID:              ".500.10.2.3.3.1.2.285278736",
				OnuTypeOID:                ".3.50.11.2.1.17.268570624",
				OnuSerialNumberOID:        ".500.10.2.3.3.1.18.285278736",
				OnuRxPowerOID:             ".500.20.2.2.2.1.10.285278736",
				OnuTxPowerOID:             ".3.50.12.1.1.14.268570624",
				OnuStatusOID:              ".500.10.2.3.8.1.4.285278736",
				OnuIPAddressOID:           ".3.50.16.1.1.10.268570624",
				OnuDescriptionOID:         ".500.10.2.3.3.1.3.285278736",
				OnuLastOnlineOID:          ".500.10.2.3.8.1.5.285278736",
				OnuLastOfflineOID:         ".500.10.2.3.8.1.6.285278736",
				OnuLastOfflineReasonOID:   ".500.10.2.3.8.1.7.285278736",
				OnuGponOpticalDistanceOID: ".500.10.2.3.10.1.2.285278736",
			},
			wantErr: false,
		},
		{
			name:    "Invalid board ID (0)",
			boardID: 0,
			ponID:   1,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid board ID (3)",
			boardID: 3,
			ponID:   1,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid PON ID (0)",
			boardID: 1,
			ponID:   0,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid PON ID (17)",
			boardID: 1,
			ponID:   17,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateBoardPonOID(tt.boardID, tt.ponID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateBoardPonOID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil {
				// Verify all OID fields match expected values
				if got.OnuIDNameOID != tt.want.OnuIDNameOID {
					t.Errorf("OnuIDNameOID = %v, want %v", got.OnuIDNameOID, tt.want.OnuIDNameOID)
				}
				if got.OnuTypeOID != tt.want.OnuTypeOID {
					t.Errorf("OnuTypeOID = %v, want %v", got.OnuTypeOID, tt.want.OnuTypeOID)
				}
				if got.OnuSerialNumberOID != tt.want.OnuSerialNumberOID {
					t.Errorf("OnuSerialNumberOID = %v, want %v", got.OnuSerialNumberOID, tt.want.OnuSerialNumberOID)
				}
				if got.OnuRxPowerOID != tt.want.OnuRxPowerOID {
					t.Errorf("OnuRxPowerOID = %v, want %v", got.OnuRxPowerOID, tt.want.OnuRxPowerOID)
				}
				if got.OnuTxPowerOID != tt.want.OnuTxPowerOID {
					t.Errorf("OnuTxPowerOID = %v, want %v", got.OnuTxPowerOID, tt.want.OnuTxPowerOID)
				}
			}
		})
	}
}

// TestInitializeBoardPonMap verifies that all 32 Board-PON combinations are generated
func TestInitializeBoardPonMap(t *testing.T) {
	boardPonMap, err := InitializeBoardPonMap()
	if err != nil {
		t.Fatalf("InitializeBoardPonMap() error = %v", err)
	}

	// Verify map has exactly 32 entries (2 boards * 16 PONs)
	expectedCount := 32
	if len(boardPonMap) != expectedCount {
		t.Errorf("InitializeBoardPonMap() generated %d entries, want %d", len(boardPonMap), expectedCount)
	}

	// Verify all expected keys exist
	for boardID := 1; boardID <= 2; boardID++ {
		for ponID := 1; ponID <= 16; ponID++ {
			key := BoardPonKey{BoardID: boardID, PonID: ponID}
			if _, ok := boardPonMap[key]; !ok {
				t.Errorf("Missing config for Board%dPon%d", boardID, ponID)
			}
		}
	}

	// Spot check a few values
	board1pon1 := boardPonMap[BoardPonKey{BoardID: 1, PonID: 1}]
	if board1pon1.OnuIDNameOID != ".500.10.2.3.3.1.2.285278465" {
		t.Errorf("Board1Pon1 OnuIDNameOID = %v, want .500.10.2.3.3.1.2.285278465", board1pon1.OnuIDNameOID)
	}

	board2pon16 := boardPonMap[BoardPonKey{BoardID: 2, PonID: 16}]
	if board2pon16.OnuTypeOID != ".3.50.11.2.1.17.268570624" {
		t.Errorf("Board2Pon16 OnuTypeOID = %v, want .3.50.11.2.1.17.268570624", board2pon16.OnuTypeOID)
	}
}

// BenchmarkGenerateBoardPonOID measures the performance of OID generation
func BenchmarkGenerateBoardPonOID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateBoardPonOID(1, 1)
	}
}

// BenchmarkInitializeBoardPonMap measures the performance of generating all 32 configs
func BenchmarkInitializeBoardPonMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = InitializeBoardPonMap()
	}
}
