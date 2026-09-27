# language: en
# @anchors
#   ref: DTCDC
#   updated_at: 2026-09-26
#   layer: feature

@DTCDC
Feature: DocTemplateCompiler — compiles documentation pages from templates that reference the specs' content

  @DTCDC-B01 @unit-level
  Scenario: The spec's content enters the compiled page
    Given a spec in layer infra with an overview section, and a template printing each infra spec's title and overview
    When the documentation is built
    Then the compiled page contains the overview's text and the spec's title

  @DTCDC-B02 @unit-level
  Scenario: The compiled page opens with the generated marker, its template path and a stamp
    Given a template "doct/x.md.tmpl"
    When the documentation is built
    Then the compiled page starts with the generated marker naming "doct/x.md.tmpl" and an inputs stamp

  @DTCDC-B03 @unit-level
  Scenario: A handwritten page is never overwritten
    Given a handwritten "docs/produto.md" and a template of the same name
    When the documentation is built
    Then the handwritten page is unchanged and "produto.md" is reported as skipped

  @DTCDC-B04 @unit-level
  Scenario: A dry run writes nothing
    Given a template
    When the documentation is built as a dry run
    Then no compiled page exists on disk

  @DTCDC-B05 @unit-level
  Scenario: A template in a subfolder is compiled
    Given a template at "doct/camadas/infra.md.tmpl"
    When the documentation is built
    Then "camadas/infra.md" is written with the selected spec's code

  @DTCDC-B06 @unit-level
  Scenario: The layer comes from the spec's header
    Given a spec whose header says infra while the map says spec
    When specs are selected by layer
    Then "layer=infra" finds it and "layer=spec" fails

  @DTCDC-B07 @unit-level
  Scenario: Only layers that have specs are offered
    Given a project whose only spec is in infra
    When the layers are listed
    Then the list is exactly infra

  @DTCDC-B08 @unit-level
  Scenario: A section runs to the next heading of its own level
    Given a spec whose rules section holds several rule headings, followed by an invariants section
    When the rules section is cut
    Then it contains every rule and nothing of the invariants section
    And a section the spec does not have is empty

  @DTCDC-B09 @unit-level
  Scenario: A spec is split into its separate rules
    Given a spec with rules B01 and B02 and invariant I01 as rule headings
    When its rules are listed
    Then there are three rules with their codes and titles, and B02's body stops before the next section

  @DTCDC-B10 @unit-level
  Scenario: A rule's title drops the HTML comment on its heading
    Given a rule heading ending in an HTML comment carrying a waiver
    When the rules index is built
    Then the rule's link label has the title's words and no comment

  @DTCDC-B11 @unit-level
  Scenario: A missing page or a page with a different stamp is stale
    Given a built documentation
    When one spec's body changes, and later a compiled page is deleted
    Then the page is reported stale in both cases, and a matching stamp is not

  @DTCDC-B12 @unit-level
  Scenario: Editing the template makes its page stale
    Given a built documentation
    When the template is edited
    Then its page is reported stale

  @DTCDC-B13 @unit-level
  Scenario: The header's date is not part of the stamp
    Given a spec whose only change is the updated_at line of its header
    When the stamp is computed before and after
    Then the stamp is the same, while a body change or a date line far below the header changes it

  @DTCDC-B14 @unit-level
  Scenario: A handwritten page is never reported stale
    Given a template whose page on disk has no generated marker
    When the stale pages are asked
    Then the page is not reported

  @DTCDC-B15 @unit-level
  Scenario: The specs no template reaches are uncovered
    Given two specs in layers infra and screen and a template selecting only infra
    When the uncovered specs are asked
    Then only the screen spec is reported

  @DTCDC-B16 @unit-level
  Scenario: Generated and handwritten markers are recognised
    Given a page with the generated marker, one with the handwritten marker and one with neither
    When each is inspected
    Then the first is generated, the second is handwritten, and the third is neither

  @DTCDC-I01 @unit-level
  Scenario: Right after a build nothing is stale
    Given templates at the root and in a subfolder
    When the documentation is built and the stale pages are asked
    Then none is reported

  @DTCDC-X01 @unit-level
  Scenario: Asking for stale pages writes nothing
    Given a compiled page made stale by a spec change
    When the stale pages are asked
    Then the page on disk is unchanged

  @DTCDC-X02 @unit-level
  Scenario: A spec file the map does not list is not loaded
    Given a spec file on disk that is not in the map
    When the compiler loads the specs
    Then no layer and no spec comes from it

  @DTCDC-E01 @unit-level
  Scenario: A wrong selection filter fails with a message that helps fix it
    Given the filters "layar=infra", "layer=screen", "infra" and "code=NAOEX"
    When specs are selected with each
    Then each fails, naming the unknown field, the existing layers, the missing "=" or the missing code

  @DTCDC-E02 @unit-level
  Scenario: A failing selection fails the build without writing the page
    Given a template selecting a layer that does not exist
    When the documentation is built
    Then the build fails and the page is not written

  @DTCDC-E03 @unit-level
  Scenario: A broken template fails the build
    Given a template with a range and no end
    When the documentation is built
    Then the build fails

  @DTCDC-E04 @unit-level
  Scenario: A spec in the map but not on disk fails the compiler
    Given a map listing a spec file that does not exist
    When the compiler is created
    Then it fails

  @DTCDC-E05 @unit-level
  Scenario: A project without templates fails the build and has nothing stale or uncovered
    Given a project with no doct folder
    When the documentation is built and the stale and uncovered pages are asked
    Then the build fails telling to run docs init, and the other two answer with nothing

  @DTCDC-E06 @unit-level
  Scenario: Asking one spec by an unknown code fails
    Given a project without the code NAOEX
    When a template asks for the spec NAOEX
    Then an error names NAOEX
