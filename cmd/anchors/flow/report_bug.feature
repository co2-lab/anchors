# language: en
# @anchors
#   code: RBFRP
#   ref: RPBUG
#   updated_at: 2026-10-03
#   layer: feature

@RPBUG
Feature: ReportBug — report a bug in Anchors itself, where it is fixed for every project

  @RPBUG-B01 @unit-level
  Scenario: The report has the bug form's sections
    Given a bug with what happened, what should happen and a minimal case
    When its issue is written
    Then the title is "[bug] " and the first line, and the body has each section, the version and the platform

  @RPBUG-B02 @unit-level
  Scenario: An open issue with the same title is told it was seen again
    Given co2-lab/anchors has an open issue titled "[bug] gate misreads a file"
    When `anchors report-bug "gate misreads a file" --expected "it reads it"` runs
    Then that issue receives a "Seen again" comment, nothing is created, and the output says "already reported to Anchors"

  @RPBUG-B03 @unit-level
  Scenario: A new bug becomes an issue at co2-lab/anchors
    Given co2-lab/anchors has no open issue with the title
    When `anchors report-bug "gate misreads a file" --expected "it reads it"` runs
    Then an issue is created there with the bug label, and the output says "reported to Anchors"

  @RPBUG-B04 @unit-level
  Scenario: A dry run prints the issue and sends nothing
    When `anchors report-bug "gate misreads a file" --expected "it reads it" --dry-run` runs
    Then the output has the title and the body, and no issue is created

  @RPBUG-B05 @unit-level
  Scenario: The minimal case comes from a file or from standard input
    Given a file with a made-up anchors.yaml
    When `anchors report-bug ... --repro <file> --dry-run` runs, and again with `--repro -`
    Then the body's minimal case holds the content each time

  @RPBUG-B06 @unit-level
  Scenario: It works without an anchors.yaml
    Given a directory with no anchors.yaml
    When `anchors report-bug "x" --expected "y" --dry-run` runs
    Then it prints the issue

  @RPBUG-B07 @unit-level
  Scenario: An escalation whose reason names the project is not sent upstream
    Given a project whose root path appears in the reason
    When `anchors escalate --bug --upstream "<reason>"` runs in local mode
    Then nothing is created and the output says "not reported to Anchors"

  @RPBUG-X01 @unit-level
  Scenario: A report that names the project is refused
    Given a project whose workflow.repo is acme/app
    When `anchors report-bug "acme/app breaks" --expected "y"` runs
    Then it is refused naming acme/app, and nothing is sent

  @RPBUG-E01 @unit-level
  Scenario: Without what should happen the report is refused
    When `anchors report-bug "x"` runs
    Then it is refused saying --expected is required, and nothing is sent

  @RPBUG-E02 @unit-level
  Scenario: A refused report leaves the link to file it
    Given the platform refuses to create the issue
    When `anchors report-bug "x" --expected "y"` runs
    Then the output carries the issues/new link and the command fails with "could not report to Anchors"

  @RPBUG-E03 @unit-level
  Scenario: An unreadable minimal case is refused
    When `anchors report-bug "x" --expected "y" --repro missing.txt` runs
    Then it is refused with the read error, and nothing is sent
