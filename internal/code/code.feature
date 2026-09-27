# language: en
# @anchors
#   ref: CDGNC
#   updated_at: 2026-09-26
#   layer: feature

@CDGNC
Feature: CodeGenerator — the short, stable identity code suggested for a unit's name

  @CDGNC-B01 @unit-level
  Scenario: Generic suffixes are dropped unless they are the whole name
    Given the names "LoginScreen", "AlertSheet", "ScrollableLayout", "Button" and "Screen"
    When the generic suffix is stripped
    Then they become "Login", "Alert", "Scrollable", "Button" and "Screen"

  @CDGNC-B02 @unit-level
  Scenario: Words are split at separators, camel case and acronyms
    Given the names "user_profile-edit" and "ABCParser" with a code length of 5
    When their codes are generated
    Then "user_profile-edit" gives "UPESR" and "ABCParser" gives "ABPRB"

  @CDGNC-B03 @unit-level
  Scenario: A long name takes the initials
    Given the five-word name "UserProfileEditFormDraft" and a code length of 5
    When its code is generated
    Then the code is "UPEFD"

  @CDGNC-B04 @unit-level
  Scenario: A two-word name takes letters from each word
    Given the name "TransactionDetail" and a code length of 5
    When its code is generated
    Then the code is "TRDTT", starting with T and holding the D of the second word

  @CDGNC-B05 @unit-level
  Scenario: A single word takes consonants before vowels
    Given the names "Spacer", "Login" and "AlertsScreen" with a code length of 5
    When their codes are generated
    Then they are "SPCRA", "LGNOI" and "LRTSA"

  @CDGNC-B06 @unit-level
  Scenario: Short codes are padded with X and existing codes are completed the same way
    Given the name "Ok", and the existing codes "mtvr" and "ABCDEFG", with a code length of 5
    When the name's code is generated and the existing codes are padded
    Then they are "KOXXX", "MTVRX" and "ABCDE"

  @CDGNC-B07 @unit-level
  Scenario: A module prefix starts the code
    Given the name "Login" with the prefix "au", and with the prefix "abcdefg", and a code length of 5
    When the codes are generated
    Then they are "AULGN" and "ABCDE"

  @CDGNC-B08 @unit-level
  Scenario: A collision varies the last position and keeps the prefix
    Given the code "AULGN" already taken
    When a unique code for "Login" with prefix "AU" is generated
    Then the code is "AULGA"
    And a unique code for a name whose code is free is the generated code itself

  @CDGNC-B09 @unit-level
  Scenario: The module prefix is the initial and the first consonant
    Given the module names "auth", "family" and "i"
    When their module prefixes are derived
    Then they are "AT", "FM" and "IX"

  @CDGNC-B10 @unit-level
  Scenario: The generated length follows the smallest declared length
    Given a project configuration declaring the code lengths 5 and 4
    When it is loaded
    Then the code of "Login" is generated with 4 characters, "LGNO"
    And declaring only lengths below 2, or none, leaves the length unchanged

  @CDGNC-I01 @unit-level
  Scenario: Generation is deterministic and a resolved code is free
    Given the name "Spacer" whose generated code is taken
    When the unique code is generated twice with the same taken codes
    Then both answers are equal, differ from the generated code, and are not taken

  @CDGNC-X01 @unit-level
  Scenario: A saturated namespace returns the generated code
    Given the code of "Spacer" and every single-position variation of it already taken
    When a unique code for "Spacer" is generated
    Then the answer is the generated code, taken as it is
