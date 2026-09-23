# Integration guide

`gs1-go` is the parsing and validation boundary between a data-capture
device and an application record. It deliberately does not own the transport,
database, EDI, or label-printing layer. Keep those responsibilities separate
so a scanner replacement does not change business logic.

## Data capture

Every input shape should be reduced to one element string before parsing.
`ParseInto` is useful for long-running consumers because the caller can reuse
the `Barcode` value for every scan.

### Keyboard wedge or HID scanner

Scanner input normally arrives as one line on stdin. The CLI already handles
this shape:

```bash
gs1 parse < scans.txt
```

For an application-owned loop:

```go
scanner := bufio.NewScanner(os.Stdin)
var barcode gs1.Barcode
for scanner.Scan() {
    barcode.Reset()
    if err := gs1.ParseInto(scanner.Text(), &barcode); err != nil {
        log.Printf("reject scan: %v", err)
        continue
    }
    store(barcode)
}
if err := scanner.Err(); err != nil {
    log.Fatal(err)
}
```

The parser normalizes common scanner noise, including CR/LF suffixes, BOMs,
NUL bytes, and FNC1 separators. Do not strip the raw scan before parsing if
the original element string is needed for audit or reprocessing.

### TCP scanner or mobile-computer stream

Read complete scanner messages from the device protocol, then pass each
message to the same parser. A TCP connection is a transport boundary, not a
barcode boundary; use the device's documented delimiter or length framing.

```go
func handleScanner(conn net.Conn) error {
    defer conn.Close()
    scanner := bufio.NewScanner(conn)
    var barcode gs1.Barcode
    for scanner.Scan() {
        barcode.Reset()
        if err := gs1.ParseInto(scanner.Text(), &barcode); err != nil {
            continue // record the rejected raw message in production
        }
        store(barcode)
    }
    return scanner.Err()
}
```

Zebra, Honeywell, and Datalogic devices can expose keyboard-wedge, TCP, or
application-managed output. Keep device configuration outside this package;
the parsed result should be identical for equivalent element strings.

### HTTP ingestion

Treat an HTTP request as untrusted input. Bound the body, authenticate the
caller, and retain the request identifier with the scan result.

```go
func scanHandler(w http.ResponseWriter, r *http.Request) {
    body := io.LimitReader(r.Body, 4<<10)
    input, err := io.ReadAll(body)
    if err != nil {
        http.Error(w, "read scan", http.StatusBadRequest)
        return
    }
    barcode, err := gs1.Parse(string(input))
    if err != nil {
        http.Error(w, "invalid scan", http.StatusUnprocessableEntity)
        return
    }
    store(barcode)
    w.WriteHeader(http.StatusAccepted)
}
```

In a production handler, add authentication, replay protection, request-size
limits, structured logging, and an idempotency key before persisting data.

## Storage and filtering

Store the original element string alongside normalized fields. A practical
scan record contains:

| Field | Source | Storage guidance |
|---|---|---|
| `raw_element_string` | `Barcode.Raw` | Preserve exactly for audit and replay. |
| `gtin14` | AI `01` | Store the trade-item GTIN as a fixed-width string, not an integer. AI `02` identifies contained trade items. |
| `lot` | AI `10` | Preserve leading zeroes and case. |
| `serial` | AI `21` | Preserve as text; it may be alphanumeric. |
| `expiry_date` | AI `17` | Convert only after successful date validation. |
| `sscc` | AI `00` | Store as text; it is an identifier, not a quantity. |
| `gln` | AI `414` or a role-specific GLN | Keep the AI role when multiple GLNs occur. |
| `captured_at` | application clock | Store UTC with the device/request ID. |

Use `Barcode.Get` or the typed accessors for the values required by the
application. Keep unknown or newly supported AIs in the raw string and
element list rather than silently discarding them. Index the normalized
columns used for operational filters such as GTIN, lot, expiry, and SSCC.

## Logistic hierarchy and dispatch advice

Represent the physical hierarchy in application records:

```text
SSCC (AI 00)            outer logistic unit
  └─ GTIN (AI 02)       contained trade item
       └─ quantity (37)  number of contained items
```

For a case-level GTIN-14, the indicator digit is prepended to the GTIN-13
item data before the check digit is recomputed. `Barcode.GTIN()` reads AI 01;
read AI 02 explicitly with `Barcode.Get("02")` when the scan describes
contained trade items.

The library parses these identifiers; it does not infer containment. Your
receiving or warehouse service should validate that a child scan is attached
to the expected parent and should retain the scan timestamp and source.

When producing a dispatch advice, map application records to the document
standard selected by the trading partners. Typical EANCOM DESADV mappings
include:

| Application value | EANCOM DESADV location |
|---|---|
| SSCC | `GIN+BJ` |
| GTIN | `LIN+1++04150000021126:SRV` |
| Additional item identifier | `PIA` |
| Quantity | `QTY+12` (despatched quantity) |
| Expiry date | `DTM+361:20250630:102` (AI 17 is YYMMDD; validate it, then convert to CCYYMMDD) |
| Lot or batch | `GIN+BX` |

GS1 XML DespatchAdvice 3.x uses equivalent structures for the same record:

| Application value | GS1 XML DespatchAdvice 3.x location |
|---|---|
| Logistic-unit SSCC | `despatchAdviceLogisticUnit/logisticUnitIdentification/sscc` |
| Trade-item GTIN | `despatchAdviceLogisticUnit/transactionalTradeItem/gtin` |
| Lot or batch | `transactionalTradeItem/lotNumber` |
| Expiration date | `transactionalTradeItem/itemExpirationDate` |
| Dispatched quantity | `despatchAdviceLogisticUnit/transactionalTradeItem/despatchedQuantity` |

The exact XML element and qualifier depend on the message version and partner
implementation guide; do not generate EANCOM or XML by concatenating scanner
text.

## Transports and standards outside this module

| Concern | Common standards or transports | Responsibility |
|---|---|---|
| Device connectivity | HID, TCP, Android intent, HTTP | Capture application |
| Business documents | GS1 XML DespatchAdvice, EANCOM DESADV | EDI/document service |
| Secure transport | AS2, AS4, HTTPS, partner VPN | Integration/platform service |
| Master-data synchronisation | GDSN | Product-data service |
| Barcode generation | GS1 symbology encoder/renderer | Label service |

See the [GS1 General Specifications](https://www.gs1.org/standards/barcodes-epcrfid-id-keys/gs1-general-specifications),
[GS1 XML](https://www.gs1.org/standards/gs1-xml),
[EANCOM](https://www.gs1.org/standards/eancom), and
[GDSN](https://www.gs1.org/standards/gdsn) documentation for the normative
requirements. See [ADR 0002](adr/0002-parser-scope.md) for this package's
boundary.
