<!-- @anchors
  code: GTSTG
  updated_at: 2026-10-01
  layer: infra
-->
# GitState — classify the project's versioning before `init` scans it, and say what to do about it

> **Code**: `GTSTG`

## Overview

Git is not a comfort detail for Anchors: the change stamp of every anchor, the diff coverage, the
pre-commit hook and, in `github` mode, the repository of the work queue all come from it. A project
without git loads and runs, but a large part of the framework is silently switched off. So `init`
classifies the project root BEFORE anything else, and tells the person what is missing.

"There is no git" has two meanings, and Anchors acts differently on each: the git program may be
missing from the machine (there is nothing to offer — offering to initialise a repository would fail
the moment it was accepted), or the program exists and the folder is simply not a repository (then
initialising is exactly the step to offer). A repository with no commit yet is a third state of its
own: without a first commit there is no history to compare against, so calling it "ready" would claim
that things work when they do not.

The unit also holds the text of the warning for each state, next to the decision, so the words that
teach the person WHY git matters are tested together with the state that triggers them, and the short
`.gitignore` seeded with the first commit.

Whether git is installed is given by the caller; the unit only reads the folder tree, so the
classification can be proven on any machine.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | an existing folder path | a path that does not exist | the caller (`init` runs in the folder it was given) |
| whether git is installed | a yes/no already measured by the caller | — | the caller, which looks the program up before asking |
| the state to explain or to offer on | one of the four states | — | this unit: an unknown state has no warning and no offer |

## Effects

| Effect | Description |
| --- | --- |
| `GTSTG-B01` | When git is not installed, the classification (`DetectGit`) is "not installed", whatever the folder holds. |
| `GTSTG-B02` | When git is installed and neither the root nor any ancestor folder has a repository, the state is "not initialised". |
| `GTSTG-B03` | A repository whose local branches hold at least one reference is "ready". |
| `GTSTG-B04` | A repository whose references were packed (no loose branch reference, a non-empty packed list) is "ready". |
| `GTSTG-B05` | A repository with no branch reference and no non-empty packed list is "no commit yet". |
| `GTSTG-B06` | A root inside an existing repository (the repository is in an ancestor folder) is "ready", so no nested repository is ever offered. |
| `GTSTG-B07` | A root whose repository marker is a file (a worktree or a submodule pointer) is "ready". |
| `GTSTG-B08` | There is an action to offer (`OfferAction`) only for "not initialised" (initialise) and "no commit yet" (commit); never for "not installed" or "ready". |
| `GTSTG-B09` | Each of the three unready states has its own warning (`AvisoGit`), in the project's language — install git, the folder is not under git, there is no commit yet — and "ready" has none. |
| `GTSTG-B10` | The seeded ignore list covers the system clutter and the Anchors working area, with no exception inside that area and nothing that guesses the stack. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `GTSTG-I01` | "Not installed" and "not initialised" are never merged: the same empty folder gives the two different states depending only on whether git is installed, and only one of them has something to offer. | classifies one isolated empty folder with and without git and compares the states and their offers |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `GTSTG-X01` | Does not guess the project's stack in the seeded ignore list. | The language has not been read yet at this point of `init`, and an ignore list with rules of the wrong stack is worse than a short one. |
| `GTSTG-X02` | Does not look up the git program itself; the caller says whether it is installed. | The suite must cover the "no git" case even on a machine that has git. |

## Errors

none — a missing repository marker, an unreadable branch folder or a missing packed list are the states themselves ("not initialised", "no commit yet"), not failures to report; the unit returns a state for every folder.

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |

none — the unit reads only the folder tree and the answer the caller passes in.

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
