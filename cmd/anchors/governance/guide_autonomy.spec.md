<!-- @anchors
  code: ATGDT
  updated_at: 2026-09-26
  layer: comando
-->
# AutonomyGuide — what an agent does with what it does not know, by the role declared locally

> **Code**: `ATGDT`

## Overview

The autonomy section is appended to the review and work guides. It tells the agent what to do when it does not know, and it CHANGES with the role declared in the project's local settings (`.anchors/settings.yaml`).

The reason it changes: `settings user-issues` closes one door — the claim does not hand an escalated card to whoever does not decide the product. But the door used most has no label: the agent ASKING the developer who is running it. The developer knows the code and answers; the answer is reasonable, and it becomes a product decision taken by someone with no authority, with no plan, no review and no trace that it was decided there. The escalation exists for that, and the difference is the RECORD: an issue stays, has an owner, and whoever decides reads it when they can.

So whoever does not decide the product reads a different instruction, in the place where it matters, and not a warning at the end. Whoever decides reads that even their decisions are written, not discussed. Every profile also reads that preparing the environment asks no authorization — the preparation commands are idempotent, and the ruler for what does ask is reversibility — because an agent without that line either asks permission for everything or for nothing.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the root | the directory whose local settings declare the role | a directory with no settings, or with settings that cannot be read | this unit: both read as "no role declared" |
| the declared role | a known role, with or without a lens, that decides the product or not | — | the settings package answers whether the role decides and what its lens and title are |

## Effects

| Effect | Description |
| --- | --- |
| `ATGDT-B01` | A role that decides the product reads "Your role (<title>) decides the direction of this product", and that the decision is still written through `anchors escalate … --for-user` so it stays recorded. |
| `ATGDT-B02` | A declared role that does not decide reads "Your role (<title>) does NOT decide the direction of this product", and that whoever decides is the `product-owner` or the `architect`. |
| `ATGDT-B03` | With no role declared, the section says "You did not declare a role", claims no declaration, and gives the closed instructions. |
| `ATGDT-B04` | Every profile, and a project with no role declared, reads "Preparing the environment does not ask for authorization": the idempotent preparation commands (`doctor --fix`, `settings role <role> --date`, `map build`), that `settings role` takes its arguments so an agent is not left waiting, and that REVERSIBILITY is the ruler for what DOES ask. |
| `ATGDT-B05` | A role with a lens reads "This role's lens (<title>)" followed by the lens; a role with no lens gets no such section. |
| `ATGDT-B06` | Whoever does not decide the product — declared or not — reads to escalate with `--for-user`, "Do not ask whoever is running you", to move on to the next card, and what is NOT to be escalated. |
| `ATGDT-B07` | A declaration by the old `user_issues` flag, with no role, is named as such — "Your declaration (the old `user_issues` flag, with no role yet …)" — saying whether it decides the direction of the product and to declare a role; no empty role is ever named. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `ATGDT-I01` | A declaration that cannot be read never opens the door: the section is the one for no role declared. | a settings file with malformed content yields "You did not declare a role" and the ban on asking |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `ATGDT-X01` | The role that decides the product is not told to stop asking. | The ban protects the product from decisions without authority; it has no object for the role that holds that authority. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `ATGDT-E01` | REF[ATGDT-I01]: a settings file that cannot be read is the one failure the unit handles, and I01 states how: it falls on the closed side | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/settings/settings.go` | `Load`, `HandlesUserIssues`, `Decided` | the local declaration and whether it decides the product |
| DEP2 | `internal/settings/roles.go` | `Title`, `Lens` | the role's name and review lens |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
