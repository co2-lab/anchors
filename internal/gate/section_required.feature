# language: en
# @anchors
#   code: SCRQF
#   ref: SCRQR
#   updated_at: 2026-10-19
#   layer: feature

@SCRQR
Feature: SectionRequired — the spec sections a unit's code calls for exist in its spec

  @SCRQR-B01 @unit-level
  Scenario: Nothing to confront leaves without a verdict
    Given a code node, a spec with no code, and a spec whose code's layer declares no requirement
    When section-required confronts each
    Then each leaves without a verdict

  @SCRQR-B02 @unit-level
  Scenario: A triggered section the spec lacks fails, named once
    Given a component layer requiring "Test IDs" when the code exposes a handle, and "Contrato de Dados" when it depends on the store layer
    And a component exposing ":login-submit" and depending on a store, whose spec has neither section
    When section-required confronts the spec
    Then it fails naming both sections, the file and what triggered each

  @SCRQR-B03 @unit-level
  Scenario: A section present satisfies the requirement
    Given the same component whose spec has "### test ids" and "## Contrato de Dados"
    When section-required confronts the spec
    Then it passes

  @SCRQR-B04 @unit-level
  Scenario: Each trigger reads what the core measures
    Given a requirement per trigger — exposes, depends_on_layer, flag "@navigates" and code_matches "showToast\("
    And code that triggers none of them: no handle, no dependency, the flag only inside a string and the pattern only inside a comment
    When section-required confronts the spec, and then code that triggers each
    Then nothing is called for at first, and each section is named once its trigger holds

  @SCRQR-B05 @unit-level
  Scenario: A trigger that cannot be measured is pending, not silent
    Given a requirement on exposed handles in a project declaring no dialect.test_handles.exposed
    When section-required confronts the spec
    Then it is pending, naming dialect.test_handles.exposed

  @SCRQR-E01 @unit-level
  Scenario: A code file that cannot be read calls for nothing
    Given a spec whose specified code file is gone from disk
    When section-required confronts the spec
    Then it names no section from the file it could not read

  @SCRQR-B06 @unit-level
  Scenario: A dependency on types only calls for nothing unless asked
    Given a component whose only dependency on the store layer is an import of its types
    When section-required confronts its spec with no data contract section
    Then nothing is called for
    And with kinds [value, type] the data contract is named, saying it depends on the store's types
