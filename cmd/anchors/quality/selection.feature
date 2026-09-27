# language: en
# @anchors
#   ref: SLCTN
#   updated_at: 2026-09-27
#   layer: feature

@SLCTN
Feature: RunSelection — a run takes only what is stale and below the minimum, unless told otherwise

  @SLCTN-B01 @unit-level
  Scenario: A file is placed in one of four boxes
    Given results fresh and stale, passing and below the minimum, and a file never measured
    When each is placed
    Then each lands in its box and the one with no result is never measured

  @SLCTN-B02 @unit-level
  Scenario: The default takes stale below the minimum and never measured, and each flag opens a side
    Given one file in each box and one never measured
    When the selection is made with no flag, with each flag, with both, and with skip-unmeasured
    Then no flag takes stale-below and never-measured, each flag adds its side, both add fresh passing, and skip-unmeasured drops the never measured

  @SLCTN-B03 @unit-level
  Scenario: A test file is stale through what it exercises and passes by its own layer
    Given a test whose exercised code changed, one that failed in integration but passed in unit, and one with only skipped cases
    When their states are read for the unit layer
    Then the first is stale, the second passes in unit, and the third is not passing

  @SLCTN-B04 @unit-level
  Scenario: A mutation result under load counts as stale and passes at the floor
    Given code files measured under load, at a report floor, below the default floor, and with no mutant run
    When their states are read
    Then the one under load is stale, the others pass or not by the floor, and the one with no mutant run passes

  @SLCTN-B05 @unit-level
  Scenario: A test file belongs to the suite that ran it
    Given a file timed by this suite, one timed only by another, one run only by this suite's layer, and one never run
    When ownership is asked for this suite
    Then the first, third and fourth are this suite's and the second is not

  @SLCTN-B06 @unit-level
  Scenario: Support files and no_signal targets never run
    Given a stale failing support file and a stale failing code file the mutation gate declares in no_signal
    When the selection is made
    Then neither is taken

  @SLCTN-B07 @unit-level
  Scenario: The run says what it selected and what it left out
    Given a suite with files in several boxes
    When it runs by default
    Then it says how many were selected and, by box, how many were left out with the flag that takes them

  @SLCTN-B08 @unit-level
  Scenario: A suite without run_changed runs whole
    Given a suite that declares no run_changed
    When it runs by default, and with a budget
    Then it runs whole saying why, and with a budget it is refused before anything runs

  @SLCTN-B09 @unit-level
  Scenario: A selection that takes nothing runs nothing
    Given a suite whose files are all fresh and passing
    When it runs by default
    Then it says there is nothing to run and no command runs

  @SLCTN-B10 @unit-level
  Scenario: The selected files run in as few batches as the ceiling allows
    Given files whose command line exceeds the ceiling together, and one longer than the ceiling alone
    When they are split into batches
    Then each batch stays under the ceiling and the long one runs alone

  @SLCTN-B11 @unit-level
  Scenario: --all runs whole and does not combine with the other choices
    Given a suite with run_changed and files in every box
    When it runs with --all, and with --all together with --changed or a state flag, and a state flag with --changed
    Then --all runs the whole command, and each combination is refused before anything runs

  @SLCTN-E01 @unit-level
  Scenario: A selective run without a map is refused
    Given a project whose map was never built
    When a suite with run_changed runs by default
    Then it is refused saying to build the map or to use --all
