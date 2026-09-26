<!-- @anchors
  code: PSDPL
  updated_at: 2026-09-26
  layer: gate
-->
# PlanSourceDeclared — a plan that NAMES a source has to declare who builds it

> **Code**: `PSDPL`

## Overview

Confronts a plan against the dependency it wrote in PROSE and never declared in `needs:`:
**the plan names a source — who is going to build the adapter for it?**

The real case, measured in blue-eyes. Plan 0008 (Frontend/Web) said `Fonte: **GA4**, and
it is the architectural exception of the project`, and declared only
`needs: plans/0005-home-e-indice.md`. The GA4 adapter came from plan 0002, and **that
dependency existed only in the prose.**

What happened: 0002 was revised (`PLTFR-R0002`) and `Ga4Adapter` was REMOVED, with a
correct argument — no document of the project sustained it. 0008 stayed intact,
depending on a source nobody was going to build any more. Both plans remained internally
coherent, and the contradiction only surfaced months later, when 0008 was started.

**Why no gate caught it, and what this one separates from its neighbours.**
`dependency-honored` confronts the `needs:` that is DECLARED; here the defect is the
`needs:` that is MISSING. `plan-seeds-valid` looks at what the plan promises to CREATE,
not at what it promises to CONSUME. This is the version, between PLANS, of step 5 of the
review guide ("the prose ages"): nobody touched 0008, and it became wrong anyway —
because another document changed.

**What it measures.** A source line (`Fonte:`/`Fontes:`) names one or more sources in
bold. For each one, the gate looks for the corresponding adapter among the seeds of ALL
plans. If the adapter exists in another plan and this one does not declare it in
`needs:`, it fails.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the confronted artifact | any map node | — (the gate does not choose the target) | the gate engine, routing by the declared `on:` — a node that is not a plan leaves without a verdict |
| the plan's content | any text, with or without a source line | — (no source line is a case, not an error) | this unit: with no source line the confrontation is skipped, never approved |
| the named source | a name in bold on the source line | prose around the bold, which explains the source and is not its name | this unit, by the bold convention the plans already use |
| the map | a built graph, or none | — | this unit: with no graph the verdict is Pending, never Pass |
| the adapter's owner | the plan that SEEDS a file ending in `Adapter.spec.md` | a seeded file outside that naming convention | the convention the three existing adapters follow; the gate fails towards the safe side when the name is off-pattern |

## Effects

| Effect | Description |
| --- | --- |
| `PSDPL-B01` | A source whose adapter another plan seeds, and which this plan does not declare in `needs:`, FAILS — it is the defect that existed only in the prose. |
| `PSDPL-B02` | The failing verdict names WHICH source and WHERE its adapter lives, so the fix is the declaration, not an investigation. |
| `PSDPL-B03` | With the owning plan declared in `needs:`, the gate passes: the prose and the declaration now say the same thing. |
| `PSDPL-B04` | Every source of the line is confronted on its own — one line may name several, and one declared does not cover the rest. |
| `PSDPL-B05` | The source name matches the adapter's file regardless of case and punctuation: `GA4` matches `Ga4Adapter.spec.md`. |
| `PSDPL-B06` | A source whose adapter NOBODY seeds is not charged: it may be the source of a future plan, and the gate cannot invent a dependency that does not exist. |
| `PSDPL-B07` | The plan that seeds the adapter itself is not charged — it does not depend on itself. |
| `PSDPL-B08` | A plan with no source line returns Skip: there is nothing to confront, and that is not approval. |
| `PSDPL-B09` | An artifact that is not a plan returns Skip: the gate has no jurisdiction over specs, code or features. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PSDPL-I01` | What was not measured is never approved. With no graph, and with no adapter seeded anywhere, the verdict is Pending — approving there would stamp a confrontation that never happened. | runs with a nil graph and with a graph carrying no adapter, and verifies neither returns Pass |
| `PSDPL-I02` | The ownership convention fails towards the SAFE side. A seeded file whose name is off the `…Adapter.spec.md` pattern owns nothing, and the gate charges nothing for it — a wrong accusation costs more than a missed one, because it teaches the reader to ignore the gate. | seeds a file off the pattern and verifies nothing is charged |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PSDPL-X01` | Does not confront the ORDER of the phases. | Whether the plan that builds the adapter comes before the plan that consumes it is another question, and `fase-ordenada` is the ruler for it. Answering it here would put the same rule in two places that would then diverge. |
| `PSDPL-X02` | Does not charge a source whose adapter nobody seeds. | It may be the source of a plan that does not exist yet. Charging it would demand a `needs:` pointing at nothing — and the gate would be asking for a lie instead of catching one. |
| `PSDPL-X03` | Does not read the plan's prose to understand WHAT the source is for. | The ruler is the bold name on the source line, which is deterministic. Judging whether the plan really consumes that source, or only mentions it, is interpretation — and interpretation is not what a blocking gate can hold. |

## Errors

Each failure the code handles is already stated as a rule of another letter; the rows below catalogue it as a failure and point at that rule.

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PSDPL-E01` | REF[PSDPL-I01]: with no map nothing was measured, and I01 answers Pending, never approval | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/mapx/model.go` | `KindPlan`, `EdgeSeeds` | core — jurisdiction comes from the node's KIND, and the adapter's owner from the seeding EDGE |
| DEP2 | `internal/config/config.go` | `Config` | core — the gate receives the project's Structure with the confronted node |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
