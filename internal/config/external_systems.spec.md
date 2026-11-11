<!-- @anchors
  code: EXSSA
  updated_at: 2026-11-10
  layer: config
-->
# ExternalSystems — the real systems the code talks to, and the suites that reach them

> **Code**: `EXSSA`

## Overview

A unit test proves a unit against doubles, and a double of a database, an API or a store accepts
whatever shape it is given. Measured in MIF: fourteen defects crossed such a boundary — an index
that did not exist, a condition outside the key schema, a field the API refused — and nine passed
against recorded responses of an API version that was gone. Each was green.

`external_systems:` names each real system, how a call to it is recognised in the code, which
suite layers reach it, and — for recorded responses — the API version the project speaks and how
long a recording lives:

```yaml
external_systems:
  - name: appsync
    calls: "\\bclient\\.graphql\\("
    proven_by: [integration, live]
  - name: stripe
    calls: "\\bstripe\\.[a-z]+\\.(create|update|retrieve)\\("
    proven_by: [live]
    version: 2025-09-30.dahlia
    max_age: 90d
```

`external-proven` (EXPRS) asks a test of each calling file in those layers; `recording-current`
(RCCRS) asks each recording to be of the version declared, within its age. The name is not
`boundaries:`, which holds the layers' boundaries.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a system | a unique name, a `calls` regular expression, one or more `proven_by` layers, an optional `version`, an optional `max_age` in days (`90d`) | no name, a repeated name, an empty or invalid `calls`, no `proven_by`, a `max_age` that is not days | this unit: refuses the configuration (`EXSSA-E01`) |

## Effects

| Effect | Description |
| --- | --- |
| `EXSSA-B01` | Each system declares its name, the pattern of a call to it, the suite layers that reach it, the API version the project speaks, and the days a recording lives. (`ExternalSystem`, `CallsRE`, `MaxAgeDays`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `EXSSA-E01` | A system with no name or a repeated one, an empty or invalid `calls`, no `proven_by`, or a `max_age` that is not a number of days. | The configuration does not load, naming the system and the field. | A pattern that does not compile recognises no call, and a system no layer proves would ask a proof nothing can give — both silent. |
