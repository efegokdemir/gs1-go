# ADR 0007: Configurable day-zero date policy

## Status
Accepted

## Date
2026-09-21

## Context
GS1 dates are `YYMMDD`. The GS1 General Specifications allow `DD = 00` for
AIs such as (17) to mean the last day of the month, and this usage is common
on pharmaceutical packs. Some national traceability systems and ERPs instead
store such dates as the first day of the month, and integrators need to
match the system of record rather than the specification. Other integrators
need the untouched six-digit string for round-tripping.

## Decision
- `ParseDate` keeps the specification behavior: day `00` resolves to the
  last day of the month.
- `ParseDateWithOptions(value, DateOptions{DayZero: …})` accepts
  `DayZeroLastDay` (default) or `DayZeroFirstDay`.
- `ParseDateRaw` validates the field and returns the original string.
- The WASM API exposes the same choice through `dateFormat` (`"raw"` default,
  `"iso"`) and `dayZero` (`"last"` default, `"first"`).
- Impossible dates (month 13, February 30) are rejected rather than
  normalized, unlike `time.Date`.

## Consequences
- Default behavior is unchanged and specification-compliant.
- Options are a struct, so further policies (for example two-digit year
  pivots if GS1 ever changes the 2000–2099 window) can be added without
  breaking the signature.
