<!-- @anchors
  code: TBYSA
  updated_at: 2026-11-10
  layer: config
-->
# TestBypass — how the code skips an effect while a test runs, and the suites that run it

> **Code**: `TBYSA`

## Overview

Code sometimes skips an effect while a test runs: an early return under an end-to-end flag so the
visual captures stay clean, a write left out under CI. Measured in MIF (MIF-1S): the one suite
against the real backend never ran the write such a return skipped, and the write reached the
store review broken.

`test_bypass:` says how a bypass is recognised and which suite layers run the path it skips:

```yaml
test_bypass:
  patterns:
    - "process\\.env\\.EXPO_PUBLIC_E2E === 'true'\\)\\s*return"
  proven_by: [integration]
```

`bypass-proven` (BYPRS) asks each file with a bypass for such a proof.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the declaration | one or more `patterns`, each a regular expression, and one or more `proven_by` layers | no pattern, an empty or invalid one, no layer | this unit: refuses the configuration (`TBYSA-E01`) |

## Effects

| Effect | Description |
| --- | --- |
| `TBYSA-B01` | The project declares the patterns of its bypasses and the suite layers that run the skipped path. (`TestBypass`, `PatternREs`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `TBYSA-E01` | No pattern, an empty or invalid pattern, or no `proven_by` layer. | The configuration does not load, naming the field. | A pattern that does not compile recognises no bypass, and no layer would ask a proof nothing can give. |
