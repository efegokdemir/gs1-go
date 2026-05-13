package gs1

import "testing"

func TestAITableIntegrity(t *testing.T) {
	for key, spec := range aiTable {
		if key != spec.AI {
			t.Errorf("map key %q does not match spec.AI %q", key, spec.AI)
		}
		if spec.Name == "" {
			t.Errorf("AI %q has empty Name", key)
		}
		if spec.FixedLen > 0 && spec.FixedLen != spec.MaxLen {
			t.Errorf("AI %q: fixed-length AI has FixedLen=%d != MaxLen=%d", key, spec.FixedLen, spec.MaxLen)
		}
		if spec.FixedLen == 0 && spec.MaxLen == 0 {
			t.Errorf("AI %q: variable-length AI has MaxLen=0", key)
		}
	}
}

func TestLookupAI(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		pos     int
		wantAI  string
		wantLen int
		wantOK  bool
	}{
		{name: "2-digit GTIN", input: "0104150000021125", pos: 0, wantAI: "01", wantLen: 2, wantOK: true},
		{name: "2-digit lot", input: "10ABC123", pos: 0, wantAI: "10", wantLen: 2, wantOK: true},
		{name: "2-digit expiry", input: "17250630", pos: 0, wantAI: "17", wantLen: 2, wantOK: true},
		{name: "3-digit additional", input: "240PROD001", pos: 0, wantAI: "240", wantLen: 3, wantOK: true},
		{name: "3-digit NHRN", input: "713BR12345", pos: 0, wantAI: "713", wantLen: 3, wantOK: true},
		{name: "4-digit weight", input: "3102123456", pos: 0, wantAI: "3102", wantLen: 4, wantOK: true},
		{name: "at offset", input: "0104150000021125172506301", pos: 16, wantAI: "17", wantLen: 2, wantOK: true},
		{name: "3-digit GLN", input: "4141234567890123", pos: 0, wantAI: "414", wantLen: 3, wantOK: true},
		{name: "3-digit GSIN", input: "40212345678901234567", pos: 0, wantAI: "402", wantLen: 3, wantOK: true},
		{name: "4-digit gross kg", input: "3302123456", pos: 0, wantAI: "3302", wantLen: 4, wantOK: true},
		{name: "4-digit gross lb", input: "3400123456", pos: 0, wantAI: "3400", wantLen: 4, wantOK: true},
		{name: "2-digit custom 90", input: "90CUSTOM", pos: 0, wantAI: "90", wantLen: 2, wantOK: true},
		{name: "2-digit custom 99", input: "99DATA", pos: 0, wantAI: "99", wantLen: 2, wantOK: true},
		{name: "unknown AI", input: "8812345", pos: 0, wantAI: "", wantLen: 0, wantOK: false},
		{name: "too short", input: "0", pos: 0, wantAI: "", wantLen: 0, wantOK: false},
		{name: "empty at pos", input: "01", pos: 2, wantAI: "", wantLen: 0, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, aiLen, ok := lookupAI(tt.input, tt.pos)
			if ok != tt.wantOK {
				t.Errorf("lookupAI() ok = %v, want %v", ok, tt.wantOK)
				return
			}
			if !ok {
				return
			}
			if spec.AI != tt.wantAI {
				t.Errorf("lookupAI() AI = %q, want %q", spec.AI, tt.wantAI)
			}
			if aiLen != tt.wantLen {
				t.Errorf("lookupAI() len = %d, want %d", aiLen, tt.wantLen)
			}
		})
	}
}
