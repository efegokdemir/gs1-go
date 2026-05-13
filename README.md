# gs1

GS1 barcode parsing for healthcare supply chain traceability.

## Install

```bash
go get github.com/galenzo17/health-interop/gs1@latest
```

Standalone module — zero external dependencies, stdlib only.

## Overview

Parses GS1-128 and GS1 DataMatrix barcode scanner output into typed elements. Supports all healthcare-relevant Application Identifiers (AIs) including GTIN, batch/lot, expiration date, serial number, and NHRN codes.

## Usage

```go
import "github.com/galenzo17/health-interop/gs1"

// Parse a barcode string
b, err := gs1.Parse("0104150000021126172506302112345ABC\x1D10LOT42X")
if err != nil {
    log.Fatal(err)
}

// Convenience methods
fmt.Println(b.GTIN())         // 04150000021126
fmt.Println(b.Lot())          // LOT42X
fmt.Println(b.SerialNumber()) // 12345ABC

expiry, _ := b.ExpirationDate()
fmt.Println(expiry) // 2025-06-30

// Generic lookup
v, ok := b.Get("17")
fmt.Println(v, ok) // 250630 true

// Iterate all elements
for _, e := range b.Elements {
    fmt.Printf("AI(%s) = %s\n", e.AI, e.Value)
}
```

## GTIN Validation

GTIN check digit validation is separate from parsing (separation of concerns):

```go
err := gs1.ValidateGTIN("04150000021126")

// Compute a check digit
check, _ := gs1.ComputeGTINCheckDigit("0415000002112")
fmt.Println(string(check)) // 6
```

Supports GTIN-8, GTIN-12, GTIN-13, and GTIN-14.

## Date Parsing

GS1 dates use YYMMDD format (years 2000-2099). Day 00 means last day of month:

```go
t, _ := gs1.ParseDate("250630") // 2025-06-30
t, _ = gs1.ParseDate("250200")  // 2025-02-28 (last day of Feb)
```

## Supported Application Identifiers

| AI | Name | Type |
|---|---|---|
| 00 | SSCC | Fixed 18N |
| 01 | GTIN | Fixed 14N |
| 02 | Content GTIN | Fixed 14N |
| 10 | Batch/Lot | Variable ..20X |
| 11 | Production Date | Fixed 6N |
| 13 | Packaging Date | Fixed 6N |
| 15 | Best Before Date | Fixed 6N |
| 17 | Expiration Date | Fixed 6N |
| 21 | Serial Number | Variable ..20X |
| 30 | Count | Variable ..8N |
| 37 | Count of Trade Items | Variable ..8N |
| 240 | Additional Product ID | Variable ..30X |
| 241 | Customer Part Number | Variable ..30X |
| 310n | Net Weight kg | Fixed 6N |
| 320n | Net Weight lb | Fixed 6N |
| 330n | Gross Weight kg | Fixed 6N |
| 340n | Gross Weight lb | Fixed 6N |
| 402 | GSIN | Fixed 17N |
| 414 | GLN | Fixed 13N |
| 710-714 | NHRN | Variable ..20X |
| 90-99 | Internal/Custom | Variable ..30-90X |

## UPC-E Expansion

```go
upca, _ := gs1.ExpandUPCE("012345") // expands 6-digit UPC-E to 12-digit UPC-A
```

## Input Formats

- GS1-128 scanner output: `0104150000021126172506301012345`
- With FNC1 separators (ASCII 29): `2112345\x1D10LOT1`
- Bracket notation: `(01)04150000021126(17)250630(10)12345`
- AIM prefix (DataMatrix): `]d2...`, `]d1...` (older)
- AIM prefix (GS1-128): `]C1...`
- AIM prefix (QR/DotCode): `]Q3...`, `]J1...`

## Scanner Resilience

`Parse()` automatically handles common 2D scanner quirks:

- **Trailing CR/LF** — stripped (scanners often append `\r\n`)
- **UTF-8 BOM** — stripped (`\xEF\xBB\xBF` from some USB configs)
- **Null bytes** — removed (USB HID scanners may inject `\x00`)
- **CR/LF as FNC1** — converted to GS (`\x1D`) when used as field separator
- **Consecutive FNC1** — collapsed to single separator

No configuration needed — resilience is the default.

## Regulatory Validation (LATAM Pharma Traceability)

Validate that a barcode contains the minimum AIs required by each country's regulator:

```go
b, _ := gs1.Parse("0104150000021126172506302112345\x1D10LOT1")

err := b.ValidateANVISA()   // Brazil: requires 01 + 17 + 10 + 21
err = b.ValidateANMAT()     // Argentina: requires 01 + 17 + 10 + 21
err = b.ValidateSNFA()      // Chile: requires 01 + 17 + 10
err = b.ValidateCOFEPRIS()  // Mexico: requires 01 + 17 + 10

// Or use the Regulator type directly
err = gs1.ANVISA.Validate(b)
```

| Regulator | Country | Required AIs |
|---|---|---|
| ANVISA | Brazil | 01 (GTIN) + 17 (Expiry) + 10 (Lot) + 21 (Serial) |
| ANMAT | Argentina | 01 (GTIN) + 17 (Expiry) + 10 (Lot) + 21 (Serial) |
| SNFA | Chile | 01 (GTIN) + 17 (Expiry) + 10 (Lot) |
| COFEPRIS | Mexico | 01 (GTIN) + 17 (Expiry) + 10 (Lot) |

### CLI

```bash
hinterop parse gs1 "<barcode>" --validate anvisa
hinterop parse gs1 "<barcode>" --validate anmat
hinterop parse gs1 "<barcode>" --validate snfa
hinterop parse gs1 "<barcode>" --validate cofepris
```
