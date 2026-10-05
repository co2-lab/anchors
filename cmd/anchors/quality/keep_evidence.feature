# language: en
# @anchors
#   code: KEFKP
#   ref: KPEVD
#   updated_at: 2026-10-05
#   layer: feature

@KPEVD
Feature: KeepEvidence — a change that proves nothing new keeps the files' evidence

  @KPEVD-B01 @unit-level
  Scenario: The evidence moves to the current content with the reason
    Given a proven spec whose rule line gained a @realizes, with the map still holding its proof
    When its evidence is kept with a reason
    Then the proof is at the current content and the reason is recorded

  @KPEVD-B02 @unit-level
  Scenario: A rebuilt map's lost evidence comes from HEAD
    Given the same change after a map build dropped the proof
    When its evidence is kept
    Then the proof comes back from the map at HEAD

  @KPEVD-B03 @unit-level
  Scenario: Nothing to keep is said
    Given a file whose evidence is already at its content
    When its evidence is kept
    Then the command says there is nothing to keep

  @KPEVD-B04 @unit-level
  Scenario: The contract stamps are refreshed under the same declaration
    Given a file kept
    When the command ends
    Then it lists the @contract stamps pointing at it under the same declaration

  @KPEVD-E01 @unit-level
  Scenario: A missing reason, file or map fails
    Given no reason, a file outside the map, and a project with no map
    When the command runs with each
    Then it fails saying which

  @KPEVD-B05 @unit-level
  Scenario: The suites a later run did not replace come back from HEAD and are carried with the rest
    Given a spec proven at HEAD by a unit suite and an e2e suite, changed, the map rebuilt and only the unit suite run again
    When keep-evidence runs on the spec
    Then the e2e proof comes back from HEAD, and both suites are at the new content
