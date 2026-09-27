<!-- @anchors
  code: SCIGS
  updated_at: 2026-09-26
  layer: scan
-->
# ScanIgnore — what the scan never sees: the built-in list, the project's `.gitignore`, and editor noise

> **Code**: `SCIGS`

## Overview

Decides which directories and files the repository scan does not look at. Everything the map, the
checks and the watcher know about the project passes through this decision, so an error here is
silent: a file wrongly ignored is not reported as missing, it simply never exists for any gate.

Three sources are composed. A small built-in list of directory names that are almost never project
material (dependency caches, build output, coverage), because descending into them is expensive and
almost always useless. The project's own root `.gitignore`, because it is the list the team already
maintains of what is disposable, and a second list in the Anchors configuration would age apart from
it. And a fixed set of editor and system ephemera (swap files, atomic-save temporaries, lock files),
which are noise of the file system and never material in any project.

The built-in list is defeatable on purpose: `build` or `dist` is compiled output in most projects and
source code in some. A layer whose pattern points inside such a directory, or a negation in the
`.gitignore`, wins over the default. Two kinds of directory are never defeatable: the machinery
(`.git`, `.anchors`) and the records Anchors itself writes (`issues`, `changes`), which are versioned
output and would feed the work queue with its own reports.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a directory, with or without a `.gitignore` at its top | — | this unit: an absent or unreadable `.gitignore` means the project declared nothing |
| the `.gitignore` lines | gitignore syntax: blank lines, `#` comments, `!` negations, leading, inner and trailing slashes, globs | nested `.gitignore` files and `.git/info/exclude`, which are not read | this unit reads only the root file |
| the layer configuration | the declared layers, or none at all | — | the caller passes the loaded configuration or nothing |
| the queried path | a root-relative path, with `/` or the native separator | an absolute path | the caller (the scan walk) always passes a root-relative path |

## Effects

| Effect | Description |
| --- | --- |
| `SCIGS-B01` | `LoadIgnore`: With nothing declared, the built-in directories (`node_modules`, `dist`, `build`, `vendor`, `.next`, `coverage`, `.expo`) are skipped. |
| `SCIGS-B02` | `LoadIgnoreFor`: A layer whose pattern names a built-in directory as a path segment re-enables that directory; a catch-all pattern that merely reaches it does not. |
| `SCIGS-B03` | A `.gitignore` negation naming a built-in directory (`!build/`) re-enables it. |
| `SCIGS-B04` | The directories `issues` and `changes`, which Anchors writes as its own records, are skipped at any depth, whatever the `.gitignore` says. |
| `SCIGS-B05` | Editor and system ephemera — atomic-save temporaries, swap, backup, lock and autosave files, `*.tmp`, `.DS_Store`, `Thumbs.db` — are never files to scan; a directory merely named like one (`tmp/`) is not affected. |
| `SCIGS-B06` | A `.gitignore` pattern with a leading or an inner slash is anchored at the root; a pattern with no slash matches the basename or any path segment at any depth. |
| `SCIGS-B07` | A pattern with a trailing slash matches directories only, and everything below an ignored directory is ignored with it; a plain file with that name is not ignored. |
| `SCIGS-B08` | The `.gitignore` rules apply in order and the last one that matches decides, so a later `!pattern` re-includes a file an earlier rule excluded. |
| `SCIGS-B09` | Without a loaded ignore set, the built-in directories, the machinery (`.git`, `.anchors`) and the Anchors records are still skipped and ephemera are still not files; nothing else is ignored. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SCIGS-I01` | The machinery directories `.git` and `.anchors` are skipped under every declaration: no negation re-enables them. | writes `!.git/` and `!.anchors/` in the `.gitignore` and verifies both are still skipped |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SCIGS-X01` | Only the root `.gitignore` is read; a nested `.gitignore` does not change what the scan sees. | The unit reproduces the essential of gitignore semantics for the list the team maintains at the root, not the whole of git's resolution. |

## Errors

none — an absent or unreadable `.gitignore` is the normal case of a project that declared nothing, and the built-in list still applies; a pattern that does not parse as a glob matches nothing, as a line git itself would not match.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | config — the declared layers whose patterns re-enable built-in directories |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
