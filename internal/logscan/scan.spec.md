<!-- @anchors
  code: LGSCL
  updated_at: 2026-10-08
  layer: apoio
-->
# LogScan — finding the occurrences of declared failures in the project's logs

> **Code**: `LGSCL`

## Overview

A spec declares the failures its unit handles as `-E` rules. This unit reads the project's logs and
counts how often each declared failure actually happened, and when (the first and last moments), so the
map can confront the declared failures with the real ones.

The log format does not matter: if the handling path records the failure's code, the code is the same
sequence of characters in JSON, plain text or syslog, and the scanner looks for the code, not the
format. The requirement belongs to the project: log the code, wrapped by default in `#[` and `]`. The
pattern is built from the codes the project's specs actually declare, never from the generic shape of a
code: measured with a bare pattern, four of five captures in a log were build ids, cache keys, SKUs and
trace ids. The delimiter removes the rest of the ambiguity.

A code the log carries and no spec declares is the layer's most valuable finding (the log shows an error
the spec never foresaw), so it is reported apart instead of discarded; without a delimiter it is only
believed on a line that says it reports a failure. Legacy logs that nobody will rewrite can be bound to
a code by a text alias, applied only to lines that carry no declared code. Anchors never guesses where
the logs are: with no declared path it reads nothing, since a log usually carries what must not leak.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the logs configuration | declared glob paths, optional delimiters, aliases and a timestamp pattern | an alias or timestamp that is not a valid pattern | this unit: the scan fails |
| the map | the project's graph, whose specs declare the `-E` rules | — | the map build |
| the log lines | any format, lines up to 4 MB | — | this unit reads them as text |

## Effects

| Effect | Description |
| --- | --- |
| `LGSCL-B01` | With no declared log path, nothing is scanned and no result is given. |
| `LGSCL-B02` | Every file matching the declared paths is read, whatever its line format; each declared code found counts one occurrence per line, and the number of files and lines read is reported. |
| `LGSCL-B03` | Only the `-E` codes the project's specs declare become occurrences; a code that merely looks like one is not captured. |
| `LGSCL-B04` | A failure code no spec declares is reported apart with its count: with delimiters, whenever it is delimited; without delimiters, only on a line carrying a failure word (error, err, fatal, severe, critical, crit, panic, exception). |
| `LGSCL-B05` | The code is looked for between `#[` and `]` by default; a declared pair of delimiters replaces them, and an empty pair looks for the bare code. |
| `LGSCL-B06` | A legacy alias binds its code to the lines its text pattern matches, only when the line carries no declared code. |
| `LGSCL-B07` | With a timestamp pattern declared, each occurrence records the earliest and latest moment it was seen. |
| `LGSCL-B08` | The occurrences come back sorted by rule code. |
| `LGSCL-B09` | The `-E` rules a spec catalogues are listed from its headings, table rows and bold bullets. (`SpecFailureCodes`) |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `LGSCL-E01` | A declared alias or timestamp is not a valid pattern. | The scan fails with the pattern's error and gives no result. | A partial scan would report fewer occurrences than happened, and the confrontation would clear failures that were never looked for. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
