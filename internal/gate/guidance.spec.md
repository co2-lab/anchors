<!-- @anchors
  code: GDGTG
  updated_at: 2026-10-11
  layer: gate
-->
# Guidance gates — the working agreement is one file, and the others point to it

> **Code**: `GDGTG`

## Overview

The canonical guidance layer holds the files at the project's root that tell whoever arrives how the
project is worked on: `README.md`, `CONTRIBUTING.md`, `AGENTS.md` and `CLAUDE.md` (DESIGN-guidance-layer.md).
The agreement lives in `CONTRIBUTING.md`; the others send the reader there, so it is said once.
`guidance-linked` measures that much. Whether what they say contradicts Anchors and the configuration is
the AI judgment `guidance-honored` asks, against `anchors guide guidance`. Neither bars a commit.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| a target | a node of the guidance layer | any other node | this unit: skipped |

## Effects

| Effect | Description |
| --- | --- |
| `GDGTG-B01` | `guidance-linked` passes `CONTRIBUTING.md`; another guidance file passes when it mentions `CONTRIBUTING.md` and that file exists beside it, fails naming the missing agreement when it does not exist, and fails naming the file when it does not point to it; a node outside the guidance layer is skipped. (`checkGuidanceLinked`) |
| `GDGTG-B02` | The catalog offers both guidance gates to every project that chose an artifact, informative, on the guidance layer's tag; `guidance-honored` is an AI judgment about the configuration too (`stale_on_config`), whose question asks for contradictions with the guide and with `anchors.yaml`, never for style, and a quoted passage with its replacement. |
| `GDGTG-B03` | The header gates leave a headless layer's files: `header-valid` and `header-layer-declared` skip them, and `--fix` writes them no header. |

## Errors

none — the checks read the target and the map; a file that cannot be read is the engine's failure, reported before the checker runs.
