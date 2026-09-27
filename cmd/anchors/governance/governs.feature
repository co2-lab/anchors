# language: en
# @anchors
#   ref: GVRNS
#   updated_at: 2026-09-26
#   layer: feature

@GVRNS
Feature: Governs — who each guide governs, and how many, read from the map

  @GVRNS-B01 @unit-level
  Scenario: The board ranks each guide by how many files it governs
    Given a map where GUIDE.md governs three files and STYLE.md governs one
    When anchors governs runs with no argument
    Then it prints "2 guide(s) with governance", GUIDE.md with 3 before STYLE.md with 1
    And the total reads "total pairs (guide, governed) = 4"

  @GVRNS-B02 @unit-level
  Scenario: A map without governance has an empty board
    Given a map with no governs edge
    When anchors governs runs with no argument
    Then it prints "no guide governs anything"

  @GVRNS-B03 @unit-level
  Scenario: The detail of a guide groups the files it governs by kind
    Given a map where GUIDE.md governs a.go, b.go and a.spec.md
    When anchors governs runs on GUIDE.md
    Then it prints "GUIDE.md governs 3 file(s):" with "[code] 2" before "[spec] 1" and each file listed

  @GVRNS-B04 @unit-level
  Scenario: A file that governs nobody is answered, not refused
    Given a map where a.go governs nothing
    When anchors governs runs on a.go
    Then it prints "a.go governs nobody" and succeeds

  @GVRNS-B05 @unit-level
  Scenario: The guide argument is resolved against the root and the map can be given
    Given a project whose map is kept outside it
    When anchors governs runs with --map pointing at that map and the absolute path of GUIDE.md
    Then it prints "GUIDE.md governs 3 file(s):"

  @GVRNS-I01 @unit-level
  Scenario: Only governs edges count as governance
    Given a map with four governs edges and one specifies edge
    When anchors governs runs with no argument
    Then the total of pairs is 4

  @GVRNS-E01 @unit-level
  Scenario: A missing map fails pointing at the map build
    Given a project with no map
    When anchors governs runs
    Then it fails with an error that says to run `anchors map build`
