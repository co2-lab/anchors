# language: en
# @anchors
#   ref: MPFRM
#   updated_at: 2026-09-26
#   layer: feature

@MPFRM
Feature: MapFormat — the map's format number decides whether this binary may read it

  @MPFRM-B01 @unit-level
  Scenario: The written format and the oldest readable format are both accepted
    Given a map in the format this binary writes, and another in the oldest format it reads
    When the format of each is checked
    Then both are accepted with no error

  @MPFRM-B02 @unit-level
  Scenario: A map from a newer binary is refused with the upgrade message
    Given a map one format above the one this binary writes
    When its format is checked
    Then it is refused with a message saying a NEWER Anchors wrote it, that data would be lost in silence including the judgment stamps, and how to upgrade

  @MPFRM-B03 @unit-level
  Scenario: A map older than the readable range asks for migration
    Given a map in format 1
    When its format is checked
    Then it is refused with a message naming anchors migrate

  @MPFRM-B04 @unit-level
  Scenario: A map with no version is format 1 and asks for migration
    Given a map whose file carries no version
    When its format is checked
    Then it is refused as needing migration, and the message says format 1 and never format 0

  @MPFRM-B05 @unit-level
  Scenario: The newer-map refusal never names the migration command
    Given one map newer than the binary and one older than the readable range
    When the format of each is checked
    Then the newer map's message does not contain anchors migrate

  @MPFRM-I01 @unit-level
  Scenario: Exactly formats 2 through 4 are readable
    Given every format from 0 to 6
    When each format is checked
    Then formats 2, 3 and 4 are accepted and every other one is refused

  @MPFRM-X01 @unit-level
  Scenario: Format 1 is migrated, not read
    Given a map in format 1, which earlier binaries wrote
    When its format is checked
    Then it is refused instead of being read
