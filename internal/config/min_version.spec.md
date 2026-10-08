<!-- @anchors
  code: MNVRM
  updated_at: 2026-10-08
  layer: config
-->
# MinVersion — whether the running binary meets the minimum version the project declares

> **Code**: `MNVRM`

## Overview

`min_version` exists for the moment an Anchors fix must reach everyone before work goes on.
A person declares it in the file the team versions, and every command compares the running
binary with it. The question is ORDINAL — does this binary meet the minimum? — so the
comparison is part by part, as numbers, never as text: as text, 0.1.9 would be newer than
0.1.84.

Some versions cannot be ordered. A local build calls itself `dev`, and it may be newer than
any release or older than all of them. A pre-release compared as if it were the final
release would answer "meets" to whoever has less than the minimum asks. Such a version is
not ordered: the comparison answers "does not meet" and gives the reason, and the caller
decides whether to warn or to block. What the unit never does is pretend to know.

The declared minimum itself is checked when the configuration loads. A value that cannot be
compared (`latest`, `0.1`, `dev`) would only fail at comparison time, and the effect there is
the warning vanishing — the exact silence the field exists to break. So such a value is
refused at load, with a message that names the format, gives an example and says what a bad
value would silence. `dev` is refused in the field even though a `dev` binary is tolerated:
a development binary is a fact one finds, the minimum is a decision someone writes, and "at
least dev" decides nothing.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the declared minimum | MAJOR.MINOR.PATCH, optionally with the tag prefix v, or nothing | any other text, pre-releases and `dev` included | this unit: the configuration load refuses it |
| the running version | any text the binary reports | — | this unit: a version that cannot be ordered does not meet the minimum, and the reason is returned |

## Effects

| Effect | Description |
| --- | --- |
| `MNVRM-B01` | Versions are ordered part by part as numbers: major over minor over patch, and 9 is smaller than 84 (`VersionOrder`). |
| `MNVRM-B02` | The tag prefix v is accepted on either side. |
| `MNVRM-B03` | With no minimum declared, any running binary meets it, whatever it reports (`AtendeMinVersion`). |
| `MNVRM-B04` | A running version that cannot be ordered (`dev`, empty, not three numeric parts) does not meet a declared minimum, and the reason is returned. |
| `MNVRM-B05` | A declared minimum that is not MAJOR.MINOR.PATCH is refused when the configuration loads, `dev` and pre-releases included. |
| `MNVRM-B06` | An absent minimum, or a well-formed one with or without the tag prefix, is accepted. |
| `MNVRM-B07` | The refusal names the expected format, gives an example, quotes the value, and says a value that cannot be compared silences the warning the field exists to give. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MNVRM-I01` | The order is antisymmetric: swapping the two versions always flips the answer, and only equal versions compare as equal. | compares every pair of a set of versions both ways |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MNVRM-X01` | A pre-release is never compared as if it were its final release. | Comparing 0.1.85-rc1 as 0.1.85 would answer "meets" to a binary that may lack the fix the minimum demands. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MNVRM-E01` | REF[MNVRM-B05]: a malformed declared minimum is the load failure B05 answers: the configuration does not load | — | — |
| `MNVRM-E02` | REF[MNVRM-B04]: a running version that cannot be ordered is the comparison failure B04 answers: it does not meet the minimum, with the reason | — | — |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
