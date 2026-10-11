<!-- @anchors
  code: AFTSF
  updated_at: 2026-10-11
  layer: comando
-->
# AffectedTests — which tests to run for what changed, by rule and region where the project declares them

> **Code**: `AFTSF`

## Overview

The gates and the CI run every test; whoever is working wants the tests a change can break, now. Reported
from the bavel adoption: a hook running 1377 end-to-end phases for six minutes on every commit, when the
change touched the region of one of them. `anchors affected-tests` reads each changed file against its
last committed version at the finest grain the project declares — a spec's rule, a code region
`#region [CODE]` — and lists the tests that carry the changed codes in their title; a file changed
outside its regions reaches every test whose evidence closure holds it. The list is printed the way the
project's runner takes it.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the files | paths relative to the root, or none — then the files that differ from the last commit | — | git, when no file is given |
| the tests | the tests the dialect's `tests` source reads, with their titles | a project with no `tests` declared | this unit: the test files are printed, and stderr says why |
| the format | `lines`, `regex`, `files`, `json` | any other | this unit: refused |

## Effects

| Effect | Description |
| --- | --- |
| `AFTSF-B01` | A spec whose rule's definition changed — added and removed rules included — reaches the tests whose title carries that rule's code; a change to the rest of the spec reaches none. (`affectedBy`, `changedRuleCodes`) |
| `AFTSF-B02` | A file whose region `#region [CODE]` changed — added and removed regions included — reaches the tests whose title carries CODE; a change outside every region reaches every test whose evidence closure holds the file, and so does a new or a deleted file. (`changedRegions`, `outsideRegions`) |
| `AFTSF-B03` | A changed test file runs whole, and a file a freeze covers reaches nothing. |
| `AFTSF-B04` | The tests are printed in the format asked: `lines` one title per line, `regex` one anchored alternation of the titles with their metacharacters escaped, `files` the test files once each, `json` each test with file, line and title; in file and line order. With no file given, the files are those that differ from the last commit. (`selectAffected`, `printAffected`) |
| `AFTSF-B05` | With no `dialect.tests` declared the titles cannot be read: the test files run whole are printed, and stderr says the tests by code are left out. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `AFTSF-E01` | A format outside the four. | The command fails naming the formats. | A list in a shape the runner does not take runs nothing, silently. |

## Open Decisions

None.
