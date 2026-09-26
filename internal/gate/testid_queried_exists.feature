# language: en
# @anchors
#   ref: TQETS
#   updated_at: 2026-09-26
#   layer: feature

@TQETS
Feature: TestidQueriedExists — every handle queried by an E2E flow must exist in code

  @TQETS-B01 @unit-level
  Scenario: The gate skips when no test handle attribute is declared
    Given a configuration without a declared test handle attribute
    When the gate confronts the artifact
    Then it returns Skip, avoiding false approval when no handle attribute was configured

  @TQETS-B02 @unit-level
  Scenario: The gate skips when no E2E surface is configured
    Given a configuration without a declared E2E surface
    When the gate confronts the artifact
    Then it returns Skip, acknowledging that no flow files were configured for confrontation

  @TQETS-B03 @unit-level
  Scenario: The gate skips when no exposed handles exist in code
    Given a codebase where no handles matching the attribute are present
    When the gate confronts the artifact
    Then it returns Skip, because there is no exposed surface to confront

  @TQETS-B04 @unit-level
  Scenario: A queried handle that exists in code passes
    Given a flow querying a test handle
    And application code exposing that handle on a component
    When the gate confronts the artifact
    Then it returns Pass, confirming the handle exists in code

  @TQETS-B05 @unit-level
  Scenario: A queried handle missing from code fails
    Given a flow querying a test handle
    And no application code exposing that handle
    When the gate confronts the artifact
    Then it returns Fail, catching an invented handle before test execution

  @TQETS-B06 @unit-level
  Scenario: The verdict names the missing handle and the querying flows
    Given a flow querying an unexposed handle
    When the gate confronts the artifact
    Then the failure verdict names the missing handle and the flow file that queries it

  @TQETS-B07 @unit-level
  Scenario: A negative assertion on a missing handle fails
    Given a flow with an assertNotVisible step querying a non-existent handle
    When the gate confronts the artifact
    Then it returns Fail, eliminating false-green passes caused by vacuity

  @TQETS-B08 @unit-level
  Scenario: An exposed template handle covers a queried instance
    Given application code exposing a dynamic template handle
    And a flow querying a concrete instance of that template
    When the gate confronts the artifact
    Then it returns Pass, matching the concrete handle against the template pattern

  @TQETS-B09 @unit-level
  Scenario: A regex pattern in a flow matches an exposed prefix head
    Given application code exposing a handle
    And a flow querying a pattern matching that handle prefix
    When the gate confronts the artifact
    Then it returns Pass, resolving the pattern against the exposed prefix

  @TQETS-B10 @unit-level
  Scenario: A flow handle with runtime interpolation does not trigger a failure
    Given a flow querying a handle with runtime interpolation syntax
    When the gate confronts the artifact
    Then it returns Pass, skipping expressions that require runtime evaluation

  @TQETS-B11 @unit-level
  Scenario: A marked handle defined in a lookup table counts as exposed
    Given application code defining a marked handle inside a constant table
    And a flow querying that handle
    When the gate confronts the artifact
    Then it returns Pass, recognizing marked literals outside standard JSX attributes

  @TQETS-B12 @unit-level
  Scenario: A template behind a nullish coalescing fallback counts as exposed
    Given application code with a template handle inside a fallback expression
    And a flow querying an instance of that template
    When the gate confronts the artifact
    Then it returns Pass, identifying handles defined within fallback branches

  @TQETS-B13 @unit-level
  Scenario: A suffix composed in a child component counts as exposed
    Given a parent component providing a handle prop
    And a child component composing a row suffix onto that prop
    And a flow querying the resulting composite handle
    When the gate confronts the artifact
    Then it returns Pass, recognizing handle composition across components

  @TQETS-B14 @unit-level
  Scenario: A terminal suffix composed from a prop counts as exposed
    Given a component composing a terminal toggle suffix onto a prop
    And a flow querying the composite terminal handle
    When the gate confronts the artifact
    Then it returns Pass, recognizing composite handles ending in dynamic suffixes

  @TQETS-B15 @unit-level
  Scenario: A handle appearing only in test files does not count as exposed
    Given a handle cited exclusively in a test file
    And a flow querying that handle
    When the gate confronts the artifact
    Then it returns Fail, because test files consume handles rather than exposing them

  @TQETS-I01 @unit-level
  Scenario: Missing configuration never reports a pass
    Given an unconfigured project environment
    When the gate confronts the artifact
    Then it returns Skip, ensuring unmeasured projects never receive approval

  @TQETS-I02 @unit-level
  Scenario: Findings are grouped by handle identifier
    Given multiple flows querying the same unexposed handle
    When the gate confronts the artifact
    Then the verdict groups all querying flow locations under that handle

  @TQETS-I03 @unit-level
  Scenario: Negative assertions are checked against code presence
    Given a flow asserting absence of an element
    When the gate confronts the artifact
    Then it verifies the element exists in code so negative checks remain meaningful

  @TQETS-X01 @unit-level
  Scenario: Handles exposed in code but unused in flows are not accused
    Given application code exposing extra test handles
    And flows that do not query those handles
    When the gate confronts the artifact
    Then it returns Pass, leaving single-direction enforcement from flows to code

  @TQETS-X02 @unit-level
  Scenario: Flow expressions requiring runtime execution are not evaluated
    Given flow files with dynamic JavaScript expressions
    When the gate confronts the artifact
    Then it restricts analysis to static segments without launching a runtime engine

  @TQETS-X03 @unit-level
  Scenario: Default test handle attributes are not inferred
    Given a project configuration without derived test handle settings
    When the gate confronts the artifact
    Then it skips execution instead of presuming testID or data-testid conventions

  @TQETS-E01 @unit-level
  Scenario: A declared E2E directory missing from disk skips the gate
    Given a project that declares an E2E surface whose directory does not exist
    When the gate confronts the artifact
    Then it returns Skip instead of Pass
