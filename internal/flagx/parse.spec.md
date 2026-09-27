<!-- @anchors
  code: FLPRF
  updated_at: 2026-09-26
  layer: apoio
-->
# FlagParse — reading a project's flag files and the scenarios their values open

> **Code**: `FLPRF`

## Overview

A project declares its feature flags in `flags/<name>.flag.md`, one file per flag, outside the code
tree. Each file holds a table of scenarios: a scenario code with the letter G, the condition under which
it applies (read with the fixed grammar, `FLGRF`), and what holds then. This unit reads those files so
the gates can confront each scenario against the code and the tests that declare themselves gated by it.

A scenario is a table row only, because its three parts are the columns of one statement. A row whose
condition the grammar refuses is kept, with the refusal attached, so it becomes a finding instead of
vanishing. Each scenario records its line so a verdict can point at it. The flag answers one question
of its own: does it declare what happens when it does not answer (the absent case), the case that most
often breaks and that almost nobody writes.

A project without `flags/` has simply not declared a flag yet, and the flags are always returned in the
same order, so two scans of the same repository report the same way.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the flag file | markdown with a table of scenario rows | scenarios written as prose or headings | this unit: only table rows with a G code are scenarios |
| the project root | a folder with or without `flags/` | — | this unit: no folder means no flags |

## Effects

| Effect | Description |
| --- | --- |
| `FLPRF-B01` | The scenarios of a flag file are its table rows whose code carries the letter G, each with its code, its parsed condition and what holds then; rows of other letters are ignored. (`ParseFile`, `ParseContent`) |
| `FLPRF-B02` | The flag's code is the unit part of its first scenario code, and its name is the file name without `.flag.md`. |
| `FLPRF-B03` | Each scenario records the line of the file it was read from. |
| `FLPRF-B04` | A row whose condition the grammar refuses is kept as a scenario carrying the refusal. |
| `FLPRF-B05` | A flag declares the absent case when one of its scenarios has the absent condition. |
| `FLPRF-B06` | Loading a project reads every `.flag.md` in `flags/`, sorted by file name, and gives nothing, without error, when `flags/` does not exist. |
| `FLPRF-B07` | The scenarios of all flags can be indexed by their code. (`ByCode`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FLPRF-E01` | REF[FLPRF-B04]: a condition the grammar refuses is the failure B04 keeps as a finding | — | — |
| `FLPRF-E02` | `flags/` exists but cannot be read as a folder, or a flag file in it cannot be read. | Loading returns the read error and no flags; only a `flags/` that does not exist gives nothing without error (B06). | Skipping it would make its scenarios vanish, and every `@gated-by` citing them would be accused of pointing at nothing. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/flagx/grammar.go` | `Parse` | apoio — the condition grammar (`FLGRF`) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
