package gs1

import "strings"

// cleanScannerInput normalizes raw barcode scanner output for parsing.
// It handles common quirks from 2D scanners (Honeywell, Zebra, Datalogic):
//   - Strips UTF-8 BOM (byte order mark)
//   - Strips leading/trailing whitespace, CR, LF, and null bytes
//   - Converts internal CR and LF to FNC1 (common scanner misconfiguration)
//   - Collapses consecutive FNC1 characters into one
//   - Removes internal null bytes
func cleanScannerInput(input string) string {
	s := input

	// Strip UTF-8 BOM.
	s = strings.TrimPrefix(s, "\xEF\xBB\xBF")

	// Strip leading noise.
	s = strings.TrimLeft(s, " \t\x00")

	// Strip trailing noise.
	s = strings.TrimRight(s, " \t\r\n\x00")

	// Fast path: if no noise characters remain, return as-is.
	if !containsScannerNoise(s) {
		return s
	}

	// Replace internal CR/LF with FNC1, remove null bytes, collapse FNC1.
	var b strings.Builder
	b.Grow(len(s))
	prevFNC1 := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\r', '\n':
			// Convert to FNC1 (scanner may use CR/LF as GS substitute).
			if !prevFNC1 {
				b.WriteByte(byte(fnc1))
				prevFNC1 = true
			}
		case byte(fnc1):
			if !prevFNC1 {
				b.WriteByte(c)
				prevFNC1 = true
			}
		case '\x00':
			// Skip null bytes.
		default:
			b.WriteByte(c)
			prevFNC1 = false
		}
	}

	return b.String()
}

// containsScannerNoise reports whether s contains characters that need cleaning.
func containsScannerNoise(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\r', '\n', '\x00':
			return true
		case byte(fnc1):
			// Check for consecutive FNC1.
			if i+1 < len(s) && s[i+1] == byte(fnc1) {
				return true
			}
		}
	}
	return false
}
