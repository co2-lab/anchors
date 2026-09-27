# language: en
# @anchors
#   ref: TSTRT
#   updated_at: 2026-09-27
#   layer: feature

@TSTRT
Feature: TestTraceable — a test linked to a feature must declare what scenario it proves

  @TSTRT-B01 @unit-level
  Scenario: Non-test artifacts skip confrontation
    Given an artifact node whose kind is not test
    When the gate confronts it
    Then it returns Skip, because traceability is charged only on test files

  @TSTRT-B02 @unit-level
  Scenario: Confrontation without a dependency graph returns Pending
    Given a test artifact and no graph built
    When the gate confronts it
    Then it returns Pending, because without the map relational edges cannot be inspected

  @TSTRT-B03 @unit-level
  Scenario: A test with no linked feature skips confrontation
    Given a test node with no incoming tested-by edge from a feature
    When the gate confronts it
    Then it returns Skip, because an unlinked test has no declared scenarios to cite

  @TSTRT-B04 @unit-level
  Scenario: A test whose linked feature cannot be read skips confrontation
    Given a test linked to a feature path that does not exist on disk
    When the gate confronts it
    Then it returns Skip, because the feature scenarios cannot be retrieved

  @TSTRT-B05 @unit-level
  Scenario: A test whose linked feature declares no scenario codes skips confrontation
    Given a test linked to a feature file containing no scenario codes
    When the gate confronts it
    Then it returns Skip, because there are no scenario codes to demand

  @TSTRT-B06 @unit-level
  Scenario: A test containing a declared scenario code passes
    Given a test whose content cites a scenario code declared in its linked feature
    When the gate confronts it
    Then it returns Pass, confirming the test declares what it proves

  @TSTRT-B07 @unit-level
  Scenario: A single declared scenario code is sufficient to pass
    Given a linked feature declaring multiple scenario codes
    And a test file citing only one of those scenario codes
    When the gate confronts it
    Then it returns Pass, satisfying the minimum declaration threshold

  @TSTRT-B08 @unit-level
  Scenario: A test containing none of the linked feature scenario codes fails
    Given a test file whose content contains none of the scenario codes of its linked feature
    When the gate confronts it
    Then it returns Fail, because the test is invisible to relational gates

  @TSTRT-B09 @unit-level
  Scenario: A test with transposed or misspelled codes fails exact substring confrontation
    Given a test file using an inverted or mistyped scenario code
    When the gate confronts it
    Then it returns Fail, because exact substring matching is required for traceability

  @TSTRT-B10 @unit-level
  Scenario: A failing verdict names the linked feature and lists expected codes
    Given a test file that omits the declared scenario codes
    When the gate confronts it
    Then it returns Fail with a defect message citing the feature path and expected codes

  @TSTRT-B11 @unit-level
  Scenario: Scenario codes appearing anywhere in test content satisfy traceability
    Given a test file citing a declared scenario code in a comment or suite title
    When the gate confronts it
    Then it returns Pass, recognizing the code anywhere within the file text

  @TSTRT-I01 @unit-level
  Scenario: Traceability is strictly charged on test artifacts
    Given a specification or code node confronted by the gate
    When the gate confronts it
    Then it returns Skip, preventing misattribution of test debt to other files

  @TSTRT-I02 @unit-level
  Scenario: Tests without a linked feature are never charged
    Given a standalone test node without any feature parent in the graph
    When the gate confronts it
    Then it returns Skip, avoiding demands for references to nonexistent specifications

  @TSTRT-I03 @unit-level
  Scenario: The acceptance threshold requires only one scenario code present
    Given a test covering one behavior among several declared scenarios
    When the gate confronts it
    Then it returns Pass, delegating exhaustive coverage checks to scenario matching gates

  @TSTRT-I04 @unit-level
  Scenario: Missing graph structure produces Pending rather than approving
    Given a test node evaluated with a nil graph pointer
    When the gate confronts it
    Then it returns Pending, preventing false approvals when relations cannot be verified

  @TSTRT-X01 @unit-level
  Scenario: The gate does not enforce one-to-one scenario coverage
    Given a feature declaring multiple scenarios and a test citing only one code
    When the gate confronts it
    Then it returns Pass, leaving per-scenario pairing to feature-test-match

  @TSTRT-X02 @unit-level
  Scenario: Test execution results are not inspected by this gate
    Given a test file containing a declared scenario code regardless of runtime pass or fail status
    When the gate confronts it
    Then it returns Pass, leaving execution validation to tests-green

  @TSTRT-X03 @unit-level
  Scenario: Scenario codes are not required in standalone test files
    Given a test file testing internal helpers without an associated feature
    When the gate confronts it
    Then it returns Skip, allowing utility tests to exist without scenario codes

  @TSTRT-X04 @unit-level
  Scenario: Test assertion semantics and quality are not evaluated
    Given a test file containing a valid scenario code without valid assertion semantics
    When the gate confronts it
    Then it returns Pass, validating relational visibility rather than test logic

  @TSTRT-B12 @unit-level
  Scenario: With a tests source a test traces only through its titles
    Given a test file whose feature code appears in a fixture string but in no test title
    When the project declares how tests are written, and when it does not
    Then with the declaration the test does not trace, and without it the code in the file traces it
    And a file the source lists no test in, such as a YAML flow under a Jest pattern, is read as without a source

  @TSTRT-E03 @unit-level
  Scenario: A failing tests source fails the gate naming the error
    Given a project whose tests script exits with an error
    When a test is confronted
    Then the gate fails naming the script's error

  @TSTRT-B13 @unit-level
  Scenario: A support file is not charged with tracing
    Given a support file a feature links to, which cites no code
    When test-traceable judges it
    Then it skips saying the file is support
