<!-- @anchors
  code: DCNAV
  updated_at: 2026-10-07
  layer: doct
-->
# Navigation — the app's navigation map and the code's dependency chain, compiled into pages

> **Code**: `DCNAV`

## Overview

The flags written beside the code — `@navigates:` on each navigation call, `@dep:` on each
import line — become edges of the map (`navigates-to`, `depends-on`), and the chain's gates
keep them equal to the code and to the specs (DESIGN-dependencies-and-navigation.md). This unit
compiles those edges into two pages a reader opens instead of the map:

- `{{ with navigation }}` — every screen (a spec of the `screen` layer), where it leads and by
  which rule, drawn as a Mermaid flowchart: an entry route (`navigation.entry`) as a stadium,
  a screen no entry reaches dashed;
- `{{ range dependencies }}` — every code file of the dependency chain, what it uses with the
  symbols, and who uses it.

`anchors docs init` seeds `navigation.md.tmpl` when a spec is of the `screen` layer, and
`dependencies.md.tmpl` when a code file declares a dependency on another.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the specs | the specs the compiler loaded; a screen is one whose header declares the `screen` layer | — | the compiler |
| the map | the `specifies`, `navigates-to` and `depends-on` edges | an edge whose ends are no screen's files, or no code files with a code of their own | this unit: the edge is left out of the page |
| the entries | `navigation.entry` in `anchors.yaml`, each a screen's name, with or without its `Screen` suffix, in any case | — | the config |

## Effects

| Effect | Description |
| --- | --- |
| `DCNAV-B01` | `navigation` lists the screen specs by name and every edge the code's navigation flags declare, each end taken to the screen whose spec specifies its file, with the rule that triggers it; a screen leading to itself is no edge. (`NavMap`, `NavScreen`, `NavEdge`, `fnNavigation`) |
| `DCNAV-B02` | A screen an entry names is marked as an entry, and a screen no entry reaches along the edges is marked unreached; with no entry declared, none is. The flowchart draws an entry as a stadium, an unreached screen dashed, each edge labelled with its rule. (`Mermaid`) |
| `DCNAV-B03` | `ScaffoldNavigation` and `ScaffoldDependencies`: the pages' templates, in the project's language, are seeded when a spec is a screen and when a code file declares a dependency, and compile into one row per screen with where it leads, and one row per file of the chain with what it uses and who uses it. (`DepFile`, `DepLink`, `fnDependencies`) |

## Errors

none — the unit reads the map and the specs the compiler already loaded; an edge it cannot place is left out of the page (Domain), not a failure.
