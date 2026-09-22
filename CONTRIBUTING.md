# Contributing

Thanks for helping improve gs1-go. This document covers the workflow and the
bar for accepted changes.

## Development

```bash
git clone https://github.com/galenzo17/gs1-go.git
cd gs1-go
make all        # vet + lint + tests
make race       # tests with the race detector
make fuzz       # 30s fuzz run of the parser
make bench      # parser benchmarks with allocation counts
make wasm       # build wasm/gs1.wasm and the browser example
```

`golangci-lint` v2 is required for `make lint`. Install it from
<https://golangci-lint.run/welcome/install/>.

## Ground rules

- **Zero dependencies.** The module depends only on the Go standard library.
  Pull requests adding a `require` line will not be accepted.
- **Specification first.** Parsing and validation behavior must be traceable
  to the GS1 General Specifications, the GS1 Syntax Dictionary, or a public
  regulator document. Cite the section in the PR description.
- **Tests are table-driven** and live next to the code. Every new parsing
  path needs positive and negative cases; add seeds to `FuzzParse` when the
  input shape is new.
- **No allocations on the hot path.** Run `make bench` and check that
  `BenchmarkParseInto` stays at 0 allocs/op.
- **Godoc everything exported.** Comments start with the identifier name and
  explain behavior, not implementation.
- **Backwards compatibility.** The public API follows semantic versioning.
  Breaking changes require a major version bump and an ADR.

## Architecture decisions

Non-trivial design choices are recorded in `docs/adr/`. Add a new ADR when a
change alters the public API shape, the parsing strategy, or the supported
targets. Copy the structure of an existing ADR.

## Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat(parser): recover from missing FNC1 between variable-length fields
fix(date): reject day 32 before month normalization
docs(readme): document ParseInto reuse contract
```

## Release process

1. Move the *Unreleased* section of `CHANGELOG.md` under a new version heading.
2. Tag: `git tag -a vX.Y.Z -m "vX.Y.Z" && git push origin vX.Y.Z`.
3. The release workflow runs the test suite, builds CLI binaries and the WASM
   bundle, and publishes a GitHub release with generated notes.

## Legal

By contributing you agree that your contribution is licensed under the MIT
license of this repository. Do not submit code copied from GS1 publications
or from software under incompatible licenses.
