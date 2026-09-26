# language: en
# @anchors
#   ref: DEPHN
#   updated_at: 2026-09-26
#   layer: feature

@DEPHN
Feature: DependencyHonored — methods promised in the dependency table are consumed in code

  @DEPHN-B01 @unit-level
  Scenario: An artifact that is not a spec leaves without a verdict
    Given a node whose kind is code, feature or test
    When the gate confronts it
    Then it returns Skip, because only specs have a Dependency Table

  @DEPHN-B02 @unit-level
  Scenario: Without a relational map the verdict is undetermined
    Given a spec node evaluated without a relational graph
    When the gate confronts it
    Then it returns Pending, because dependency and specification edges cannot be traversed without a graph

  @DEPHN-B03 @unit-level
  Scenario: A spec declaring no confrontable symbols leaves without a verdict
    Given a spec whose Dependency Table declares no backticked symbols
    When the gate confronts it
    Then it returns Skip, because prose descriptions and empty tables promise no verifiable identifiers

  @DEPHN-B04 @unit-level
  Scenario: A spec governing no code leaves the verdict undetermined
    Given a spec that promises dependency symbols but specifies no code files
    When the gate confronts it
    Then it returns Pending, because there is no governed code to confront yet

  @DEPHN-B05 @unit-level
  Scenario: When every promised symbol appears in governed code, the gate passes
    Given a spec whose promised dependency symbols all appear in the governed code
    When the gate confronts it
    Then it returns Pass with an empty message

  @DEPHN-B06 @unit-level
  Scenario: When a promised symbol is absent from governed code, the gate fails
    Given a spec declaring a promised dependency symbol that does not appear in governed code
    When the gate confronts it
    Then it returns Fail, naming the unused symbol and dependency target file

  @DEPHN-B07 @unit-level
  Scenario: When an absent symbol resembles an identifier in code, the verdict suggests the rename
    Given a spec promising a symbol absent from code while code uses a closely related identifier name
    When the gate confronts it
    Then it returns Fail and includes a suggestion to rename the symbol

  @DEPHN-B08 @unit-level
  Scenario: Symbols appearing only in comments do not fulfill the promise
    Given a governed code file where a promised symbol appears only inside line comments
    When the gate confronts the spec
    Then it returns Fail, because comments are stripped before checking dependency usage

  @DEPHN-I01 @unit-level
  Scenario: Prose descriptions in dependency methods are never treated as contracts
    Given a spec with dependency methods written as plain prose without backticks
    When the gate extracts promised symbols
    Then no symbols are extracted, because prose conveys intent rather than a verifiable contract

  @DEPHN-I02 @unit-level
  Scenario: Symbol presence is matched strictly on token word boundaries
    Given a governed code file that contains a longer identifier embedding the promised symbol as a substring
    When the gate verifies symbol presence
    Then it considers the symbol unused, because substring occurrences do not match token word boundaries

  @DEPHN-I03 @unit-level
  Scenario: Near-symbol rename suggestions are strictly conservative
    Given an absent symbol and candidate code identifiers that do not share a prefix or suffix extension or have insufficient length
    When the gate searches for a nearest symbol
    Then no suggestion is produced, avoiding arbitrary or noisy edit-distance guesses

  @DEPHN-X01 @unit-level
  Scenario: The gate performs static textual confrontation without runtime execution
    Given governed code containing the promised symbol identifier statically
    When the gate confronts the unit
    Then it evaluates token presence in source text without executing the program or validating runtime call graphs

  @DEPHN-X02 @unit-level
  Scenario: The gate does not interpret dependency semantics or parameter signatures
    Given a dependency table promising an identifier name without type annotations or signatures
    When the gate confronts the unit
    Then it verifies only exact identifier token presence, delegating type semantics to the compiler

  @DEPHN-E01 @unit-level
  Scenario: A specified code file missing from disk is left out of the confrontation
    Given a spec whose map lists two specified code files, one of them no longer on disk
    And the file still on disk uses every promised symbol
    When the gate confronts the spec
    Then it passes, confronting only the file that is there
    And a symbol used by no file still on disk is charged as unused
