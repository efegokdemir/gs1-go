# ADR 0000: Record architecture decisions

## Status
Accepted

## Context
Parsing GS1 element strings involves judgment calls that are not obvious
from the code: how tolerant to be of scanner noise, which ambiguities to
resolve heuristically, what stays out of scope. Contributors need the
reasoning, not just the result.

## Decision
Use lightweight Architecture Decision Records in `docs/adr/`, numbered
sequentially, with Status, Context, Decision and Consequences sections.
Records are immutable once accepted; later records supersede earlier ones.

## Consequences
- Every non-trivial design change in a pull request is accompanied by an ADR.
- Reviewers can evaluate alternatives that were considered and rejected.
