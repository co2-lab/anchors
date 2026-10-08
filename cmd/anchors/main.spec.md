<!-- @anchors
  code: CLMNC
  updated_at: 2026-10-08
  layer: comando
-->
# CliMain — the entry point that stamps the build identity, prints a failure once and turns it into the exit code the hooks read

> **Code**: `CLMNC`

## Overview

The entry point wires what no package can wire for itself, then runs the root command. It
hands the build's version to every place that records it: the version the CLI reports, and
the generator the map records, so an older binary is accused instead of silently undoing what
a newer one wrote. It hands the configuration loader the migration registry's knowledge of
renamed keys, which the loader cannot import without a cycle, so an unknown key in an
old-format configuration is told apart from a typo: the advice is to migrate, not to update the
binary.

When the command fails, the entry point prints the error once, prefixed with the localized
word for error, and chooses the exit code. A file the project does not govern is not a
failure: it exits with its own code, which the hooks read to let a commit of configuration
files through, while any other failure exits 1.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the build identity | the version, commit and date set at link time, or the development defaults | — | the release build |
| the command's outcome | success, the not-governed signal, or any other error | — | the commands; this unit classifies |

## Effects

| Effect | Description |
| --- | --- |
| `CLMNC-B01` | A failing command prints its error once on standard error, prefixed with the localized word for error, and the process exits 1. |
| `CLMNC-B02` | A command that signals a file the project does not govern exits with the not-governed code, 3. |
| `CLMNC-B03` | The build version reaches the version the CLI reports and the generator the map records; a build without link-time values reports the development identity. |
| `CLMNC-B04` | An unknown key that the migration renames, in a configuration of an older format, is reported as renamed, with the advice to run `anchors migrate`. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CLMNC-I01` | REF[CLMNC-B02]: the not-governed signal and a real failure never share an exit code | — |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CLMNC-X01` | REF[CLMNC-B01]: the entry point, not the root, prints the error, so it comes out once | — |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CLMNC-E01` | REF[CLMNC-B01]: any command failure is handled here by printing it once and exiting non-zero | — | — |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
