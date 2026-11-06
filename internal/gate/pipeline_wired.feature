# language: en
# @anchors
#   code: PWFPP
#   ref: PWSPP
#   updated_at: 2026-11-06
#   layer: feature

@PWSPP
Feature: PipelineWired — the stages the workflows call are the stages the configuration declares

  @PWSPP-B01 @unit-level
  Scenario: The gate passes when every call names a declared stage, every stage is called and recorded, and every workflow calls one or is unguarded
    Given a review and a production workflow each asking and recording its stage, and an unguarded build workflow
    When the gate confronts them with the two declared stages
    Then it passes, naming three workflows and two stages

  @PWSPP-B02 @unit-level
  Scenario: A call to an undeclared stage, a stage no workflow calls, a stage asked and never recorded, and a workflow calling none fail the gate
    Given a review workflow recording deploy instead of its stage, no production workflow, and an admin workflow calling none
    When the gate confronts them with the two declared stages
    Then it fails naming each of the four gaps

  @PWSPP-B03 @unit-level
  Scenario: With no stage declared, or no workflow found, the gate is pending
    Given a project with no pipeline, and one with stages and no workflow
    When the gate runs
    Then it is pending, naming the globs
