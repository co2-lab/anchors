<!-- @anchors
  code: HDGDH
  updated_at: 2026-09-26
  layer: infra
-->
# HeaderGuide — render the project's header guide in the stack's comment dialect, passing the gate that init itself declares

> **Code**: `HDGDH`

## Overview

`init` seeds a header guide: the built-in ruler for the annotation block at the top of every file
(`anchors guide header`), instantiated for the project in front of it. The examples are written in the
comment dialect of the chosen stack — a Python or Ruby project reads `#`, a C-like stack reads `//` —
and the grouping example uses the project's real first module, so the reader sees their own code, not a
generic sample. The dialect only shapes the example: Anchors reads any comment at run time.

The seeded guide must pass the `guide-checklist` gate that the same `init` declares. A guide seeded by
`init` that fails `init`'s own gate is the worst first impression possible — measured in a real project,
it was the first blocking finding. So the guide always carries a compliance-points section, under the
title of the project's language, with five points checkable one by one. Beyond the gate, the points are
what makes the guide confrontable rather than merely read.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the preset | any catalogued preset, or none (an empty preset) | — | this unit: an unknown or empty preset gets the `//` dialect and the title "project" |
| the module names | the project's module names, possibly none | — | this unit: with none, the example module is `auth` and no module list is printed |
| the language of the project | any language of the translation catalog | — | the translation catalog, which carries the section title in every supported language |

## Effects

| Effect | Description |
| --- | --- |
| `HDGDH-B01` | In the rendered guide (`RenderHeaderGuide`), the examples use `#` comments for the django, fastapi, python-lib, rails and phoenix presets, and `//` comments for every other preset. |
| `HDGDH-B02` | The grouping example names the first module of the project; with no modules it names `auth`. |
| `HDGDH-B03` | The project's modules are listed, comma-separated, only when there are modules. |
| `HDGDH-B04` | The guide always states the essentials — the `code:` and `updated_at:` annotations, the `header-valid` gate and a pointer to `anchors guide header` — even with no preset. |
| `HDGDH-B05` | The guide always carries the compliance-points section under the title of the current language, with the points CK1 to CK5. |
| `HDGDH-B06` | The guide is titled after the preset's title, or "project" when the preset has none. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `HDGDH-I01` | The guide seeded in any supported language is recognised by the compliance-section heading the `guide-checklist` gate looks for, with at least one checkable point. | renders the guide in every supported language and matches it against the gate's heading, built from all translations of the title |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `HDGDH-X01` | Does not write the guide to disk; it returns the text for `init` to write. | Rendering stays pure, so every dialect and language is provable without a project folder. |

## Errors

none — rendering only concatenates text from its inputs and the translation catalog; there is nothing that can fail.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/i18n/i18n.go` | `TIn`, `Current` | core — the compliance-section title in the project's language |
| DEP2 | `internal/initx/presets.go` | `Preset` | infra — the stack's name and title |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
