package gs1

import (
	"errors"
	"testing"
)

func TestEncode(t *testing.T) {
	elements := []Element{
		{AI: "10", Value: "LOT42"},
		{AI: "01", Value: "04150000021126"},
		{AI: "17", Value: "250630"},
		{AI: "21", Value: "SERIAL"},
	}
	want := "\x1D01041500000211261725063010LOT42\x1D21SERIAL"
	got, err := Encode(elements)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if got != want {
		t.Fatalf("Encode() = %q, want %q", got, want)
	}
	b, err := Parse(got)
	if err != nil {
		t.Fatalf("Parse(Encode()) error = %v", err)
	}
	if got := b.HRI(); got != "(01)04150000021126(17)250630(10)LOT42(21)SERIAL" {
		t.Errorf("HRI() = %q", got)
	}
	if got := b.String(); got != b.HRI() {
		t.Errorf("String() = %q, want HRI %q", got, b.HRI())
	}
}

func TestEncodeErrors(t *testing.T) {
	tests := []struct {
		name     string
		elements []Element
		wantErr  error
	}{
		{name: "empty", wantErr: ErrEmptyInput},
		{name: "unknown AI", elements: []Element{{AI: "999", Value: "x"}}, wantErr: ErrUnknownAI},
		{name: "wrong fixed length", elements: []Element{{AI: "17", Value: "2506"}}, wantErr: ErrInvalidData},
		{name: "non-numeric", elements: []Element{{AI: "01", Value: "0415000002112A"}}, wantErr: ErrInvalidData},
		{name: "FNC1 in value", elements: []Element{{AI: "10", Value: "LOT\x1D42"}}, wantErr: ErrInvalidData},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Encode(tt.elements)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Encode() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncodeSeparatesFixedLengthNonPredefinedAI(t *testing.T) {
	got, err := Encode([]Element{
		{AI: "402", Value: "12345678901234567"},
		{AI: "10", Value: "LOT1"},
	})
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	want := "\x1D40212345678901234567\x1D10LOT1"
	if got != want {
		t.Errorf("Encode() = %q, want %q", got, want)
	}
}
