<!-- @anchors
  code: OBHNB
  updated_at: 2026-09-26
  layer: gate
-->
# ObligationHonored — the cross-cutting duty that lives OUTSIDE the unit

> **Code**: `OBHNB`

## Overview

Confronts the CROSS-CUTTING OBLIGATIONS a project declares against reality: **a node whose
header carries the trigger attribute must appear in the files the obligation demands.**

It is the class of defect no per-unit spec catches, because the duty lives OUTSIDE the
unit. The case that motivated it is real: a new data model, carrying free-text annotations
written by the user, was left out of the account-deletion script — personal data that would
never be erased. The project had been bitten by this before, and the script itself carries a
comment calling it a "silent LGPD violation". It repeated, and none of the nine existing
gates saw it, because every one of them looks INSIDE the unit.

**The honest exception, and the third state.** A node can waive itself with
`obligation_waived: <name> — <reason>` in the header, and the reason is MANDATORY: a waiver
without a justification is treated as absent. But waiving is not the only real case — the
most common one is that the duty is REAL and will be paid in another phase, when the handler
that consumes the table does not exist yet. Waiving would be a lie (the duty did not stop
existing), and leaving it red confuses ACKNOWLEDGED DEBT with forgetfulness, which is exactly
the distinction the pillar exists to preserve. `obligation_pending: <name> — <when>` asserts
three things — that the duty is known, that it still holds, and when it will be paid — and
the verdict is Pending, visible in the report and never Pass.

What separates it from its neighbours: `doc-required` charges a document the Structure
declared for a LAYER; this one charges the presence of a TOKEN in the files a named duty
points at, triggered by an attribute the node itself declares.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the declared obligations | what the project's Structure declares, or nothing | — (no obligation is a case, not an error) | the Structure: what is not declared is not charged |
| the trigger attribute | a `key: value` pair present in the node's header | a mention of the same text in the document's body | this unit: only the beginning of the file counts as header |
| the declaration of waiver or debt | the marker followed by the obligation name and a WRITTEN reason after an em dash or a spaced hyphen | a bare marker, and a hyphen glued to the name | this unit: without a reason the declaration is worth nothing, and a glued hyphen would let the obligation's own name supply the separator |
| the destination files | whatever the obligation's globs match | — (a glob matching nothing is a case, not an error) | this unit: with no file to read there is no violation to declare |

## Effects

| Effect | Description |
| --- | --- |
| `OBHNB-B01` | A node that carries the trigger and does not appear in the demanded file fails, and the verdict carries the declared REASON for the duty. |
| `OBHNB-B02` | A node that carries the trigger and does appear passes — the duty is fulfilled. |
| `OBHNB-B03` | A node without the trigger attribute contracts no obligation: the duty is charged by what the node declares about itself, not by what it might resemble. |
| `OBHNB-B04` | A waiver WITH a written reason exempts the node; the same waiver without one does not, because that is what separates the honest exception from silence. |
| `OBHNB-B05` | A project that declares no obligation is skipped: there is nothing to confront, and inventing duties would charge what nobody committed to. |
| `OBHNB-B06` | An acknowledged DEBT, with the when written down, yields Pending — it is a record, visible in the report, never an exemption. |
| `OBHNB-B07` | A bare debt marker, with no when, keeps failing: it assumes no debt, it only hides better. |
| `OBHNB-B08` | Waiver and debt stay distinct: only the waiver resolves the duty, because the debt is still owed. |
| `OBHNB-B09` | The failing verdict OFFERS the three ways out — fulfil, waive with a reason, or acknowledge the debt with a when. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `OBHNB-I01` | The token searched for is derived through the declared form — `screaming-snake`, `snake`, `kebab`, a free template, or the raw name. Guessing the shape in the engine would put one project's mess inside the framework. | applies each declared form to the same node name and verifies the resulting token |
| `OBHNB-I02` | A glob that matches no file produces no violation. Accusing where there was nothing to read would stamp what was never measured. | points the obligation at a glob matching nothing and verifies that nothing is accused |
| `OBHNB-I03` | The node's own `identified_as` wins over the obligation's automatic form. It is the only source that knows the project's real irregularity — one model becomes a plural env var, a sibling becomes a singular one, with no derivable rule; inverting this order accuses 28 correct models (measured). | declares `identified_as` against a form that would derive another token and verifies which one is searched for |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `OBHNB-X01` | Does not decide WHICH obligations exist, nor which files satisfy them. | That is the project's decision, declared in the Structure. A gate that invented cross-cutting duties would charge what nobody committed to, and the team would turn it off. |
| `OBHNB-X02` | Does not understand what the destination file DOES with the token. | The ruler is presence, which is deterministic. A purge script that names the table and then never erases it passes here — separating "forgotten" from "remembered" is already the defect this gate was built for; judging the implementation is another ruler. |
| `OBHNB-X03` | Does not read a declaration written in the body of the document. | Only the beginning of the file counts as header. Without that cut, a mention in the prose — an example, a quotation — would waive an obligation nobody meant to waive. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `Node` | core — the node's identity is what the token is derived from |
| DEP2 | `internal/config/config.go` | `Config` | core — the obligations, their triggers and their destinations are declared in the Structure |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
