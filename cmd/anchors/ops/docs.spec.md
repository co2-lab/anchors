<!-- @anchors
  code: DCCMD
  updated_at: 2026-10-08
  layer: comando
-->
# DocsCommand — the documentation is compiled from the specs through templates, against a map rebuilt from the tree

> **Code**: `DCCMD`

## Overview

The documentation has three layers and its content lives in one. The specs are the source;
the templates are the frame that references excerpts of the specs; the compiled pages are
generated and nobody edits them. Markdown has no native inclusion, so without a build the
only way to show a spec's content in a page would be to copy it by hand, and a hand-written
copy goes stale unseen.

`docs build` compiles the templates into pages. It compiles against a map rebuilt from the
working tree, not the map on disk: a map that had lost nodes in a merge once produced
compiled pages missing four whole units, with no error, and the commit looked normal. The
rebuilt map stays in memory. Whoever needs a specific map asks for it explicitly, and is
warned that what that map does not know does not get in. A page that exists without the
compiler's marker was written by hand; it is never overwritten, and the command names it
and says how to resolve the double ownership. A dry run compiles without writing.

`docs init` writes the skeleton of the templates once, never overwriting an existing one
unless forced. `docs duties` answers which documents the project requires, for the whole
project, for a layer, or for one unit (resolved to its layer and code), and teaches the
known kinds when the project declares none.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a governed project with `anchors.yaml` | a directory with no config | this unit: build and duties fail with the load error |
| the map on disk | a loadable map, required by `init` and by `build --no-map-rebuild` | a missing map in those two cases | this unit: fails pointing at `anchors map build` |
| the templates | any number of templates in the template directory, including none | — | the compiler; this unit reports the result |
| the unit or layer | a path to a unit, or a layer name | — | this unit: an unknown layer simply has no duty |

## Effects

| Effect | Description |
| --- | --- |
| `DCCMD-B01` | `docs build` compiles against a map rebuilt from the tree, so a unit the map on disk does not know still reaches the compiled page. |
| `DCCMD-B02` | With `--no-map-rebuild`, the build compiles against the map on disk, and warns on standard error that what the map does not know does not get in. |
| `DCCMD-B03` | With `--dry-run`, the pages are compiled and listed as not written, and nothing is written. |
| `DCCMD-B04` | A page that exists without the compiler's marker is skipped and kept as it is, and the output names it and says to delete either the template or the page. |
| `DCCMD-B05` | With no template, the build says there is nothing to compile. |
| `DCCMD-B06` | `docs init` writes the template skeleton; a second run skips every existing template and says so, and `--force` rewrites them. |
| `DCCMD-B07` | `docs duties` lists every required document; with `--layer` only those the layer triggers, or says none is required; with `--unit` those of the unit's layer. |
| `DCCMD-B08` | When the project declares no required document, `docs duties` says so and lists the kinds Anchors knows how to instruct. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DCCMD-I01` | REF[DCCMD-B04]: a hand-written page is never overwritten by a build | — |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DCCMD-X01` | `docs build` never writes the map: the rebuilt map lives only for the compilation. | Writing the map is the map command's job, with its own merge rules; a side effect here would surprise every other reader. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DCCMD-E01` | `docs build` runs where `anchors.yaml` cannot be loaded. | It fails with the load error, pointing at `anchors init`. | There is no tree definition to rebuild the map from. |
| `DCCMD-E02` | `docs duties` runs where `anchors.yaml` cannot be loaded. | It fails with the load error. | The duties are declared in the config; answering "none" would lie. |
| `DCCMD-E03` | `docs build --no-map-rebuild` or `docs init` runs with no map on disk. | It fails pointing at `anchors map build`. | Both were asked to use the map on disk, and there is none. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
