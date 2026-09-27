# language: en
# @anchors
#   ref: DLCTI
#   updated_at: 2026-09-26
#   layer: feature

@DLCTI
Feature: Dialect — the lexicon of the project's language, between an agnostic gate and concrete code

  @DLCTI-B01 @unit-level
  Scenario: A project that declares no dialect gets only the naming defaults
    Given a project with no configuration, and one with no dialect block
    When the effective dialect is resolved for each
    Then both have the default set-promise and set-slice conventions and no exported-function, loop, cursor or handling pattern

  @DLCTI-B02 @unit-level
  Scenario: The family fills every field the project left empty, and a declared field wins
    Given a dialect with the family "go", a declared loop pattern and a declared handling pattern
    When the effective dialect is resolved
    Then the loop and handling patterns are the declared ones, the other fields are the Go family's, and the declaration itself is unchanged

  @DLCTI-B03 @unit-level
  Scenario: The family name is matched ignoring case
    Given a dialect with the family "TS"
    When the effective dialect is resolved
    Then its exported-function pattern is the ts family's

  @DLCTI-B04 @unit-level
  Scenario: An unknown family contributes nothing, and the known ones are listed in order
    Given a dialect with the family "cobol"
    When the effective dialect is resolved and the known families are listed
    Then the lexicon stays empty, and the list holds every built-in family in alphabetical order

  @DLCTI-B05 @unit-level
  Scenario: The naming conventions apply to any family and a declared one replaces them
    Given a python dialect declaring only a set-promise pattern, and a dialect declaring only a set-slice pattern
    When the effective dialects are resolved
    Then each keeps its declared convention and gets the default for the other

  @DLCTI-B06 @unit-level
  Scenario: The Gherkin language defaults to English, and one outside the table keeps its code with English keywords
    Given the Gherkin languages "", "pt", "PT" and "eo"
    When the keywords are resolved for each
    Then they give en with "Scenario", pt with "Cenário" twice, and eo with "Scenario"

  @DLCTI-B07 @unit-level
  Scenario: Every way to open a scenario, in every language, deduplicated and longest first
    Given the Gherkin keyword table
    When the recognised scenario openings are listed
    Then each appears once, longer before shorter, including "Example", "Cenario" and "シナリオ", with "Esquema do Cenário" before "Cenário"

  @DLCTI-B08 @unit-level
  Scenario: Every result keyword, in every language, sorted
    Given the Gherkin keyword table
    When the recognised result keywords are listed
    Then they are exactly the table's result keywords, in alphabetical order

  @DLCTI-B09 @unit-level
  Scenario: An empty or invalid pattern compiles to nothing
    Given an empty pattern, a pattern with a lookbehind the engine rejects, and a valid loop pattern
    When each is compiled
    Then the first two give nothing and the third matches "for x"

  @DLCTI-B10 @unit-level
  Scenario: The opt-out is read by the YAML field name, ignoring case and spaces
    Given an opt-out list naming " Collection_Query " and another naming "HTTPStatus"
    When the fields collection_query, exported_func and http_status are asked about
    Then only collection_query is waived

  @DLCTI-B11 @unit-level
  Scenario: The set-promise verb is recognised after a provider prefix and never inside a word
    Given the names listUsers, GetAllItems, cognitoListDevices, dynamo_query_all, allocateSlot, callbackUrl and enlistment
    When the default set-promise convention is applied
    Then the first four are recognised and the last three are not

  @DLCTI-B12 @unit-level
  Scenario: The set-slice convention is a query verb opening the name, followed by a slice word
    Given the names listRecent, get_first, fetchTopItems, listUsers and cachedListRecent
    When the default set-slice convention is applied
    Then the first three are recognised and the last two are not

  @DLCTI-I01 @unit-level
  Scenario: Every pattern a family or a naming default brings compiles
    Given every built-in family
    When each family's effective dialect is resolved and its patterns compiled
    Then no non-empty pattern compiles to nothing

  @DLCTI-I02 @unit-level
  Scenario: The keywords written for any language are among those every reader recognises
    Given every language of the Gherkin table
    When its scenario, outline and result keywords are resolved
    Then each is among the recognised scenario openings or result keywords

  @DLCTI-X01 @unit-level
  Scenario: No family brings a collection query, which is the project's to declare
    Given every built-in family
    When each family's effective dialect is resolved
    Then its collection query pattern is empty
