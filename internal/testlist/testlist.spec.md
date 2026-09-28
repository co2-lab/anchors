<!-- @anchors
  code: TSTLS
  updated_at: 2026-09-28
  layer: infra
-->
# TestList — the project's tests, read the way the project says they are written

> **Code**: `TSTLS`

## Overview

Reads the project's tests: which tests exist, in which file, on which line, and under which title.
How a test is written belongs to the project's language and test library, not to Anchors, so the
project declares it, choosing one of two sources:

- a **pattern**: the regular expression of the call that opens a test, up to its opening parenthesis.
  The title is the string literal right after it.
- a **script**: a command the project provides, run at the project root, which prints the tests on
  stdout under a versioned contract. The script can ask the test library itself, so when the library
  changes its answer changes with it.

The script's contract, version 1, is one JSON object:
`{"version": 1, "tests": [{"file": "src/a.test.ts", "line": 12, "title": "…"}]}` — `file` relative
to the project root, `line` counted from one, `title` as written.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the source | a pattern, a script, or neither | both at once | this unit: refuses, naming the conflict |
| the pattern | a regular expression that compiles | one that does not compile | this unit: refuses, naming the error |
| the files | paths relative to the root | — | the caller, from the map's test nodes |
| the script's output | the contract's JSON object | anything else | this unit: refuses, naming what is wrong |

## Effects

| Effect | Description |
| --- | --- |
| `TSTLS-B01` | With a pattern, each call it opens that is followed by a string literal is a test: the literal is its title, and the line is the line where the call starts. |
| `TSTLS-B02` | A title may be quoted with single quotes, double quotes or backticks; a backslash escapes the next character inside the first two, and only a backtick title may span lines. |
| `TSTLS-B03` | A call whose next token is not a string literal is not a test this reading can name, and is left out. |
| `TSTLS-B04` | The pattern's tests come in file order, then in the order they appear in each file. |
| `TSTLS-B05` | With a script, the command runs at the project root and the tests are what it prints under the contract. |
| `TSTLS-B06` | Output outside the contract is refused, naming what is wrong: not a single JSON object, an unknown field, a version other than 1, no `tests`, an empty or non-relative `file`, a `line` below 1, an `end` before its `line`, or an empty `title`. (`Parse`) |
| `TSTLS-B07` | A `file` in the script's output is read with forward slashes and without redundant segments, so it matches the map's IDs. |
| `TSTLS-B08` | A source that declares neither a pattern nor a script lists no test and raises no error. |
| `TSTLS-B09` | The script may say where each test ends (`end`, its last line), and the test carries it; a test without `end`, and every test a pattern reads, carries zero — the end is not known. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `TSTLS-E01` | The source declares both a pattern and a script | an error naming the conflict | Two sources could disagree, and the project must say which one speaks |
| `TSTLS-E02` | The pattern does not compile | an error naming the compile error | A pattern that matches nothing would read as a project with no tests |
| `TSTLS-E03` | The script exits with an error | an error naming the command, the exit and the last line it wrote to stderr | A failed script listed nothing, which is not the same as no tests |
| `TSTLS-E04` | REF[TSTLS-B06]: the script answers outside the contract | — | — |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
