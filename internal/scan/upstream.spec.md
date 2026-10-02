<!-- @anchors
  code: UPOWP
  updated_at: 2026-10-02
  layer: scan
-->
# UpstreamOwnership — which files Anchors still owns in a project, and where a file's `@anchors` header begins and ends

> **Code**: `UPOWP`

## Overview

Anchors seeds pipelines into a project's workflow directory. While such a file still carries the
template marker it belongs to the Anchors project, not to the one that received it: its rules, spec
and tests live upstream, and the local copy is replaced whole when the doctor fixes the environment.
Measured in a real project, the map asked the team to own and specify files they did not own, and
gave one of them an identity inferred from an example in its comments. The marker is the whole test,
on purpose: a team that edits the file removes the marker, and from then on the file is theirs and
governed like any other.

The unit also delimits a file's `@anchors` header block. Reading a header key over the WHOLE file
reads prose and code as declarations — a workflow's jq program with a line starting `parent:` became
a node's parent in the map. The block opens on the first comment line that carries the `@anchors`
token and ends where that comment ends, so consumers read header keys only inside it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the file path | a root-relative path, with `/` or `\` separators | an absolute path | the scan passes the path it walked |
| the file content | any bytes | — | this unit: content without the marker or without a header yields "not owned" or no header |
| the header's comment syntax | HTML comments, block comments, and line comments starting with `//`, `#`, `*` or `--` | a header written outside a comment, which is prose | this unit: only a comment line opens the header |

## Effects

| Effect | Description |
| --- | --- |
| `UPOWP-B01` | `IsUpstreamOwned`: A file is owned upstream only when it lives under `.github/workflows/` AND its content carries the marker `anchors:template`; the marker elsewhere, or a workflow without it, is the project's. |
| `UPOWP-B02` | A Windows-style path under the workflow directory is recognised the same as its slash form. |
| `UPOWP-B03` | `AnchorsHeader`: The header opens only on a comment line carrying the `@anchors` token followed by whitespace, the end of an HTML comment, or the end of the line: `@anchors-shared-code` and a mention of `@anchors` in prose open nothing, and a file with no such line has no header. |
| `UPOWP-B04` | An HTML header ends at the line carrying `-->`, and a block-comment header at the line carrying `*/`; a header opened on a line that already closes its comment is that single line. |
| `UPOWP-B05` | A line-comment header ends at the first line that is not a comment. |
| `UPOWP-B06` | A header comment that never closes runs to the end of the file. |
| `UPOWP-B07` | The header is at the TOP: only blank lines, comments — line comments and the insides of `/* */` and `<!-- -->` blocks — and a shebang may precede it. A block below any other line is not the header unless it declares inside `@fixed-header: <why>`; a bare `@fixed-header` declares nothing. `HeaderOffTop` says a file carries such an undeclared block and no header. |
| `UPOWP-B08` | A comment opens the header only when `@anchors` is its first word, right after the comment marker; prose that names the header does not. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `UPOWP-I01` | Nothing after the header's comment is part of the header, whichever of the comment forms opened it. | reads the header of a block-comment file followed by a body line and verifies the body line is absent |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `UPOWP-X01` | Ownership is decided by the marker and the directory alone; the file's name and contents otherwise are not consulted. | Removing the marker is how a team takes a file over, so the marker must be the whole test. |

## Errors

none — every input has an answer: a path or content that does not qualify is simply not owned upstream, and a file without a header comment has no header, which is the normal case of an unannotated file.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
