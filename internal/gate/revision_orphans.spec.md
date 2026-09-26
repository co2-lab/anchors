<!-- @anchors
  code: RVORP
  updated_at: 2026-09-26
  layer: gate
-->
# RevisionOrphans — the rules a revision changed the meaning of, without saying so

> **Code**: `RVORP`

## Overview

Confronts a revision declared inside a spec against the OTHER rules of the same unit: **when a revision rewrites a rule, the sibling rules that speak of the same thing must be named — either as revised too, or as read and still valid.**

A rule does not live alone. It shares vocabulary with its siblings, and it is that vocabulary a revision changes — not only the text of the rule it rewrites. A revision that abolishes a concept leaves every other rule still asserting it, and the spec then claims two contradictory things at once.

Measured in the reference app, and it is what produced this gate. The `NTCNN-R0002` changed the notification badge from a COUNT to a DOT, and named the rules it rewrote: `B03`, `B04`, `B07`. The invariant `I02` was not named — and it is titled *"the badge never COUNTS what the list does not show"*, with a body reading *"the NUMBER on the bell matches what appears on opening"*. The invariant governed arithmetic the revision had abolished.

**Nothing accused it.** The `B03` was correct, the `I02` was well-formed, the triad complete, the suite green. The contradiction surfaced MONTHS later, when another agent went to implement and could not tell which of the two to follow — and it became a decision that had to escalate to the user, with nobody left remembering the context. Seven contradictions of this exact shape surfaced in a single batch.

The ruler is ONE shared domain word, and the threshold was measured in both directions. Requiring two failed the very case that produced the gate: the `I02` shares exactly one word with what the revision rewrote — `badge` — and that word carries the whole contradiction. What makes one word enough is CLEANING the title, not counting: with negations and waiver comments in, the same spec accused three rules (`B01` and `B05` entered on "não" alone); with them out, it accuses one — the target.

A ubiquity filter was also tried, discarding terms appearing in more than two thirds of the rules. It discarded `lista` and `badge` — precisely the subject — and made the real case accuse nothing. In a well-written spec the domain vocabulary repeats on purpose: discarding what repeats is discarding the subject.

The ruler is CO-CITATION, not meaning. Asking whether two rules contradict each other requires reading them, and that is judgement — it would make this a judge, not a gate. What a machine decides alone is narrower and sufficient: which rules of this unit share the vocabulary of the rules the revision touched, and were not mentioned.

Leaving the finding is cheap, and deliberately so: `Checked: I02` asserts that somebody read it, not that it is correct. It is cheap to write AFTER reading and impossible to write honestly without reading, and that asymmetry is what makes the ruler work. Demanding that every revision check every rule of the unit would fail always on a nine-rule spec, and satisfying it would become theatre — the whole list pasted in unread.

This gate operates in distinct territory from neighbouring gates:
- Unlike `plan-revised`, which governs the relationship between two distinct PLANS across the map, this gate stays inside one spec, between a revision and the sibling rules of its own unit.
- Unlike `plan-change-justified`, which asks whether an edit carries a written justification, this gate asks whom that justification forgot.
- Unlike `spec-feature-match`, which confronts declared requirements against written scenarios, this gate never leaves the spec: both sides of its confrontation are rules of the same file.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | nodes of kind spec | nodes of kind plan, feature, code, or test | this unit: skips non-spec nodes |
| revisions | `{CODE}-R000N` blocks declaring `Revises:` | a spec with no revision | this unit: skips when the spec declares none |
| the revision keywords | `Revises:` / `Checked:`, in any language of the translation catalog | keywords hardcoded in one language | the translation catalog (`revision.keyword.*`) |
| sibling rules | the rules this spec DEFINES | codes the spec merely cites from other units | `definedRequirements`, shared with `spec-feature-match` |
| the vocabulary | significant terms of the rule titles | stopwords, negations, and waiver comments | this unit: cleans the title before comparing |

## Effects

| Effect | Description |
| --- | --- |
| `RVORP-B01` | When the confronted node is not a spec, the gate skips confrontation. |
| `RVORP-B02` | When the spec declares no revision, the gate skips confrontation. |
| `RVORP-B03` | When a revision declares no `Revises:`, the gate abstains with a pending verdict: the field is new, and 439 revisions written before it exist in the reference app — accusing all of them at once produces noise, not a queue. |
| `RVORP-B04` | When a revision names a rule the spec does not define, the gate fails, reporting the unknown code. |
| `RVORP-B05` | When a sibling rule shares significant vocabulary with a revised rule and appears in neither `Revises:` nor `Checked:`, the gate reports it as an orphan, naming the shared terms. |
| `RVORP-B06` | When every vocabulary-sharing sibling appears in `Revises:` or `Checked:`, the gate passes. |
| `RVORP-B07` | When a rule appears in `Checked:`, it leaves the accusation without asserting that it is correct — only that somebody read it. |
| `RVORP-B08` | The gate is of the BLOCKING class: a new project is born with it blocking, and an existing one takes it through the same maturation as every structural gate — informative until the project promotes it. It never depends on an ingested signal, so it can block from day one. |
| `RVORP-I01` | A rule never accuses itself: the revised rule is excluded from its own orphan candidates. |
| `RVORP-X01` | Negations and waiver comments are stripped from a rule title before comparison: they are not what the rule asserts. |

## Open Decisions

| Code | Question |
| --- | --- |
| `RVORP-Q01` <!-- @no-scenario: a decided question is history, not behaviour; the rule it became has the scenario --> | Should the gate block or inform? Decided by the user: it blocks — became `RVORP-B08`. |
