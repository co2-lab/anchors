# language: en
# @anchors
#   code: APFTP
#   ref: APRCP
#   updated_at: 2026-10-03
#   layer: feature

@APRCP
Feature: ApprovalReachable — the doctor says when the required approval can never be given, and how to get out

  @APRCP-B01 @unit-level
  Scenario: Only an administrator whose protection spares administrators can bypass it
    Given the platform says the current account is an administrator of the repository
    And the branch protection does not enforce its rules on administrators
    When the bypass is asked
    Then it answers yes with no reason, and answers no when the protection enforces on administrators

  @APRCP-B02 @unit-level
  Scenario: Each refusal of the bypass names its reason
    Given the platform CLI is missing, or the user, the permission or the protection cannot be read, or the account is only a writer
    When the bypass is asked
    Then it answers no with the reason of the condition that failed

  @APRCP-B03 @unit-level
  Scenario: Zero required approvals or no configuration reports nothing
    Given a workflow that requires zero approvals, or no configuration at all
    When the reachability is checked
    Then there is no finding

  @APRCP-B04 @unit-level
  Scenario: Without the platform CLI the reachability check stays silent
    Given one required approval and no platform CLI on the PATH
    When the reachability is checked
    Then there is no finding, because another finding already names the missing tool

  @APRCP-B05 @unit-level
  Scenario: A required approval the account cannot bypass gives one warning on the repository
    Given one required approval on repository acme/x
    And the current account is a writer, not an administrator
    When the reachability is checked
    Then there is exactly one warning aprovacao-inalcancavel on acme/x carrying the count 1
    And an administrator with an escape gets no finding

  @APRCP-B06 @unit-level
  Scenario: The warning names both ways out
    Given one required approval the current account cannot bypass
    When the reachability is checked
    Then the warning mentions a service account, required_approvals: 0 and doctor --fix

  @APRCP-B07 @unit-level
  Scenario: Turning the requirement off sends a protection with zero approvals
    Given a platform that accepts the protection update
    When the requirement is turned off for acme/x on main
    Then one update of the main branch protection is sent
    And its body requires zero approving reviews and does not enforce on administrators

  @APRCP-I01 @unit-level
  Scenario: The warning appears exactly when the bypass is refused
    Given the same scripted platform answers for an administrator with an escape and for a writer
    When the bypass and the reachability are both asked
    Then the administrator gets a yes and no warning, and the writer gets a no and one warning

  @APRCP-X01 @unit-level
  Scenario: The reachability check never changes the branch protection
    Given one required approval the current account cannot bypass
    When the reachability is checked
    Then every call made to the platform is a read, and none updates the protection

  @APRCP-E01 @unit-level
  Scenario: A refused protection update names the branch and carries the platform output
    Given a platform that refuses every call with the output "no rule"
    When the requirement is turned off for main
    Then the error reads "main: no rule"
