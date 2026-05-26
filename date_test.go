package gs1

import (
	"errors"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr error
	}{
		// Happy path
		{name: "normal date", input: "250630", want: time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)},
		{name: "jan 1 2000", input: "000101", want: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{name: "dec 31 2099", input: "991231", want: time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC)},
		{name: "feb 28 non-leap", input: "250228", want: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)},
		{name: "feb 29 leap year", input: "240229", want: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)},

		// Day=00 (last day of month)
		{name: "day 00 jan", input: "250100", want: time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)},
		{name: "day 00 feb non-leap", input: "250200", want: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)},
		{name: "day 00 feb leap", input: "240200", want: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)},
		{name: "day 00 apr", input: "250400", want: time.Date(2025, 4, 30, 0, 0, 0, 0, time.UTC)},
		{name: "day 00 dec", input: "251200", want: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)},

		// Invalid dates
		{name: "month 00", input: "250001", wantErr: ErrInvalidDate},
		{name: "month 13", input: "251301", wantErr: ErrInvalidDate},
		{name: "day 32", input: "250132", wantErr: ErrInvalidDate},
		{name: "feb 30", input: "250230", wantErr: ErrInvalidDate},
		{name: "feb 29 non-leap", input: "250229", wantErr: ErrInvalidDate},
		{name: "apr 31", input: "250431", wantErr: ErrInvalidDate},

		// Invalid format
		{name: "too short", input: "25063", wantErr: ErrInvalidDate},
		{name: "too long", input: "2506301", wantErr: ErrInvalidDate},
		{name: "letters", input: "25063A", wantErr: ErrInvalidDate},
		{name: "empty", input: "", wantErr: ErrInvalidDate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.input)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ParseDate(%q) = %v, want error %v", tt.input, got, tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ParseDate(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseDate(%q) unexpected error: %v", tt.input, err)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("ParseDate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseDateWithOptions_DayZeroFirstDay(t *testing.T) {
	opts := DateOptions{DayZero: DayZeroFirstDay}
	tests := []struct {
		name  string
		input string
		want  time.Time
	}{
		{name: "jan first", input: "250100", want: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
		{name: "feb first non-leap", input: "250200", want: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)},
		{name: "feb first leap", input: "240200", want: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
		{name: "apr first", input: "250400", want: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)},
		{name: "dec first", input: "251200", want: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDateWithOptions(tt.input, opts)
			if err != nil {
				t.Fatalf("ParseDateWithOptions(%q) unexpected error: %v", tt.input, err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("ParseDateWithOptions(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseDateWithOptions_DayZeroLastDay(t *testing.T) {
	opts := DateOptions{DayZero: DayZeroLastDay}
	got, err := ParseDateWithOptions("250200", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseDateWithOptions_NormalDateUnaffected(t *testing.T) {
	opts := DateOptions{DayZero: DayZeroFirstDay}
	got, err := ParseDateWithOptions("250630", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseDateRaw(t *testing.T) {
	got, err := ParseDateRaw("250630")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "250630" {
		t.Errorf("got %q, want %q", got, "250630")
	}

	// Day 00 is valid
	got, err = ParseDateRaw("250200")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "250200" {
		t.Errorf("got %q, want %q", got, "250200")
	}

	// Invalid input
	_, err = ParseDateRaw("999999")
	if err == nil {
		t.Error("expected error for invalid date")
	}
}
