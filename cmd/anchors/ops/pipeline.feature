# language: en
# @anchors
#   code: PCFPA
#   ref: PCSPA
#   updated_at: 2026-11-06
#   layer: feature

@PCSPA
Feature: Pipeline — ask whether a stage may run, and record that it passed

  @PCSPA-B01 @unit-level
  Scenario: can-run passes a stage on its branch when every requirement holds, naming each, and refuses it otherwise
    Given a review stage requiring the flow and the ci gates, and a sha merged from develop into review
    When can-run asks on review, on production, with failing gates, and for a sha merged from feat/x
    Then the first passes naming each requirement, and the others fail naming the branch, the gates and the flow

  @PCSPA-B02 @unit-level
  Scenario: passed: asks whether the head the sha was merged from passed that stage; from: whether it came from that head
    Given a production sha merged from review whose head has no review record, then has one, then a sha merged from hotfix/login
    When can-run asks for the production stage
    Then the first is refused, and the second and the hotfix may run

  @PCSPA-B03 @unit-level
  Scenario: done runs verify: and records the stage's outcome on the sha as a commit status
    Given a review stage whose verify passes, then fails
    When done runs for it
    Then the sha carries anchors/pipeline/review-release as success, then as failure, and the second done fails

  @PCSPA-B04 @unit-level
  Scenario: run confronts the stage, runs its run: and its verify:, and records the outcome
    Given a lint stage with run and verify, whose run then fails
    When run runs for it, and for a stage can-run refuses
    Then run and verify ran and success is recorded, then only run ran and failure is recorded, and the refused stage does not run

  @PCSPA-E01 @unit-level
  Scenario: A stage the configuration does not declare does not run
    Given a pipeline without the stage deploy
    When can-run asks for deploy
    Then it fails naming the declared stages

  @PCSPA-E02 @unit-level
  Scenario: A stage with no branch or sha, by flag or environment, does not run
    Given no branch or sha by flag or environment
    When can-run asks for a stage
    Then it fails naming the flags and the variables

  @PCSPA-E03 @unit-level
  Scenario: A requirement GitHub cannot answer refuses the stage; a stage with no verify: is not done
    Given GitHub that cannot be asked, and a stage with only run
    When can-run asks for the production stage, and done for the stage with only run
    Then the stage is refused naming why, and done points to pipeline run

  @PCSPA-B05 @unit-level
  Scenario: passed: holds when the sha itself carries the stage's success, before asking the head it was merged from
    Given a production sha merged from review, and a stage requiring passed: production-release
    When can-run asks before and after the production release is recorded on that sha
    Then it is refused before, and may run after, naming this commit
