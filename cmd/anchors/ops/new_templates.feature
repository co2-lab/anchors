# language: en
# @anchors
#   code: NTFNW
#   ref: NWTMN
#   updated_at: 2026-10-03
#   layer: feature

@NWTMN
Feature: NewTemplates — the catalog of artifact skeletons: which kinds exist, their headers, their sections and the presets that pick them

  @NWTMN-B01 @unit-level
  Scenario: The catalog holds seven kinds and each is born by new
    Given the template catalog
    When new runs for every kind with --code ABCDE
    Then the kinds are action, feature, flow, plan, product, spec, test
    And each file is written at its path holding ABCDE

  @NWTMN-B02 @unit-level
  Scenario: A markdown header is an HTML comment with the identity and placeholders
    Given a spec with code ABCDE
    When its header is written
    Then it is an HTML comment with code ABCDE, placeholder values for updated_at and layer

  @NWTMN-B03 @unit-level
  Scenario: A feature header opens with the Gherkin language line
    Given a feature with ref ABCDE
    When its header is written with no language, then with pt
    Then it opens with "# language: en" then the ref and layer feature
    And with pt it opens with "# language: pt"

  @NWTMN-B04 @unit-level
  Scenario: A test header uses the line comment of the output file
    Given a test with ref ABCDE at calc_test.py, calc.test.ts and calc_test.go
    When its header is written
    Then the Python one uses "#", the TypeScript and Go ones use "//", each with layer test

  @NWTMN-B05 @unit-level
  Scenario: Every section realizes only canonical rule letters
    Given the template catalog
    When the letters its sections realize are collected
    Then each is among the engine's canonical rule letters

  @NWTMN-B06 @unit-level
  Scenario: The test body is idiomatic for each known family
    Given the unit calcTotal with code XXXXX
    When the test body is written for python, go, java, kotlin, csharp, rust, ruby, php and ts
    Then each holds its language's construct, such as "def test_calc_total(" or "func TestCalcTotal(t *testing.T)"
    And none holds another language's syntax

  @NWTMN-B07 @unit-level
  Scenario: Every test body carries the scenario code
    Given the code ABCDX
    When the test body is written for every family, including an unknown one
    Then each body carries ABCDX joined to B01 by a dash, in the test's name

  @NWTMN-B08 @unit-level
  Scenario: An unknown family gets an instruction, not guessed syntax
    Given the family cobol
    When the test body is written
    Then it holds an instruction to the author and neither describe( nor def

  @NWTMN-B09 @unit-level
  Scenario: Each spec preset is an ordered set of catalog sections
    Given the spec presets
    When they are listed
    Then they are the twelve named presets
    And each names only catalog sections, opens with title and closes with open

  @NWTMN-B10 @unit-level
  Scenario: Product doctrine has its own sections and no layer
    Given the product kind
    When a doctrine CreditLimit with code CRLMT is rendered
    Then every section key starts with doctrine_
    And it emits CRLMT-R01 and no spec rule
    And its header has no layer

  @NWTMN-B11 @unit-level
  Scenario: A unit name becomes snake_case with one separator and whole acronyms
    Given the unit names "My-Name", "HTTPServer", "getHTTPCode" and "calc -- total"
    When they are turned into snake_case
    Then they are "my_name", "http_server", "get_http_code" and "calc_total"
    And the Python test body for "HTTPServer" is "def test_http_server("

  @NWTMN-B12 @unit-level
  Scenario: The catalog ties rules to what they use
    Given the spec catalog and its presets
    When the three sections and the presets are read, and a screen spec is rendered in English
    Then validations realizes V, presentation-validations P and rule-uses no letter, rule-uses is in every preset, and the screen reads Validations, Presentation validations and Rule uses with the rule code first
