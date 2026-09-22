# ADR 0003: Scanner resilience by default

## Status
Accepted

## Context
Real 2D scanners (Honeywell, Zebra, Datalogic and others) in keyboard-wedge
or USB HID mode produce output that a strict parser rejects: trailing CR/LF,
NUL bytes, a UTF-8 BOM from some USB configurations, and CR or LF used in
place of the FNC1 group separator.

## Decision
`Parse` and `ParseInto` normalize input unconditionally before parsing:

1. Strip a leading UTF-8 BOM.
2. Trim leading and trailing whitespace, CR, LF and NUL.
3. Convert internal CR/LF to FNC1 (ASCII 29).
4. Collapse consecutive FNC1 into one.
5. Drop internal NUL bytes.

Clean input takes a fast path that only scans for noise characters.
Recognized AIM identifiers: `]C1`, `]d1`, `]d2`, `]e0`, `]Q3`, `]J1`.

Alternatives rejected:
- A `ParseWithOptions` toggle. Nobody wants strict parsing of scanner noise,
  and the option would double the test matrix.
- Cleaning only in the CLI. Library users integrate the same scanners.

## Consequences
- Existing callers gain resilience without code changes.
- Overhead is a few nanoseconds on clean input and tens of nanoseconds when
  cleaning is needed.
- A `ParseStrict` can be added later without changing current behavior.
