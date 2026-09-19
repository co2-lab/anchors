<!-- @anchors
  code: CSDCN
  updated_at: 2026-09-19
  layer: gate
-->
# ContractStatusDeclared — the output contract lists the status codes the code really returns, and only those

> **Code**: `CSDCN`

## Overview

Confronts the `Output Contract` table of an INTERFACE spec — a handler, a route — against the
status codes the governed code actually emits: **the table was written, but does the code
still keep it?**

Why the gate exists, measured: in an audit of 51 spec-versus-code divergences in the
reference app (2026-08), this was the MOST REPEATED pattern — eight handlers declared a
contract the code did not honour. And the omitted status was almost always the SECURITY
one: the 403 of ownership, the 409 of conflict. Whoever writes the table thinks about the
happy path and about the "business" errors, not about the refusals of access.

**The two sides of the error are different, and both matter.** A status EMITTED and not
declared leaves the client — programmed from the table — unable to handle the refusal: the
user sees a generic error where there was a specific reason. That was `accept-org-invite`,
which declared 3 status codes and emitted 8; the two missing 403s and the 409 were the
defence against invite hijacking. A status DECLARED and never emitted is worse in another
way: it is dead code in the client, and it disappears with nobody noticing.
`reanalyse-metadata` declared 402 for exceeded quota and no path of the handler emits 402
(the quota answers 429) — a client treating 402 as "needs to pay" would never fire that
branch.

What separates it from its neighbours: `dependency-honored` confronts the symbols the spec
promises to consume; this one confronts the status codes the spec promises to return. And
unlike `spec-complete`, it never charges the EXISTENCE of the section — with no table there
is nothing to confront.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any node of the map | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` |
| the artifact's kind | only a spec is confronted | code, test, feature, plan | this unit: anything that is not a spec leaves without a verdict |
| the contract section | the title in the catalogue's languages (`Contrato de Saída`, `Output Contract`) | a section the catalogue does not name | this unit: without the section the confrontation is skipped, never approved |
| the declared status codes | the concrete three-digit numbers of the table | the generic ranges `4xx`/`5xx`, which are read as coverage and not as status | this unit: charging `500` because the spec says `5xx` would be a false positive |
| the lexicon that reads the code | `dialect.http_status` of the project's Structure | — (an absent dialect is a case, not an error) | the Structure: without the declared lexicon the verdict is Pending, never approval |
| the map | a built graph, or none | — | this unit: with no map there is no edge to the code, and the verdict is Pending |

## Effects

| Effect | Description |
| --- | --- |
| `CSDCN-B01` | A status EMITTED by the code and absent from the table fails, and the verdict names it — the client programmed from the table does not handle the refusal. |
| `CSDCN-B02` | A status DECLARED in the table and emitted by no path fails too, as dead code in the client that disappears unnoticed. |
| `CSDCN-B03` | A faithful table passes: the gate that only accuses is a noise generator, and nobody keeps one. |
| `CSDCN-B04` | The 500 of the top-level try/catch is not charged for absence: it is infrastructure every handler carries, not a decision of this one. |
| `CSDCN-B05` | A declared `5xx` covers the 5xx codes the code emits; a generic `4xx` covers nothing, because it would hide exactly the access refusals this gate hunts. |
| `CSDCN-B06` | A status that only appears inside a COMMENT is not an emitted status: the comment describes what the function used to do. |
| `CSDCN-B07` | Without the contract section the confrontation is skipped — charging the section's existence belongs to `spec-complete`. |
| `CSDCN-B08` | Code that returns no status at all is skipped: a cron or trigger handler has `void` as its contract, and there is nothing to confront. |
| `CSDCN-B09` | A status passed as a literal to a locally defined helper counts as emitted — `fail(400, …)` is a 400, wherever the envelope is built. |
| `CSDCN-B10` | When the code carries a DYNAMIC status — a helper that takes the code by parameter — the gate stops asserting the phantom side, and keeps charging the literals it did find. |
| `CSDCN-B11` | Without a declared `http_status` in the dialect the verdict is Pending, and names the known families — the meter does not fake conformity nor guess the stack. |
| `CSDCN-B12` | A project that explicitly waives the `http_status` field is skipped: the opt-out is declared, and it is honoured. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CSDCN-I01` | The lexicon that reads the status comes from the project's dialect, never from the gate. Embedding one stack's syntax would make the gate silent on every other one — and silence reads as conformity. | confronts the same defect in a Go `net/http` handler with `dialect.family: go` and verifies the same accusation |
| `CSDCN-I02` | A dialect declared by hand, with no family, teaches the gate its own lexicon. Without this the agnosticism would be a promise limited to the built-in families. | declares `http_status` as a raw pattern for a Rails handler and verifies the status is charged |
| `CSDCN-I03` | A named constant is worth the number it means: `http.StatusForbidden` is a declared 403, and reading only digits would approve every handler written with constants. | emits the status by named constant and verifies the verdict names the number |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CSDCN-X01` | Does not demand the generic ranges, and does not invent their semantics. | `4xx` and `5xx` are what the spec writes for "any failure". Treating them as status codes would require guessing which numbers they cover, and the guess would be charged as if it were a declaration. |
| `CSDCN-X02` | Does not judge WHEN each status is right — only whether the number appears on both sides. | Whether the 403 belongs on that branch is judgement about the design. Here the ruler is the correspondence between two sets of numbers, which is deterministic and does not depend on reading intent. |
| `CSDCN-X03` | Does not charge the phantom side when the code builds the status dynamically. | With a helper receiving the code by parameter, a declared value may well be emitted through a call textual reading cannot reach. Accusing would be a false positive, and mass false positives are what makes a team turn the gate off. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindSpec`, `EdgeSpecifies` | core — the kind gives jurisdiction, and the edge says which code realises the contract |
| DEP2 | `internal/config/config.go` | `KnownDialectFamilies` | core — the lexicon that reads the status comes from the Structure, and the Pending verdict names the families |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
