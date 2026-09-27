<!-- @anchors
  code: MGFLM
  updated_at: 2026-09-26
  layer: apoio
-->
# MigrateFile — takes one project file from its declared format to the current one, renaming only what is a key

> **Code**: `MGFLM`

## Overview

Anchors' YAML keys and gate names are identifiers, and some of them were born in Portuguese or with
names that lied about what they hold. Renaming them without a migration would break every existing
project, and in the worst way — silently: a map whose judgment stamps sit under the old key would
lose them, and the check would redo every judgment, charging again what someone already answered.
Reading both names forever does not close either: the file never gets fixed, every reader must know
both, and the next rename inherits the problem of the previous ones. The migration runs once and
leaves the project on the new format.

This unit migrates ONE file. It reads the format the file declares in its top-level `version:`
line — without the YAML parser, on purpose, because the file may be in a shape the current parser
refuses, which is exactly the case to detect and fix. A file with no such line is format 1, the
format of the first files, which did not declare it. It then applies, in order, the steps between
that format and the target, each on the result of the previous one, renaming keys and values only
where they are keys and values: the old names also appear in comments and prose inside Anchors'
files, and a substring replacement would rewrite text that is no key at all. Each rename belongs to
one file, chosen by the file's name, so a key of the map is never renamed in the configuration.

The `version:` line is raised to the target even when nothing else changed, because the format
belongs to the file, not to what it happens to contain; a map with no judgments is still format 2
after migrating, and leaving it at 1 would rerun the migration on every command. A map from before
the field existed gets the line inserted right after its leading comment block, where the map
writer puts it. A dry run measures what would change without writing, which is what lets the
doctor warn without touching the repository of someone who did not ask.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the file | an Anchors map or configuration file, readable, in any format up to the target | an unreadable or missing file | this unit returns the read error |
| the declared format | a top-level `version:` line holding a number, or no such line | a number too large to represent | this unit returns an error saying the version is not a number |
| the target format | the binary's current format | a target with no registered chain of steps | the step chain refuses it with an error naming the missing step |

## Effects

| Effect | Description |
| --- | --- |
| `MGFLM-B01` | `FormatOf` answers that the file's format is the number in its top-level `version:` line; a file without that line is format 1. |
| `MGFLM-B02` | `MigrateFile`: A file already at or above the target format is left untouched and reported as unchanged. |
| `MGFLM-B03` | An old key is renamed only where it is a key at the start of a line — at any indentation, list items included — keeping its value; the same word in a comment or inside a value is not touched. |
| `MGFLM-B04` | Each rename applies only to the file it was declared for: a key of the map is never renamed in the configuration file, nor the reverse. |
| `MGFLM-B05` | A value is renamed only under the key the step names and only when the whole value matches an old name. |
| `MGFLM-B06` | The steps are applied in order, each on the result of the previous, so a key renamed twice across formats ends under its latest name. |
| `MGFLM-B07` | The `version:` line is raised to the target, even when no key or value changed. |
| `MGFLM-B08` | A map or a configuration without a `version:` line gets one inserted right after its leading comment block. |
| `MGFLM-B09` | The result counts each renamed key, and each renamed `key: value`, by its old form, and reports whether the text changed. |
| `MGFLM-B10` | A dry run reports what would change, counts included, without writing the file. |
| `MGFLM-B11` | A top-level `version:` line followed by a comment is read by its number, and raised in place keeping the comment. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `MGFLM-I01` | Migrating twice is the same as migrating once: the second run changes nothing and reports no change. | migrates a file, migrates it again, and compares the text and the reported change |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `MGFLM-X01` | Does not load the file through the YAML parser: text the parser would refuse is still read and migrated. | The file to fix may be in a shape the current parser refuses — exactly the case the migration exists for. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `MGFLM-E01` | The file cannot be read. | The read error is returned and nothing is written. | There is no format to read, and nothing to migrate. |
| `MGFLM-E02` | The top-level `version:` line holds something that is not a number (text such as `abc`, or a number too large to represent). | An error naming the file and saying the version is not a number; nothing is written. | Guessing a format would apply the wrong steps, and reading it as format 1 inserted a second `version:` line. |
| `MGFLM-E03` | A step between the file's format and the target is missing. | The chain's error is returned and the file is left untouched. | Writing the new number with part of the conversion missing would make the file lie about its own format. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/migra/step.go` | `StepsFrom` | apoio — the ordered steps between two formats |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
