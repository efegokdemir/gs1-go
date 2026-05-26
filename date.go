package gs1

import (
	"fmt"
	"time"
)

// DayZeroPolicy controls how day=00 in YYMMDD dates is resolved.
type DayZeroPolicy int

const (
	// DayZeroLastDay resolves day 00 to the last day of the month (GS1 default).
	DayZeroLastDay DayZeroPolicy = iota
	// DayZeroFirstDay resolves day 00 to the first day of the month.
	DayZeroFirstDay
)

// DateFormat controls the output style for date fields.
type DateFormat int

const (
	// DateFormatTime returns parsed time.Time values (default).
	DateFormatTime DateFormat = iota
	// DateFormatRaw returns the raw YYMMDD string without parsing to time.Time.
	DateFormatRaw
)

// DateOptions configures date parsing behavior.
type DateOptions struct {
	DayZero DayZeroPolicy
	Format  DateFormat
}

// ParseDate parses a 6-digit YYMMDD date string as used in GS1 barcodes.
// Years map to 2000-2099. Day 00 means the last day of the specified month
// (deprecated since 2025 but supported for backward compatibility).
func ParseDate(yymmdd string) (time.Time, error) {
	return ParseDateWithOptions(yymmdd, DateOptions{})
}

// ParseDateWithOptions parses a YYMMDD date with configurable day-zero policy.
// When opts.Format is DateFormatRaw, the returned time is zero-valued; use
// ParseDateRaw instead for a string result.
func ParseDateWithOptions(yymmdd string, opts DateOptions) (time.Time, error) {
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
		switch opts.DayZero {
		case DayZeroFirstDay:
			return time.Date(year, time.Month(mm), 1, 0, 0, 0, 0, time.UTC), nil
		default:
			// Last day of month: first day of next month minus one day.
			t := time.Date(year, time.Month(mm)+1, 1, 0, 0, 0, 0, time.UTC)
			return t.AddDate(0, 0, -1), nil
		}
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

// ParseDateRaw validates a YYMMDD string and returns it as-is without
// converting to time.Time. Useful when the caller wants the raw GS1 value.
func ParseDateRaw(yymmdd string) (string, error) {
	// Validate by parsing, then discard the time.
	_, err := ParseDateWithOptions(yymmdd, DateOptions{})
	if err != nil {
		return "", err
	}
	return yymmdd, nil
}
