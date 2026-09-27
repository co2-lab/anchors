# language: en
# @anchors
#   ref: SRRGS
#   updated_at: 2026-09-26
#   layer: feature

@SRRGS
Feature: SourceRegion — the declared interval of a rule code in source, and the composition of test scripts

  @SRRGS-B01 @unit-level
  Scenario: A nested region closes before the outer one
    Given the region "MLETX-A03" on lines 2 to 8 containing the region "MLETX-B05" on lines 4 to 6
    When the regions are extracted
    Then there is no defect, "MLETX-B05" comes first with lines 4 to 6, and "MLETX-A03" follows with lines 2 to 8

  @SRRGS-B02 @unit-level
  Scenario: A deeply indented one-line region does not swallow the file
    Given a deeply indented region of three lines followed by 500 other lines
    When the regions are extracted
    Then there is one region covering exactly 3 lines

  @SRRGS-B03 @unit-level
  Scenario: The region revision ignores changes outside it
    Given two texts with the same region and a different line after it, and a third with a different line inside it
    When the regions are extracted from each
    Then the first two revisions are equal and the third differs

  @SRRGS-B04 @unit-level
  Scenario: A swapped close does not cascade into the following regions
    Given a region "CODEX-A01" closed by "CODEX-B02" followed by a well-formed region "CODEX-A02"
    When the regions are extracted
    Then there is exactly one "fecho-trocado" defect and both regions come out

  @SRRGS-B05 @unit-level
  Scenario: A close without a code closes the open region
    Given the region "CODEX-A01" closed by a bare "#endregion"
    When the regions are extracted
    Then there is one region "CODEX-A01" from line 1 to line 3 and no defect

  @SRRGS-B06 @unit-level
  Scenario: A file without regions has no defect
    Given a text with an old-style code comment and no region marker
    When the regions are extracted
    Then there is no region and no defect

  @SRRGS-B07 @unit-level
  Scenario: Composition resolves relative to the script's directory
    Given the script "apps/mobile/.maestro/suites/smoke/SS-03.yaml" running "../../utils/login.yaml", "../../utils/dismissOsDialogs.yaml" through "file:", and a runScript
    When its composition is read
    Then it is "apps/mobile/.maestro/utils/login.yaml" and "apps/mobile/.maestro/utils/dismissOsDialogs.yaml", in that order

  @SRRGS-B08 @unit-level
  Scenario: Duplicates and self-references are not dependencies
    Given the script "a/b/flow.yaml" running "../utils/login.yaml" twice, a runScript, and "./flow.yaml"
    When its composition is read
    Then it is only "a/utils/login.yaml", once

  @SRRGS-B09 @unit-level
  Scenario: A script with no composition has no dependency
    Given a script that only taps an element
    When its composition is read
    Then there is no dependency

  @SRRGS-I01 @unit-level
  Scenario: One pairing defect is reported once
    Given a region closed by the wrong code followed by a well-formed region
    When the regions are extracted
    Then exactly one defect is reported and both regions come out

  @SRRGS-X01 @unit-level
  Scenario: Composition comes only from declared paths
    Given a script whose steps mention other scripts only by name in a tap step
    When its composition is read
    Then there is no dependency

  @SRRGS-E01 @unit-level
  Scenario: A region never closed is reported
    Given a region "CODEX-A01" opened and never closed
    When the regions are extracted
    Then there is exactly one defect of kind "sem-fecho"

  @SRRGS-E02 @unit-level
  Scenario: A close with no open region is reported
    Given a close "CODEX-A01" with no region open
    When the regions are extracted
    Then there is exactly one defect of kind "fecho-orfao"

  @SRRGS-E03 @unit-level
  Scenario: A close naming another code is reported
    Given the region "CODEX-A01" closed by "CODEX-B02"
    When the regions are extracted
    Then there is exactly one defect of kind "fecho-trocado"
