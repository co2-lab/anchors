# language: en
# @anchors
#   code: KNFTK
#   ref: DCKND
#   updated_at: 2026-10-03
#   layer: feature

@DCKND
Feature: DocKinds — what each kind of project documentation must answer, told to the agent that writes it

  @DCKND-B01 @unit-level
  Scenario: Every known kind has a title, items to cover and a trap
    Given the six known kinds openapi, c4, schema, components, adr and runbook
    When the instruction of each one is looked up
    Then each has a non-empty title, at least one item and a non-empty trap

  @DCKND-B02 @unit-level
  Scenario: The kind is found ignoring case and surrounding spaces
    Given a document declared with the kind " OpenAPI "
    When its instruction is looked up
    Then it is the openapi instruction

  @DCKND-B03 @unit-level
  Scenario: An unknown kind gets a minimal instruction titled by the kind or the path
    Given a document of kind "glossary" and another with no kind at path "docs/x.md"
    When their instructions are looked up
    Then the first is titled "glossary" and the second "docs/x.md"
    And neither has items nor a trap

  @DCKND-B04 @unit-level
  Scenario: The duty text lists title and path, why, items and trap in order
    Given an openapi document at "docs/api.yaml" declared with a why
    When its duty is rendered
    Then the text starts with the title and the path in backticks
    And the why comes next, then one bulleted line per item, then the trap line

  @DCKND-I01 @unit-level
  Scenario: The known kinds and the instructions are the same set
    Given the list of known kinds
    When each is looked up and the sets are compared
    Then none falls back to the minimal instruction and both sets have the same size

  @DCKND-X01 @unit-level
  Scenario: No documentation kind is refused
    Given a document of a kind Anchors has never heard of
    When its duty is rendered
    Then a duty text naming the document is returned
