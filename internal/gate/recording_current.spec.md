<!-- @anchors
  code: RCCRS
  updated_at: 2026-11-10
  layer: gate
-->
# RecordingCurrent — a recorded response is of the system the project speaks to today

> **Code**: `RCCRS`

## Overview

A double built from a real response is right on the day it is recorded. Measured in MIF: nine
defects passed against recordings whose system had moved on — an API version that changed
(Stripe's `dahlia`), a capability added, a presign condition the old response did not have —,
and the double kept answering the old shape.

Each recording says where it came from, on a line of its own: `@recorded: <system> <version>
<date>`. The project declares, per system (`external_systems:`, EXSSA), the version it speaks and
how many days a recording lives. A recording of another version, of an undeclared system, or
whose mark does not say the three, fails; one past its system's age diverges until it is
recorded again.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a recording's mark | `@recorded:` followed by the system, the version and the date (`2026-10-02`), in any comment | a mark missing any of the three, or a date that is no date | this unit: fails, naming the mark |
| the systems | `external_systems:` with `version` and `max_age` | — | the configuration (EXSSA) |

## Effects

| Effect | Description |
| --- | --- |
| `RCCRS-B01` | A recording of a version other than its system's declared `version` fails, saying the two versions; a recording of the declared one, or of a system that declares none, passes. (`checkRecordingCurrent`) |
| `RCCRS-B02` | A recording older than its system's `max_age` diverges, saying its date and age; with no `max_age`, a recording does not age. |
| `RCCRS-B03` | A file with no recording is skipped. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RCCRS-E01` | A mark that does not say the system, the version and a date, or that names a system `external_systems:` does not declare. | Fails, naming the mark or the system. | A recording no one can date or place cannot be told current, and passing it would vouch for a double nobody can check. |
