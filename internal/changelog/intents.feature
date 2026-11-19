# language: en
# @anchors
#   code: INTNB
#   ref: INTNA
#   updated_at: 2026-11-18
#   layer: feature

@INTNA
Feature: Intents — the change intents, read from their files and rendered as a release

  @INTNA-B01 @unit-level
  Scenario: An intent reads its product, bump and type from the front matter, and its text after it
    Given an intent for app, bump minor, type feat, with CRLF endings
    When it is parsed
    Then it holds app, minor, feat and its text

  @INTNA-E01 @unit-level
  Scenario: An intent with no front matter, one not closed, or a key it does not know is an error naming the file
    Given intents with no front matter, a front matter not closed, a key scope, and YAML that does not parse
    When they are parsed
    Then each fails naming the file

  @INTNA-B02 @unit-level
  Scenario: A directory's intents are read sorted by name, each one that does not parse is an error of its own, and no directory holds none
    Given b.md, a.md, c.md with no front matter, a text file and an archived intent under released
    When the directory is read, and one that does not exist
    Then a.md and b.md are read in order, c.md is one error, and the missing directory holds nothing

  @INTNA-B03 @unit-level
  Scenario: Intents are grouped by type in the order declared, another type after them, and rendered with a heading per type
    Given a fix, a perf and a feat, with feat, fix and docs declared
    When they are grouped and rendered as the candidate 1.1.0
    Then the sections are feat, fix and perf, under 1.1.0 — Unreleased, and a template that does not parse fails
