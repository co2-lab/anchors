# language: en
# @anchors
#   ref: FLSCF
#   updated_at: 2026-09-30
#   layer: feature

@FLSCF
Feature: FlagScenarios — the scenarios a feature flag declares are written, complete, cited and tested

  @FLSCF-B01 @unit-level
  Scenario: The flag gates skip what is not a flag, and a flag with no scenario
    Given a spec node, a flag file with no scenario table, and, for the citation gate, a flag node
    When each flag gate confronts them
    Then every gate returns Skip

  @FLSCF-B02 @unit-level
  Scenario: A flag whose every condition is in the grammar passes
    Given the checkout flag with the conditions `= "off"`, `= "on"` and `absent`
    When flag-scenario-grammar confronts it
    Then it returns Pass

  @FLSCF-B03 @unit-level
  Scenario: A condition in prose fails the grammar, naming the scenario
    Given the checkout flag whose scenario CHKUT-G01 has the condition "when the user is a beta tester"
    When flag-scenario-grammar confronts it
    Then it returns Fail, and the message names CHKUT-G01

  @FLSCF-B04 @unit-level
  Scenario: A flag that declares the absent case passes completeness
    Given the checkout flag with a scenario whose condition is `absent`
    When flag-scenarios-complete confronts it
    Then it returns Pass

  @FLSCF-B05 @unit-level
  Scenario: A flag without the absent case fails, naming the waiver
    Given the checkout flag with only an `= "on"` scenario
    When flag-scenarios-complete confronts it
    Then it returns Fail, and the message names `@no-absent`

  @FLSCF-B06 @unit-level
  Scenario: The absent waiver needs a written reason, outside backticks
    Given the checkout flag with only an `= "on"` scenario, and in turn a waiver with a reason, a bare waiver, and a waiver quoted in backticks
    When flag-scenarios-complete confronts each
    Then only the waiver with a reason passes; the bare and the quoted ones fail

  @FLSCF-B07 @unit-level
  Scenario: Only a citation of a G code is confronted
    Given a spec with no citation, and a spec citing `@gated-by CRED-B03`
    When flag-scenario-exists confronts each
    Then both return Skip

  @FLSCF-B08 @unit-level
  Scenario: A citation of a declared scenario passes
    Given a project whose checkout flag declares CHKUT-G02, and a spec citing `@gated-by CHKUT-G02`
    When flag-scenario-exists confronts the spec
    Then it returns Pass

  @FLSCF-B09 @unit-level
  Scenario: Citations of scenarios that do not exist fail, each named once and sorted
    Given a project whose checkout flag declares G01 to G03, and a spec citing CHKUT-G98 twice and CHKUT-G97 once
    When flag-scenario-exists confronts the spec
    Then it returns Fail, and the message lists "CHKUT-G97, CHKUT-G98" with CHKUT-G98 appearing once

  @FLSCF-B10 @unit-level
  Scenario: Without a map the governance gate is Pending
    Given the checkout flag and no map
    When flag-scenario-governs confronts it
    Then it returns Pending

  @FLSCF-B11 @unit-level
  Scenario: A scenario no rule cites fails governance
    Given the checkout flag whose scenario CHKUT-G01 alone has an incoming citation
    When flag-scenario-governs confronts it
    Then it returns Fail naming CHKUT-G02 and CHKUT-G03, and not CHKUT-G01

  @FLSCF-B12 @unit-level
  Scenario: A scenario with a reasoned governance waiver is not charged
    Given a flag whose only scenario, uncited, ends its outcome with a governance waiver and a reason, and in turn with a bare waiver
    When flag-scenario-governs confronts each
    Then the reasoned waiver passes and the bare one fails

  @FLSCF-B13 @unit-level
  Scenario: Without a map the coverage gate is Pending
    Given the checkout flag and no map
    When flag-covered confronts it
    Then it returns Pending

  @FLSCF-B14 @unit-level
  Scenario: Coverage is judged per scenario
    Given the checkout flag whose ingested proven codes hold only CHKUT-G01, and in turn all three
    When flag-covered confronts it
    Then the first fails naming CHKUT-G02 and CHKUT-G03 and not CHKUT-G01, and the second passes

  @FLSCF-B15 @unit-level
  Scenario: A scenario no test names fails as having no test
    Given the checkout flag, and a test file that names no scenario, and in turn a test file naming CHKUT-G01 only in a comment
    When flag-covered confronts it
    Then it returns Fail naming every scenario, and the commented code is not read as a written test

  @FLSCF-B16 @unit-level
  Scenario: A written test not yet proven is told apart: not ingested, or ingested and not green
    Given a test file naming CHKUT-G01, first with no ingested execution and then with an ingestion that proved only G02 and G03
    When flag-covered confronts the checkout flag
    Then the first message says the test was written but no execution ingested, and the second says it did not pass

  @FLSCF-I01 @unit-level
  Scenario: A gate that could not measure never answers Pass
    Given the checkout flag and no map
    When flag-scenario-governs and flag-covered confront it
    Then neither returns Pass

  @FLSCF-X01 @unit-level
  Scenario: No waiver is accepted without a written reason
    Given a flag with a bare `@no-absent`, and a flag with a bare `@no-govern`
    When their gates confront them
    Then both fail

  @FLSCF-B17 @unit-level
  Scenario: A gated-by citation is read at the code length the project declares
    Given a project that declares code length 7 and a flag that declares no scenario "CHKUTXY-G99"
    When flag-scenario-exists confronts a spec citing "@gated-by CHKUTXY-G99"
    Then it returns Fail naming "CHKUTXY-G99"

  @FLSCF-E02 @unit-level
  Scenario: Flags that cannot be read leave the citation Pending, naming the flags folder
    Given a spec citing "CHKUT-G02" and a "flags" folder that cannot be read
    When flag-scenario-exists confronts the spec
    Then it returns Pending with a message that names the flags folder, not the missing map

  @FLSCF-B18 @unit-level
  Scenario: With a tests source a flag scenario is written only when a title cites it
    Given a flag scenario code that appears in a test's fixture but in no title
    When flag-covered runs with the project's tests declared, and without
    Then with the declaration the scenario has no test, and without it the scenario is written
    And a file the source lists no test in, such as a YAML flow under a Jest pattern, is read as without a source

  @FLSCF-E03 @unit-level
  Scenario: A failing tests source fails flag-covered naming the error
    Given a project whose tests script exits with an error
    When flag-covered runs
    Then it fails naming the script's error

  @FLSCF-B19 @unit-level
  Scenario: A flag scenario is green only by its own proof
    Given a flag whose scenarios are proven only through a variant of the first one
    When flag coverage confronts it
    Then the first scenario is not green, as no case proved exactly it
