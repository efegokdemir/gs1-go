# ADR 0006: Missing FNC1 recovery heuristic

## Status
Accepted

## Date
2026-05-25

## Context
Some scanner configurations drop the FNC1 separator between two
variable-length fields, so `10LOT42X21SERIAL` arrives as a single run of
characters. A strict parser reads the serial into the lot until the 20
character limit is exceeded and fails. The data is recoverable when a valid
AI code marks the boundary, but AI codes such as `21` also occur naturally
inside lot numbers (`HC23L25212800`), so the first match is often wrong.

## Decision
When a variable-length value exceeds its maximum length, search the
overflowing region for candidate boundaries where:

1. a known AI code begins,
2. the data following it is plausible for that AI (length and character class),
3. a dry-run parse of the remainder succeeds with no further recovery.

Among valid candidates, pick the one closest to the midpoint of the search
range, preferring the later position on ties. This *balanced split* favors
two reasonably sized fields over one tiny field and one huge one, which
matches how real lot and serial numbers are sized.

Alternatives rejected:
- First valid boundary: fails on lot numbers containing `21` or `10`.
- Last valid boundary: fails symmetrically on serials containing AI codes.
- Refusing to recover: rejects otherwise usable scans at the point of care.

## Consequences
- Scans with dropped separators parse in the common case.
- Recovery is heuristic. Applications that require certainty should check
  that the raw input contains FNC1 where expected, or validate the recovered
  lot and serial against a master record.
- The fuzz target exercises the recovery path; it must never panic or loop.
