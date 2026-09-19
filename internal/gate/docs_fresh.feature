# language: en
# @anchors
#   ref: DCFRD
#   updated_at: 2026-09-19
#   layer: feature

@DCFRD
Feature: DocsFresh — the compiled document has to reflect the spec

  @DCFRD-B01 @unit-level
  Scenario: A compiled document that no longer matches the spec fails
    Given a project whose documentation was compiled from its specs
    And a spec that changed afterwards with nobody recompiling
    When the gate confronts the spec
    Then it returns Fail, because the page still asserts the old rule with real content and
      is therefore convincing

  @DCFRD-B02 @unit-level
  Scenario: The verdict names the stale documents
    Given a project with one compiled document left behind by a spec change
    When the gate confronts the spec
    Then the message carries that document's name and the directory it lives in, so the
      author does not have to diff the whole output

  @DCFRD-B03 @unit-level
  Scenario: A compiled document that matches the templates passes
    Given a project whose documentation was compiled from the specs as they stand now
    When the gate confronts the spec
    Then it returns Pass

  @DCFRD-B04 @unit-level
  Scenario: An artifact that is not a spec leaves without a verdict
    Given a node whose kind is code, test or feature
    When the gate confronts it
    Then it returns Skip, because only the source of the documentation is confronted

  @DCFRD-B05 @unit-level
  Scenario: A project with no template directory is not charged
    Given a project with specs and no template directory at all
    When the gate confronts a spec
    Then it returns Skip, because the project never opted into the compiled-documentation
      mechanism

  @DCFRD-B06 @unit-level
  Scenario: A template whose document was never produced counts as stale
    Given a template declared in the project and no compiled document beside it
    When the gate confronts the spec
    Then it returns Fail, because absence is staleness — there is nothing asserting what the
      spec says today

  @DCFRD-B07 @unit-level
  Scenario: A document written by hand is not charged
    Given a compiled path holding a page written by hand, with no generation marker
    When the gate confronts the spec
    Then it returns Pass, because the build refuses to overwrite it and demanding a command
      that changes nothing is a warning nobody can act on

  @DCFRD-B08 @unit-level
  Scenario: A template that cannot be compiled fails
    Given a template the compiler cannot evaluate
    When the gate confronts the spec
    Then it returns Fail carrying the compiler's own error, because a broken template is a
      defect and not a reason to go quiet

  @DCFRD-B09 @unit-level
  Scenario: A project whose specs cannot be read leaves without a verdict
    Given a map that declares a spec the disk does not hold
    When the gate confronts the project
    Then it returns Skip carrying the reason, because a gate that could not look must not
      approve

  @DCFRD-B10 @unit-level
  Scenario: The comparison happens in memory
    Given a project whose compiled document is stale
    When the gate confronts the spec
    Then no compiled document is produced, because the comparison never leaves memory

  @DCFRD-I01 @unit-level
  Scenario: The gate never repairs what it points at
    Given a project whose compiled document is stale
    When the gate confronts the spec twice in a row
    Then both runs return the same verdict, because a gate that fixed the defect would pass
      on the second run and the defect would only surface for whoever cloned the repository

  @DCFRD-I02 @unit-level
  Scenario: The charge starts from the spec and never from the compiled document
    Given a project whose compiled document is stale
    When the gate runs over every node of the project
    Then the verdict lands on the spec, because the compiled document is a build product and
      not a declared unit

  @DCFRD-X01 @unit-level
  Scenario: The gate does not judge whether the compiled document is good
    Given a template that selects a passage nobody would have chosen
    And a compiled document equal to what that template produces now
    When the gate confronts the spec
    Then it returns Pass, because the ruler is equality with the templates — what the
      template selects is the author's decision

  @DCFRD-X02 @unit-level
  Scenario: The gate does not run the build even knowing the fix
    Given a project whose compiled document is stale
    When the gate confronts the spec
    Then the document on disk is left exactly as it was, because the check runs in a hook
      and in CI and must not depend on having run before

  @DCFRD-X03 @unit-level
  Scenario: The gate does not charge documents written by hand
    Given a compiled path holding a hand-written page and a template that would produce
      something else entirely
    When the gate confronts the spec
    Then it returns Pass, because the compiler will not write that file and the author would
      be sent to run a command that changes nothing

  @DCFRD-X04 @unit-level
  Scenario: The compiled document is not confronted as an artifact of its own
    Given a compiled document sitting in the output directory
    When the gate runs over the project
    Then that document receives no verdict, because giving it a node would demand a spec, a
      feature and a test for a build product
