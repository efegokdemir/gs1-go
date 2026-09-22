# ADR 0002: Parser scope

## Status
Accepted

## Context
The package was audited against a GS1 LATAM solution-provider certification
checklist covering Application Identifiers, symbologies, data-capture
capabilities and business reporting. Deciding what a *library* should own,
versus the application built on top of it, keeps the API small and the
dependency count at zero.

## Decision

### In scope
- Application Identifiers used in healthcare and general distribution:
  GTIN, SSCC, batch/lot, dates, serial, counts, weights, GLN, GSIN, NHRN and
  the 90–99 internal range.
- Element strings from GS1-128, GS1 DataMatrix, GS1 QR Code, GS1 DataBar and
  composite symbols, plus bare EAN-13 / UPC-A / GTIN-14 and UPC-E expansion.
- AI separation, FNC1 handling, AIM symbology identifiers, check digits,
  date validation, and regulator-mandated AI presence checks.

### Out of scope
- Scanner hardware concerns (connectivity, configuration).
- Barcode generation, print-quality verification (ISO/IEC 15415/15416).
- GS1 Digital Link URI construction and resolution; the
  [GS1 Barcode Syntax Engine](https://github.com/gs1/gs1-syntax-engine)
  covers this area.
- Business documents (purchase orders, dispatch advices, packing lists) and
  EPCIS event capture. Applications implement these using the parsed data.

### GS1 DataBar
DataBar is a symbology family. Scanners decode it to the same AI element
string as GS1-128, so no symbology-specific code is needed.

## Consequences
- The public API stays centered on `Parse`, `Barcode` and validators.
- Feature requests for generation or Digital Link are redirected to the GS1
  reference implementations rather than duplicated here.
