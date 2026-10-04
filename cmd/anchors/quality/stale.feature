# language: en
# @anchors
#   code: STFTA
#   ref: STEDS
#   updated_at: 2026-10-03
#   layer: feature

@STEDS
Feature: StaleEdges — lists the confrontation debt: expired test evidence and stale edges

  @STEDS-B01 @unit-level
  Scenario: Expired test evidence is listed before the stale edges
    Given a map with a test whose dependency advanced since ingestion and a stale edge
    When the stale command runs
    Then the heading "1 EXPIRED test evidence(s)" comes before the stale edges heading
    And the note that expired evidence is not a test defect is printed

  @STEDS-B02 @unit-level
  Scenario: Each expired evidence names why it expired
    Given a test whose own file changed and a test whose own file and a dependency changed
    When the stale command runs
    Then the first reads "the test file itself changed"
    And the second reads "1 dependencie(s) changed, e.g.: a.go (and the test itself)"

  @STEDS-B03 @unit-level
  Scenario: Stale edges are split into never validated and drifted
    Given a map with an edge stamped at an older spec revision and an edge never stamped
    When the stale command runs
    Then the first edge is listed with "(advanced rev)" and the second with "(never validated)"
    And the summary reads "1 never validated, 1 with rev drift"

  @STEDS-B04 @unit-level
  Scenario: A map with every edge validated prints the clean message
    Given a map whose one edge is stamped at the current revisions
    When the stale command runs
    Then it prints "no stale edges (1 edges, all validated)"

  @STEDS-I01 @unit-level
  Scenario: An edge stamped at the current revisions is never listed
    Given a map with an edge stamped at the current revisions of both ends
    When the stale command runs
    Then that edge is not in the stale listing

  @STEDS-X01 @unit-level
  Scenario: Listing the stale edges leaves the map unchanged
    Given a map with a stale edge
    When the stale command runs
    Then the map file has the same content as before

  @STEDS-E01 @unit-level
  Scenario: The stale command without a map points at the map build
    Given a project with configuration and no map
    When the stale command runs
    Then it fails with an error naming "anchors map build"
