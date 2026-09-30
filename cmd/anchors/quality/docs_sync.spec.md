<!-- @anchors
  code: DCSYN
  updated_at: 2026-09-30
  layer: comando
-->
# DocsSyncForCommit — the commit carries the pages its specs produce

> **Code**: `DCSYN`

## Overview

A commit that changed a spec was barred by `docs-fresh` until someone ran `anchors docs build`
by hand, and a project could not make that gate blocking without turning every spec change
into a two-step commit. The pre-commit phase now compiles the pages of `docs/` the commit
leaves out of date — from the commit's specs, templates and map — and stages them with it, as
it does the map.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project | a git work tree with a `doct/` folder of templates and a map | no `doct/`, no map, or `docs.pre_commit: false` | this unit: nothing is compiled or staged |
| the specs and templates | their content as the index holds it | an unstaged edit, an untracked file | `scan.IndexReader`: they are read as the commit records them |
| the pages | pages the compiler generated, carrying its marker | a page written by hand, a page git ignores | this unit: never written, never staged |

## Effects

| Effect | Description |
| --- | --- |
| `DCSYN-B01` | A page the commit's specs leave out of date is compiled from the staged specs and templates, written and staged, when the tree holds nothing the index does not; a page already up to date is not compiled. (`syncDocsForCommit`, `Compiler.Compiled`) |
| `DCSYN-B02` | When the tree holds governed changes the index does not, the pages go straight into the index and the files on disk stay as they are; the staged page is compiled from the staged spec, not from the tree's. (`stageBlob`) |
| `DCSYN-B03` | With `docs.pre_commit: false`, without a `doct/` folder, or without a map, nothing is compiled or staged. |

## Constraints

| Code | Must not | Why |
| --- | --- | --- |
| `DCSYN-X01` | Write or stage a page written by hand, or a page git ignores. | A page with no generator marker is someone's work, and a page git ignores is not the commit's; the build refuses the first too. |
| `DCSYN-X02` | Write or stage anything outside `docs/`. | The hook rewrites what the gates would otherwise bar; anything else it touched would be a change nobody asked for. |

## Invariants

| Code | Invariant |
| --- | --- |
| `DCSYN-I01` | After the sync, `docs-fresh` confronted with the index finds no page out of date: what the hook staged is what the gate reads. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DCSYN-E01` | The index cannot be read, a template does not compile, or a page cannot be written or staged. | The error comes back; the hook warns and does not block. | `docs-fresh` still holds the commit to its pages; a sync that failed must not be what blocks it. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/doct/doct.go` | `NewWith`, `Compiled` | doct — the pages the templates produce |
| DEP2 | `internal/scan/scan.go` | `IndexReader`, `GovernedTreeChanges` | scan — the files as the index has them |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
