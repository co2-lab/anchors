# language: en
# @anchors
#   ref: DTAUI
#   updated_at: 2026-09-26
#   layer: feature

@DTAUI
Feature: Audit — the dossier of everything pending on one file, for fixing it in one pass

  @DTAUI-B01 @unit-level
  Scenario: A file the map does not know is refused
    Given a project whose map does not know the file new.go
    When anchors audit runs on new.go
    Then it fails with an error saying "new.go" is not in the map

  @DTAUI-B02 @unit-level
  Scenario: Without impact the dossier covers the file alone
    Given a project where a.go has a failing gate and other.go is also in the map
    When anchors audit runs on a.go without --impact
    Then the header reads "audit: a.go — the file"
    And nothing about other.go is printed

  @DTAUI-B03 @unit-level
  Scenario: With impact the dossier covers the unit on the impact path
    Given a project where a.spec.md specifies a.go and other.go is unrelated
    When anchors audit runs on a.spec.md with --impact
    Then the header announces "the unit (" with the node count
    And a.go is listed as "○ a.go (impact)" while other.go stays out

  @DTAUI-B04 @unit-level
  Scenario: A gate result shows by verdict with the first line of its detail
    Given a failing gate whose detail has two lines and a gate that awaits judgment
    When anchors audit runs on the file
    Then the failing gate reads "✗ [gate] always-fails — broken line one" without the second line
    And the judgment gate is marked ⏳

  @DTAUI-B05 @unit-level
  Scenario: A doctor finding shows by severity when it cites a node in scope
    Given a warning finding and an informational finding about a.go
    When the dossier is printed
    Then the warning reads "⚠ [doctor:orphan] no edges" and the information reads "ℹ [doctor:hint] just a note"

  @DTAUI-B06 @unit-level
  Scenario: Only a failing gate and a warning finding count as actionable
    Given one pending gate, one warning finding and one informational finding in scope
    When the dossier is printed
    Then it counts "1 actionable pending item(s)"

  @DTAUI-B07 @unit-level
  Scenario: A scope with nothing pending says so
    Given only passing gate results and no finding
    When the dossier is printed
    Then it says "✓ nothing pending"

  @DTAUI-B08 @unit-level
  Scenario: The target prints first and the impact nodes after it
    Given pending items on the target a.go and on a.spec.md
    When the dossier is printed
    Then "● a.go" appears before "○ a.spec.md (impact)"

  @DTAUI-I01 @unit-level
  Scenario: Nothing outside the audited scope reaches the dossier
    Given a warning finding about elsewhere.go, which is not in the scope
    When the dossier is printed
    Then the finding "out of scope" is not printed

  @DTAUI-E01 @unit-level
  Scenario: A project without configuration fails loading it
    Given a directory with no anchors.yaml
    When anchors audit runs there
    Then it fails with "load config"

  @DTAUI-E02 @unit-level
  Scenario: A project without a map fails pointing at the map build
    Given a project with anchors.yaml and no map
    When anchors audit runs there
    Then it fails with "load map" and the hint "run `anchors map build`"
