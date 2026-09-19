# language: en
# @anchors
#   ref: VLANV
#   updated_at: 2026-09-19
#   layer: feature

@VLANV
Feature: ValueAnchored — every value of a closed set points at the rule that justifies it

  @VLANV-B01 @unit-level
  Scenario: An artifact that is not a spec leaves without a verdict
    Given a node whose kind is code, test or feature
    When the gate confronts it
    Then it returns Skip, because the gate reaches the code through the spec

  @VLANV-B02 @unit-level
  Scenario: A value of a closed set with no anchor is failed
    Given a closed set declaring the values "15m" and "1h" with no anchor comment on any of them
    When the gate confronts it
    Then it returns Fail, because each value is a domain decision that left no address

  @VLANV-B03 @unit-level
  Scenario: The verdict names the unanchored value
    Given a closed set whose value "15m" carries no anchor
    When the gate confronts it
    Then the verdict names "15m", so the reader does not have to hunt for which value it was

  @VLANV-B04 @unit-level
  Scenario: A value whose anchor carries rule key and value passes
    Given a closed set where each value is preceded by an anchor holding the rule key and that same value
    When the gate confronts it
    Then it returns Pass, because the address exists and it checks out

  @VLANV-B05 @unit-level
  Scenario: An anchor that asserts one value while the line says another is failed
    Given an anchor asserting "15m" written above a line whose literal is "5m"
    When the gate confronts it
    Then it returns Fail, because a lying anchor looks like traceability while pointing at the wrong place

  @VLANV-B06 @unit-level
  Scenario: The verdict of a lying anchor shows both sides of the divergence
    Given an anchor asserting "15m" written above a line whose literal is "5m"
    When the gate confronts it
    Then the verdict carries both "15m" and "5m", because one side alone does not show the drift

  @VLANV-B07 @unit-level
  Scenario: Lying anchors are reported before the unanchored ones
    Given a closed set carrying one lying anchor and one value with no anchor at all
    When the gate confronts it
    Then the lying anchor appears first in the verdict, because the absent anchor can be seen
      and the lying one cannot

  @VLANV-B08 @unit-level
  Scenario: Without a declared value anchor pattern the gate skips
    Given a project that declares no value anchor pattern
    When the gate confronts a closed set
    Then it neither approves nor fails, because it cannot read and will not stamp what it did not measure

  @VLANV-B09 @unit-level
  Scenario: The skip names the setting that enables the gate
    Given a project that declares no value anchor pattern
    When the gate confronts a closed set
    Then the verdict names the value_anchor setting, so the reader learns how to turn the gate on

  @VLANV-B10 @unit-level
  Scenario: A declaration that opens no list is not a closed set
    Given a public symbol declared as a single scalar value on one line
    When the gate confronts it
    Then it returns Pass, because a scalar is not a closed set and there is nothing to anchor

  @VLANV-I01 @unit-level
  Scenario: An anchor pattern with a single capture group does not enable the gate
    Given a declared anchor pattern that captures only the rule key
    When the gate confronts a closed set
    Then it neither approves nor fails, because without the second group the anchor asserts
      no value and the confrontation cannot happen

  @VLANV-I02 @unit-level
  Scenario: A line carrying an anchor is never read as the end of the list
    Given a closed set whose second value carries an anchor that lies, after a first anchored value
    When the gate confronts it
    Then it returns Fail, because the brackets inside the anchor must not close the set early

  @VLANV-I03 @unit-level
  Scenario: With no built map the verdict is pending
    Given no graph built
    When the gate confronts a spec
    Then it does not approve, because approving without being able to look would stamp
      what was never measured

  @VLANV-X01 @unit-level
  Scenario: The gate does not judge whether the value is a good one
    Given a closed set whose value "999y" is anchored to a rule key and matches it exactly
    When the gate confronts it
    Then it returns Pass, because the ruler is that the decision has an address and the
      address does not lie — whether the value belongs is judgement

  @VLANV-X02 @unit-level
  Scenario: The gate does not accuse a line whose literal it cannot read
    Given a closed set line carrying a computed expression instead of a quoted literal
    When the gate confronts it
    Then it returns Pass, because a false negative is better than a mass of false positives
      that would train the team to ignore the gate

  @VLANV-X03 @unit-level
  Scenario: The gate does not decide what an anchor or a public symbol looks like
    Given a project whose declared export pattern does not match the way this file writes its set
    When the gate confronts it
    Then it returns Pass, because both shapes are declared by the project — inventing them
      would charge a convention nobody adopted
