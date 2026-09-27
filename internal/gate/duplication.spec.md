<!-- @anchors
  code: DUPLC
  updated_at: 2026-09-27
  layer: gate
-->
# Duplication — no code file holds a block copied from somewhere else

> **Code**: `DUPLC`

## Overview

The native check behind `no-duplication`. It runs jscpd once per scan and reads its JSON report,
which lists every clone with both files and their lines, instead of reading jscpd's exit code. jscpd
exits 0 with any duplication unless a `threshold` is configured, so a gate that read the exit code
approved what it did not measure: in the reference app, 35 clones (1.17% of the lines), exit 0, and
`check --all` announced the gate "clean, ready to become blocking".

Each file gets its own verdict. A file that holds a copy fails, naming the other side of each clone
and the lines; every other file passes. With `--changed` only the files of the change are judged,
so a commit answers for the clones it touches; jscpd still scans the whole project, because a new
copy only exists relative to its original.

The calibration is the project's `.jscpd.json`, which jscpd reads by itself (`minLines`,
`ignore`, `jscpd:ignore-start` markers). Its `threshold` keeps jscpd's meaning: the percentage of
duplicated lines the project tolerates.

## Signature

| Parameter | Type | Description |
| --- | --- | --- |
| `content` | `string` | Text of the file (not read: the report names the clones) |
| `n` | `mapx.Node` | The file being judged |
| `root` | `string` | Project root, where jscpd runs and `.jscpd.json` lives |
| `g` | `*mapx.Graph` | The map; with the root, the key of the one run per scan |
| `cfg` | `*config.Config` | Not read |

**Returns**: `(Verdict, string)` — `Pass` for a file in no clone, `Fail` listing its clones,
`Pending` when the clones are within the declared threshold or when jscpd produced no report.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the tool | `npx` on the PATH, which fetches jscpd | no `npx` | the gate runner: `needs_tool: npx` turns the gate into Skip |
| the report | jscpd's `jscpd-report.json`, paths relative to the root or absolute | no report, or one that is not JSON | this unit: Pending naming why, never an approval |
| `.jscpd.json` | jscpd's configuration, with or without `threshold` | a file that is not JSON | this unit: read as declaring no threshold, as jscpd itself cannot use it |

## Effects

| Effect | Description |
| --- | --- |
| `DUPLC-B01` | A file that takes part in no clone passes. |
| `DUPLC-B02` | A file that takes part in a clone fails, naming for each clone its own lines, the other file and that file's lines, whichever side of the clone it is on. |
| `DUPLC-B03` | A clone inside a single file is described by its two line ranges alone. |
| `DUPLC-B04` | When `.jscpd.json` declares a `threshold` and the project's duplicated percentage is at or under it, the file's clones are reported as Pending, not failed; over it they fail. |
| `DUPLC-B05` | Without a declared `threshold`, any clone fails, whatever jscpd's exit code. |
| `DUPLC-B06` | jscpd runs once per scan: the same root and map reuse the report, and a new map runs it again. |
| `DUPLC-B07` | An absolute path in the report is read relative to the project root. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DUPLC-E01` | jscpd writes no report, or one that is not JSON | Pending, naming the last line jscpd printed or the parse error | The duplication was not measured; approving would say it was, failing would blame the code for the tool |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Config` | core — the gate's configuration, passed and not read |
| DEP2 | `internal/i18n/i18n.go` | `T` | core — the verdict messages |
| DEP3 | `internal/mapx/model.go` | `Graph`, `Node` | core — the file being judged and the scan's map |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
