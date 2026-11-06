# language: en
# @anchors
#   code: PCFPP
#   ref: PCSPP
#   updated_at: 2026-11-06
#   layer: feature

@PCSPP
Feature: Pipeline — the stages a workflow asks about, what each requires, and what validates it

  @PCSPP-B01 @unit-level
  Scenario: A stage declares its branches, its requirements and the command that runs or validates it
    Given a pipeline with a review stage requiring the flow and the ci gates, and a production stage requiring the flow and one of a review pass or a hotfix head
    When the configuration is loaded
    Then the stages, their requirements and the unguarded workflows are read as declared, and the workflows default to GitHub's

  @PCSPP-E01 @unit-level
  Scenario: A stage with nothing that runs or validates it, a passed: naming no stage, or an unknown requirement refuses the configuration
    Given stages with no run or verify, a passed naming no stage, an unknown requirement and bad branch patterns
    When the configuration is loaded
    Then it does not load, naming the stage
