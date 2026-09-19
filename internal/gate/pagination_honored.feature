# language: en
# @anchors
#   ref: PGNHN
#   updated_at: 2026-09-19
#   layer: feature

@PGNHN
Feature: PaginationHonored — what promises a SET does not return the first page in silence

  @PGNHN-B01 @unit-level
  Scenario: A limit received from the caller passes
    Given an exported function whose signature takes the limit as a parameter
    When the gate confronts it
    Then it returns Pass, because the page is deliberate and whoever asked for it
      knows there is more

  @PGNHN-B02 @unit-level
  Scenario: A limit hidden in a default value is accused
    Given an exported function whose name promises the whole set
    And the limit lives in a default value the caller never sees
    When the gate confronts it
    Then it returns Fail, because the hundred-and-first row is never processed and
      nobody is told

  @PGNHN-B03 @unit-level
  Scenario: The NAME bounds the promise
    Given an exported function whose name promises no set at all
    And it returns a partial result
    When the gate confronts it
    Then it returns Pass, because there is no promise to break

  @PGNHN-B04 @unit-level
  Scenario: Sibling functions that paginate are the proof by asymmetry
    Given a module whose sibling functions loop until the cursor is exhausted
    And one exported function returns a single page
    When the gate confronts it
    Then it returns Fail, because the author knew the pattern — the one that does not
      paginate is forgetfulness, not decision

  @PGNHN-B05 @unit-level
  Scenario: A waiver with a written reason leaves the report
    Given an exported function carrying a waiver marker followed by the reason
    When the gate confronts it
    Then it returns Pass, and the function no longer appears in the report

  @PGNHN-B06 @unit-level
  Scenario: The verdict offers the way out
    Given an exported function accused of hiding the limit
    When the gate confronts it
    Then the verdict names the waiver marker, so whoever reads it learns the declared
      way out instead of guessing

  @PGNHN-B07 @unit-level
  Scenario: Without a declared dialect the verdict is undetermined
    Given a project that declares no dialect for its stack
    When the gate confronts any code
    Then it returns Pending, because approving without being able to read the code
      would stamp what was never checked

  @PGNHN-I01 @unit-level
  Scenario: The ruler is agnostic across languages
    Given the same hidden-limit defect written in two different stacks
    And each project declares its own dialect
    When the gate confronts both
    Then both are accused, because the confronted truth — the name promises a set, the
      return is partial — belongs to no language

  @PGNHN-I02 @unit-level
  Scenario: Where the construct is not recognised, the gate stays silent
    Given code whose pattern the declared dialect does not reach
    When the gate confronts it
    Then it accuses nothing, because a false positive here teaches the team to ignore
      the gate

  @PGNHN-I03 @unit-level
  Scenario: A cursor with no loop does not count as pagination
    Given an exported function that returns the cursor and never walks it
    When the gate confronts it
    Then it returns Fail, because the consumer is left with the same slice, now wearing
      the appearance of completeness

  @PGNHN-I04 @unit-level
  Scenario: A provider prefix in the name does not hide the promise
    Given an exported function whose name carries the provider prefix before the promise
    When the gate confronts it
    Then it returns Fail, because what the name says holds wherever it comes from

  @PGNHN-X01 @unit-level
  Scenario: The gate does not invent a cursor the provider does not offer
    Given an exported function calling a dependency with no pagination mechanism
    When the gate confronts it
    Then it returns Pass, because accusing the absence of a mechanism the dependency
      lacks would hand the author a defect that is not theirs

  @PGNHN-X02 @unit-level
  Scenario: The gate does not measure performance or page size
    Given an exported function that exposes its limit and returns a small page
    When the gate confronts it
    Then it returns Pass, because judging whether a hundred is many depends on the
      domain, and that is the project's decision
