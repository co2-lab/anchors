<!-- @anchors
  code: MGCMM
  updated_at: 2026-09-26
  layer: comando
-->
# MigrateCommand — the command the format error promises, bringing the map and the config up to this binary's format

> **Code**: `MGCMM`

## Overview

When a binary meets a map or a config written in an older format, its error tells the user
to run `anchors migrate`. This unit is that command: an error that promises a command which
does not exist is worse than no error, because whoever reads it tries, fails, and stops
trusting the next message.

The command rewrites the two versioned files of a project, the map and the config, into the
format the running binary reads, through the migration steps of the migration package. It
reports each file on its own line: already current, migrated (with the format it went from
and to), or skipped with the reason it could not be read. For every file that changed it
lists the renamed keys, one per line with the number of occurrences, so whoever reviews the
diff knows what to expect before opening it.

A dry run says what would change and writes nothing. A real run ends by reminding that the
files are versioned and the migration must be committed, since a migration left on one
machine makes the next agent migrate again and produce the same diff.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | a directory; either file may be missing | — | this unit: a missing file is reported and skipped |
| each file's format | a format at or below the binary's | a file whose format cannot be determined | the migration package reports it; this unit prints it and moves on |

## Effects

| Effect | Description |
| --- | --- |
| `MGCMM-B01` | One run brings both the map and the config to the binary's current format, each reported as migrated from its old format. |
| `MGCMM-B02` | A file that is missing or cannot be read is reported on a `·` line with the reason, and the other file still migrates. |
| `MGCMM-B03` | A file already in the current format is reported as such and left untouched. |
| `MGCMM-B04` | For each migrated file, the renamed keys are listed in alphabetical order, each with its number of occurrences. |
| `MGCMM-B05` | With `--dry-run` the command reports what would be migrated, writes nothing, and says that nothing was written. |
| `MGCMM-B06` | After a real migration the command tells the user to commit it. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MGCMM-I01` | The migration is idempotent: a second run reports the files as current, asks for no commit, and leaves them byte for byte as the first run wrote them. | migrates, runs again, and compares the output and the file |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MGCMM-X01` | The command does not know which keys each format renamed: the steps belong to the migration package, and this unit only applies them and reports. | A format change is declared once, beside its reason; a second copy here would drift. |
| `MGCMM-X02` | The command never commits: it only reminds. | The commit belongs to whoever ran it, who reviews the diff first. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MGCMM-E01` | REF[MGCMM-B02]: a file that cannot be read is the failure this command handles, by reporting it and migrating the other file | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/migra/migrate.go` | `MigrateFile` | apoio — the migration steps |
| DEP2 | `internal/mapx/model.go` | `DefaultPath`, `FormatoAtual` | mapa — the map's path and the current format |
| DEP3 | `internal/config/config.go` | `AbsRoot`, `DefaultFile` | config — the config's path |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
