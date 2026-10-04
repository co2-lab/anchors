# language: en
# @anchors
#   code: GPFGN
#   ref: GNPTG
#   updated_at: 2026-10-03
#   layer: feature

@GNPTG
Feature: GeneratedPaths — the product names the files it derives, so a conflict in them is rebuilt, not merged

  @GNPTG-B01 @unit-level
  Scenario: The patterns match the map, the compiled docs and plan progress files only
    Given a governed project
    When the generated patterns are read
    Then anchors.graph.yaml, docs/arquitetura.md and plans/0001-progress.md match
    And sub/anchors.graph.yaml, internal/gate/x.spec.md and plans/0001-progress.md.bak do not

  @GNPTG-B02 @unit-level
  Scenario: The default output is one pattern per line
    Given a governed project
    When generated-paths runs with no format
    Then it prints three lines, one pattern each

  @GNPTG-B03 @unit-level
  Scenario: The re format is one valid alternation
    Given a governed project
    When generated-paths runs with --format re
    Then it prints one line that compiles as a regular expression

  @GNPTG-B04 @unit-level
  Scenario: The dot of a path matches only a literal dot
    Given the alternation of the generated patterns
    When it is matched against anchorsXgraph.yaml
    Then it does not match

  @GNPTG-I01 @unit-level
  Scenario: Both forms carry the same patterns
    Given a governed project
    When the line form and the alternation form are both printed
    Then the lines joined by | equal the alternation

  @GNPTG-X01 @unit-level
  Scenario: The command only names the derived files
    Given a governed project
    When generated-paths runs
    Then its output is exactly the three patterns

  @GNPTG-E01 @unit-level
  Scenario: A directory without anchors.yaml is refused
    Given a directory with no anchors.yaml
    When generated-paths runs there
    Then it fails
