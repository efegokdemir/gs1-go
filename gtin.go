package gs1

import "fmt"

// NewSSCC builds an 18-digit Serial Shipping Container Code from its
// extension digit, company prefix, and serial reference. The company prefix
// and serial reference must contain exactly 16 digits together.
func NewSSCC(extensionDigit byte, companyPrefix, serialReference string) (string, error) {
	if extensionDigit < '0' || extensionDigit > '9' {
		return "", fmt.Errorf("%w: extension digit must be a digit", ErrInvalidData)
	}
	if err := validateDigits("company prefix", companyPrefix); err != nil {
		return "", err
	}
	if err := validateDigits("serial reference", serialReference); err != nil {
		return "", err
	}
	if len(companyPrefix)+len(serialReference) != 16 {
		return "", fmt.Errorf("%w: company prefix and serial reference must total 16 digits", ErrInvalidData)
	}

	partial := string(extensionDigit) + companyPrefix + serialReference
	check, _ := computeCheckDigit(partial)
	return partial + string(check), nil
}

// NewGTIN14 builds a case-level GTIN-14 from a valid GTIN-13 and an
// indicator digit. The indicator is prepended to the GTIN-13 data digits and
// the check digit is recomputed for the resulting 14-digit key.
func NewGTIN14(indicator byte, gtin13 string) (string, error) {
	if indicator < '0' || indicator > '9' {
		return "", fmt.Errorf("%w: indicator must be a digit", ErrInvalidData)
	}
	if err := ValidateGTIN(gtin13); err != nil {
		return "", fmt.Errorf("invalid GTIN-13: %w", err)
	}
	if len(gtin13) != 13 {
		return "", fmt.Errorf("%w: GTIN-13 must be 13 digits", ErrInvalidData)
	}

	partial := string(indicator) + gtin13[:12]
	check, _ := computeCheckDigit(partial)
	return partial + string(check), nil
}

// NewGTIN13 builds a GTIN-13 from a company prefix and item reference. The
// two components must contain exactly 12 digits together.
func NewGTIN13(companyPrefix, itemReference string) (string, error) {
	if err := validateDigits("company prefix", companyPrefix); err != nil {
		return "", err
	}
	if err := validateDigits("item reference", itemReference); err != nil {
		return "", err
	}
	if len(companyPrefix)+len(itemReference) != 12 {
		return "", fmt.Errorf("%w: company prefix and item reference must total 12 digits", ErrInvalidData)
	}

	partial := companyPrefix + itemReference
	check, _ := computeCheckDigit(partial)
	return partial + string(check), nil
}

// SplitSSCC validates an SSCC and splits it into its extension digit,
// company prefix, and serial reference. gcpLen specifies the company-prefix
// length because GS1 does not encode that length in an SSCC.
func SplitSSCC(sscc string, gcpLen int) (extension byte, companyPrefix, serialReference string, err error) {
	if len(sscc) != 18 {
		return 0, "", "", fmt.Errorf("%w: SSCC must be 18 digits", ErrInvalidData)
	}
	if err := validateDigits("SSCC", sscc); err != nil {
		return 0, "", "", err
	}
	if err := ValidateCheckDigit(sscc); err != nil {
		return 0, "", "", err
	}
	if gcpLen < 1 || gcpLen > 15 {
		return 0, "", "", fmt.Errorf("%w: GCP length must be between 1 and 15", ErrInvalidData)
	}

	return sscc[0], sscc[1 : 1+gcpLen], sscc[1+gcpLen : 17], nil
}

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
	if len(partial) == 0 {
		return 0, fmt.Errorf("%w: partial length must be positive", ErrInvalidData)
	}
	return computeCheckDigit(partial)
}

func computeCheckDigit(partial string) (byte, error) {
	n := len(partial)
	if n > 30 {
		return 0, fmt.Errorf("%w: partial length %d exceeds 30 digits", ErrInvalidData, n)
	}
	if err := validateDigits("partial key", partial); err != nil {
		return 0, err
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

func validateDigits(label, value string) error {
	for i := 0; i < len(value); i++ {
		if value[i] < '0' || value[i] > '9' {
			return fmt.Errorf("%w: %s contains non-digit character at position %d", ErrInvalidData, label, i)
		}
	}
	return nil
}

// ValidateCheckDigit validates the modulo-10 check digit of a numeric GS1
// key whose final digit is the check digit.
func ValidateCheckDigit(key string) error {
	if len(key) < 2 || len(key) > 31 {
		return fmt.Errorf("%w: key length %d not valid", ErrInvalidData, len(key))
	}
	if err := validateDigits("key", key); err != nil {
		return err
	}
	expected, _ := computeCheckDigit(key[:len(key)-1])
	if expected != key[len(key)-1] {
		return fmt.Errorf("%w: expected '%c', got '%c'", ErrInvalidCheckDigit, expected, key[len(key)-1])
	}
	return nil
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
