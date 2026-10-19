# language: en
# @anchors
#   code: DLFTA
#   ref: DLCTI
#   updated_at: 2026-10-19
#   layer: feature

@DLCTI
Feature: Dialect — the lexicon of the project's language, between an agnostic gate and concrete code

  @DLCTI-B01 @unit-level
  Scenario: A project's dialect is what it declares, and nothing fills it
    Given a project with no configuration, one with no dialect block, and one naming only the family "go"
    When the effective dialect is resolved for each
    Then none has a naming convention, an exported-function, loop, cursor, handling, tests, definition or env-read pattern
    And a declared loop pattern is read as declared, with nothing beside it

  @DLCTI-B02 @unit-level
  Scenario: Each family has an example of the lexicon, never applied by itself
    Given the family "go"
    When its example is asked for
    Then it carries the Go lexicon with its statement end, signature end, limits and guards
    And changing the copy handed out changes no family

  @DLCTI-B03 @unit-level
  Scenario: The family name is matched ignoring case
    Given the family "TS"
    When its example is asked for
    Then its exported-function pattern is the ts family's

  @DLCTI-B04 @unit-level
  Scenario: An unknown family has no example, and the known ones are listed in order
    Given the family "cobol"
    When its example is asked for and the known families are listed
    Then there is no example, and the list holds every built-in family in alphabetical order

  @DLCTI-B05 @unit-level
  Scenario: Every family's example carries the two naming conventions
    Given every known family, and a dialect declaring only a set-slice pattern
    When the examples are asked for and the dialect is resolved
    Then each example has the set-promise and set-slice conventions, and the dialect only its own

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
    When each family's example is asked for and its patterns compiled
    Then no non-empty pattern compiles to nothing

  @DLCTI-I02 @unit-level
  Scenario: The keywords written for any language are among those every reader recognises
    Given every language of the Gherkin table
    When its scenario, outline and result keywords are resolved
    Then each is among the recognised scenario openings or result keywords

  @DLCTI-X01 @unit-level
  Scenario: No family brings a collection query, which is the project's to declare
    Given every built-in family
    When each family's example is asked for
    Then its collection query pattern is empty

  @DLCTI-B13 @unit-level
  Scenario: The Go family recognises both shapes of error handling
    Given the Go family's handle patterns
    When they read "if err != nil {" and "if err := os.Remove(p); err != nil {"
    Then both lines are recognised as handling a failure

  @DLCTI-B14 @unit-level
  Scenario: The Go and TS families say how a test is written
    Given the go and ts families' examples, and a ts project that declares its own tests script
    When they read a test
    Then go reads t.Run, ts reads it, test and describe with their modifiers, each with its title
    And the project's script is read as declared

  @DLCTI-B15 @unit-level
  Scenario: The TS family recognises catch with or without its binding
    Given a try block closed by catch (e) and one closed by catch with no binding
    When the ts family's handle patterns read them
    Then both are recognised as handling, and a word merely containing catch is not

  @DLCTI-B16 @unit-level
  Scenario: The families say what an assertion is, and a project may declare only its own
    Given the Go family's example, and a Go project declaring only its assertion
    When each is read
    Then the example asserts with t.Errorf, and the project has its assertion and no way to open a test

  @DLCTI-B17 @unit-level
  Scenario: The families say how code defines a name
    Given Go, TS and Python code defining functions and types, and a project declaring its own definition
    When each family's definition reads the code
    Then it captures the names defined, and the declared one is read as declared

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
    When each family's example, and the pattern of a project that declares none, reads them
    Then each family reads its own, and the project that declares none reads nothing


  @DLCTI-B21 @unit-level
  Scenario: The fallible patterns are the project's, and the family's only when the project declares none
    Given a project declaring its own fallible pattern, and a project declaring none with no family
    When the dialect is resolved
    Then the first keeps its pattern and window, and the second has none

  @DLCTI-B22 @unit-level
  Scenario: A family knows its fallible patterns and does not impose them on a project that declared none
    Given the ts, go and python families, an unknown family, and a ts project that declared no fallible patterns
    When the family's patterns and the project's dialect are read
    Then each known family has its patterns, the unknown has none, and the project is measured with none

  @DLCTI-B23 @unit-level
  Scenario: import_resolve and navigation_call are read from the dialect
    Given a dialect declaring import extensions, an alias, a separator and a navigation call with its route group, and a ts project that declared no navigation call
    When the config is loaded and the dialect read
    Then the dialect has the extensions, the alias, the separator and the navigation pattern, and the ts project has no navigation call while its family has an example

  @DLCTI-B24 @unit-level
  Scenario: The OCaml family reads OCaml
    Given a project declaring the ocaml family and nothing else
    When its dialect is read against an OCaml module
    Then the top-level let and let rec are exported, a local let is not, Sys.getenv_opt names the variable, the family's example reads the modules opened and qualified, and the project has no import pattern

  @DLCTI-B25 @unit-level
  Scenario: The import resolution declares imports by name
    Given an anchors.yaml whose import_resolve says by_name with name_case lower_first
    When it is loaded
    Then the dialect resolves imports by name, lowering the first letter

  @DLCTI-B26 @unit-level
  Scenario: A family's fallible patterns say when a call propagates
    Given the ts family's patterns and go's HTTP pattern
    When return useQuery, an arrow whose body is an awaited fetch, return http.Get and a query kept in a variable are matched
    Then the first three match their propagates group, and the query kept in a variable does not

  @DLCTI-B27 @unit-level
  Scenario: The test handle patterns are read and checked on load
    Given a dialect declaring exposed, mentioned and queried patterns, one with a pattern that does not compile, and one whose pattern captures with no group
    When the config is loaded
    Then the first is read, and the others refuse the config naming the pattern

  @DLCTI-B28 @unit-level
  Scenario: The ts family has an example of the test handle patterns, never applied by itself
    Given the ts family with the attribute data-testid, the go family, the ts family with no attribute, a ts project that declared none, and a nested conditional handle whose branch is a template
    When the examples and the project's dialect are read
    Then the ts example reads data-testid and every branch of the conditional, none of its condition, the others give none, and the project has no patterns

  @DLCTI-B29 @unit-level
  Scenario: A dialect pattern runs over the whole content, as an expression or a region
    Given a pattern spanning lines, a region of three names, a pattern narrowed to Go files, and patterns written as text and as a map
    When they are matched and loaded, and five broken patterns are loaded
    Then each capture keeps its line, the region's names share its start and read as listed, the Go pattern skips a TS file, and each broken pattern refuses the config naming it

  @DLCTI-B30 @unit-level
  Scenario: The chain patterns are the dialect's, and a family gives none by itself
    Given a dialect declaring imports, import names and declarations, one whose import pattern has no path group, and a ts project that declared none
    When the configs are loaded and the dialect read
    Then the first is read, the second refuses naming the field, and the ts project has no import patterns nor import resolution

  @DLCTI-B31 @unit-level
  Scenario: The families have an example of the chain patterns
    Given the ts, go, python and ocaml examples and an unknown family
    When each reads a source of its language
    Then ts reads a list across lines, a side effect, a destructured dynamic import and a flow, go reads an import and a block, python reads from-imports, ocaml reads opened modules, and the unknown family has none

  @DLCTI-B32 @unit-level
  Scenario: The families' examples of the mock shapes and the route pattern
    Given the ts, python and go families
    When their examples of mocks and of a route pattern are asked for
    Then ts reads jest.mock, a factory passed and a typed one, and name="X", a param list and addResource
    And python reads patch, and go has none

  @DLCTI-B33 @unit-level
  Scenario: Aliases scoped to the importers a glob matches
    Given `@/` declared for apps/mobile/, and scoped to apps/landing-page/src/ for apps/landing-page/**
    When the aliases of an importer in each app are asked for
    Then the mobile importer reads apps/mobile/ and the landing importer apps/landing-page/src/, the other aliases kept
    And a scoped entry with no glob or no aliases refuses the config

  @DLCTI-B34 @unit-level
  Scenario: The exposed patterns read a handle whose head is the caller's
    Given the ts example for testID, and a project declaring an exposed pattern with a suffix group
    When `${testID}-option-${v}` and `${prefix}${key}` are read, and the project's config loads
    Then the first gives the segment "-option-", the second an empty one, and the config loads
