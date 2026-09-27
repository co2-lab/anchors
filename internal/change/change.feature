# language: en
# @anchors
#   ref: CHRCC
#   updated_at: 2026-09-26
#   layer: feature

@CHRCC
Feature: ChangeRecord — the delivery record an agent leaves when it finishes a stage

  @CHRCC-B01 @unit-level
  Scenario: The key joins the stage and the normalised unit
    Given a delivery of stage "code" over the unit "internal/gate/mock stamped.go"
    When its key is computed
    Then the key is "code--internal-gate-mock-stamped-4dc4c23f"

  @CHRCC-B02 @unit-level
  Scenario: A pending record lives in the changes folder under its key
    Given a delivery of stage "code" over the unit "internal/gate/mock stamped.go"
    When its path under a project root is computed
    Then the path is "changes/code--internal-gate-mock-stamped-4dc4c23f.md" under that root

  @CHRCC-B03 @unit-level
  Scenario: The header carries the stage, the unit, the date and the agent only when named
    Given a delivery of stage "test" over "internal/x/y.go" dated 2026-09-26 by agent "worker-2"
    And a second delivery that names no agent
    When both are rendered
    Then the first opens with a change-layer header carrying stage, unit, date and "agent: worker-2"
    And the second has no agent line

  @CHRCC-B04 @unit-level
  Scenario: The record states the intent and the touched files
    Given a delivery whose intent is "  Proves the Y rules.  " touching two files
    When it is rendered
    Then "What was done" holds "Proves the Y rules." and "Files" lists both files

  @CHRCC-B05 @unit-level
  Scenario: Empty decision and proof sections are still written
    Given a delivery with no decisions and nothing unproven
    When it is rendered
    Then the decisions section says "none — …" and the proof section says "nothing — …"

  @CHRCC-B06 @unit-level
  Scenario: Saving the same stage and unit again replaces the record
    Given a project with no changes folder
    And a delivery of stage "spec" over "b.go" saved once
    When the same stage and unit is saved again with a new intent
    Then it is written at the same path and the pending list still has one record for it

  @CHRCC-B07 @unit-level
  Scenario: The pending list is the sorted markdown files of the changes folder
    Given a changes folder holding two records, a text file and a subfolder named "sub.md"
    When the pending records are listed
    Then only the two records come back, in sorted order

  @CHRCC-B08 @unit-level
  Scenario: Marking a record reviewed moves it to the history under the same name
    Given a saved pending record
    When it is marked reviewed
    Then it now lives at "changes/reviewed/" with the same file name

  @CHRCC-B09 @unit-level
  Scenario: A second review of the same key keeps the first in the history
    Given a record of stage "code" over "a.go" already marked reviewed
    When a second delivery of the same stage and unit is saved and marked reviewed
    Then the history holds both, the second under a numbered name

  @CHRCC-B10 @unit-level
  Scenario: A line break or a comment close in a field cannot break the header
    Given a delivery whose unit is "a.go", a line break, "layer: forged", a line break and "-->"
    When it is rendered
    Then the header has no "layer: forged" line, the unit line carries it escaped, and the date line follows intact

  @CHRCC-I01 @unit-level
  Scenario: A reviewed record leaves the pending list and stays in the history
    Given two saved records
    When one of them is marked reviewed
    Then it is no longer in changes, it exists in the history, and only the other is pending

  @CHRCC-I02 @unit-level
  Scenario: Two different units never share a key
    Given the units "a/b.go", "a-b.go", "a/b.ts", "a.b.go" and "a b.go"
    When their keys under the same stage are computed
    Then the five keys are distinct, and "a\b.go" has the key of "a/b.go"

  @CHRCC-X01 @unit-level
  Scenario: A subfolder of changes is never listed as pending
    Given a changes folder holding a subfolder named "sub.md"
    When the pending records are listed
    Then the subfolder is not in the list

  @CHRCC-E01 @unit-level
  Scenario: A changes folder that cannot be read is an error
    Given a project where "changes" is a file, not a folder
    When the pending records are listed
    Then an error is returned

  @CHRCC-E02 @unit-level
  Scenario: A changes folder that cannot be created fails the save
    Given a project where "changes" is a file, not a folder
    When a delivery is saved
    Then an error is returned

  @CHRCC-E03 @unit-level
  Scenario: Marking a missing record reviewed fails
    Given no record named "missing.md" in changes
    When it is marked reviewed
    Then an error is returned
