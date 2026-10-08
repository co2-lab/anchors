<!-- @anchors
  code: CNGDC
  updated_at: 2026-10-08
  layer: infra
-->
# ContributingGuide — render the project's CONTRIBUTING.md from the configuration init writes

> **Code**: `CNGDC`

## Overview

Whoever arrives at a project under Anchors — a person or an agent — does not know that the spec comes first,
which command gives the next task, or which gate bars a commit. All of that is already decided in
`anchors.yaml`, in a form a newcomer does not read. `init` seeds a `CONTRIBUTING.md` that says it in prose,
rendered from the configuration it just wrote, so the guide names this project's layers, artifacts and
blocking gates rather than an example's.

The guide proposes no structure. The code layers it lists are the ones the project declared; the layer kinds
it names (entry points, use cases, domain, repositories, infrastructure, presentation) are illustration of
keeping layers apart, never a layout to move files into. When the project already has a `CONTRIBUTING.md`,
init never touches it: the section this unit renders is what init shows for the project to add by hand.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the configuration | the configuration init wrote, possibly with no layer, no gate and no workflow, or none | — | this unit: a missing part is said to be missing, never invented |
| the guide folder | the project's guide folder, or none | — | this unit: with none, no folder is named |

## Effects

| Effect | Description |
| --- | --- |
| `CNGDC-B01` | The whole guide (`RenderContributing`) is a `# Contributing` title, a paragraph saying it was seeded from `anchors.yaml`, and the section (`ContributingSection`). |
| `CNGDC-B02` | The section opens with `## Working with Anchors` and states the order of the work: spec, then feature, then test, then code, and that a bug fix needs its rule and scenario. |
| `CNGDC-B03` | The section lists where each piece lives: the pattern of every spec, feature and test layer, in that order, and the colocation templates when the project declares colocation. |
| `CNGDC-B04` | The section lists every declared code layer with its pattern, by name; with none, it says no code layer is declared yet. The layer kinds it names are given as illustration. |
| `CNGDC-B05` | The daily commands name the GitHub repository as the queue in `github` mode, and the local queue otherwise. |
| `CNGDC-B06` | The gates whose failure blocks are named as barring the commit and the rest as informing; with no blocking gate, it says every gate only informs. Divergences and pending items are said to inform unless `severity` says otherwise. |
| `CNGDC-B07` | The guides part names the project's guide folder when there is one. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CNGDC-I01` | The guide never names a layer, gate or folder the configuration does not hold. | renders a configuration with one layer and one gate and finds no other name in the layer and gate lists |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CNGDC-X01` | Does not write the guide to disk; it returns the text for `init` to write or show. | Whether to write is init's decision — an existing guide is the project's — and rendering stays pure, provable without a folder. |

## Errors

none — rendering only concatenates text from the configuration; there is nothing that can fail.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
