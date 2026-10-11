# language: en
# @anchors
#   code: GGFGD
#   ref: GDGTG
#   updated_at: 2026-10-11
#   layer: feature

@GDGTG
Feature: Guidance gates — the working agreement is one file, and the others point to it

  @GDGTG-B01 @unit-level
  Scenario: The guidance points to the agreement
    Given a project with CONTRIBUTING.md, an AGENTS.md pointing to it, a README.md that does not, and one with no CONTRIBUTING.md
    When guidance-linked runs on each
    Then CONTRIBUTING.md and AGENTS.md pass, the README fails naming itself, the file with no agreement fails naming it, and a code file is skipped

  @GDGTG-B02 @unit-level
  Scenario: The catalog offers the guidance gates
    Given a project that chose any artifact
    When the default gates are seeded
    Then guidance-linked and guidance-honored are there, informative, on the guidance tag, and the judgment is about the configuration and asks for contradictions and replacements

  @GDGTG-B03 @unit-level
  Scenario: A headless layer is asked no header
    Given a README.md with no header in a headless layer
    When header-valid and header-layer-declared run, and the header is fixed
    Then both skip it, and the fix writes nothing
