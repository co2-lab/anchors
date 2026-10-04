# language: en
# @anchors
#   code: GEFGT
#   ref: GHEGT
#   updated_at: 2026-10-03
#   layer: feature

@GHEGT
Feature: GitHubEnvironment — the doctor warns, before the work starts, about the pieces the GitHub flow silently needs

  @GHEGT-B01 @unit-level
  Scenario: Local mode checks nothing
    Given a configuration without GitHub mode and an empty project
    When the GitHub environment is checked
    Then there is no finding

  @GHEGT-B02 @unit-level
  Scenario: Each missing pipeline gives a warning that names what stops happening
    Given a project in GitHub mode with none of the flow's pipelines on disk
    When the pipelines are checked
    Then there is one warning pipeline-ausente per flow pipeline
    And each one says what does not happen and points at --fix

  @GHEGT-B03 @unit-level
  Scenario: A pipeline without serialization is its own finding, not a missing one
    Given the flow's pipelines seeded and the claim pipeline rewritten without concurrency
    When the pipelines are checked
    Then there is exactly one finding, pipeline-sem-serializacao, saying the same card can reach two agents

  @GHEGT-B04 @unit-level
  Scenario: A pipeline behind its template is reported only while it carries the marker
    Given a seeded pipeline whose content drifted from the template with the marker intact
    When the pipelines are checked
    Then a pipeline-desatualizado warning is given
    And once the marker is removed the same drift is not reported

  @GHEGT-B05 @unit-level
  Scenario: An unreadable branch protection is an unprotected branch
    Given a platform that fails the read of the branch protection of acme/exemplo
    When the branch protection is checked
    Then there is one warning main-sem-protecao on acme/exemplo saying a direct push skips the card

  @GHEGT-B06 @unit-level
  Scenario: A protection without required reviews is partially protected
    Given a platform that answers the protection does not require pull request reviews
    When the branch protection is checked
    Then there is one warning main-sem-protecao with the partially protected message
    And a protection that requires reviews gives no finding

  @GHEGT-B07 @unit-level
  Scenario: Without the platform CLI the branch protection is not asked
    Given no platform CLI on the PATH
    When the branch protection is checked
    Then there is no finding

  @GHEGT-B08 @unit-level
  Scenario: The protection is read on the declared integration branch
    Given a GitHub-mode project whose integration branch is "develop", on a platform where "develop" requires reviews and "main" does not
    When the branch protection is checked, and again after "develop" loses its protection
    Then the first check gives no finding, and the second gives one whose message names "`develop`"

  @GHEGT-I01 @unit-level
  Scenario: After the fix seeds the pipelines no pipeline finding remains
    Given an empty project in GitHub mode that has pipeline findings
    When the flow's pipelines are seeded and the pipelines are checked again
    Then there is no pipeline finding

  @GHEGT-X01 @unit-level
  Scenario: The board is never charged
    Given a project in GitHub mode with the flow's pipelines seeded
    When the GitHub environment is checked
    Then no finding names the board or the project mirror
