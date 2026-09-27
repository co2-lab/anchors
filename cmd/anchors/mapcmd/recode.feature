# language: en
# @anchors
#   ref: RCDEO
#   updated_at: 2026-09-26
#   layer: feature

@RCDEO
Feature: Recode — renames an identity code and carries the change to every textual surface of the project

  @RCDEO-B01 @unit-level
  Scenario: Codes typed in lower case are renamed in upper case
    Given a project whose login spec has the code LOGIN
    When the recode runs from "login" to "signn"
    Then the report names "recode LOGIN → SIGNN:"

  @RCDEO-B02 @unit-level
  Scenario: The plan counts each file's occurrences by kind
    Given a project whose login spec has the code LOGIN in its header and a scenario code
    When the recode runs from LOGIN to SIGNN without applying
    Then the plan lists the login spec with a header line and a scenario-code line

  @RCDEO-B03 @unit-level
  Scenario: A recode without the apply switch writes nothing
    Given a project whose login spec has the code LOGIN
    When the recode runs from LOGIN to SIGNN without applying
    Then the login spec is unchanged on disk
    And the report ends with "(dry-run — nothing was written"

  @RCDEO-B04 @unit-level
  Scenario: Applying rewrites the header and the scenario codes of the spec and the test
    Given a project whose login spec and login test carry LOGIN-B01
    When the recode runs from LOGIN to SIGNN with the apply switch
    Then the spec declares "code: SIGNN" and "SIGNN-B01", the test carries "SIGNN-B01"
    And the report says "file(s) rewritten"

  @RCDEO-B05 @unit-level
  Scenario: Applying rebuilds the map from the headers
    Given a project of five files whose map gives the login spec the code LOGIN
    When the recode runs from LOGIN to SIGNN with the apply switch
    Then the rebuilt map gives the login spec the code SIGNN
    And the report says "map rebuilt (5 nodes)"

  @RCDEO-B06 @unit-level
  Scenario: Applying keeps the judgments, stamps and flow of the previous map
    Given a project whose map has a judged and stamped edge on the login spec and a flow graph
    When the recode runs from LOGIN to SIGNN with the apply switch
    Then the rebuilt map still has that edge's judgment and stamp, and the flow graph

  @RCDEO-I01 @unit-level
  Scenario: After applying, neither the spec nor the map carries the old code
    Given a project whose login spec has the code LOGIN
    When the recode runs from LOGIN to SIGNN with the apply switch
    Then the word LOGIN appears nowhere in the login spec and no map node has the code LOGIN

  @RCDEO-X01 @unit-level
  Scenario: The map is rebuilt from the files, not edited as text
    Given a project whose map was built before the recode
    When the recode runs from LOGIN to SIGNN with the apply switch
    Then the map's login spec node carries the revision of the rewritten file

  @RCDEO-E01 @unit-level
  Scenario: A project with no configuration is refused naming the configuration file
    Given a directory with no anchors.yaml
    When the recode runs from LOGIN to SIGNN
    Then it fails with a message naming "anchors.yaml"

  @RCDEO-E02 @unit-level
  Scenario: Recoding a code no file carries is refused
    Given a project where no file carries the code NOPEX
    When the recode runs from NOPEX to SIGNN
    Then it fails saying the code does not appear in any file

  @RCDEO-E03 @unit-level
  Scenario: A write failure after some files changed says the project is half converted
    Given a project whose login spec is writable and whose login test is read-only
    When the recode runs from LOGIN to SIGNN with the apply switch
    Then it fails
    And the report says "1 file(s) had ALREADY been changed" and "half converted"

  @RCDEO-E04 @unit-level
  Scenario: A recode with a single code is refused
    Given a project whose login spec has the code LOGIN
    When the recode runs with only "LOGIN"
    Then it is refused for the number of arguments
