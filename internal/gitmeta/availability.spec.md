<!-- @anchors
  code: GTAVG
  updated_at: 2026-09-26
  layer: infra
-->
# GitAvailability — why a git operation cannot happen, named with its fix

> **Code**: `GTAVG`

## Overview

"Could not use git" has causes with opposite fixes: git may not be installed, or the project may not be
under a repository. A raw error such as `exit status 128` sends the user to investigate the wrong layer.
This unit classifies the project root into one of three answers (git available, no git binary, no
repository) and turns the answer into the sentence that names what the command was about to do, the
cause and the fix, for the command to embed in its error.

The classification is cheap (a lookup on the PATH and a look for a `.git` entry), so it is called at the
point of use, right before the operation, where the command knows which action is left incomplete. A
subfolder of a repository is under git: the `.git` entry may be in any ancestor folder, and complaining
there would send the user to create a nested repository. When git is available, the failure is
something else, and no git cause is invented.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the project root | any folder path | — | the caller; a folder with no `.git` up to the filesystem root is "no repository" |
| the action | the words of what the command was about to do | an empty action | the calling command |

## Effects

| Effect | Description |
| --- | --- |
| `GTAVG-B01` | When `git` is not on the PATH, the answer is "no git binary", whatever the folder holds. (`Check`) |
| `GTAVG-B02` | With git on the PATH, a `.git` entry in the root or in any ancestor folder makes the answer "available". |
| `GTAVG-B03` | With git on the PATH and no `.git` entry up to the filesystem root, the answer is "no repository". |
| `GTAVG-B04` | The explanation for "no git binary" names the action and says git is not installed (or not on the PATH), and never sends the user to `git init`. (`Explain`) |
| `GTAVG-B05` | The explanation for "no repository" names the action and sends the user to `git init`, and never says git is missing. |
| `GTAVG-B06` | When git is available the explanation is empty, so no git cause is invented for a failure of another kind. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `GTAVG-E01` | REF[GTAVG-B01]: git missing from the PATH is the failure B01 classifies, with its fix in B04 | — | — |
| `GTAVG-E02` | REF[GTAVG-B03]: a root outside any repository is the failure B03 classifies, with its fix in B05 | — | — |

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
