<!-- @anchors
  code: DCSCD
  updated_at: 2026-10-08
  layer: apoio
-->
# DocScaffolds — the starting templates `anchors docs init` proposes, one page per question and per layer

> **Code**: `DCSCD`

## Overview

The pages of a documentation answer different questions for different readers, and the
organisation is what lets each reader find theirs without reading the others. This unit
proposes that organisation as ordinary templates the project then edits: an architecture page
(how the system is built, as a C4 model), a behaviour index (what happens, every scenario), a
rules index (every rule of the system) and one page per layer, where the content lives.

The content lives in ONE place, the layer page, and the cross-cutting views are indexes. The
first version repeated the text in both views, and in the reference project the rules page came
out with 9,608 lines — the cross-cut of 84 units is the whole document once more, growing with
the project. An index has one line per rule, and its link leads to the text on the layer page.

The scaffolds PROPOSE, they do not impose: once written, a template belongs to the team, and
running init again never erases the frame the team wrote. Each written template opens with a
comment saying what question the page answers, where whoever edits it will read it.

## Effects

| Effect | Description |
| --- | --- |
| `DCSCD-B01` | `Scaffolds`: There are three fixed page templates — architecture, behaviour and rules — each with a template name, a body and a sentence saying what the page answers. |
| `DCSCD-B02` | `ScaffoldLayer`: A layer's page template is named after the layer inside the layer folder of the project's language (`LayerDir`, DCLND-B09). |
| `DCSCD-B03` | Init writes the three fixed templates plus one layer page for each layer that has specs, each one opening with a template comment carrying the sentence of what the page answers. |
| `DCSCD-B05` | `ScaffoldOpenAPI`: When a spec has an `Endpoint` section, init also writes `openapi.yaml.tmpl`, compiling the project's OpenAPI under the project folder's name and version `0.1.0`, for the team to edit. |
| `DCSCD-B04` | An existing template is not rewritten unless forced; it is reported as skipped. Forced, it is rewritten. |
| `DCSCD-B05` | On a small layer, the layer page carries each scenario's steps under a heading of its own, and the behaviour index links to that heading. |
| `DCSCD-B06` | On a big layer, the layer page carries the layout's summary sentence and each unit's overview and rule list, without the scenarios' steps. |
| `DCSCD-B07` | The architecture page shows the protocol of every declared conversation, lists external containers at level 2, and gives a level 3 section to each internal container only. |
| `DCSCD-B09` | The templates are in the project's language: their names (`architecture`, `arquitetura`, `arquitectura`…), what they show — headings, paragraphs, diagram labels, empty-state notes — and the section titles the layer page asks for, which are the language's titles in the section catalog; a language with no table gets English, and no template keeps an unfilled text. The notes to whoever edits a template are in English. |
| `DCSCD-B08` | When no container is declared, the architecture page says so instead of leaving empty diagrams. |

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a directory where the templates folder can be created | a root where the templates folder cannot be created or written | this unit: the write error is returned, see `DCSCD-E01` |
| the layers | the layers that have specs, as the compiler lists them | — | the compiler (`DTCDC`) |
| the force flag | on or off | — | the caller (`anchors docs init --force`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCSCD-I01` | The skeleton init writes compiles with the compiler, producing one page per template written. | runs init over a project with a spec and a feature, builds, and verifies every template became a page |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCSCD-X01` | Init only writes templates; it never writes a compiled page. | Compiling is `docs build`'s job, with its own protection of handwritten pages; init writing pages would bypass it. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DCSCD-E01` | The templates folder cannot be created or a template cannot be written. | Init stops and returns the error, with the templates written so far. | A partial skeleton reported as complete would leave the team building from pages that are not there. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
