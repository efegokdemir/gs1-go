package gs1

import (
	"fmt"
	"time"
)

// ParseDate parses a 6-digit YYMMDD date string as used in GS1 barcodes.
// Years map to 2000-2099. Day 00 means the last day of the specified month
// (deprecated since 2025 but supported for backward compatibility).
func ParseDate(yymmdd string) (time.Time, error) {
	if len(yymmdd) != 6 {
		return time.Time{}, ErrInvalidDate
	}

	for i := 0; i < 6; i++ {
		if yymmdd[i] < '0' || yymmdd[i] > '9' {
			return time.Time{}, ErrInvalidDate
		}
	}

	yy := int(yymmdd[0]-'0')*10 + int(yymmdd[1]-'0')
	mm := int(yymmdd[2]-'0')*10 + int(yymmdd[3]-'0')
	dd := int(yymmdd[4]-'0')*10 + int(yymmdd[5]-'0')

	year := 2000 + yy

	if mm < 1 || mm > 12 {
		return time.Time{}, fmt.Errorf("%w: month %d", ErrInvalidDate, mm)
	}

	if dd == 0 {
		// Last day of month: first day of next month minus one day.
		t := time.Date(year, time.Month(mm)+1, 1, 0, 0, 0, 0, time.UTC)
		return t.AddDate(0, 0, -1), nil
	}

	if dd > 31 {
		return time.Time{}, fmt.Errorf("%w: day %d", ErrInvalidDate, dd)
	}

	t := time.Date(year, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
	// Go normalizes invalid dates (e.g., Feb 30 → Mar 2). Detect this.
	if t.Month() != time.Month(mm) || t.Day() != dd {
		return time.Time{}, fmt.Errorf("%w: day %d invalid for month %d", ErrInvalidDate, dd, mm)
	}

	return t, nil
}
