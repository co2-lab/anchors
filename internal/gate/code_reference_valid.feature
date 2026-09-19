# language: en
# @anchors
#   ref: CRVCD
#   updated_at: 2026-09-19
#   layer: feature

@CRVCD
Feature: CodeReferenceValid — cross-referenced requirement codes must resolve to existing units

  @CRVCD-B01 @unit-level
  Scenario: Non-specification artifacts skip confrontation
    Given an artifact whose kind is not a specification
    When the gate confronts it
    Then it returns Skip, avoiding confrontation on non-specification artifacts

  @CRVCD-B02 @unit-level
  Scenario: Confronting without a map graph returns pending
    Given a specification citing external requirement codes and a nil map graph
    When the gate confronts it
    Then it returns Pending, because identity owners cannot be determined without the map

  @CRVCD-B03 @unit-level
  Scenario: Confronting with an empty identity universe returns pending
    Given a specification and a map graph containing no known identity codes
    When the gate confronts it
    Then it returns Pending, refusing to approve against an empty identity universe

  @CRVCD-B04 @unit-level
  Scenario: Citations resolving to existing units in the map pass
    Given a specification citing external requirement codes that exist in the map
    When the gate confronts it
    Then it returns Pass, confirming that every cited unit exists

  @CRVCD-B05 @unit-level
  Scenario: Citations pointing to non-existent units fail
    Given a specification citing an external requirement code not present in the map
    When the gate confronts it
    Then it returns Fail, reporting the unresolvable code

  @CRVCD-B06 @unit-level
  Scenario: Citations matching the specification's own identity code pass
    Given a specification citing requirement codes that share its own identity code
    When the gate confronts it
    Then it returns Pass, recognizing self-references as valid

  @CRVCD-B07 @unit-level
  Scenario: Multiple orphaned requirement citations are reported sorted
    Given a specification citing multiple requirement codes absent from the map
    When the gate confronts it
    Then it returns Fail, listing all unresolvable codes in alphabetical order

  @CRVCD-B08 @unit-level
  Scenario: Identity ownership is resolved from graph nodes and header metadata
    Given a map graph where specification identities are declared in graph nodes or disk headers
    When the gate confronts a specification referencing those units
    Then it discovers the declared identities and validates the references

  @CRVCD-B09 @unit-level
  Scenario: Non-requirement tokens are ignored
    Given a specification containing general prose and identifier names that do not match requirement syntax
    When the gate confronts it
    Then it ignores those tokens without attempting to resolve them as requirement codes

  @CRVCD-I01 @unit-level
  Scenario: Self-references to a specification's own requirements never fail
    Given a specification referencing internal requirement codes defined in its own text
    When the gate confronts it
    Then it upholds the invariant that a specification owns its own code identifiers

  @CRVCD-I02 @unit-level
  Scenario: Missing map or empty identity universe returns pending rather than pass
    Given a specification evaluated with no graph or without known identities
    When the gate confronts it
    Then it returns Pending, never stamping approval over unmeasured territory

  @CRVCD-I03 @unit-level
  Scenario: Unresolvable external citations always produce a blocking fail verdict
    Given a specification containing dangling cross-references to absent units
    When the gate confronts it
    Then it produces a blocking Fail verdict, preventing empty traceability links

  @CRVCD-X01 @unit-level
  Scenario: Implementation correctness of referenced requirements is not evaluated
    Given a specification with valid cross-references to existing units
    When the gate confronts it
    Then it validates reference existence without inspecting downstream requirement semantics

  @CRVCD-X02 @unit-level
  Scenario: Non-specification artifacts are not inspected by this gate
    Given code or test artifacts containing requirement code citations
    When the gate confronts them
    Then it skips evaluation, leaving non-specification checks to other gates

  @CRVCD-X03 @unit-level
  Scenario: External requirement citations are not mandatory
    Given a specification that contains no external requirement code references
    When the gate confronts it
    Then it passes without demanding external citations
