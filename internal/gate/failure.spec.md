<!-- @anchors
  code: FLRAI
  updated_at: 2026-10-08
  layer: gate
-->
# Failure — the failure a spec declares must be handled, recorded, and every handling declared

> **Code**: `FLRAI`

## Overview

A spec catalogues how its unit fails, with its own letter (`-E`) and a table of condition and
result. Until these gates, nothing confronted that catalogue with the code: a declared failure
crossed the whole pipeline without anyone asking whether the code handles it, whether the handling
records it, or whether a handling that exists answers any declared failure at all.

Three gates close that static layer, and none depends on production or on a log format:

- `failure-handled` asks whether the governed code has any path that handles a failure, when the
  spec declares failures;
- `failure-logged` asks whether that handling also records the occurrence;
- `failure-declared` is the inverse of the first, and catches the commonest case: somebody wrote a
  defence and never declared what it prevents.

Handling is not "having a catch". What a check that refuses, a recovered panic and a caught
exception share is the effect, not the syntax, so the project's dialect says what a handling and a
record look like in its language. Without those patterns the gates have measured nothing and say
so. The code is judged as a set: the code does not cite the failure's code, so no rule can be tied
to one specific path, and the gates assert the case that matters, a unit that declares failures
with no handling at all.

A failure can leave the charge only with knowledge written beside it: `@resilient: <reason>` says
"it happens, I know why, and the flow absorbs it". A bare marker exempts nothing. The same rows
carry the conclusions the observation layer reads: `@resilient` and `@observing: <what was ruled
out>`, each with its reason read whole.

A unit whose handling matches are all normal flow (a lazy map initialisation, a pattern that did
not match) closes its failure section with `none — <why>`, and that satisfies `failure-declared`.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | a node of kind `spec` | any other kind | this unit: every gate skips it |
| the declared failures | failure rules in the catalogued forms: a heading, a table row, a bold bullet | a failure code mentioned in prose | this unit: prose is not read as a declaration |
| the governed code | the files the spec specifies, reached through the map | a spec with no specified file, or files that cannot be read | this unit: the verdict stays Pending |
| the dialect | handling and recording patterns, declared by the project or its family | a project with no patterns | this unit: the verdict stays Pending |

## Effects

### Shared by the three gates

| Effect | Description |
| --- | --- |
| `FLRAI-B01` | Every failure gate skips an artifact that is not a spec. |
| `FLRAI-B02` | A failure rule is read in the three catalogued forms — heading, table row, bold bullet; a failure code cited in prose declares nothing. |
| `FLRAI-B03` | A spec that declares no failure is skipped by `failure-handled` and `failure-logged`. |
| `FLRAI-B04` | Without handling patterns in the dialect every failure gate is Pending, and `failure-logged` is also Pending without recording patterns: a Pass would stamp what was never measured. |
| `FLRAI-B05` | Without governed code that can be read — no map, no specifies edge, or no specified file on disk — every failure gate is Pending. |
| `FLRAI-B06` | The governed code is read without its comment lines, so a handling written only in a comment does not count. |

### failure-handled

| Effect | Description |
| --- | --- |
| `FLRAI-B07` | A spec that declares failures, over code with no handling path at all, fails, naming the charged failures in order. |
| `FLRAI-B08` | Any handling path in the governed code passes the gate. |
| `FLRAI-B09` | A failure marked `@resilient` with a written reason is not charged by `failure-handled` nor by `failure-logged`; when every declared failure is resilient, both pass. |
| `FLRAI-B10` | A bare `@resilient` marker, with no reason, exempts nothing. |

### failure-logged

| Effect | Description |
| --- | --- |
| `FLRAI-B11` | Code that handles but records nothing fails, naming the charged failures. |
| `FLRAI-B12` | Code that handles and records the occurrence passes. |

### failure-declared

| Effect | Description |
| --- | --- |
| `FLRAI-B13` | Code with handling paths and a spec that declares no failure fails. |
| `FLRAI-B14` | Code with no handling path is skipped: there is no defence to declare. |
| `FLRAI-B15` | Code with handling paths and a spec that declares at least one failure passes. |
| `FLRAI-B16` | A spec whose failure section opens with `none` and a written reason passes; a bare `none`, a `none` outside the section, or a section that opens with prose does not close it. |

### Conclusions of the observation layer

| Effect | Description |
| --- | --- |
| `FLRAI-B17` | For each declared failure (`FailureConclusions`), the reasons of `@resilient` and `@observing` are read whole; a failure with neither carries no conclusion. |
| `FLRAI-B18` | A conclusion's reason ends at its table cell: the next column is never read into it. |
| `FLRAI-B19` | A failure rule and its conclusion are read at the code lengths the project declares (`code_lengths`), not a fixed range: with a declared length of 7, a 7-character `-E` rule is a declared failure. |
| `FLRAI-B20` | `failure-declared`: each fallible source of a unit — a call of its code that a `dialect.fallible_patterns` entry recognises (comment lines and trailing comments do not count) — is named by a declared failure (`-E`): its row, or its row of the rules' uses, cites the name called as a whole word; a failure that does not name it ("not found") does not answer it. Unanswered, the gate fails naming each source by file and line, unless the Errors section is closed with `none — <reason>` or the spec waives with `@no-failure: <reason>`. (`fallibleSources`, `fallibleCalls`, `FallibleCall`, `FallibleSource`, `uncoveredSources`) |
| `FLRAI-B21` | A file of a layer the project marks `fallible: true` that the unit's code imports — by its `@dep:` flag — is a fallible source too, named by its path and answered by a failure citing the file's stem (`useBudget`) or a name the code imports from it; while a spec still declares the dependency in a table, its `DEPn` names it as well. A dependency on any other layer is not a source. (`FallibleSource`) |
| `FLRAI-B22` | `failure-handled`: each fallible call of the unit's code has its handling in its window — the call's statement from its first line (a destructuring above it that reads the error; the climb stops at a blank line or at one ending a statement, `;`, `}` or `)`), and the pattern's `window` after it (`DefaultFallibleWindow` when it declares none) —, matched by the pattern's `handled` or, when it declares none, the project's `handle_patterns`; a call with none fails, named by file, line and text, unless its line or the line above waives it with `@no-handle: <reason>`. (`unhandledCalls`) |
| `FLRAI-B23` | A fallible call whose result is the unit's own — `return`ed, or the body of an arrow (`=>`), with or without `await` — hands its failure to the caller: `failure-handled` does not charge it, and `failure-declared` still asks the unit's spec to name it. (`propagated`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `FLRAI-I01` | `failure-handled` and `failure-logged` never both fail the same spec: with no handling the first charges and the second steps aside, and with handling that records nothing it is the other way round. | runs both gates over code with no handling and over code that handles silently, and verifies Fail/Skip and Pass/Fail |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `FLRAI-X01` | Does not tie a declared failure to a specific handling path: one handling path answers for every declared failure. | The code does not cite the failure's code, so no rule can be matched to one path; what can be asserted is the unit with failures and no handling at all. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FLRAI-E01` | REF[FLRAI-B05]: a specified file that cannot be read is left out of the governed code, and when none can be read B05 answers with Pending | — | — |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
