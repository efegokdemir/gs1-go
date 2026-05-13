package gs1

import (
	"fmt"
	"strings"
	"time"
)

const fnc1 = '\x1D' // GS (Group Separator), used as FNC1 in scanner output

// Element represents a single AI-value pair extracted from a GS1 barcode.
type Element struct {
	AI    string // application identifier code, e.g., "01"
	Value string // raw data value
}

// Barcode represents a fully parsed GS1 barcode (GS1-128 or DataMatrix).
type Barcode struct {
	Raw      string    // original input string
	Elements []Element // parsed AI-value pairs in scan order
	index    map[string]int
}

// GTIN returns the GTIN value (AI 01), or "" if not present.
func (b Barcode) GTIN() string {
	v, _ := b.Get("01")
	return v
}

// Lot returns the batch/lot number (AI 10), or "" if not present.
func (b Barcode) Lot() string {
	v, _ := b.Get("10")
	return v
}

// SerialNumber returns the serial number (AI 21), or "" if not present.
func (b Barcode) SerialNumber() string {
	v, _ := b.Get("21")
	return v
}

// SSCC returns the SSCC value (AI 00), or "" if not present.
func (b Barcode) SSCC() string {
	v, _ := b.Get("00")
	return v
}

// Count returns the item count (AI 30), or "" if not present.
func (b Barcode) Count() string {
	v, _ := b.Get("30")
	return v
}

// ExpirationDate returns the parsed expiration date (AI 17).
func (b Barcode) ExpirationDate() (time.Time, error) {
	v, ok := b.Get("17")
	if !ok {
		return time.Time{}, fmt.Errorf("%w: AI (17) not present", ErrInvalidData)
	}
	return ParseDate(v)
}

// ProductionDate returns the parsed production date (AI 11).
func (b Barcode) ProductionDate() (time.Time, error) {
	v, ok := b.Get("11")
	if !ok {
		return time.Time{}, fmt.Errorf("%w: AI (11) not present", ErrInvalidData)
	}
	return ParseDate(v)
}

// BestBeforeDate returns the parsed best-before date (AI 15).
func (b Barcode) BestBeforeDate() (time.Time, error) {
	v, ok := b.Get("15")
	if !ok {
		return time.Time{}, fmt.Errorf("%w: AI (15) not present", ErrInvalidData)
	}
	return ParseDate(v)
}

// Get returns the value for the given AI code and whether it was found.
// If the AI appears multiple times, the first occurrence is returned.
func (b Barcode) Get(ai string) (string, bool) {
	idx, ok := b.index[ai]
	if !ok {
		return "", false
	}
	return b.Elements[idx].Value, true
}

// Parse parses a GS1 barcode string (GS1-128 or DataMatrix scanner output)
// into a Barcode with typed elements. It handles AIM symbology identifiers
// (e.g., ]C1, ]d2), FNC1 separators (ASCII 29), and bracket notation
// (e.g., "(01)04150000021126(17)250630").
func Parse(input string) (Barcode, error) {
	if strings.TrimSpace(input) == "" {
		return Barcode{}, ErrEmptyInput
	}

	// Clean scanner noise, then convert bracket notation.
	data := cleanScannerInput(input)
	data = stripBracketNotation(data)

	b := Barcode{
		Raw:      input,
		Elements: make([]Element, 0, 8),
		index:    make(map[string]int, 8),
	}

	pos := skipPrefix(data)

	for pos < len(data) {
		if data[pos] == byte(fnc1) {
			pos++
			continue
		}

		spec, aiLen, ok := lookupAI(data, pos)
		if !ok {
			return Barcode{}, fmt.Errorf("%w: at position %d", ErrUnknownAI, pos)
		}
		pos += aiLen

		value, newPos, err := extractData(data, pos, spec)
		if err != nil {
			return Barcode{}, err
		}
		pos = newPos

		if err := validateData(value, spec); err != nil {
			return Barcode{}, err
		}

		b.Elements = append(b.Elements, Element{AI: spec.AI, Value: value})
		if _, exists := b.index[spec.AI]; !exists {
			b.index[spec.AI] = len(b.Elements) - 1
		}
	}

	if len(b.Elements) == 0 {
		return Barcode{}, ErrEmptyInput
	}

	return b, nil
}

// skipPrefix skips leading FNC1 and AIM symbology identifiers.
func skipPrefix(data string) int {
	pos := 0
	if pos < len(data) && data[pos] == byte(fnc1) {
		pos++
	}
	if pos < len(data) && data[pos] == ']' && pos+3 <= len(data) {
		sym := data[pos+1]
		// ]C1 = GS1-128, ]d1/]d2 = DataMatrix, ]e0 = GS1 composite,
		// ]Q3 = GS1 QR, ]J1 = GS1 DotCode
		if sym == 'C' || sym == 'd' || sym == 'e' || sym == 'Q' || sym == 'J' {
			pos += 3
		}
	}
	return pos
}

// extractData reads the data field for an AI starting at pos.
func extractData(data string, pos int, spec aiSpec) (string, int, error) {
	if spec.FixedLen > 0 {
		if pos+spec.FixedLen > len(data) {
			return "", pos, fmt.Errorf("%w: AI (%s) needs %d chars, got %d",
				ErrTruncatedData, spec.AI, spec.FixedLen, len(data)-pos)
		}
		return data[pos : pos+spec.FixedLen], pos + spec.FixedLen, nil
	}

	end := pos
	for end < len(data) && data[end] != byte(fnc1) {
		end++
	}
	value := data[pos:end]
	if len(value) > spec.MaxLen {
		return "", pos, fmt.Errorf("%w: AI (%s) data length %d exceeds max %d",
			ErrInvalidData, spec.AI, len(value), spec.MaxLen)
	}
	if len(value) == 0 {
		return "", pos, fmt.Errorf("%w: AI (%s) has empty data", ErrInvalidData, spec.AI)
	}
	newPos := end
	if newPos < len(data) && data[newPos] == byte(fnc1) {
		newPos++
	}
	return value, newPos, nil
}

// validateData checks that the value conforms to the AI's data type.
func validateData(value string, spec aiSpec) error {
	if spec.DataType != dataNumeric {
		return nil
	}
	for i := 0; i < len(value); i++ {
		if value[i] < '0' || value[i] > '9' {
			return fmt.Errorf("%w: AI (%s) expects numeric data, got %q",
				ErrInvalidData, spec.AI, value)
		}
	}
	return nil
}

// stripBracketNotation converts bracket notation "(01)0415...(17)250630"
// to raw AI string with FNC1 separators between fields. If no brackets are
// found, returns the input unchanged.
func stripBracketNotation(input string) string {
	if !strings.Contains(input, "(") {
		return input
	}
	var b strings.Builder
	b.Grow(len(input))
	first := true
	i := 0
	for i < len(input) {
		if input[i] == '(' {
			// Insert FNC1 before each AI except the first, so the parser
			// can detect field boundaries for variable-length AIs.
			if !first {
				b.WriteByte(byte(fnc1))
			}
			first = false
			i++ // skip '('
			for i < len(input) && input[i] != ')' {
				b.WriteByte(input[i])
				i++
			}
			if i < len(input) {
				i++ // skip ')'
			}
		} else {
			b.WriteByte(input[i])
			i++
		}
	}
	return b.String()
}
