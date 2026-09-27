# language: en
# @anchors
#   ref: FLPRF
#   updated_at: 2026-09-26
#   layer: feature

@FLPRF
Feature: FlagParse — reading a project's flag files and the scenarios their values open

  @FLPRF-B01 @unit-level
  Scenario: The G rows of the table are the scenarios
    Given the flag file "flags/new-checkout.flag.md" with four G rows, and a file with one B row and one G row
    When they are parsed
    Then the first gives four scenarios, G01 being eq "off" and G03 gte "50", and the second only "CHKUT-G01"

  @FLPRF-B02 @unit-level
  Scenario: The flag's code and name
    Given the flag file "flags/new-checkout.flag.md" whose scenarios are CHKUT-G01 to CHKUT-G04
    When it is parsed
    Then its code is "CHKUT" and its name "new-checkout"

  @FLPRF-B03 @unit-level
  Scenario: Each scenario records its line
    Given the flag file "flags/new-checkout.flag.md"
    When it is parsed
    Then each scenario has a line, increasing down the table

  @FLPRF-B04 @unit-level
  Scenario: A refused condition becomes a finding
    Given a row "CHKUT-G01" whose condition is "when the user is a beta tester"
    When it is parsed
    Then the scenario is kept and carries the refusal

  @FLPRF-B05 @unit-level
  Scenario: The flag says whether it declares the absent case
    Given the new-checkout flag with an "absent" row, and a flag with only `= "on"`
    When each is asked whether it declares the absent case
    Then the first says yes and the second no

  @FLPRF-B06 @unit-level
  Scenario: Flags load in a stable order, and none without a folder
    Given a project without "flags/", and a project with "zeta", "alpha" and "middle" flag files
    When each is loaded
    Then the first gives nothing and no error, and the second gives "alpha", "middle", "zeta"

  @FLPRF-B07 @unit-level
  Scenario: Scenarios are indexed by code
    Given the new-checkout flag
    When its scenarios are indexed
    Then "CHKUT-G03" resolves to the gte scenario and "CHKUT-G99" to nothing
