# language: en
# @anchors
#   ref: INCTN
#   updated_at: 2026-10-01
#   layer: feature

@INCTN
Feature: InitCatalogs — the language dialect catalog and the @TBD instruction that init seeds into judgment gates

  @INCTN-B04 @unit-level
  Scenario: The @TBD instruction forbids pass, orders a waiver naming the absence, and names the piece asked about
    Given a judgment gate asking about the code
    When its @TBD instruction is produced
    Then the English text mentions @TBD, the `waived` answer and the naming of the absence, and forbids `pass`
    And the instruction produced for the test names the test

  @INCTN-B06 @unit-level
  Scenario: A file is a test when its name carries a convention's prefix and suffix
    Given the names foo_test.go, test_foo.py, Login.spec.tsx, UserTest.java, _test.go and foo.go
    When each is matched against the catalog
    Then foo_test.go, test_foo.py, Login.spec.tsx and UserTest.java are tests
    And _test.go, with nothing besides the suffix, and foo.go are not

  @INCTN-B07 @unit-level
  Scenario: A convention gives its glob and its template
    Given the conventions `_test.go` and `test_*.py`
    When their glob and template are asked
    Then they are **/*_test.go with {{dir}}/{{name}}_test.go, and **/test_*.py with {{dir}}/test_{{name}}.py

  @INCTN-B08 @unit-level
  Scenario: A test file's unit name drops the convention's prefix and suffix
    Given the test files test_foo.py and foo_test.go
    When their unit names are asked
    Then both are foo

  @INCTN-B09 @unit-level
  Scenario: With no convention read, the family default is used
    Given a Python project with no test file
    When its conventions are asked
    Then the convention is `test_*.py`
    And a project with no test file and no family has no convention and no template

  @INCTN-B10 @unit-level
  Scenario: A family's coverage hint names the reports ingest reads
    Given the families go and python
    When their coverage hints are asked
    Then the go hint names lcov and junit, and the python hint names pytest
    And a family with no hint, and no family, get none

  @INCTN-I02 @unit-level
  Scenario: The @TBD instruction demands checking that the @TBD is still true
    Given a judgment gate asking about the code
    When its @TBD instruction is produced
    Then the text tells the judge that a marker whose piece already exists is out of date

  @INCTN-I03 @unit-level
  Scenario: No convention of the catalog is shadowed by a shorter one
    Given the convention catalog
    When every pair of conventions is compared
    Then a longer form always comes before a shorter form it contains

  @INCTN-I04 @unit-level
  Scenario: Every family default is a convention of the catalog
    Given the family defaults
    When each is looked up in the catalog
    Then each is found, with its own family

  @INCTN-X02 @unit-level
  Scenario: The catalog names no folder
    Given the conventions, the family defaults and the manifests
    When each is inspected
    Then none carries a directory separator
