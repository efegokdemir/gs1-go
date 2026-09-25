package gs1

import (
	"errors"
	"testing"
)

func TestValidateGTIN(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// Valid GTINs (verified against GS1 spec)
		{name: "GTIN-14", input: "04150000021126", wantErr: nil},
		{name: "GTIN-13", input: "5901234123457", wantErr: nil},
		{name: "GTIN-12", input: "036000291452", wantErr: nil},
		{name: "GTIN-8", input: "96385074", wantErr: nil},
		{name: "GTIN-14 all zeros", input: "00000000000000", wantErr: nil},
		{name: "GTIN-13 example 2", input: "4006381333931", wantErr: nil},
		{name: "GTIN-14 example 2", input: "10614141000415", wantErr: nil},
		{name: "GTIN-13 example 3", input: "0734567890124", wantErr: nil},

		// Invalid check digit
		{name: "wrong check GTIN-14", input: "04150000021127", wantErr: ErrInvalidCheckDigit},
		{name: "wrong check GTIN-13", input: "5901234123458", wantErr: ErrInvalidCheckDigit},
		{name: "wrong check GTIN-8", input: "96385075", wantErr: ErrInvalidCheckDigit},

		// Invalid format
		{name: "too short", input: "1234567", wantErr: ErrInvalidData},
		{name: "length 9", input: "123456789", wantErr: ErrInvalidData},
		{name: "length 10", input: "1234567890", wantErr: ErrInvalidData},
		{name: "too long", input: "123456789012345", wantErr: ErrInvalidData},
		{name: "letters", input: "041500000211AB", wantErr: ErrInvalidData},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGTIN(tt.input)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("ValidateGTIN(%q) = %v, want nil", tt.input, err)
				}
				return
			}
			if err == nil {
				t.Errorf("ValidateGTIN(%q) = nil, want %v", tt.input, tt.wantErr)
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateGTIN(%q) = %v, want %v", tt.input, err, tt.wantErr)
			}
		})
	}
	if _, err := ComputeGTINCheckDigit("1234567890"); !errors.Is(err, ErrInvalidData) {
		t.Errorf("ComputeGTINCheckDigit() error = %v, want ErrInvalidData", err)
	}
}

func TestComputeGTINCheckDigit(t *testing.T) {
	tests := []struct {
		name    string
		partial string
		want    byte
	}{
		{name: "GTIN-14 partial", partial: "0415000002112", want: '6'},
		{name: "GTIN-13 partial", partial: "590123412345", want: '7'},
		{name: "GTIN-12 partial", partial: "03600029145", want: '2'},
		{name: "GTIN-8 partial", partial: "9638507", want: '4'},
		{name: "all zeros 13", partial: "0000000000000", want: '0'},
		{name: "GTIN-13 partial 2", partial: "073456789012", want: '4'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComputeGTINCheckDigit(tt.partial)
			if err != nil {
				t.Fatalf("ComputeGTINCheckDigit(%q) error = %v", tt.partial, err)
			}
			if got != tt.want {
				t.Errorf("ComputeGTINCheckDigit(%q) = '%c', want '%c'", tt.partial, got, tt.want)
			}
		})
	}
}

func TestIdentifierBuilders(t *testing.T) {
	if got, err := NewSSCC('1', "0614141", "123456789"); err != nil || got != "106141411234567897" {
		t.Fatalf("NewSSCC() reference vector = %q, %v; want %q", got, err, "106141411234567897")
	}

	sscc, err := NewSSCC('3', "7801234", "000000009")
	if err != nil {
		t.Fatalf("NewSSCC() error = %v", err)
	}
	if sscc != "378012340000000095" {
		t.Fatalf("NewSSCC() = %q, want %q", sscc, "378012340000000095")
	}
	if gotExt, gotPrefix, gotSerial, err := SplitSSCC(sscc, 7); err != nil || gotExt != '3' || gotPrefix != "7801234" || gotSerial != "000000009" {
		t.Fatalf("SplitSSCC() = %q, %q, %q, %v", gotExt, gotPrefix, gotSerial, err)
	}

	gtin13, err := NewGTIN13("7801234", "00005")
	if err != nil || gtin13 != "7801234000056" {
		t.Fatalf("NewGTIN13() = %q, %v", gtin13, err)
	}
	gtin14, err := NewGTIN14('1', gtin13)
	if err != nil || gtin14 != "17801234000053" {
		t.Fatalf("NewGTIN14() = %q, %v", gtin14, err)
	}
}

func TestIdentifierBuildersRejectInvalidInput(t *testing.T) {
	if _, err := NewSSCC('x', "7801234", "000000009"); !errors.Is(err, ErrInvalidData) {
		t.Errorf("NewSSCC() error = %v, want ErrInvalidData", err)
	}
	if _, err := NewSSCC('3', "7801234", "9"); !errors.Is(err, ErrInvalidData) {
		t.Errorf("NewSSCC() short components error = %v, want ErrInvalidData", err)
	}
	if _, err := NewGTIN14('1', "1234567890123"); !errors.Is(err, ErrInvalidCheckDigit) {
		t.Errorf("NewGTIN14() error = %v, want ErrInvalidCheckDigit", err)
	}
	if err := ValidateCheckDigit("378012340000000097"); !errors.Is(err, ErrInvalidCheckDigit) {
		t.Errorf("ValidateCheckDigit() error = %v, want ErrInvalidCheckDigit", err)
	}
	if _, _, _, err := SplitSSCC("378012340000000095", 3); !errors.Is(err, ErrInvalidData) {
		t.Errorf("SplitSSCC() error = %v, want ErrInvalidData", err)
	}
	if _, err := NewGTIN14('1', "04150000021126"); !errors.Is(err, ErrInvalidData) {
		t.Errorf("NewGTIN14() length error = %v, want ErrInvalidData", err)
	}
}

func TestExpandUPCE(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "d6 is 0", input: "123450", want: "012000003455"},
		{name: "d6 is 1", input: "123451", want: "012100003454"},
		{name: "d6 is 2", input: "123452", want: "012200003453"},
		{name: "d6 is 3", input: "123453", want: "012300000451"},
		{name: "d6 is 4", input: "123454", want: "012340000053"},
		{name: "d6 is 5", input: "012345", want: "001234000057"},
		{name: "d6 is 9", input: "012349", want: "001234000095"},
		{name: "8-digit with NS and check", input: "00123457", want: "001234000057"},
		{name: "letters", input: "01234A", wantErr: true},
		{name: "wrong length", input: "01234", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandUPCE(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ExpandUPCE(%q) = %q, want error", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ExpandUPCE(%q) error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ExpandUPCE(%q) = %q, want %q", tt.input, got, tt.want)
			}
			// Validate the expanded UPC-A
			if err := ValidateGTIN(got); err != nil {
				t.Errorf("expanded UPC-A %q has invalid check digit: %v", got, err)
			}
		})
	}
}

func BenchmarkValidateGTIN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ValidateGTIN("04150000021126")
	}
}
