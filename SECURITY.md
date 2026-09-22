# Security policy

## Supported versions

The latest minor release receives security fixes. Older releases are not
patched; upgrade to the current release.

## Threat model

gs1-go parses untrusted scanner input. Bugs that matter most are:

- panics or unbounded memory use on crafted input (denial of service in a
  scanning loop),
- element strings that parse to a *different* GTIN, lot, or serial than the
  one encoded (mis-identification of a medicine or device),
- check-digit or date validation accepting invalid data.

The parser runs under continuous fuzzing in CI to catch the first class.

## Reporting a vulnerability

Please do not open a public issue. Use GitHub's private vulnerability
reporting on this repository (*Security* tab, *Report a vulnerability*), or
email the maintainer listed in the repository profile.

Include the exact input (escape FNC1 as `\x1D`), the observed behavior, and
the expected behavior. You will receive an acknowledgement within 72 hours
and a fix or mitigation plan within 14 days for confirmed reports.
