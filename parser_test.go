package gs1

import (
	"errors"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   error
		wantCount int
		wantAIs   []string
	}{
		// Single AI — fixed length
		{
			name:      "GTIN only",
			input:     "0104150000021126",
			wantCount: 1,
			wantAIs:   []string{"01"},
		},
		// Single AI — variable length
		{
			name:      "lot only",
			input:     "10ABC123",
			wantCount: 1,
			wantAIs:   []string{"10"},
		},
		// Multiple fixed-length AIs (no FNC1 needed)
		{
			name:      "GTIN + expiry",
			input:     "0104150000021126" + "17250630",
			wantCount: 2,
			wantAIs:   []string{"01", "17"},
		},
		// Fixed + variable with FNC1
		{
			name:      "GTIN + expiry + lot",
			input:     "0104150000021126" + "17250630" + "10BATCH42",
			wantCount: 3,
			wantAIs:   []string{"01", "17", "10"},
		},
		// Variable + variable with FNC1 separator
		{
			name:      "serial + lot with FNC1",
			input:     "21SN12345\x1D10LOT999",
			wantCount: 2,
			wantAIs:   []string{"21", "10"},
		},
		// Full healthcare barcode (GTIN + expiry + serial + lot)
		{
			name:      "full healthcare barcode",
			input:     "0104150000021126172506302112345ABC\x1D10LOT42X",
			wantCount: 4,
			wantAIs:   []string{"01", "17", "21", "10"},
		},
		// AIM prefix DataMatrix
		{
			name:      "AIM prefix ]d2",
			input:     "]d20104150000021126172506301012345",
			wantCount: 3,
			wantAIs:   []string{"01", "17", "10"},
		},
		// AIM prefix GS1-128
		{
			name:      "AIM prefix ]C1",
			input:     "]C10104150000021126172506301012345",
			wantCount: 3,
			wantAIs:   []string{"01", "17", "10"},
		},
		// Bare GTIN (EAN-13, EAN-8, UPC-A — no AI prefix)
		{
			name:      "bare EAN-13",
			input:     "7800038041425",
			wantCount: 1,
			wantAIs:   []string{"01"},
		},
		{
			name:      "bare UPC-A",
			input:     "036000291452",
			wantCount: 1,
			wantAIs:   []string{"01"},
		},
		// Leading FNC1
		{
			name:      "leading FNC1",
			input:     "\x1D0104150000021126",
			wantCount: 1,
			wantAIs:   []string{"01"},
		},
		// 3-digit AI
		{
			name:      "AI 240 additional product",
			input:     "240PRODUCTCODE001",
			wantCount: 1,
			wantAIs:   []string{"240"},
		},
		// 4-digit AI (weight)
		{
			name:      "AI 3102 net weight kg",
			input:     "3102001500",
			wantCount: 1,
			wantAIs:   []string{"3102"},
		},
		// 3-digit AI NHRN Brazil
		{
			name:      "AI 713 NHRN Brazil",
			input:     "713BR12345678",
			wantCount: 1,
			wantAIs:   []string{"713"},
		},
		// SSCC
		{
			name:      "SSCC",
			input:     "00123456789012345675",
			wantCount: 1,
			wantAIs:   []string{"00"},
		},
		// GLN
		{
			name:      "GLN",
			input:     "4141234567890123",
			wantCount: 1,
			wantAIs:   []string{"414"},
		},
		// GSIN
		{
			name:      "GSIN",
			input:     "40212345678901234567",
			wantCount: 1,
			wantAIs:   []string{"402"},
		},
		// Gross weight
		{
			name:      "gross weight kg",
			input:     "3302001500",
			wantCount: 1,
			wantAIs:   []string{"3302"},
		},
		// Count AI
		{
			name:      "count",
			input:     "3025",
			wantCount: 1,
			wantAIs:   []string{"30"},
		},
		// Production date
		{
			name:      "production date",
			input:     "11250101",
			wantCount: 1,
			wantAIs:   []string{"11"},
		},
		// Trailing FNC1 (some scanners emit it)
		{
			name:      "trailing FNC1",
			input:     "10LOT1\x1D",
			wantCount: 1,
			wantAIs:   []string{"10"},
		},

		// Bracket notation
		{
			name:      "bracket notation simple",
			input:     "(01)04150000021126(17)250630",
			wantCount: 2,
			wantAIs:   []string{"01", "17"},
		},
		{
			name:      "bracket notation full",
			input:     "(02)17795678901213(10)AAB123(17)201231(37)20",
			wantCount: 4,
			wantAIs:   []string{"02", "10", "17", "37"},
		},

		// Scanner resilience
		{
			name:      "trailing CRLF",
			input:     "0104150000021126172506301012345\r\n",
			wantCount: 3,
			wantAIs:   []string{"01", "17", "10"},
		},
		{
			name:      "trailing LF",
			input:     "0104150000021126\n",
			wantCount: 1,
			wantAIs:   []string{"01"},
		},
		{
			name:      "BOM prefix",
			input:     "\xEF\xBB\xBF0104150000021126",
			wantCount: 1,
			wantAIs:   []string{"01"},
		},
		{
			name:      "CR as FNC1",
			input:     "2112345\r10LOT1",
			wantCount: 2,
			wantAIs:   []string{"21", "10"},
		},
		{
			name:      "null bytes stripped",
			input:     "01041500000211261725063010\x0012345",
			wantCount: 3,
			wantAIs:   []string{"01", "17", "10"},
		},
		{
			name:      "AIM prefix ]d1",
			input:     "]d10104150000021126",
			wantCount: 1,
			wantAIs:   []string{"01"},
		},
		{
			name:      "AIM prefix ]Q3",
			input:     "]Q30104150000021126",
			wantCount: 1,
			wantAIs:   []string{"01"},
		},
		{
			name:      "BOM + AIM + data + CRLF",
			input:     "\xEF\xBB\xBF]d20104150000021126172506301012345\r\n",
			wantCount: 3,
			wantAIs:   []string{"01", "17", "10"},
		},
		{
			name:      "double FNC1 between fields",
			input:     "2112345\x1D\x1D10LOT1",
			wantCount: 2,
			wantAIs:   []string{"21", "10"},
		},

		// Error cases
		{
			name:    "empty string",
			input:   "",
			wantErr: ErrEmptyInput,
		},
		{
			name:    "only spaces",
			input:   "   ",
			wantErr: ErrEmptyInput,
		},
		{
			name:    "unknown AI",
			input:   "8812345",
			wantErr: ErrUnknownAI,
		},
		// Custom AI 90-99
		{
			name:      "custom AI 90",
			input:     "90CUSTOMDATA\x1D",
			wantCount: 1,
			wantAIs:   []string{"90"},
		},
		{
			name:    "truncated GTIN",
			input:   "010415000002",
			wantErr: ErrTruncatedData,
		},
		{
			name:    "non-numeric in numeric field",
			input:   "01ABCDEFGHIJKLMN",
			wantErr: ErrInvalidData,
		},
		{
			name:    "variable data exceeds max",
			input:   "10AAAAABBBBBCCCCCDDDDDE",
			wantErr: ErrInvalidData,
		},
		{
			name:    "only FNC1",
			input:   "\x1D",
			wantErr: ErrEmptyInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := Parse(tt.input)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Parse(%q) = %v, want error %v", tt.input, b.Elements, tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Parse(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			if len(b.Elements) != tt.wantCount {
				t.Errorf("Parse(%q) got %d elements, want %d", tt.input, len(b.Elements), tt.wantCount)
			}
			for i, wantAI := range tt.wantAIs {
				if i >= len(b.Elements) {
					break
				}
				if b.Elements[i].AI != wantAI {
					t.Errorf("Parse(%q) element[%d].AI = %q, want %q", tt.input, i, b.Elements[i].AI, wantAI)
				}
			}
		})
	}
}

func TestParseConvenienceMethods(t *testing.T) {
	input := "0104150000021126" + "17250630" + "2112345ABC\x1D" + "10LOT42X"
	b, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if got := b.GTIN(); got != "04150000021126" {
		t.Errorf("GTIN() = %q, want %q", got, "04150000021126")
	}
	if got := b.Lot(); got != "LOT42X" {
		t.Errorf("Lot() = %q, want %q", got, "LOT42X")
	}
	if got := b.SerialNumber(); got != "12345ABC" {
		t.Errorf("SerialNumber() = %q, want %q", got, "12345ABC")
	}

	expiry, err := b.ExpirationDate()
	if err != nil {
		t.Fatalf("ExpirationDate() error = %v", err)
	}
	if !expiry.Equal(time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("ExpirationDate() = %v, want 2025-06-30", expiry)
	}

	v, ok := b.Get("01")
	if !ok || v != "04150000021126" {
		t.Errorf("Get(01) = (%q, %v), want (%q, true)", v, ok, "04150000021126")
	}

	_, ok = b.Get("00")
	if ok {
		t.Error("Get(00) should return false for missing AI")
	}
}

func TestParseConvenienceMethodsMissing(t *testing.T) {
	b, err := Parse("0104150000021126")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if got := b.Lot(); got != "" {
		t.Errorf("Lot() = %q, want empty", got)
	}
	if got := b.SerialNumber(); got != "" {
		t.Errorf("SerialNumber() = %q, want empty", got)
	}

	_, err = b.ExpirationDate()
	if err == nil {
		t.Error("ExpirationDate() should error when AI 17 is missing")
	}
}

func BenchmarkParse(b *testing.B) {
	input := "0104150000021126172506302112345ABC\x1D10LOT42X"
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input)
	}
}

func BenchmarkParseMinimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Parse("0104150000021126")
	}
}
