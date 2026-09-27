# language: en
# @anchors
#   ref: FLRSA
#   updated_at: 2026-09-26
#   layer: feature

@FLRSA
Feature: Failures — the observed failures that the spec has not explained yet

  @FLRSA-B01 @unit-level
  Scenario: Only the ingested failures whose rule carries no conclusion are listed, with the spec that declares them
    Given a login spec declaring three failures, two of them concluded as resilient or under observation
    And a log whose occurrences of the three were ingested
    When the failures review runs
    Then only the unconcluded failure is listed, under the header of one open failure
    And the login spec is named as the one that declares it

  @FLRSA-B02 @unit-level
  Scenario: The open failures are listed from the most frequent to the least
    Given a map holding three open failures in the order one, two and five occurrences
    When the failures review runs
    Then the failure with five occurrences comes first, then the one with two, then the one with one

  @FLRSA-B03 @unit-level
  Scenario: With all the concluded failures are listed too, with their conclusion
    Given a login spec whose ingested failures include one resilient and one under observation
    When the failures review runs with the option to show everything
    Then the resilient reason and what is under observation are printed with their failures

  @FLRSA-B04 @unit-level
  Scenario: An occurrence measured against another version of the spec is flagged
    Given an ingested failure whose recorded spec revision differs from the spec's current revision
    When the failures review runs
    Then the occurrence is flagged saying the spec changed since the ingestion

  @FLRSA-B05 @unit-level
  Scenario: With nothing open the review says so instead of printing an empty list
    Given a login spec in which every ingested failure carries a conclusion
    When the failures review runs
    Then it prints that no failure is under observation

  @FLRSA-I01 @unit-level
  Scenario: What the log ingestion binds to a spec is what the failures review lists
    Given a log ingested with the real log ingestion into a project with a login spec
    When the failures review runs
    Then the open failure bound by the ingestion is listed with the login spec

  @FLRSA-X01 @unit-level
  Scenario: The failures review leaves the map and the specs unchanged
    Given a project with ingested failures, its map and its login spec
    When the failures review runs with the option to show everything
    Then the map and the login spec are byte for byte what they were

  @FLRSA-E01 @unit-level
  Scenario: Without a map the failures review fails
    Given a directory with no map
    When the failures review runs on it
    Then the command fails

  @FLRSA-E02 @unit-level
  Scenario: A spec whose file can no longer be read is skipped by the review
    Given a map with ingested failures on a login spec whose file was deleted afterwards
    When the failures review runs
    Then none of its failures is listed and the review prints that no failure is under observation
