<!-- @anchors
  code: LCINL
  updated_at: 2026-09-26
  layer: infra
-->
# LcovIngest — line coverage per file, and the uncovered lines of a change, read from an lcov report

> **Code**: `LCINL`

## Overview

Line coverage is read, not measured: many tools write the lcov text format, and reading it keeps Anchors
agnostic of the stack. This unit reads an lcov report into one coverage per source file — the covered and
total line counts, and, per instrumented line, whether it ran — and answers the diff coverage question when
crossed with the lines a change added: which of those lines are instrumented, and which of them no test ran.

A line the tool did not instrument (a comment, a blank line) is not executable, so it never counts against a
change, neither as uncovered nor in the denominator.

## Effects

| Effect | Description |
| --- | --- |
| `LCINL-B01` | `ParseLCOV`: Each record, opened by its source-file line and closed by the end-of-record line, becomes one file's coverage, in report order. |
| `LCINL-B02` | A line-hit entry with a positive hit count marks that line covered, and one with zero hits marks it uncovered. |
| `LCINL-B03` | When the record states its line totals, those totals win over the count of line-hit entries. |
| `LCINL-B04` | The uncovered lines of a change are the changed lines that are instrumented and not covered. |
| `LCINL-B05` | The instrumented count of a change is the number of changed lines the report instrumented. |
| `LCINL-B06` | The coverage percentage is covered over total lines times one hundred, and zero when the file has no line. |
| `LCINL-B07` | A record left without its end-of-record line is closed by the next source-file line or by the end of the report. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the report file | an lcov text report | any other format, whose lines are ignored | the project, declaring the report path in its configuration |
| the changed lines | the added lines of a file, as the diff reading gives them | — | the diff reading |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `LCINL-I01` | An uncovered changed line is always an instrumented changed line: the uncovered lines never exceed the instrumented count. | crosses changed lines — covered, uncovered and not instrumented — with a file's coverage and compares both answers |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `LCINL-X01` | Only line coverage is read: branch and function entries do not change the line counts. | The gates that consume the report ask about lines; mixing branch counts in would change the denominator they reason about. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `LCINL-E01` | The report file cannot be opened. | The error is returned and no report is produced. | A missing report must not look like a project with no instrumented line. |

## Dependencies

none — the unit reads a text file and depends on nothing of the project.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
