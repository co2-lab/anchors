<!-- @anchors
  code: DLVRE
  updated_at: 2026-10-08
  layer: comando
-->
# Deliver — record what a stage delivered, where the reviewer reads it, and send the author to the review

> **Code**: `DLVRE`

## Overview

`anchors deliver` closes a stage of the cycle (spec, code, feature, test or plan) by recording what
was delivered: the unit, the files, the declared intent, the decisions the ruler did not make and the
author made alone, and what the author knows is not proved. The last two look optional and are not:
a silent decision is where most defects that cross the gates come from, and declaring them separates
assumed debt from forgetting. Measured in three rounds of a real end-to-end run: seven serious defects
passed with every gate green, and all seven came from adversarial review — which only happened because
someone remembered to ask for it.

Where the record lives depends on the workflow mode, and the two are exclusive. In local mode it is a
file under `changes/`, which the watcher sees and turns into a review task. In github mode the reviewer
reads the card, so the record is a comment on the card — found by the number given, or by the unit's
code in the map — and the command fails rather than falling back to the file, because a record on disk
is invisible to whoever reviews the card. A vendored pipeline owned upstream has no card in github mode,
so nothing is recorded and the output says why.

The command never reads the clock: the date is given by whoever records, so the record is reproducible.
After recording, it instructs the next step — the review of the unit — with the command ready, and then
confronts the declaration against the disk (see the delivery confrontation).

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the stage | spec, code, feature, test, plan | anything else; none | this unit: refuses |
| the unit | a path, relative to the project root or absolute, of a file that exists or of a unit one of whose pieces exists | a unit with no piece on disk | this unit: refuses and asks to check the path |
| the intent | non-blank text | empty or whitespace-only | this unit: refuses |
| the date | a date given by the caller | none | this unit: refuses; it never reads the clock |
| the decisions and gaps | free prose, commas included, repeatable | — | this unit: each flag value is one item |
| the card | a card number, in github mode | a closed card | the board: only an open card is found |

## Effects

| Effect | Description |
| --- | --- |
| `DLVRE-B01` | The command refuses without `--stage` or `--unit`, with a stage outside spec, code, feature, test and plan, with a blank `--intent`, and without `--date`, saying the tool does not read the clock. |
| `DLVRE-B02` | A unit that does not exist is accepted when a piece of the same unit exists (its spec, feature, test or code by the same stem); the output says the unit does not exist yet, names the piece found, and that the record keeps the unit given — which it does. The unit and every file given as an absolute path are recorded relative to the root. |
| `DLVRE-B03` | Each `--decision` and `--uncovered` value is one item even when it holds commas, while `--file` splits on commas. |
| `DLVRE-B04` | In local mode the record is written under `changes/` with the unit and the intent, and the output names the file. |
| `DLVRE-B05` | In github mode with `--card`, the record is a comment on that open card, and the output names the card. |
| `DLVRE-B06` | In github mode without `--card`, the card is the open card titled with the unit's code, which the map gives for the unit's exact file or, failing that, for any piece of the same unit. |
| `DLVRE-B07` | In github mode, a vendored pipeline owned upstream delivered without `--card` records nothing and succeeds, saying it is owned upstream. |
| `DLVRE-B08` | After recording, the output instructs the next step: `anchors work review --for <unit>`. |
| `DLVRE-B09` | The hint to start the watcher is printed only when the watcher is not running in the project. |
| `DLVRE-B10` | The note that zero free decisions and zero proof gaps were declared is printed only when both lists are empty. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `DLVRE-I01` | A refused delivery records nothing. | every refusal is run on the same project, and `changes/` does not exist afterwards |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `DLVRE-X01` | In github mode the record is never written under `changes/`. | A record that should be on the card and landed on disk is invisible to the reviewer — measured: 73 records in a repository, none reviewed. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `DLVRE-E01` | In github mode without `--card`, the map gives the unit no code. | The command fails saying the code was not found, to run `anchors map build`, or to name the card with `--card <n>`. | The code is what identifies the card; guessing would record on the wrong one. |
| `DLVRE-E02` | In github mode, no open card carries the unit's code. | The command fails naming the code, and says to reopen the card or to name another with `--card <n>`. | The record must land on an open card the reviewer reads. |
| `DLVRE-E03` | In github mode, the card given with `--card` is closed. | The command fails saying the card is closed. | A comment on a closed card reaches no reviewer. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
