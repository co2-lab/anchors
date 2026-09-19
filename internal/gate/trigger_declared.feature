# language: en
# @anchors
#   ref: TRDCT
#   updated_at: 2026-09-19
#   layer: feature

@TRDCT
Feature: TriggerDeclared — cited compliance triggers and obligations must exist in the declared vocabulary

  @TRDCT-B01 @unit-level
  Scenario: Non-spec artifacts skip confrontation
    Given an artifact node whose kind is not spec
    When the gate confronts it
    Then it returns Skip, because instructional trigger citations reside in specifications

  @TRDCT-B02 @unit-level
  Scenario: Specifications without cited triggers skip confrontation
    Given a specification file containing no cited obligation triggers
    When the gate confronts it
    Then it returns Skip, because citing compliance triggers is optional

  @TRDCT-B03 @unit-level
  Scenario: Cited triggers return Pending when no vocabulary is declared
    Given a specification citing an obligation trigger
    And a project configuration with no declared obligations or packs
    When the gate confronts it
    Then it returns Pending, because there is no vocabulary against which to verify

  @TRDCT-B04 @unit-level
  Scenario: Non-trigger key-value citations are ignored
    Given a specification citing ordinary metadata keys such as layer or code
    When the gate confronts it
    Then it ignores those keys and returns Skip, avoiding false alarms on general metadata

  @TRDCT-B05 @unit-level
  Scenario: Natural language mentions without backticks are ignored
    Given a specification mentioning compliance concepts in unquoted prose
    When the gate confronts it
    Then it ignores the text and returns Skip, evaluating only backticked symbol citations

  @TRDCT-B06 @unit-level
  Scenario: Valid declared trigger values pass confrontation
    Given a specification citing a trigger value declared in configuration or packs
    When the gate confronts it
    Then it returns Pass, confirming the cited trigger belongs to the active vocabulary

  @TRDCT-B07 @unit-level
  Scenario: Valid declared obligation names pass confrontation
    Given a specification citing an obligation name declared in configuration or packs
    When the gate confronts it
    Then it returns Pass, confirming the cited obligation exists in the active vocabulary

  @TRDCT-B08 @unit-level
  Scenario: An undeclared trigger value fails with a nearest-match suggestion
    Given a specification citing a trigger value absent from the declared vocabulary
    When the gate confronts it
    Then it returns Fail, offering a closest matching declared trigger as a fix hint

  @TRDCT-B09 @unit-level
  Scenario: An undeclared trigger without close match lists declared triggers
    Given a specification citing an unknown trigger with no substring similarity to declared triggers
    When the gate confronts it
    Then it returns Fail, listing available declared triggers

  @TRDCT-B10 @unit-level
  Scenario: An undeclared obligation name fails confrontation
    Given a specification citing an obligation name absent from declared obligations
    When the gate confronts it
    Then it returns Fail, reporting the obligation was not found in the vocabulary

  @TRDCT-B11 @unit-level
  Scenario: Duplicate citations of triggers or obligations are deduplicated
    Given a specification repeating the same undeclared trigger multiple times
    When the gate confronts it
    Then it reports that undeclared trigger only once in the error list

  @TRDCT-B12 @unit-level
  Scenario: Multiple vocabulary errors are sorted deterministically
    Given a specification citing multiple unknown triggers and obligations
    When the gate confronts it
    Then it returns Fail with defect messages sorted alphabetically and aggregated

  @TRDCT-B13 @unit-level
  Scenario: The list of declared triggers is capped at four
    Given a project declaring more triggers than the cap
    And none of them is close enough to the written one to be suggested
    When the gate confronts the artifact
    Then the verdict lists at most four of them, because a list of every trigger a large
      project declares would bury the advice it exists to give

  @TRDCT-I01 @unit-level
  Scenario: Trigger confrontation applies exclusively to specification artifacts
    Given code, feature, or test artifacts confronted by the gate
    When the gate confronts them
    Then it returns Skip, keeping verification scoped to where instructional triggers are taught

  @TRDCT-I02 @unit-level
  Scenario: Missing compliance packs produce Pending rather than Pass
    Given a specification citing compliance triggers in a project lacking pack configuration
    When the gate confronts it
    Then it returns Pending, avoiding false approvals of unverified compliance claims

  @TRDCT-I03 @unit-level
  Scenario: Trigger keys are limited to recognized compliance predicates
    Given a specification containing backticked key-value pairs outside the compliance predicate set
    When the gate confronts it
    Then it ignores them, restricting confrontation to recognized compliance triggers

  @TRDCT-I04 @unit-level
  Scenario: Undeclared triggers produce a blocking Fail verdict
    Given a specification citing an invalid trigger value
    When the gate confronts it
    Then it returns Fail, preventing erroneous vocabulary instructions from propagating

  @TRDCT-X01 @unit-level
  Scenario: The gate does not require every specification to cite triggers
    Given a standard specification that does not touch regulated data
    When the gate confronts it
    Then it returns Skip, treating compliance trigger citations as optional

  @TRDCT-X02 @unit-level
  Scenario: Code and test artifacts are not checked for trigger citations
    Given source code or test files containing trigger-like strings
    When the gate runs
    Then it leaves non-specification files outside its evaluation scope

  @TRDCT-X03 @unit-level
  Scenario: Code obligation implementation is not evaluated by this gate
    Given a specification citing valid triggers
    When the gate confronts it
    Then it evaluates vocabulary validity without inspecting code compliance fulfillment

  @TRDCT-X04 @unit-level
  Scenario: Unquoted prose text is not evaluated as symbol citations
    Given prose sentences containing trigger keywords without backticks
    When the gate confronts it
    Then it treats the text as unquoted narrative rather than syntax citations
