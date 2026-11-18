<!-- @anchors
  code: RNMBR
  updated_at: 2026-10-08
  layer: comando
-->
# Renumber — moves the revisions a branch added when the base already took their number

> **Code**: `RNMBR`

## Overview

Two open pull requests revising the same spec each take the next free revision number, and both take the
same one. The second to merge would land a revision code that already means something else on the base.
The number is kept — it says how many times the document changed — and it is the branch that arrives
second that renumbers its own revisions, never the base's: the base's revisions are history other people
have already read and cited.

This command finds the revision blocks the branch added whose number the base already uses, moves them to
the next free numbers, and rewrites the citations of the old code in everything the branch changed —
specs, plans, features, tests and code comments. A citation is rewritten only on a line the branch added:
a line that already existed where the branch forked cites what the base means by that code. It works
after a rebase as well as before it, and it is the sibling of `recode`: that one renames an identity code
across the project, this one renames revision codes across what the branch changed.

The decision of which revisions collide and to which number they move belongs to the revision planner;
this unit gathers the three versions of each file from git — the working copy, the one at the merge base
and the one at the base — hands them to the planner, and applies what it answers.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the base | a ref git can resolve, or nothing (then the default base) | a ref with no merge base with the current branch | this unit: fails naming the base |
| the spec files | paths inside the project, or none (then every changed Markdown file) | a path that cannot be read | this unit: fails naming the file |
| the changed files | any file the branch changed | a deleted, unreadable or binary file | this unit: skips it when rewriting citations |
| the project | a git repository with a configuration file | a directory with no configuration | this unit: fails naming the configuration file |

## Effects

| Effect | Description |
| --- | --- |
| `RNMBR-B01` | With no base given, the base is the remote copy of the integration branch when it exists, and the local integration branch otherwise. |
| `RNMBR-B02` | With no file given, the specs examined are the Markdown files the branch changed since the merge base, including changes not yet committed and files not yet tracked. |
| `RNMBR-B03` | Only a revision the branch added whose number the base already uses is moved, to the next free number; a revision the base has never moves. |
| `RNMBR-B04` | The old code is rewritten in every text file the branch changed, but only on the lines the branch added; a line that existed where the branch forked keeps citing the base's revision. |
| `RNMBR-B05` | A changed file that is binary is not rewritten. |
| `RNMBR-B06` | The dry run shows what would move and where it is cited, and writes nothing. |
| `RNMBR-B07` | When no added revision collides with the base, the command says there is nothing to renumber and writes nothing. |
| `RNMBR-B08` | A rewritten file keeps its permission bits. |
| `RNMBR-B09` | With files given, only those specs' revisions are examined. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `RNMBR-I01` | After a renumber, every revision code the base had still means what it meant on the base, in every file. | renumbers a rebased branch and confirms the base's revision and its citation are untouched |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `RNMBR-X01` | Does not stage or commit what it rewrites. | The rewrite of a citation the branch wrote about the base's revision would move too; the user reviews the diff before it enters history. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `RNMBR-E01` | The project has no configuration file. | The command fails naming the configuration file. | The default base comes from the configured integration branch. |
| `RNMBR-E02` | The base has no merge base with the current branch. | The command fails naming the base, before reading any file. | Without the fork point the command cannot tell what the branch added. |
| `RNMBR-E03` | A spec file given on the command line cannot be read. | The command fails naming the file, and nothing is written. | A renumber that skipped the file asked for would report success on work it did not do. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
