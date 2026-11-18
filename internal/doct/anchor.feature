# language: en
# @anchors
#   code: ANFTN
#   ref: DCLND
#   updated_at: 2026-10-03
#   layer: feature

@DCLND
Feature: DocLinks — the anchors, links, sizes and layer arrows the documentation templates are given

  @DCLND-B01 @unit-level
  Scenario: A heading's anchor follows the GitHub convention and keeps accents
    Given the heading "A régua e a dívida" and the heading "`código` entre crases"
    When their anchors are computed
    Then they are "a-régua-e-a-dívida" and "código-entre-crases"

  @DCLND-B02 @unit-level
  Scenario: A rule's link points at the rule on a small layer and at its unit on a big one
    Given the rule "GLCGL-B01 — r" of the unit "GLCGL — GoLive" in layer "infra"
    When its link is built with the layer small and then big
    Then the first ends in "layers/infra.md#glcgl-b01--r"
    And the second ends in the anchor of the unit's heading
    And a rule written as a table row links to the unit's heading even on the small layer

  @DCLND-B03 @unit-level
  Scenario: A scenario's link falls back to the unit, then to the bare page, on a big layer
    Given a scenario of the unit "GLCGL" and a scenario whose spec is unknown, in layer "infra"
    When their links are built with the layer small and then big
    Then on the small layer each points at its own scenario heading
    And on the big layer the first points at the unit's heading and the second at "layers/infra.md"

  @DCLND-B04 @unit-level
  Scenario: A selection's size counts units, rules, lines and scenarios, and is remembered per selection
    Given one spec with three rules and a feature with two scenarios, and two layers of different sizes
    When the size of each selection is asked twice
    Then the spec's size is one unit, three rules and two scenarios
    And each selection returns the same answer the second time, different from the other selection's

  @DCLND-B05 @unit-level
  Scenario: A selection that matches no spec has size zero and is not big
    Given a filter naming a layer with no spec
    When its size and its bigness are asked
    Then the size is zero, no error is returned, and the selection is not big

  @DCLND-B06 @unit-level
  Scenario: Arrows between layers aggregate dependency edges between specs and carry their count
    Given two screen specs needing the same shared spec, a plan needing a spec, a tested-by edge, and an edge inside one layer
    When the arrows between layers are deduced
    Then there is exactly one arrow, from screen to shared, with count 2

  @DCLND-B07 @unit-level
  Scenario: The arrows are ordered by source layer then target layer
    Given dependency edges producing the arrows b→a, a→c and a→b
    When the arrows are deduced
    Then they come in the order a→b, a→c, b→a

  @DCLND-B08 @unit-level
  Scenario: A diagram node identifier has only ASCII letters, digits and underscores
    Given the names "feature-hook" and "padrão"
    When their node identifiers are computed
    Then they are "n_feature_hook" and "n_padr_o"

  @DCLND-I01 @unit-level
  Scenario: Every generated link resolves on a small layer and on a big one
    Given one layer below the cut-off and one layer above it
    When the default documentation skeleton is built
    Then every link in the generated pages points at a heading anchor that exists on its target page

  @DCLND-X01 @unit-level
  Scenario: The bigness asked by a template has no threshold of its own
    Given a layout whose cut-off is changed after the compiler was created
    When a template asks whether a selection is big
    Then the answer follows the compiler's layout

  @DCLND-B09 @unit-level
  Scenario: A layer's page is in the folder of its template, or the project language's folder
    Given a project with no template for the layer infra and no language
    When the page of infra is asked
    Then it is layers/infra.md
    And with the template doct/paginas/infra.md.tmpl it is paginas/infra.md
    And the language folders are camadas for pt-BR, capas for es and layers for en
