# language: en
# @anchors
#   ref: VLANV
#   updated_at: 2026-09-23
#   layer: feature

@VLANV
Feature: ValueAnchored — a replicated key is declared where it is used, and every copy carries the same value

  @VLANV-B01 @unit-level
  Scenario: A declaration matching its code line passes
    Given a code file declaring `COLOR-OK` as `#1F8A5B`
    And the code line below it holds `#1F8A5B`
    When the gate confronts the file
    Then it passes

  @VLANV-B02 @unit-level
  Scenario: A declaration whose code line says another value fails
    Given a code file declaring `COLOR-OK` as `#1F8A5B`
    And the code line below it holds `#2A9D6B`
    When the gate confronts the file
    Then it fails, showing the key, the declared value and the line

  @VLANV-B03 @unit-level
  Scenario: Comment and blank lines between declaration and code are skipped
    Given two stacked declarations and an explanatory comment above one code line
    When the gate confronts the file
    Then both declarations are checked against that code line

  @VLANV-B04 @unit-level
  Scenario: Copies of the same key with different values are reported
    Given `COLOR-OK` declared as `#2A9D6B` in one file and as `#1F8A5B` in another
    And each file is locally consistent
    When the gate confronts either file
    Then it fails, listing every place and value of the key

  @VLANV-B05 @unit-level
  Scenario: The value a rule declares in the spec is the source
    Given a spec row `TKNSX-R01` declaring `#2A9D6B`
    And code declarations of `TKNSX-R01` as `#1F8A5B` that agree with each other
    When the gate confronts a file holding one of them
    Then it fails, naming the spec and the value it declares

  @VLANV-B06 @unit-level
  Scenario: Without a declared pattern the gate skips and says how to enable it
    Given a project that does not declare `derived.value_anchor`
    When the gate confronts a code file
    Then it skips, naming the setting

  @VLANV-B07 @unit-level
  Scenario: A declaration with nothing below annotates nothing
    Given a declaration on the last line of a file
    When the gate confronts the file
    Then it fails

  @VLANV-B08 @unit-level
  Scenario: An anchor inside prose is a mention, not a declaration
    Given a comment that explains the anchor syntax in a sentence
    When the gate confronts the file
    Then the mention is not charged

  @VLANV-X01 @unit-level
  Scenario: A prose rule declares no value
    Given rules whose lines carry backticked identifiers only in headings or prose cells
    When declarations point at those rules
    Then they are not charged against the spec

  @VLANV-X02 @unit-level
  Scenario: A literal nobody declared is not charged
    Given a closed set with no declaration
    When the gate confronts the file
    Then it skips
