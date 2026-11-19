# language: en
# @anchors
#   code: PCFPP
#   ref: PCSPP
#   updated_at: 2026-11-18
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

  @PCSPP-B02 @unit-level
  Scenario: A stage asks of the version its sha carries, and of a branch whose head the sha contains
    Given a review stage requiring contains production and version of app fresh, listed for target review and expected, and a production stage requiring same_as_passed review-release
    When the configuration is loaded
    Then each requirement is read and written back as declared, the product left out where it was

  @PCSPP-E02 @unit-level
  Scenario: A version asked of no product, with nothing asked, of a stage or a listing that do not exist, or a contains: pattern refuses the configuration
    Given version of web, version asking nothing, same_as_passed qa, a key latest, a listing key match, contains hotfix/*, listed with no file nor released:, a version with no versioning block, and two products unnamed
    When the configuration is loaded
    Then it does not load, naming the stage and the field

  @PCSPP-B03 @unit-level
  Scenario: A stage names the product it releases
    Given production-release with release app
    When the configuration is loaded
    Then the stage releases app

  @PCSPP-E03 @unit-level
  Scenario: A stage that releases a product versioning: does not declare refuses the configuration
    Given production-release with release web, and only app declared
    When the configuration is loaded
    Then it does not load, naming the stage and web
