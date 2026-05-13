package gs1

import "testing"

func TestCleanScannerInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// No-op cases
		{name: "clean input", input: "0104150000021126", want: "0104150000021126"},
		{name: "with FNC1", input: "2112345\x1D10LOT1", want: "2112345\x1D10LOT1"},

		// Trailing CRLF (most common scanner quirk)
		{name: "trailing LF", input: "0104150000021126\n", want: "0104150000021126"},
		{name: "trailing CRLF", input: "0104150000021126\r\n", want: "0104150000021126"},
		{name: "trailing CR", input: "0104150000021126\r", want: "0104150000021126"},

		// Leading/trailing whitespace
		{name: "leading spaces", input: "  0104150000021126", want: "0104150000021126"},
		{name: "trailing spaces", input: "0104150000021126  ", want: "0104150000021126"},
		{name: "both spaces", input: "  0104150000021126  ", want: "0104150000021126"},

		// Null bytes
		{name: "trailing null", input: "0104150000021126\x00", want: "0104150000021126"},
		{name: "leading null", input: "\x000104150000021126", want: "0104150000021126"},
		{name: "internal null", input: "01041500\x0000021126", want: "0104150000021126"},

		// BOM
		{name: "UTF-8 BOM", input: "\xEF\xBB\xBF0104150000021126", want: "0104150000021126"},
		{name: "BOM + trailing CRLF", input: "\xEF\xBB\xBF0104150000021126\r\n", want: "0104150000021126"},

		// CR/LF as FNC1
		{name: "CR as FNC1", input: "2112345\r10LOT1", want: "2112345\x1D10LOT1"},
		{name: "LF as FNC1", input: "2112345\n10LOT1", want: "2112345\x1D10LOT1"},
		{name: "CRLF as FNC1", input: "2112345\r\n10LOT1", want: "2112345\x1D10LOT1"},

		// Consecutive FNC1 collapse
		{name: "double FNC1", input: "2112345\x1D\x1D10LOT1", want: "2112345\x1D10LOT1"},
		{name: "triple FNC1", input: "2112345\x1D\x1D\x1D10LOT1", want: "2112345\x1D10LOT1"},

		// Mixed noise
		{name: "BOM + AIM + data + CRLF", input: "\xEF\xBB\xBF]d20104150000021126\r\n", want: "]d20104150000021126"},
		{name: "null + CR combo", input: "2112345\x00\r10LOT1", want: "2112345\x1D10LOT1"},

		// Empty / noise only
		{name: "empty", input: "", want: ""},
		{name: "only CRLF", input: "\r\n", want: ""},
		{name: "only BOM", input: "\xEF\xBB\xBF", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanScannerInput(tt.input)
			if got != tt.want {
				t.Errorf("cleanScannerInput(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func BenchmarkCleanScannerInput(b *testing.B) {
	// Typical scanner output with trailing CRLF
	input := "0104150000021126172506302112345ABC\x1D10LOT42X\r\n"
	for i := 0; i < b.N; i++ {
		_ = cleanScannerInput(input)
	}
}

func BenchmarkCleanScannerInputClean(b *testing.B) {
	// Already clean input — should hit fast path
	input := "0104150000021126172506302112345ABC\x1D10LOT42X"
	for i := 0; i < b.N; i++ {
		_ = cleanScannerInput(input)
	}
}
