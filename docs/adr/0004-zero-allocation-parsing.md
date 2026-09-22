# ADR 0004: Zero-allocation parsing hot path

## Status
Accepted

## Date
2026-05-13

## Context
An earlier design allocated three times per `Parse` call: the element slice,
a `map[string]int` index for `Get`, and the heap escape triggered by
`append`. In continuous scanner loops and batch pipelines this GC pressure
capped throughput at roughly 2.5M parses per second per core. The target was
0–1 allocations and sustained throughput above 5M parses per second.

## Decision

### Replace the index map with a linear scan
Healthcare barcodes carry one to six elements. A linear scan over a slice of
32-byte structs beats a map lookup for n ≤ 8 because there is no hashing and
the data sits in one or two cache lines. The field was unexported, so the
change is invisible to callers.

| Approach | Allocs saved | `Get` cost | Complexity |
|---|---|---|---|
| Keep map, `clear()` on reuse | 0 | O(1) | Low |
| **Linear scan** | **1** | **O(n), n ≤ 8** | **Low** |
| Fixed-size array index | 1 | O(1) | Medium |

### Expose `ParseInto` with a caller-owned `Barcode`
A caller that owns the `Barcode` and calls `Reset` between scans gets zero
allocations. `sync.Pool` was rejected: its Get/Put overhead is a measurable
fraction of a ~200 ns parse, and callers who want pooling can wrap
`ParseInto` themselves.

### Strings need no action
`Element.AI` references the literal in the AI table and `Element.Value` is a
substring of the input, so neither allocates. The only allocation left is
zero-padding a bare EAN-13/UPC-A to GTIN-14, a fallback path.

## Consequences
- `Parse` performs one allocation; `ParseInto` performs none on the happy path.
- `Get` is O(n); revisit if a future format exceeds ~20 elements.
- `BenchmarkParseInto` must stay at 0 allocs/op; CI contributors check this.
