<!-- @anchors
  code: CNSTC
  updated_at: 2026-10-08
  layer: comando
-->
# StampCommand — writes the missing contract stamps on test doubles, and refreshes them after a change

> **Code**: `CNSTC`

## Overview

A test double reproduces a real module, and the mock-stamped gate recomputes, for each double, a stamp that
records the block of the module it copies. This command writes those stamps.

In its plain form it only ADDS: above every double of a governed module that has no stamp it writes one,
and it lists, per test, what it wrote and which doubles it left unstamped and why. An existing stamp is
never rewritten, even when it diverges: a divergence is the gate saying the double may have drifted, and a
tool refreshing the hash would let the stamp certify itself. With no test named, every test of the map is
considered. A dry run says what it would write and writes nothing.

After a function changes, the refresh form is the other half of the loop: given the changed file it lists
every double that was stamped against the previous version — the test and line, the member, the old and new
hash, and the lines of the stamped block that left and that came compared to the last commit — and then
updates those stamps. The listing is the point: each double named reproduces the old contract and must be
adjusted in the same commit. A stamp whose member is gone is reported and left alone.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the test files | paths of tests, relative to the root or not; none means every test of the map | a path that cannot be read | this unit: an unreadable test is reported and skipped (`CNSTC-E04`) |
| the refresh files | paths of changed modules | a module no double is stamped against | this unit: it says no double is stamped against it (`CNSTC-B05`) |
| the configuration and the map | the project's configuration and map | a project with neither | this unit: it refuses (`CNSTC-E01`, `CNSTC-E02`) |

## Effects

| Effect | Description |
| --- | --- |
| `CNSTC-B01` | With no test named, every test of the map is considered, in sorted order; a test with nothing to stamp is not listed. |
| `CNSTC-B02` | Writes the missing stamp above each double and lists, under the test, each stamp written and each double left unstamped with its reason, closing with the totals of stamps, files and doubles left unstamped. |
| `CNSTC-B03` | The dry run says what it would write and leaves every file as it was. |
| `CNSTC-B04` | The refresh lists each double stamped against the previous version of the file (test and line, member, old and new hash, and the block's lines that left and came since the last commit), updates their stamps and closes telling the user to adjust each double in the same commit. |
| `CNSTC-B05` | The refresh of a file no double is stamped against says so and writes nothing. |
| `CNSTC-B06` | When the stamped block has no version in the last commit, the refresh says there is nothing to compare instead of a diff. |
| `CNSTC-B07` | The refresh in dry run lists the doubles and says nothing was written, leaving the stamps as they were. |
| `CNSTC-B08` | The block diff lists the lines of the old block missing from the new one, marked as left, then the lines of the new block missing from the old one, marked as came, counting repeated lines. |
| `CNSTC-B09` | Handed a test file that holds stamps instead of a module, the refresh names the modules its stamps point at and the command to refresh them, and changes nothing. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CNSTC-I01` | Stamping is idempotent: a second run over the same tests writes nothing. | stamps a fixture, runs again, and the second run reports zero stamps written |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CNSTC-X01` | The plain form never rewrites an existing stamp, even one whose hash diverges from the module. | A divergent stamp is the gate's signal that the double may be stale; rewriting it would erase that signal without anyone looking at the double. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CNSTC-E01` | The project has no configuration. | The command fails with the configuration error. | The stamp's hash follows the project's detection rules; without them there is nothing to compute. |
| `CNSTC-E02` | The map does not exist or cannot be read. | Error pointing at `anchors map build`. | The map says which files are tests and which modules they double. |
| `CNSTC-E03` | During a refresh, the member a stamp names no longer exists in the changed file. | The stamp is listed as NOT refreshed with the instruction to adjust the double, delete the stamp and stamp again; it is left unchanged. | A renamed or removed member cannot be re-hashed: guessing the new member would certify a double against code it does not copy. |
| `CNSTC-E04` | A test to stamp cannot be read. | It is reported as not readable and the run goes on with the other tests. | One missing file in the map must not stop the stamps of every other test. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
