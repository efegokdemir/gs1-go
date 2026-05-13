package gs1

import "fmt"

// ValidateGTIN validates the modulo-10 check digit of a GTIN string.
// It accepts GTIN-8, GTIN-12, GTIN-13, and GTIN-14.
func ValidateGTIN(gtin string) error {
	n := len(gtin)
	if n != 8 && n != 12 && n != 13 && n != 14 {
		return fmt.Errorf("%w: length %d not valid (expected 8, 12, 13, or 14)", ErrInvalidData, n)
	}

	for i := 0; i < n; i++ {
		if gtin[i] < '0' || gtin[i] > '9' {
			return fmt.Errorf("%w: non-digit character at position %d", ErrInvalidData, i)
		}
	}

	expected, _ := ComputeGTINCheckDigit(gtin[:n-1])
	if expected != gtin[n-1] {
		return fmt.Errorf("%w: expected '%c', got '%c'", ErrInvalidCheckDigit, expected, gtin[n-1])
	}

	return nil
}

// ComputeGTINCheckDigit computes the GS1 modulo-10 check digit for a partial
// GTIN (all digits except the check digit). Returns the check digit as a byte
// ('0'-'9').
func ComputeGTINCheckDigit(partial string) (byte, error) {
	n := len(partial)
	if n != 7 && n != 11 && n != 12 && n != 13 {
		return 0, fmt.Errorf("%w: partial length %d not valid", ErrInvalidData, n)
	}

	for i := 0; i < n; i++ {
		if partial[i] < '0' || partial[i] > '9' {
			return 0, fmt.Errorf("%w: non-digit character at position %d", ErrInvalidData, i)
		}
	}

	// Weights alternate 3,1,3,1... counting from the rightmost position of
	// the partial string. The rightmost digit gets weight 3.
	sum := 0
	for i := n - 1; i >= 0; i-- {
		d := int(partial[i] - '0')
		if (n-1-i)%2 == 0 {
			sum += d * 3
		} else {
			sum += d
		}
	}

	check := (10 - (sum % 10)) % 10
	return byte('0' + check), nil
}

// ExpandUPCE expands a 6-digit UPC-E code (without check digit) or an
// 8-digit UPC-E (with number system and check digit) to a 12-digit UPC-A.
// The input must be 6 or 8 digits.
func ExpandUPCE(upce string) (string, error) {
	for i := 0; i < len(upce); i++ {
		if upce[i] < '0' || upce[i] > '9' {
			return "", fmt.Errorf("%w: non-digit character in UPC-E", ErrInvalidData)
		}
	}

	var digits string
	switch len(upce) {
	case 8:
		// Format: NS + 6 digits + check. Extract the 6 core digits.
		digits = upce[1:7]
	case 6:
		digits = upce
	default:
		return "", fmt.Errorf("%w: UPC-E must be 6 or 8 digits, got %d", ErrInvalidData, len(upce))
	}

	// UPC-E to UPC-A expansion rules based on last digit of the 6-digit code.
	// Given digits d1..d6, expand to number system 0 + manufacturer(5) + product(5).
	var expanded string
	switch digits[5] {
	case '0', '1', '2':
		expanded = "0" + string(digits[0:2]) + string(digits[5]) + "0000" + string(digits[2:5])
	case '3':
		expanded = "0" + string(digits[0:3]) + "00000" + string(digits[3:5])
	case '4':
		expanded = "0" + string(digits[0:4]) + "00000" + string(digits[4:5])
	default: // 5-9
		expanded = "0" + string(digits[0:5]) + "0000" + string(digits[5])
	}

	partial := expanded
	check, _ := ComputeGTINCheckDigit(partial)
	return partial + string(check), nil
}
