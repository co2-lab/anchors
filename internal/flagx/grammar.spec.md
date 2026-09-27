<!-- @anchors
  code: FLGRF
  updated_at: 2026-09-26
  layer: apoio
-->
# FlagGrammar — the fixed grammar of a feature-flag scenario's condition

> **Code**: `FLGRF`

## Overview

A feature-flag scenario states WHEN it applies: the flag's value is equal to something, above a
threshold, contains a word, or the flag does not answer at all. Prose in that column covers any case and
can be confronted with none, so the condition follows a fixed grammar that refuses what it did not
foresee, and refuses loudly, saying what it expected.

The operators are the union of what the market's flag tools offer (LaunchDarkly, Unleash, Flagsmith).
The tools disagree only on spelling, so the grammar keeps one operator per meaning and accepts every
house style as an alias: a project writing `>=` and one writing `NUM_GTE` make the same statement. The
absent case (the flag service is down, a new environment, a local test) is the one that breaks in
production and the one nobody writes, so "absent" and "present" are first-class operators that compare
against nothing, and their words come from the translation catalog, so a flag file written in any
supported language is read. Operators that evaluate against an external service (segments, modulo) are
deliberately not accepted: a scenario that depends on one is written as the value it produces.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the condition cell | a comparison (operator and operand) or an absent/present word, optionally in backticks | prose, service-dependent operators, an operator without operand | this unit: it is refused with a specific error |

## Effects

| Effect | Description |
| --- | --- |
| `FLGRF-B01` | The symbolic, LaunchDarkly-style and Flagsmith/Unleash-style spellings of a comparison give the same operator, ignoring case, spaces, underscores and hyphens. (`Parse`) |
| `FLGRF-B02` | The operand loses one layer of surrounding double or single quotes, and keeps any inner quotes. |
| `FLGRF-B03` | The absent and present cases, in their English forms or in the word the catalog gives in any supported language, compare against nothing and carry no operand. |
| `FLGRF-B04` | The longest operator wins: a two-character symbol before its one-character prefix, and the longest word prefix before a shorter one. |
| `FLGRF-B05` | Only the absent and present operators need no operand. |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `FLGRF-X01` | Does not accept operators that evaluate against a flag service (segment match, modulo), nor prose. | A segment lives in the flag service's database, not in the repository; the grammar must not pretend to confront what it cannot see. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `FLGRF-E01` | The condition is empty or blank. | Refused: a scenario states when it applies. | An empty condition applies to nothing and would hide a missing statement. |
| `FLGRF-E02` | An operator has nothing to compare against. | Refused, naming the operator. | A half-written comparison would be read as some other statement. |
| `FLGRF-E03` | The condition is not a known comparison nor an absent/present word. | Refused, quoting the condition and showing the forms that are accepted. | A grammar that refuses without saying what it wanted moves the guessing to whoever writes the scenario. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/i18n/i18n.go` | `AllTranslations` | apoio — the absent/present words in every language (`INCTA`) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
