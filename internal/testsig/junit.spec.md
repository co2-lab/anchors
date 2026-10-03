<!-- @anchors
  code: JUIJN
  updated_at: 2026-10-03
  layer: infra
-->
# JUnitIngest — the run's outcome per test case, and the scenario codes each case proves, read from a JUnit report

> **Code**: `JUIJN`

## Overview

Anchors never runs a project's tests: it reads the report the project's runner already writes. JUnit XML is
the execution format nearly every runner emits, so reading it keeps the tool agnostic of the stack.

This unit reads the report into one result per test case — its name, its class, its file, and whether it
failed or was skipped — and then answers two questions about scenario codes. Which codes are PROVEN: those
named by a case that neither failed nor was skipped. Which codes were SEEN: those named by any case, whatever
the outcome, so that a partial run knows which scenarios it actually measured. The codes are recognized in the
case's NAME by the rule code grammar, in the project's vocabulary.

## Effects

| Effect | Description |
| --- | --- |
| `JUIJN-B01` | `ParseJUnit`: A report whose root lists suites is read into one result per case of every suite. |
| `JUIJN-B02` | A report whose root is a single suite is read the same way. |
| `JUIJN-B03` | Suites nested inside suites are flattened: their cases are read as well. |
| `JUIJN-B04` | A case without its own file takes its suite's file; a case that names its file keeps it. |
| `JUIJN-B05` | A case with a failure or an error is failed; a case marked skipped is skipped. |
| `JUIJN-B06` | The proven codes are the codes named by the cases that neither failed nor were skipped. |
| `JUIJN-B07` | The seen codes are the codes named by every case, passed, failed or skipped. |
| `JUIJN-B08` | `CodesInCase`: Every code a case's name mentions is extracted, in the vocabulary the project declared. |
| `JUIJN-B09` | A file that is not a JUnit report is read as a report with no cases: nothing is proven by it. |
| `JUIJN-B10` | Each case carries its run time in seconds from its `time` attribute; a missing, malformed or negative time reads as 0. |
| `JUIJN-B12` | A case naming a visual-regression code of a state — `BUTTN-VR-S01` — also proves (and is seen for) the state it captures, `BUTTN-S01`: a green capture shows the state under its condition, looking as the spec says. |
| `JUIJN-B11` | The codes a report proves and sees carry the variant a case names (`CODE-B02#02`): a case of one variant never stands for its sibling. (`ScenarioCodesInCase`, `PassedCodes`, `SeenCodes`) |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the report file | a JUnit XML file, with a suites or a suite root | any other file, which yields no case | the project, declaring the report path in its configuration |
| the case name | any text, with or without scenario codes | — | this unit extracts only what the grammar recognizes |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `JUIJN-I01` | Every proven code is also a seen code: a run cannot prove what it did not measure. | reads a report of passed, failed and skipped cases and verifies the proven codes are contained in the seen codes |
| `JUIJN-I02` | A scenario is proven only by a case that ran, passed, and names exactly that scenario: a skipped or failed case proves nothing, and a variant is never proven by its sibling or by its rule. | reads a report where one variant passes and its sibling is skipped, and verifies only the passing variant is proven |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `JUIJN-X01` | Only the case's name carries codes: a code in the suite's or the class's name proves nothing. | A suite groups many cases; a code on the group would count every case of it as proof of one rule. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `JUIJN-E01` | The report file cannot be read. | The read error is returned and no report is produced. | A missing report must not look like a run with no cases. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/testsig/code.go` | `mustCodeRE` | infra — the scenario code grammar, in the project's vocabulary |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
