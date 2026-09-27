# language: en
# @anchors
#   ref: ARCGN
#   updated_at: 2026-09-26
#   layer: feature

@ARCGN
Feature: AgentRoleCLI — who this agent is, and the role it declared, as the commands show and ask it

  @ARCGN-B01 @unit-level
  Scenario: The agent identity falls back from the session to the user to default
    Given ANCHORS_SESSION is "  worker-2 " and USER is "alice"
    When the agent identity is read, then read again with ANCHORS_SESSION empty, then with USER empty too
    Then it is "<host>/worker-2", then "<host>/alice", then "<host>/default"

  @ARCGN-B02 @unit-level
  Scenario: The role list shows every known role
    Given the catalogue of known roles
    When the role list is built
    Then every role name and what it does appears in it

  @ARCGN-B03 @unit-level
  Scenario: Showing a role says whether it decides the product and shows its lens
    Given the product owner, the developer and the security reviewer roles
    When each role is shown
    Then the product owner is told it acts on "`needs-user`" cards, the developer is told to use "anchors escalate", and the security reviewer's lens is printed

  @ARCGN-B04 @unit-level
  Scenario: An unrecognised answer is echoed back and the question is asked again
    Given the terminal answers "wizard" and then "po"
    When the role is asked
    Then the output repeats "wizard" in quotes and the role is the product owner

  @ARCGN-B05 @unit-level
  Scenario: Three unrecognised answers give up
    Given the terminal answers "a", "b", "c" and then "po"
    When the role is asked
    Then an error pointing at "anchors settings role" is returned and "po" is never read

  @ARCGN-B06 @unit-level
  Scenario: A pipe or the null device is not an interactive terminal
    Given the input is a pipe, and then the null device
    When the terminal is checked for interactivity
    Then both answers are false

  @ARCGN-B07 @unit-level
  Scenario: Only a role that handles escalated cards decides the product
    Given a project with no settings, then with the product owner role, then with the developer role
    When the agent is asked whether it decides the product
    Then the answers are false, true and false

  @ARCGN-X01 @unit-level
  Scenario: The role list is the settings catalogue, not a copy
    Given the catalogue of known roles held by the settings unit
    When the role list is built
    Then each catalogue entry, and the description the catalogue gives it, is what appears

  @ARCGN-E01 @unit-level
  Scenario: A closed input while asking is an error
    Given the input is closed before any answer
    When the role is asked
    Then an error is returned and no role

  @ARCGN-E02 @unit-level
  Scenario: An unreadable settings file does not unlock the product decision
    Given a project whose settings file is not valid YAML
    When the agent is asked whether it decides the product
    Then the answer is false
