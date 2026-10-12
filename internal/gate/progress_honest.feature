# language: en
# @anchors
#   code: PHFPR
#   ref: PRHNP
#   updated_at: 2026-10-12
#   layer: feature

@PRHNP
Feature: ProgressHonest — the progress file tells the truth about the disk

  @PRHNP-B01 @unit-level
  Scenario: An artifact that is not a plan leaves without a verdict
    Given a node whose kind is spec, code or test
    When the gate confronts it
    Then it returns Skip, because the gate is anchored on the plan, which is what the map holds

  @PRHNP-B02 @unit-level
  Scenario: An untracked plan is skipped, not failed
    Given a plan no demand of anchors.progress.md tracks
    When the gate confronts it
    Then it returns Skip, because "is the progress tracked" is a different question
      from "is the progress true"

  @PRHNP-B03 @unit-level
  Scenario: The skip for an untracked plan says how to track it
    Given a plan no demand of anchors.progress.md tracks
    When the gate confronts it
    Then the verdict names the command that tracks it, so the reader does not have to look
      it up

  @PRHNP-B04 @unit-level
  Scenario: A ticked item whose file does not exist is failed
    Given a progress item ticked as done and citing a path that is not on disk
    When the gate confronts the plan
    Then it returns Fail, because it declares done what is not

  @PRHNP-B05 @unit-level
  Scenario: An open item whose file already exists is failed
    Given a progress item left open and citing a path that is already on disk
    When the gate confronts the plan
    Then it returns Fail, because it produces rework — somebody redoes what is done

  @PRHNP-B06 @unit-level
  Scenario: A spec the plan seeds and the progress does not list is failed
    Given a plan whose checkbox items seed two specs
    And a progress file that lists only one of them
    When the gate confronts the plan
    Then it returns Fail, because a gate that only looks inside the file never sees
      what is missing from it

  @PRHNP-B07 @unit-level
  Scenario: A checkbox item promising no file at all is failed
    Given a progress item that is an untouched template marker with no path
    When the gate confronts the plan
    Then it returns Fail, because an eternal open box makes the plan look unfinished forever

  @PRHNP-B08 @unit-level
  Scenario: The ticked-but-absent finding is reported first
    Given a progress carrying both a ticked item with no file and an open item whose file exists
    When the gate confronts the plan
    Then the ticked-but-absent finding appears first, because whoever reads the board
      decides on it

  @PRHNP-B09 @unit-level
  Scenario: An item in prose citing no path is not charged
    Given progress items reading "review with the team" and "agreed at the daily"
    When the gate confronts the plan
    Then it returns Pass, because there is nothing to confront

  @PRHNP-B10 @unit-level
  Scenario: A progress that agrees with the disk on every item passes
    Given a ticked item whose file is on disk and an open item whose file is not
    When the gate confronts the plan
    Then it returns Pass

  @PRHNP-B11 @unit-level
  Scenario: The verdict names each offending path
    Given a progress item ticked as done and citing a path that is not on disk
    When the gate confronts the plan
    Then the verdict names that path, so the reader does not have to diff the file
      against the disk by hand

  @PRHNP-B12 @unit-level
  Scenario: The verdict carries only the directions that found something
    Given a progress whose only finding is of one direction
    When the gate confronts the plan
    Then the verdict carries that direction's finding and no heading for the directions
      that found nothing, because a zero-item finding reads as a defect that is not there

  @PRHNP-I01 @unit-level
  Scenario: The progress file's name has one definition, derived from the scanner
    Given the progress file's name as the gate reads it
    When it is compared with the scanner's
    Then they are the same, because a second constant here would diverge in silence and the
      gate would read a file nobody writes

  @PRHNP-I02 @unit-level
  Scenario: A seed is matched by path, never by the item's text
    Given a plan seeding a spec with a long description
    And a progress listing the same path under a shorter wording
    When the gate confronts the plan
    Then the seed is not accused as missing, because what identifies the item is the file

  @PRHNP-I03 @unit-level
  Scenario: A spec mentioned in the plan's prose is not a seed
    Given a plan whose revision paragraph names a spec without a checkbox
    And a progress listing every spec the plan actually seeds
    When the gate confronts the plan
    Then it returns Pass, because the checkbox is the promise and the prose speaks of
      what already exists

  @PRHNP-I04 @unit-level
  Scenario: A template file is never a seeded spec
    Given a plan whose checkbox item cites a template spec path
    And a progress that does not list it
    When the gate confronts the plan
    Then it returns Pass, because the mould is the shape work is poured into, not work to be done

  @PRHNP-X01 @unit-level
  Scenario: The gate does not charge that the plan is tracked
    Given a plan that predates the progress file and no demand tracks
    When the gate confronts it
    Then it does not return Fail, because merging "does it exist" with "is it true"
      would report two different debts as one finding

  @PRHNP-X02 @unit-level
  Scenario: The gate does not judge the content of an item beyond the path
    Given a progress item whose description contradicts what the cited file contains
    And that file is on disk and the item is ticked
    When the gate confronts the plan
    Then it returns Pass, because the ruler here is the disk, which needs no opinion

  @PRHNP-X03 @unit-level
  Scenario: The gate does not put the progress file into the map
    Given a progress file that changes on every delivery
    When the map is built
    Then the progress file is not a node, because a gate reaching it would charge every
      edit of the artifact that exists in order to change

  @PRHNP-E01 @unit-level
  Scenario: An unreadable progress file is pending
    Given an anchors.progress.md edited by hand
    When the gate confronts a plan
    Then it is pending, naming why, and does not fail — progress-intact fails the file

  @PRHNP-B13 @unit-level
  Scenario: The progress file is the one Anchors wrote
    Given a project with no progress file, one written by anchors progress, and one edited by hand
    When progress-intact runs on each
    Then it skips the first, passes the second and fails the third naming the hand edit

