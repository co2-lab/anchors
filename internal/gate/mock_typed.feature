# language: en
# @anchors
#   ref: MCTYM
#   updated_at: 2026-09-26
#   layer: feature

@MCTYM
Feature: MockTyped — every test double must derive from the module it replaces

  @MCTYM-B01 @unit-level
  Scenario: A double with no tie fails and the verdict names the loose module
    Given a test doubling a governed module with a factory carrying no tie
    When the gate confronts it
    Then it returns Fail naming that module, because the double is a frozen copy of a
      contract that may no longer exist

  @MCTYM-B02 @unit-level
  Scenario: A double whose factory carries the declared tie passes
    Given a test doubling a governed module
    And the factory is annotated with the tie shape the project declared
    When the gate confronts it
    Then it returns Pass, because the annotation is what makes the compiler check name,
      signature and return against the real module

  @MCTYM-B03 @unit-level
  Scenario: The charge is per module, not per file
    Given a test with one tied double and one loose double
    When the gate confronts it
    Then it returns Fail counting one loose double, because a verdict per file would give
      the hole as resolved

  @MCTYM-B04 @unit-level
  Scenario: A double with no factory is not charged
    Given a test declaring a double of a governed module with no factory at all
    When the gate confronts it
    Then it returns Skip, because the automock derives from the real module by construction
      and cannot drift

  @MCTYM-B05 @unit-level
  Scenario: Without the tie shape declared the gate goes quiet
    Given a project that declares no tie shape
    And a test carrying a loose double
    When the gate confronts it
    Then it returns Skip naming what has to be declared

  @MCTYM-B06 @unit-level
  Scenario: An artifact that is not a test leaves without a verdict
    Given a node whose kind is spec, code or feature
    When the gate confronts it
    Then it returns Skip, because the double lives in the test and charging elsewhere
      would accuse the wrong file

  @MCTYM-B07 @unit-level
  Scenario: A test that doubles nobody leaves without a verdict
    Given a test that asserts on plain values and declares no double
    When the gate confronts it
    Then it returns Skip and not Pass, because nothing was checked and a Pass would
      inflate the count of greens with nothing

  @MCTYM-B08 @unit-level
  Scenario: Another runner of the same ecosystem is recognised the same way
    Given a test doubling a governed module through a different runner of the same ecosystem
    And the factory carries no tie
    When the gate confronts it
    Then it returns Fail, because the double is the double regardless of which library
      spells it

  @MCTYM-B09 @unit-level
  Scenario: A third-party library double is not charged
    Given a test whose only doubles are of libraries the graph does not govern
    When the gate confronts it
    Then it returns Skip, because an external dependency has its version pinned and its
      double swaps a component for a stub instead of reproducing a contract

  @MCTYM-B10 @unit-level
  Scenario: The third-party exemption is not an escape hatch
    Given a test with a third-party double and one loose double of a governed module
    When the gate confronts it
    Then it returns Fail counting one loose double and naming only the own one

  @MCTYM-B11 @unit-level
  Scenario: A relative import that resolves in the map is an own module
    Given a test doubling a governed module through a relative path, with no tie
    When the gate confronts it
    Then it returns Fail, because the criterion is resolving to a node of the map and not
      the shape of the specifier

  @MCTYM-B12 @unit-level
  Scenario: Another ecosystem's dialect is charged the same way
    Given a project declaring another ecosystem's detection pattern and tie shape
    And one test doubling a governed module without the tie and another with it
    When the gate confronts each of them
    Then the first returns Fail naming the module and the second returns Pass

  @MCTYM-I01 @unit-level
  Scenario: What is governed is decided by the graph, never by a prefix list
    Given a test doubling the same governed module through an alias, a relative path and a bare name
    And none of the three carries a tie
    When the gate confronts each of them
    Then all three are charged, because the graph already exists, assumes no alias
      convention, and follows the project on its own

  @MCTYM-I02 @unit-level
  Scenario: An undeclared tie shape skips instead of guessing one
    Given a test carrying a loose double of a governed module
    And a project that declares no tie shape
    When the gate confronts it
    Then it returns Skip and not Pass, because inferring the TypeScript form would report
      green over what was never checked in any other ecosystem

  @MCTYM-I03 @unit-level
  Scenario: The verdict counts the loose doubles, not just names them
    Given a test with one tied double and one loose double
    When the gate confronts it
    Then the verdict says how many are loose, because naming one without saying how many
      leaves the reader unable to tell whether the others are on the list too

  @MCTYM-X01 @unit-level
  Scenario: The gate does not check whether the annotated type matches the real module
    Given a factory annotated with a type that has nothing to do with the doubled module
    When the gate confronts it
    Then it returns Pass, because the ruler here is that the tie WAS WRITTEN — checking
      that it matches is the compiler's job, and it does it better

  @MCTYM-X02 @unit-level
  Scenario: The gate does not cover drift of behaviour
    Given a tied double whose factory returns a value with the right shape and the wrong meaning
    When the gate confronts it
    Then it returns Pass, because this gate covers drift of SHAPE — promising behaviour
      would sell a green it does not give

  @MCTYM-X03 @unit-level
  Scenario: The gate carries no built-in tie shape and no built-in ecosystem
    Given a project in an ecosystem where a double is a satisfied interface and no call to detect
    And neither a tie shape nor a detection dialect declared
    When the gate confronts a test of that project
    Then it returns Skip, because the right verdict is that the gate does not apply —
      not a green over an unchecked file

  @MCTYM-X04 @unit-level
  Scenario: The gate does not charge third-party library doubles
    Given a test doubling several third-party libraries and no governed module
    When the gate confronts it
    Then it returns Skip, because a gate that accuses everything is switched off and takes
      the legitimate findings with it

  @MCTYM-X05 @unit-level
  Scenario: A double with no factory is not charged
    Given the project declares mock_detect
    And a test with an automock and a spy mock of governed modules
    When the gate confronts it
    Then neither is charged, and an unannotated factory still is

  @MCTYM-E02 @unit-level
  Scenario: With no map the verdict is pending, not an external-only skip
    Given a test that mocks a module and no map built
    When the gate confronts the test
    Then it returns Pending with the no-map message
