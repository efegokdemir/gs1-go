# ADR 0001: Standalone repository and module path

## Status
Accepted

## Date
2026-09-22

## Context
The parser started as the `gs1` sub-module of
[health-interop](https://github.com/galenzo17/health-interop), a broader
LATAM healthcare interoperability toolkit. It was already a separate Go
module with zero dependencies, but its import path
(`github.com/galenzo17/health-interop/gs1`), issue tracker, release tags and
CI were shared with unrelated national-identifier packages and a Cobra CLI.

GS1 tooling has a natural audience outside healthcare interoperability:
retail, logistics, and anyone integrating a barcode scanner. A repository
whose name and README are about GS1 is easier to discover, cite and
contribute to.

## Decision
- Extract the module to its own repository, `github.com/galenzo17/gs1-go`,
  preserving the commit history of the `gs1/` directory via `git subtree split`.
- Module path `github.com/galenzo17/gs1-go`, package name `gs1`. The
  repository name follows the GS1 organization's convention of a `gs1-`
  prefix plus a descriptive suffix (`gs1-syntax-engine`,
  `gs1-syntax-dictionary`), with the language as suffix as in
  `gs1encoders-swift`.
- Versioning restarts at `v0.1.0` with plain `vX.Y.Z` tags; there is no
  longer a `gs1/` tag prefix.
- Governance mirrors the GS1 organization repositories: `CONTRIBUTING.md`,
  `SECURITY.md`, a changelog and a `docs/` folder. The license remains MIT.
- The CLI is rewritten on the standard library `flag` package so the module
  keeps zero dependencies. The Bubble Tea TUI stays in health-interop.

## Consequences
- Consumers of `github.com/galenzo17/health-interop/gs1` should migrate the
  import path; only the path changes, the API is unchanged apart from the
  added `LookupAI`.
- health-interop can depend on this module instead of vendoring it, once it
  updates its `go.work` and import paths.
- The repository can adopt a trademark notice clarifying that it is not a
  GS1 product, which matters more for a repository named after GS1.
