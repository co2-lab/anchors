# language: en
# @anchors
#   code: CDFTB
#   ref: RCGRL
#   updated_at: 2026-10-03
#   layer: feature

@RCGRL
Feature: RuleCodeGrammar — the grammar that recognizes a scenario code in a test's name, in the project's vocabulary

  @RCGRL-B01 @unit-level
  Scenario: A rule code of each canonical letter is recognized in a test name
    Given test names carrying the codes ABCDX-B01, ABCDX-E01, ABCDX-I01 and ABCDX-Q01 followed by prose
    When the codes are read from each name
    Then each name yields exactly its code

  @RCGRL-B02 @unit-level
  Scenario: A rule code with a lowercase slug is recognized with the slug
    Given the test name "ABCDX-B01-some-slug proves it"
    When the codes are read from it
    Then the only code is "ABCDX-B01-some-slug"

  @RCGRL-B03 @unit-level
  Scenario: Design-system and visual-regression codes are recognized
    Given the test names "ABCDX-DS-Button-primary renders" and "ABCDX-VR matches"
    When the codes are read from them
    Then the codes are "ABCDX-DS-Button-primary" and "ABCDX-VR"

  @RCGRL-B04 @unit-level
  Scenario: Declared rule letters replace the vocabulary
    Given the project declares the rule letters "BZ"
    When the codes are read from "ABCDX-Z01" and from "ABCDX-I01"
    Then "ABCDX-Z01" is recognized and "ABCDX-I01" is not

  @RCGRL-B05 @unit-level
  Scenario: A declared code length replaces the accepted identity length
    Given the default code length, under which "ABC-B01" is not a code
    When the project declares a code length of 3
    Then "ABC-B01" is recognized

  @RCGRL-B06 @unit-level
  Scenario: An empty declaration keeps the vocabulary in place
    Given the project declared the letters "BZ" and a code length of 3
    When an empty value is declared for the letters and for the length
    Then the letters stay "BZ" and the length stays 3

  @RCGRL-I01 @unit-level
  Scenario: The default rule letters equal the configuration's default letters
    Given the package's default rule letters
    When they are compared with the configuration's default rule letters
    Then both are the same string

  @RCGRL-X01 @unit-level
  Scenario: An identity too long or glued to a longer word is not a code
    Given the test names "XABCDXY-B01 too long" and "abcABCDX-B01 glued"
    When the codes are read from them
    Then no code is recognized

  @RCGRL-B07 @unit-level
  Scenario: A scenario code keeps its variant
    Given a case name and a feature citing a scenario with a variant and one without
    When their codes are read
    Then the variant is kept, the rule is the code without it, and each tag comes once

  @RCGRL-B08 @unit-level
  Scenario: A rule is proven only when each of its scenarios is
    Given a rule with two variants of which one is proven, a rule with both proven, and a rule no feature declares
    When the proven rules are read
    Then the first is not proven and names its missing variant, the second is, and the third is proven by its own code
