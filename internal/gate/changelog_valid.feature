# language: en
# @anchors
#   code: CHVLB
#   ref: CHVLA
#   updated_at: 2026-11-18
#   layer: feature

@CHVLA
Feature: ChangelogValid — a product's changelog outputs hold what they claim, whoever wrote them

  @CHVLA-B01 @unit-level
  Scenario: An output whose blocks read as versions, newest first, one per tag, in its format, passes
    Given app tagged app-v1.0.0 and app-v1.1.0, and CHANGELOG.md with an unreleased block, then 1.1.0 and 1.0.0, each with Features
    When the gate runs with a heading of version and date, Features required and 200 characters at most
    Then it passes

  @CHVLA-B02 @unit-level
  Scenario: Blocks out of order, twice, an unreleased one below another, or a marker no version, fail
    Given a file with 1.0.0 before 1.1.0, 1.1.0 twice, an unreleased block below them, and a marker 1.x
    When the gate runs
    Then it fails naming each

  @CHVLA-B03 @unit-level
  Scenario: A tag with no block, a block with no tag, a heading off the format or a date that is not the tag's fail
    Given a file with 1.2.0 that no tag records, 1.1.0 dated 2001-01-01, 1.0.5 headed by one hash, and no block for 1.0.0; and a project with no file at all
    When the gate runs
    Then it fails naming each block, and with no file says two tags need a block

  @CHVLA-B04 @unit-level
  Scenario: A required section missing, a section left empty, or a block above max_chars fails
    Given 1.1.0 with an empty Fixes section, and 1.0.0 with a long Fixes section and no Features
    When the gate runs with Features required and 60 characters at most
    Then it fails naming the empty section, the missing one and the length

  @CHVLA-B05 @unit-level
  Scenario: The block of a release that changed against HEAD fails; the unreleased one may change
    Given the committed file, then 1.1.0 and the unreleased block edited
    When the gate runs
    Then it fails naming 1.1.0, and not the unreleased block

  @CHVLA-B06 @unit-level
  Scenario: required: changed fails an output unchanged since the latest release while the sources hold a version not released
    Given the sources at 1.2.0 and the file as app-v1.1.0 left it
    When the gate runs, then with the file edited, then with the sources back at 1.1.0
    Then it fails naming 1.2.0 and app-v1.1.0, then passes twice

  @CHVLA-B07 @unit-level
  Scenario: A per_version output has a file per tag named by format.file, plus an unreleased one at most
    Given notes with v1.1.0.txt, unreleased.md, README and v1.3.0.txt, and then v1.0.0.txt and v1.1.0.txt alone
    When the gate runs with the file v{version}.txt
    Then it fails naming README, the missing 1.0.0 and 1.3.0 with no tag, then passes

  @CHVLA-B08 @unit-level
  Scenario: No product with an output skips the gate
    Given a product with no changelog outputs
    When the gate runs
    Then it is skipped, saying no product declares changelog.outputs
