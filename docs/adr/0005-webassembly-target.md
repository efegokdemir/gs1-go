# ADR 0005: WebAssembly build target

## Status
Accepted

## Date
2026-05-13

## Context
Browser-based pharmacy, dispensing and warehouse applications receive scanner
output as keystrokes and need the same parsing rules as the back end. The
module is pure Go with no dependencies, so compiling it to WebAssembly is
cheap and keeps one implementation.

## Decision

### Go's native `GOOS=js GOARCH=wasm`
| Approach | Binary | Stdlib | Maintenance |
|---|---|---|---|
| **Go native WASM** | **~3 MB stripped** | **Full** | **Same source** |
| TinyGo | ~0.5–1 MB | Partial | Second toolchain |
| Rewrite in Rust | ~0.2 MB | n/a | Second implementation |

Binary size is acceptable for line-of-business applications. TinyGo can be
revisited if size becomes a constraint.

### Global `gs1.*` functions
`gs1.parse(input, options)`, `gs1.validateGTIN(gtin)` and
`gs1.validateRegulatory(input, regulator)` map directly onto scanner
`onScan` callbacks. A class-based API adds ceremony; a Web Worker adds
latency to a sub-millisecond operation.

### Thin loader and hand-written types
`wasm/gs1-loader.js` wraps Go's `wasm_exec.js` in an `async loadGS1(url)`
that returns the namespace. `wasm/gs1.d.ts` is hand-written; the surface is
three functions.

### Dates default to raw `YYMMDD`
JavaScript callers choose `dateFormat: "iso"` and `dayZero: "first" | "last"`
explicitly, mirroring `DateOptions` in Go.

## Consequences
- One parser, one test suite, one set of correctness guarantees.
- Supported: Chrome 57+, Firefox 52+, Safari 11+, Edge 79+. Not supported: IE 11.
- `wasm_exec.js` must match the Go toolchain used to build; the build script
  copies it from `GOROOT`.
