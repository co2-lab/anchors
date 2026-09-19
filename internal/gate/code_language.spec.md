<!-- @anchors
  code: CDLNG
  updated_at: 2026-09-19
  layer: gate
-->
# CodeLanguage — the code does not go back to mixing languages

> **Code**: `CDLNG`

## Overview

Defends a decision that, without a gate, undoes itself: **what is WRITTEN in the code is
English; what is READ is translated.**

Anchors was born with the code in Portuguese and migrated. An external contributor should not
need Portuguese to read a function's name. But that work is lost with a
single PR from someone who does not know the rule — and the rule is nowhere the compiler
reads. This gate is where it comes to be.

**What it looks at, and what it ignores on purpose:**

- IDENTIFIERS — yes. They are what the contributor needs to read in order to work.
- COMMENTS — no. They carry the measurements and the whys, in the team's language.
- USER-FACING TEXT — no. That goes through the translation catalogue, and charging it here would
  duplicate the ruler in two places that would diverge.

The detection **is not by dictionary**, and the reason was measured: comparing each word against the
system dictionary accused compound identifiers in English, which no common dictionary
has — 491 accusations for 47 real cases, 90% false positive. A gate like that is
turned off on the first day.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, which routes by the declared `on:` |
| the content | any text, empty included | — (absent text is a case, not an error) | this unit: without a declaration there is nothing to accuse |
| the candidate word | a declaration identifier, with length above the floor | a word too short to carry a language | this unit, by the length floor |

## Effects

| Effect | Description |
| --- | --- |
| `CDLNG-B01` | An identifier in Portuguese is ACCUSED; an identifier in English passes without noise. |
| `CDLNG-B02` | The verdict RETURNS the word that accused — without it, whoever reads looks for the needle in the whole file. |
| `CDLNG-B03` | Only a DECLARATION is the subject: what does not declare an identifier is not read. |
| `CDLNG-B04` | The declarations are truly found, in every form the language offers — not only in the most common one. `PortugueseIdentifiers` sweeps the content and returns what accused. |
| `CDLNG-B05` | `WordIsPortuguese` decides ONE word, and `IdentifierIsPortuguese` decides a whole identifier by breaking it into the words that compose it — it is the separation that lets the length floor hold per word, and not for the whole identifier. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CDLNG-I01` | A short word does not count. Below the length floor there is no language to infer, and accusing there would be noise — which is how a gate loses the trust of whoever reads it. | confronts short identifiers and verifies that none is accused |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CDLNG-X01` | Does not read a COMMENT. | The comment carries the measurement and the why, in the language of whoever decided. Translating it would lose what it exists to record, and the ruler is about what the contributor needs to READ in order to work. |
| `CDLNG-X02` | Does not read USER-FACING TEXT. | A message to the user goes through the translation catalogue, which already resolves it by the project's language. Charging it here would duplicate the ruler in two places that would diverge. |
| `CDLNG-X03` | Does not use a dictionary to decide the language. | Measured: the system dictionary accuses a compound identifier in English, and produced 90% false positive. A gate that is wrong nine times out of ten is turned off, and then it defends nothing. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindCode` | core — the gate only has jurisdiction over code |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
