<!-- @anchors
  code: DCLND
  updated_at: 2026-10-02
  layer: apoio
-->
# DocLinks — the anchors, links, sizes and layer arrows the documentation templates are given

> **Code**: `DCLND`

## Overview

An index is only useful if each entry leads to the text. The target of an entry is a heading
anchor, and no documentation generator invents one: GitHub, MkDocs and Starlight follow the same
convention — the heading text in lower case, punctuation removed, spaces turned into hyphens,
accents kept. An anchor written by hand fails in silence: the link exists, is clickable and goes
nowhere. This unit computes the anchor once, so no template author has to.

It also builds the complete link to a rule or a scenario. The index cannot guess the shape of
the page it points at: a layer page shows each rule and scenario as its own heading when the
layer is small, and only summarizes each unit when it is big. The first version always pointed
at the rule's own anchor, and on the big layers that anchor did not exist — 483 of 812 links
broken. Here the rule is one: whoever decides the page's shape builds the link. Big or small is
asked of the layout decision (`DCOXX`), the same object the page templates ask, so link and page
cannot disagree.

It measures a selection of specs (units, rules, lines, scenarios) for that decision, deduces the
arrows between layers from the map's dependency edges for the architecture diagram, and gives
diagram nodes identifiers the diagram language can parse.

## Effects

| Effect | Description |
| --- | --- |
| `DCLND-B01` | `GitHubAnchor`: A heading's anchor is its text trimmed and lower-cased, keeping letters and digits (accented ones included), turning each space, hyphen and underscore into a hyphen, and dropping every other character. |
| `DCLND-B02` | A rule's link goes to its layer's page, at the rule's own heading when the layer is small and the rule is written as a heading, and at the heading of the unit that holds it when the layer is big or the rule is a table row or a bullet. |
| `DCLND-B09` | A layer's page (`layerPage`) is in the folder that holds the layer's template under `doct/`; with no template for the layer, in the folder of the project's language (`LayerDir`: `layers`, `camadas`, `capas`). |
| `DCLND-B03` | A scenario's link goes to its layer's page, at the scenario's heading when the layer is small, at its unit's heading when the layer is big, and at the bare page when the layer is big and the scenario's spec is unknown. |
| `DCLND-B04` | The size of a selection counts its units, their rules, their spec lines and their scenarios, and is remembered per selection, so asking again returns the same answer and two selections never share one. |
| `DCLND-B05` | The size of a selection that matches no spec is zero, not an error, and such a selection is never big. |
| `DCLND-B06` | The arrows between layers are the map's `needs` and `depends-on` edges between two specs of different layers, one arrow per pair of layers carrying the number of edges behind it; edges of other types, edges touching a non-spec node and edges inside one layer are left out. |
| `DCLND-B07` | The arrows are ordered by source layer and then by target layer, so the diagram does not change shape between builds. |
| `DCLND-B08` | `MermaidID`: A diagram node identifier is `n_` followed by the name with every character outside ASCII letters and digits replaced by an underscore. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a heading or name | any text | — | this unit: every character is either kept, turned into a hyphen or underscore, or dropped |
| a selection filter | the same `field=value` filter the spec selection accepts | a malformed filter | the spec selection (`DTCDC`) rejects it; here it is answered as an empty selection |
| the map's edges | edges of any type between any nodes | — | this unit: only dependency edges between two specs are counted |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCLND-I01` | Every link the generated index writes resolves to a heading that exists on the page it points at, on a small layer and on a big one alike. | builds the default skeleton over one small and one big layer, collects every generated link and every heading anchor, and verifies each link's anchor exists on its target page |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCLND-X01` | The unit does not decide what "big" means: it asks the layout decision, and the template function exposing it takes no threshold. | A template that picked its own threshold would recreate the defect: the page summarizing by one criterion and the links built by another. |

## Errors

none — a selection that matches nothing is answered as an empty size (`DCLND-B05`), and an edge to a node that is not a spec is filtered out (`DCLND-B06`); the unit reads nothing that can fail.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/doct/layout.go` | `Layout`, `Big` | apoio — the one big-or-small decision |
| DEP2 | `internal/doct/doct.go` | `Compiler`, `Spec`, `Rule` | apoio — the loaded specs and the spec selection |
| DEP3 | `internal/doct/scenarios.go` | `Scenario` | apoio — the scenarios counted and linked |
| DEP4 | `internal/mapx/model.go` | `Edges` | mapa — the dependency edges between units |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
