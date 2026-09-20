<!-- @anchors
  code: SAIDA
  updated_at: 2026-09-20
-->
# Declared ways out — how a gate goes quiet, and what the silence means

> **Code**: `SAIDA`

## Overview

Every Anchors gate needs a way out: there is always the legitimate case the ruler did not
foresee, and a gate with no honest exit is a gate the team learns to work around — by
turning it off, ignoring it, or writing what it wants to hear.

But "silencing the gate" is not one thing. It is **two different assertions**, and
confusing them costs the gate the most valuable information it holds: the difference
between *I decided it is not needed* and *I have not done it yet*.

This doctrine is cross-cutting by construction. It belongs to neither `triad-complete` nor
`plan-doctrine-exists` nor `doctrine-realized` — it holds for all three, and for every
gate that comes to offer a way out. Until now it lived duplicated in each one's comments,
which is exactly the divergence the product axis exists to end.

## Rules

### SAIDA-R01 — the PERMANENT waiver is declared with `@no-<thing>: <reason>`

Asserts that the requirement will **never** apply to this unit: the code is pure
configuration, the behaviour is not observable, the proof lives elsewhere. It is a
decision taken, and the gate passes — `Pass`, for good.

### SAIDA-R02 — debt is declared with `@TBD: <what> — <reason>`

Asserts the piece **does not exist yet**, and that someone will write it. It is pending
work, not a decision: the gate returns `Pending`, and the finding keeps showing up until
it is paid.

### SAIDA-R03 — debt NEVER becomes `Pass`

Measured in this repository before the fix: `triad-complete` threw `@TBD` into the same
bucket as `@no-*`, and a spec with `@TBD: code,feature,test` came out **green**,
indistinguishable from a complete triad. The work that remained disappeared from the radar
because of the honest declaration of whoever assumed it — the worst possible incentive.

### SAIDA-R04 — a bare marker waives nothing    @TBD: the reason is checked per gate, and no spec catalogues it as its own rule yet

`@no-code` on its own, with no `:` and no written reason, is not a way out: it is the
silence the gates exist to end. The reason is mandatory in both markers, and it is checked
by pattern — whoever reads the unit finds there why it is waived, without hunting for the
decision somewhere else.

### SAIDA-R05 — the way out holds where it is written    @TBD: holds today by construction (each gate reads its own line), and no spec asserts it

A marker on a line waives that line; in the header, the unit. There is no waiver that
holds for the whole project written in a distant file — the decision stays where whoever
reads will find it.

### SAIDA-R06 — a marker inside backticks is a MENTION, not a declaration    @TBD: implemented in `triad-complete` and in the doctrine axis, not yet catalogued as a spec rule

A revision explaining the removal of a waiver cites the marker (*"the `@TBD: code` waiver
is gone"*), and without this distinction the citation **reactivates** the waiver the text
says has ended. An active marker is never inside backticks.

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SAIDA-X01` @TBD: a constraint nobody has had to assert yet | No gate infers a way out from the wording of a message. | The way out is a declared field or marker, never a recognised phrase — deducing intent from prose ages at the first rewrite of the message. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
