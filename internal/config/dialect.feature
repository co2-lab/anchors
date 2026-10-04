# language: en
# @anchors
#   code: DLFTA
#   ref: DLCTI
#   updated_at: 2026-10-03
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
  Scenario: The Gherkin language defaults to English, is found in any case, and one outside the table keeps its code with English keywords
    Given the Gherkin languages "", "pt", "PT", "eo", "zh-CN", "ZH-cn" and "en-AU"
    When the keywords are resolved for each
    Then they give en with "Scenario", pt with "Cenário" twice, eo with "Scenario", zh-CN with "场景" twice, and en-AU kept as declared with "Scenario"

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

  @DLCTI-B13 @unit-level
  Scenario: The Go family recognises both shapes of error handling
    Given the Go family's handle patterns
    When they read "if err != nil {" and "if err := os.Remove(p); err != nil {"
    Then both lines are recognised as handling a failure

  @DLCTI-B14 @unit-level
  Scenario: The Go and TS families say how a test is written
    Given the go and ts families, and a ts project that declares its own tests script
    When the effective dialect is read
    Then go reads t.Run, ts reads it, test and describe with their modifiers, and the project's script wins

  @DLCTI-B15 @unit-level
  Scenario: The TS family recognises catch with or without its binding
    Given a try block closed by catch (e) and one closed by catch with no binding
    When the ts family's handle patterns read them
    Then both are recognised as handling, and a word merely containing catch is not

  @DLCTI-B16 @unit-level
  Scenario: The families say what an assertion is, and a project may declare only its own
    Given the Go family, a Go project declaring only its assertion, and one declaring its own tests pattern
    When each dialect is resolved
    Then the family asserts with t.Errorf, the second keeps t.Run with its own assertion, and the third takes no assertion from the family

  @DLCTI-B17 @unit-level
  Scenario: The families say how code defines a name
    Given Go, TS and Python code defining functions and types, and a project declaring its own definition
    When each family's definition reads the code
    Then it captures the names defined, and the declared one wins

  @DLCTI-B18 @unit-level
  Scenario: Every examples keyword, in every language, sorted
    Given the Gherkin table
    When the examples keywords are listed
    Then each language's keyword appears once, in alphabetical order

  @DLCTI-B19 @unit-level
  Scenario: The Go family sees an error in a field and a sentinel error
    Given the Go family
    When its handle and log patterns are matched
    Then "if result.Error != nil {" and "return ErrWalletLinkNotFound" are handling
    And "return nil, result.Error" and "return nil, ErrWalletLinkNotFound" are recording
    And "return e.Error()" is not recording

  @DLCTI-B20 @unit-level
  Scenario: Each family reads environment variables its own way
    Given the reads os.Getenv("A"), process.env.A, os.environ["A"], System.getenv("A"), env::var("A"), ENV["A"], getenv('A') and Environment.GetEnvironmentVariable("A")
    When each family's pattern, and the pattern of a project with no family, reads them
    Then each family reads its own, and the project with no family reads them all

