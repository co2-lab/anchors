# language: en
# @anchors
#   code: DCFTA
#   ref: DCCMD
#   updated_at: 2026-10-03
#   layer: feature

@DCCMD
Feature: DocsCommand — the documentation is compiled from the specs through templates, against a map rebuilt from the tree

  @DCCMD-B01 @unit-level
  Scenario: Build compiles from the tree, even what the map on disk does not know
    Given a spec pkg/Coisa.spec.md, a template over the spec layer, and a map on disk with no node
    When docs build runs
    Then docs/camadas.md holds "O que a unidade faz"

  @DCCMD-B02 @unit-level
  Scenario: No-map-rebuild compiles against the map on disk and warns
    Given the same project
    When docs build runs with --no-map-rebuild
    Then stderr says the compiled output comes from the map on disk
    And the spec the map does not know does not reach the page

  @DCCMD-B03 @unit-level
  Scenario: A dry run compiles without writing
    Given a template doct/guia.md.tmpl
    When docs build runs with --dry-run
    Then it lists "docs/guia.md compiled (not written)"
    And docs/guia.md does not exist

  @DCCMD-B04 @unit-level
  Scenario: A hand-written page is skipped, kept and named
    Given a template doct/manual.md.tmpl and a hand-written docs/manual.md
    When docs build runs
    Then it reports "1 file(s) NOT generated" naming docs/manual.md
    And docs/manual.md still reads "written by hand"

  @DCCMD-B05 @unit-level
  Scenario: Build with no template says there is nothing to compile
    Given an empty template directory
    When docs build runs
    Then it says there is no template in doct/

  @DCCMD-B06 @unit-level
  Scenario: Init writes the skeleton once and needs a map
    Given a project with a map
    When docs init runs twice after one template was edited, then with --force
    Then the first run writes the templates
    And the second says the edited one already exists and keeps it
    And --force rewrites it
    And and with no map init points at anchors map build

  @DCCMD-B07 @unit-level
  Scenario: Duties answers for the project, a layer or a unit
    Given an openapi document triggered by layer api and a c4 document
    When duties runs, then with --layer api, --layer web and --unit api/handler.go
    Then all duties list both documents
    And api lists only docs/api.yaml
    And web says no documentation is required
    And the unit lists docs/api.yaml

  @DCCMD-B08 @unit-level
  Scenario: Duties without a declaration teaches the known kinds
    Given a config with no docs section
    When duties runs
    Then it says the project declares no required documentation and lists openapi, c4

  @DCCMD-X01 @unit-level
  Scenario: Build never writes the map
    Given a map on disk with no node
    When docs build runs and rebuilds the map from the tree
    Then the map on disk is byte for byte what it was

  @DCCMD-E01 @unit-level
  Scenario: Build without a config points at init
    Given a directory with no anchors.yaml
    When docs build runs
    Then it fails pointing at anchors init

  @DCCMD-E02 @unit-level
  Scenario: Duties without a config fails
    Given a directory with no anchors.yaml
    When duties runs
    Then it fails

  @DCCMD-E03 @unit-level
  Scenario: No-map-rebuild without a map says so
    Given a project with no map on disk
    When docs build runs with --no-map-rebuild
    Then it fails with "load map"
