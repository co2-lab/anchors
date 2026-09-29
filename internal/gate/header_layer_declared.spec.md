<!-- @anchors
  code: HDLYD
  updated_at: 2026-09-28
  layer: gate
-->
# HeaderLayerDeclared — the layer a header declares is one the Estrutura has

> **Code**: `HDLYD`

## Overview

The header's `layer:` decides which templates derive the unit's siblings and where the
documentation files it. A layer the Estrutura does not have is read as nothing: the unit falls back
to the default templates, lands in a page of its own, and shows up — if at all — as a layer
outside every container on the architecture page. In the reference app a section declared
`layer: landing-component` for months while the Estrutura only had `landing-feature`, and it
surfaced by chance.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the header | an `@anchors` header with `layer:` | no header, no layer, a `TODO` placeholder | this unit: Skip (the placeholder is `placeholder-filled`'s) |
| the Estrutura | the project's `layers:` | none declared | this unit: Skip |

## Effects

| Effect | Description |
| --- | --- |
| `HDLYD-B01` | A header layer the Estrutura declares passes; one it does not fails, naming the layer and the declared layers in order. Only the `@anchors` header is read: a `layer:` line in the body is not a declaration. |
| `HDLYD-B02` | A header layer that differs from a declared one only in case fails naming both: the lookup is exact, so the unit is read as unknown. |
| `HDLYD-B03` | A file whose header declares no layer, a placeholder layer (`TODO…`), and a project with no layers are skipped. |

## Errors

none — the gate reads the file's header and the configuration

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/scan/scan.go` | `HeaderLayerOf` | scan — the layer the header declares |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
