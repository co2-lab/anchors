# language: en
# @anchors
#   ref: CDLNG
#   updated_at: 2026-09-19
#   layer: feature

@CDLNG
Feature: CodeLanguage — the code does not go back to mixing languages

  @CDLNG-B01 @unit-level
  Scenario: An identifier in the wrong language is accused, and an English one passes
    Given a file declaring a function named in Portuguese
    And another declaring a function named in English
    When the gate confronts the file
    Then the Portuguese identifier is accused
    And the English one passes without noise

  @CDLNG-B02 @unit-level
  Scenario: The verdict returns the word that accused
    Given a file declaring an identifier in the wrong language
    When the gate confronts it
    Then the verdict names that exact word, so the reader does not hunt the whole file

  @CDLNG-B03 @unit-level
  Scenario: Only a DECLARATION is the subject
    Given a file whose Portuguese words appear outside any declaration
    When the gate confronts it
    Then it accuses nothing, because what declares no identifier is not read

  @CDLNG-B04 @unit-level
  Scenario: The declarations are found in every form the language offers
    Given a file declaring identifiers as function, type, variable and constant
    When the gate confronts it
    Then every one of those declarations is read, not only the most common form

  @CDLNG-B05 @unit-level
  Scenario: Deciding one word is separate from deciding a whole identifier
    Given a compound identifier made of several words
    When the gate decides its language
    Then it breaks the identifier into its words and decides each one
    And that separation is what lets the length floor apply per word instead of to
      the identifier as a whole

  @CDLNG-I01 @unit-level
  Scenario: A short word does not count
    Given a file declaring identifiers below the length floor
    When the gate confronts it
    Then none is accused, because below the floor there is no language to infer and
      accusing there is the noise that costs a gate its credibility

  @CDLNG-X01 @unit-level
  Scenario: The gate does not read comments
    Given a file whose comments are written in the team's own language
    And every identifier is in English
    When the gate confronts it
    Then it returns Pass, because the comment carries the measurement and the why

  @CDLNG-X02 @unit-level
  Scenario: The gate does not read user-facing text
    Given a file whose user-facing strings are written in the project's language
    And every identifier is in English
    When the gate confronts it
    Then it returns Pass, because that text goes through the translation catalog

  @CDLNG-X03 @unit-level
  Scenario: The gate does not use a dictionary to decide the language
    Given a file declaring compound identifiers in English that no common dictionary holds
    When the gate confronts it
    Then none is accused, because the dictionary approach was measured at ninety percent
      false positives — and a gate that wrong is switched off, defending nothing
