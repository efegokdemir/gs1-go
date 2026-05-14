package gs1

import (
	"errors"
	"testing"
)

// helper to build a Barcode with given AIs (values don't matter for regulatory checks)
func testBarcode(ais ...string) Barcode {
	b := Barcode{
		Elements: make([]Element, len(ais)),
	}
	for i, ai := range ais {
		b.Elements[i] = Element{AI: ai, Value: "test"}
	}
	return b
}

func TestValidateANVISA(t *testing.T) {
	tests := []struct {
		name    string
		ais     []string
		wantErr error
	}{
		{name: "compliant", ais: []string{"01", "17", "10", "21"}, wantErr: nil},
		{name: "compliant with extras", ais: []string{"01", "17", "10", "21", "240", "713"}, wantErr: nil},
		{name: "missing GTIN", ais: []string{"17", "10", "21"}, wantErr: ErrMissingRequiredAI},
		{name: "missing expiry", ais: []string{"01", "10", "21"}, wantErr: ErrMissingRequiredAI},
		{name: "missing lot", ais: []string{"01", "17", "21"}, wantErr: ErrMissingRequiredAI},
		{name: "missing serial", ais: []string{"01", "17", "10"}, wantErr: ErrMissingRequiredAI},
		{name: "empty barcode", ais: []string{}, wantErr: ErrMissingRequiredAI},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := testBarcode(tt.ais...)
			err := b.ValidateANVISA()
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("ValidateANVISA() = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Errorf("ValidateANVISA() = nil, want %v", tt.wantErr)
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateANVISA() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateANMAT(t *testing.T) {
	tests := []struct {
		name    string
		ais     []string
		wantErr error
	}{
		{name: "compliant", ais: []string{"01", "17", "10", "21"}, wantErr: nil},
		{name: "missing serial", ais: []string{"01", "17", "10"}, wantErr: ErrMissingRequiredAI},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testBarcode(tt.ais...).ValidateANMAT()
			if tt.wantErr == nil && err != nil {
				t.Errorf("ValidateANMAT() = %v, want nil", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateANMAT() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSNFA(t *testing.T) {
	tests := []struct {
		name    string
		ais     []string
		wantErr error
	}{
		{name: "compliant", ais: []string{"01", "17", "10"}, wantErr: nil},
		{name: "compliant with serial", ais: []string{"01", "17", "10", "21"}, wantErr: nil},
		{name: "missing lot", ais: []string{"01", "17"}, wantErr: ErrMissingRequiredAI},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testBarcode(tt.ais...).ValidateSNFA()
			if tt.wantErr == nil && err != nil {
				t.Errorf("ValidateSNFA() = %v, want nil", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateSNFA() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCOFEPRIS(t *testing.T) {
	tests := []struct {
		name    string
		ais     []string
		wantErr error
	}{
		{name: "compliant", ais: []string{"01", "17", "10"}, wantErr: nil},
		{name: "missing GTIN", ais: []string{"17", "10"}, wantErr: ErrMissingRequiredAI},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testBarcode(tt.ais...).ValidateCOFEPRIS()
			if tt.wantErr == nil && err != nil {
				t.Errorf("ValidateCOFEPRIS() = %v, want nil", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateCOFEPRIS() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegulatorValidateGeneric(t *testing.T) {
	custom := Regulator{Name: "Custom", Country: "XX", RequiredAIs: []string{"01", "10"}}
	b := testBarcode("01", "10", "21")
	if err := custom.Validate(b); err != nil {
		t.Errorf("custom regulator validation failed: %v", err)
	}
}
