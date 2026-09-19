<!-- @anchors
  code: RLTYR
  updated_at: 2026-09-19
  layer: gate
-->
# RuleTypes — the rule VOCABULARY is extensible, but it must be DECLARED

> **Code**: `RLTYR`

## Overview

Confronts a spec against the alphabet that makes traceability possible: **each letter of a
rule code (`{CODE}-<letter><NN>`) is the initial of the term that names a section — and a
letter nobody declared is INVISIBLE.**

That is the worst kind of hole, and it is the reason the gate exists. The rule appears in
the spec, in the `.feature` and in the test, and even so `feature-test-match` does not see
it: the code regex does not match the letter. **It looks covered and is not** — a green
that certifies a link nothing ever traversed.

It confronts three things, in this order of gravity: an UNDECLARED LETTER used in the
file; a SECTION cataloguing rules under a title no letter claims; and a CONFLICT in the
vocabulary itself, where two different sections claim the same letter. A fourth finding —
a section the project declared as rule-cataloguing but filled WITHOUT a code — is
Pending rather than Fail, because it is the inverse gap and the ruler is opt-in.

With no vocabulary declared the gate does NOT go quiet: it confronts the CANONICAL
letters. It used to Skip, and the effect was a canonical gate — seeded by `init` in every
project — that measured nothing: it took a row in the `check` table and reported
indeterminate forever. That is the same silence Anchors fights everywhere else, a
declared gate that confronts nothing giving the impression of a defence that does not
exist.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the artifact's content | any text, empty included | — (text with no code at all is a case, not an error) | this unit: with no code there is nothing to confront, and that is Pass |
| the declared vocabulary | what the Structure declares, or nothing | — (no declaration falls back to the canonical letters, never to silence) | the Structure, and this unit for the fallback |
| the accepted code length | what the project declared in `code_lengths` | — | `internal/config/config.go`, consulted PER CALL and never frozen in a package variable |

## Effects

| Effect | Description |
| --- | --- |
| `RLTYR-B01` | A letter used in the spec and NOT declared in the vocabulary fails, and the verdict names the letter. |
| `RLTYR-B02` | A declared letter used under a claimed section passes. |
| `RLTYR-B03` | A section that CATALOGUES rules under a title no letter claims fails, and the verdict names the section. |
| `RLTYR-B04` | A CONFLICT in the vocabulary itself — the same letter claimed by two different terms — fails before the file is even read: each letter belongs to ONE term. |
| `RLTYR-B05` | With no vocabulary declared the gate confronts the CANONICAL letters instead of going quiet. |
| `RLTYR-B06` | A heading that IS the rule code itself is not a category section: it is the rule's own header, and is not charged. |
| `RLTYR-B07` | A section that merely CITES codes belonging to other sections does not catalogue rules, and claims no letter. |
| `RLTYR-B08` | A section that DEFINES a code in the first cell of a table does catalogue rules, and is charged. |
| `RLTYR-B09` | A section the project declared as `sections_require_code` that is FILLED and carries no code returns Pending, naming the section. |
| `RLTYR-B10` | A section whose table already carries the code is not charged. |
| `RLTYR-B11` | A declared section that is NOT in `sections_require_code` is not charged: it merely enumerates values, and demanding a rule of an index would invent a duty. |
| `RLTYR-B12` | A project that does not use `sections_require_code` changes no behaviour — the ruler is born opt-in, or it would accuse an entire existing base at once. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RLTYR-I01` | A spec with no rule code at all is not this gate's problem. Demanding a catalogue here would duplicate `spec-complete`, which is who charges the existence of a catalogued rule. | confronts a spec of pure prose and verifies Pass |
| `RLTYR-I02` | The verdict NAMES the letter AND where to declare it. A gate that fails without saying what transfers the diagnostic work to whoever reads it. | confronts an undeclared letter and verifies the letter and `rule_types` appear in the verdict |
| `RLTYR-I03` | A section without a code borrows a neighbour's. Measured in a real project: 48 specs with "Events / Callbacks" filled in and not one code, and the scenarios proving those events borrowed the code of the neighbouring state — an `-S` governing behaviour. | confronts a filled `requires_code` section with no code and verifies the finding |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RLTYR-X01` | Does not decide WHICH letters exist. | The vocabulary is extensible by design and belongs to the project. The gate charges that it be DECLARED, not that it be a particular set — the canonical letters are the fallback, not a ceiling. |
| `RLTYR-X02` | Without a declared vocabulary it charges only the letter, not sections or terms. | Sections and terms only exist once the project declares the vocabulary. Charging them against an implicit vocabulary would invent a rule nobody wrote. |
| `RLTYR-X03` | Does not judge whether the letter is the RIGHT one for that rule. | Whether a behaviour was catalogued as `-B` or as `-S` is editorial judgement. The ruler here is that the letter be recognisable by the traceability, which is deterministic. |
| `RLTYR-X04` | Does not charge format — only traceability. | A section with no code, or one that only cites other people's codes, is ignored. What the gate defends is that a code the tooling cannot see does not exist. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `RuleType` | core — the vocabulary, its sections and `sections_require_code` are declared in the Structure |
| DEP2 | `internal/config/config.go` | `DefaultRuleLetters` | core — the canonical letters, the ruler used when the project declares no vocabulary |
| DEP3 | `internal/config/config.go` | `CodeLengthPattern` | core — the accepted code length comes from the project and is read PER CALL, never frozen |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
