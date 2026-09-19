# language: en
# @anchors
#   ref: FTMFT
#   updated_at: 2026-09-19
#   layer: feature

@FTMFT
Feature: FeatureTestMatch — scenarios in feature must be implemented in test by code and description

  @FTMFT-B01 @unit-level
  Scenario: Non-feature artifacts skip confrontation
    Given an artifact that is not of kind feature
    When the gate confronts the artifact
    Then it returns Skip, restricting scope strictly to feature files

  @FTMFT-B02 @unit-level
  Scenario: A nil graph returns pending without approving
    Given an execution context with no built graph
    When the gate confronts the artifact
    Then it returns Pending, because linked tests cannot be resolved without a map

  @FTMFT-B03 @unit-level
  Scenario: A feature declaring no scenarios skips confrontation
    Given a feature file containing no scenarios with requirement codes
    When the gate confronts the artifact
    Then it returns Skip, because there are no scenarios to confront

  @FTMFT-B04 @unit-level
  Scenario: A feature with no linked tests returns pending
    Given a feature with scenarios and no outgoing tested-by edges in the map
    When the gate confronts the artifact
    Then it returns Pending, waiting for test files to be linked

  @FTMFT-B05 @unit-level
  Scenario: Scenarios belonging to non-test surfaces are skipped
    Given a feature containing scenarios tagged exclusively for e2e or vr surfaces
    When the gate confronts the artifact
    Then it returns Pass, skipping scenarios evaluated on other surfaces

  @FTMFT-B06 @unit-level
  Scenario: A scenario code completely absent from tests fails
    Given a feature scenario whose code does not appear in any linked test body
    When the gate confronts the artifact
    Then it returns Fail, treating the missing code as a blocking defect

  @FTMFT-B07 @unit-level
  Scenario: A scenario code appearing only in comments fails
    Given a linked test containing the scenario code only within comment lines
    When the gate confronts the artifact
    Then it returns Fail, because comments do not constitute executable implementation

  @FTMFT-B08 @unit-level
  Scenario: Tests implementing scenario codes with exact titles pass
    Given linked tests implementing all feature codes with exact matching titles
    When the gate confronts the artifact
    Then it returns Pass, confirming full implementation and title correspondence

  @FTMFT-B09 @unit-level
  Scenario: Tests with matching codes but drifting descriptions issue a warning
    Given linked tests implementing feature codes with divergent description text
    When the gate confronts the artifact
    Then it returns Pending, issuing an informative warning without blocking

  @FTMFT-B10 @unit-level
  Scenario: Test titles containing quotes are parsed without truncation
    Given a test title with nested quotation marks
    When the gate confronts the artifact
    Then it extracts the entire title without premature truncation

  @FTMFT-B11 @unit-level
  Scenario: Sibling scenario codes in composite test titles are extracted
    Given a test title listing multiple sibling scenario codes separated by slashes
    When the gate confronts the artifact
    Then it correctly attributes the test title to each cited code

  @FTMFT-B12 @unit-level
  Scenario: Shared test titles verify scenario presence through test body
    Given multiple scenarios sharing a single compound test title
    When the gate confronts the artifact
    Then it evaluates descriptive correspondence against the test body content

  @FTMFT-B13 @unit-level
  Scenario: Test comments contribute to descriptive match but not code presence
    Given a test with scenario vocabulary present in preceding comments
    When the gate confronts the artifact
    Then it accepts the descriptive alignment while requiring the code in test statements

  @FTMFT-B14 @unit-level
  Scenario: Scenario codes match with exact word boundaries
    Given a scenario code that forms a prefix of another code
    When the gate confronts the artifact
    Then it matches only the exact identifier without false substring collisions

  @FTMFT-B15 @unit-level
  Scenario: Go t.Run declarations are recognized as valid test titles
    Given a Go test file using t.Run declarations with scenario codes
    When the gate confronts the artifact
    Then it recognizes the subtest title and validates correspondence

  @FTMFT-B16 @unit-level
  Scenario: Script comment markers are stripped when verifying code presence
    Given a test file with hash comment markers
    When the gate confronts the artifact
    Then comment lines are stripped prior to checking code presence

  @FTMFT-B17 @unit-level
  Scenario: RootCode returns the root requirement code without scenario sub-index
    Given a scenario code carrying a numeric sub-index
    When RootCode processes the identifier
    Then it strips the sub-index suffix and returns the base requirement code

  @FTMFT-I01 @unit-level
  Scenario: Missing scenario code is always a failure
    Given a feature scenario missing from linked test code
    When the gate confronts the artifact
    Then it produces a Fail verdict, enforcing that missing code is never tolerated

  @FTMFT-I02 @unit-level
  Scenario: Descriptive divergence is always an informative warning
    Given a scenario with matching code and differing description
    When the gate confronts the artifact
    Then it produces a Pending verdict, ensuring text variation does not block delivery

  @FTMFT-I03 @unit-level
  Scenario: Code presence ignores comments while description matching reads them
    Given a test file with comments
    When the gate confronts the artifact
    Then it uses stripped content for code presence and commented content for description

  @FTMFT-X01 @unit-level
  Scenario: Static analysis does not run tests or inspect execution results
    Given a linked test file
    When the gate confronts the artifact
    Then it inspects source text statically without invoking execution runners

  @FTMFT-X02 @unit-level
  Scenario: Non-unit surfaces are left to their respective gates
    Given scenarios mapped to e2e or vr surfaces
    When the gate confronts the artifact
    Then it bypasses them, delegating enforcement to surface-specific checkers

  @FTMFT-X03 @unit-level
  Scenario: Minor description drift does not block promotion
    Given minor phrasing variations between scenario and test
    When the gate confronts the artifact
    Then it emits a warning instead of failing the check
