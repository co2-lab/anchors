# language: en
# @anchors
#   ref: CLMNC
#   updated_at: 2026-09-26
#   layer: feature

@CLMNC
Feature: CliMain — the entry point that stamps the build identity, prints a failure once and turns it into the exit code the hooks read

  @CLMNC-B01 @unit-level
  Scenario: A failing command prints its error once and exits 1
    Given an empty directory
    When anchors no-such-command runs
    Then it exits 1
    And stderr starts with "error: " and holds the unknown-command message once

  @CLMNC-B02 @unit-level
  Scenario: A file the project does not govern exits with the not-governed code
    Given a governed project with a map and a package.json no layer matches
    When anchors check --changed package.json --no-record runs
    Then it exits 3
    And stderr says the file is not governed

  @CLMNC-B03 @unit-level
  Scenario: The build version reaches the reported version and the map's generator
    Given a build with no link-time values
    When anchors --version runs, then anchors map build
    Then it reports "dev (commit none, built unknown)"
    And the map records generated_by: dev

  @CLMNC-B04 @unit-level
  Scenario: A renamed key in an older config points at migrate
    Given an anchors.yaml in format 1 holding trinca_opcional
    When anchors generated-paths runs
    Then it exits 1 saying `trinca_opcional` was renamed and to run anchors migrate
