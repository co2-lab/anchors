<!-- @anchors
  code: OPQSP
  updated_at: 2026-09-26
  layer: gate
-->
# OpenQuestions — a spec with an open question is not ready to implement

> **Code**: `OPQSP`

## Overview

Confronts a spec against the decisions it has NOT yet taken, and keeps it out of "ready"
while there is an open question.

The defect class is UNRESOLVED AMBIGUITY — the cheapest to avoid and the most expensive
to discover late. The path is always the same: the spec does not decide something the code
needs; whoever implements picks a defensible reading and moves on; the choice is never
confronted with the one who had the answer; the product ships with the wrong reading. No other
gate catches it, because all the pieces exist and reference one another — the defect is a decision
nobody took.

What this gate adds to the advice "don't guess, record and report" is a declared
PLACE for the record. Without a place, recording becomes a PR comment that dies in the merge.
With a place, the question is a visible work item, and the spec only becomes implementable when
the section empties.

The intended cycle: whoever writes notices what they do not know and writes it in the section; the gate accuses
while there is an item; the question is taken to whoever decides; the answer BECOMES A RULE, with
a code, and the item leaves the section.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any map node | — (the gate does not choose the target) | the gate engine, which routes by the declared `on:` |
| the decisions section | the title in any variation the catalogue and the project name, including `## Open Decisions`, `## Open Questions`, and `## Decisões em Aberto` | a title neither of the two names | this unit, by the vocabulary of accepted titles |
| the section's content | catalogued items, prose, or nothing | — (an absent section is a case, not an error) | this unit: whoever did not open the section is not demanded |
| the title's lexicon | declared in the project's Structure | — (omitted falls back to the framework's) | the Structure, with the framework's catalogue as the floor |

## Effects

| Effect | Description |
| --- | --- |
| `OPQSP-B01` | An artifact that is not a spec leaves without a verdict: only the spec has an open decision to demand. |
| `OPQSP-B02` | Whoever OPENED the section is confronted by its content: an open item blocks, a closed section releases. |
| `OPQSP-B03` | An open item BLOCKS: while there is a question, the spec does not pass as ready. |
| `OPQSP-B04` | A section closed honestly — opened under accepted titles such as `## Open Decisions`, `## Open Questions`, or `## Decisões em Aberto` and with no item — releases. Saying "there is no question" is different from not having looked. |
| `OPQSP-B05` | An item marked as RESOLVED does not block: the question stays in the trace, and what closed it is the rule that was born from it. |
| `OPQSP-B06` | Every question needs a CODE. Without identity it does not become a traceable item nor survive a rewrite of the spec. |
| `OPQSP-B07` | `OpenDecisions` COUNTS a spec's pending decisions, for whoever needs the number instead of the verdict — it is what allows reporting the pendency as a systemic lead, in the same standing as an absent signal. The count reads the project's lexicon by the same route as the confrontation: counting zero in a spec whose section is called something else would assert "there is no pending decision" about a spec full of them, which is the silence this unit exists to eliminate. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `OPQSP-I01` | Prose is not an item. Explanatory text inside the section does not count as a question — otherwise the author would learn to explain nothing. | writes prose in the section with no catalogued item and verifies that it does not block |
| `OPQSP-I02` | The section's boundary is respected: what comes after it is not read as a question. Without that, the whole spec would turn into a decisions section. | writes items in a following section and verifies that only those in the section count |
| `OPQSP-I03` | The column that says what the question BECOMES is not its identity, and filling it does not close the question. They are two things: the foreseen destination and the answer given. | fills the destination column without resolving and verifies that it still blocks |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `OPQSP-X01` | Does not judge whether the question is GOOD nor whether the answer is right. | The ruler is deterministic: either an open item exists, or it does not. Evaluating the merit of a doubt is judgement, and judgement belongs to another class of gate. |
| `OPQSP-X02` | Does not FAIL the spec that does not have the section — it records the pendency and says how to close it. | The absence does not distinguish "everything was decided" from "the section was deleted", and the two ask for opposite things. Failing would be treating migration as a defect; going quiet would be the silence the gate exists to eliminate. The verdict stays undetermined and TEACHES the way out: close with the declaration that there is no question, or write down what was not decided. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — the section's title may come from the project's lexicon, and the gate reads the Structure to find out |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
