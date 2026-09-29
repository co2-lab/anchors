# language: en
# @anchors
#   ref: MGFLM
#   updated_at: 2026-09-29
#   layer: feature

@MGFLM
Feature: MigrateFile — takes one project file from its declared format to the current one, renaming only what is a key

  @MGFLM-B01 @unit-level
  Scenario: The format is the top-level version line, and its absence means format 1
    Given a map whose only lines are comments and "nodes: []"
    And a map declaring "version: 3"
    When their formats are read
    Then the first is format 1 and the second is format 3

  @MGFLM-B02 @unit-level
  Scenario: A file already at the target is left untouched
    Given a map declaring "version: 2" that still has the key "gerado_por"
    When it is migrated to format 2
    Then the text is unchanged and the result reports no change

  @MGFLM-B03 @unit-level
  Scenario: Old map keys are renamed where they are keys and keep their values
    Given a format 1 map with "gerado_por: 0.1.83", an indented "code_declarado: true" and a list item's "julgamentos:" holding a "doc-self-contained" judgment
    When it is migrated to format 2
    Then it has "generated_by:", "code_declared:", "judgments:" and "version: 2" and none of the old keys
    And the values "0.1.83" and "doc-self-contained" are still there
    And the old key names written in a comment or inside a quoted title are unchanged

  @MGFLM-B04 @unit-level
  Scenario: Each rename applies only to its own file
    Given a format 1 configuration with "trinca_opcional: [covered-by, tested-by]"
    And a format 1 configuration with the map's keys "julgamentos" and "gerado_por"
    When both are migrated to format 2
    Then the first has "triad_optional: [covered-by, tested-by]"
    And the second still has "julgamentos" and "gerado_por"

  @MGFLM-B05 @unit-level
  Scenario: A value is renamed only under its key and when it matches whole
    Given a format 1 configuration with "- name: regra-cumprida", "  when: regra-cumprida" and "- name: regra-cumprida-extra"
    When it is migrated to format 2
    Then the first becomes "- name: rule-fulfilled"
    And the other two are unchanged

  @MGFLM-B06 @unit-level
  Scenario: Steps are applied in order, each on the result of the previous
    Given a step producing format 2 that renames "a" to "b" and a step producing 3 that renames "b" to "c"
    And a format 1 configuration with "a: 1"
    When it is migrated to format 3
    Then it has "c: 1" and neither "a:" nor "b:"

  @MGFLM-B07 @unit-level
  Scenario: The version is raised even when nothing else changes
    Given a map with "version: 1" and no old key
    When it is migrated to format 2
    Then it has "version: 2" and the result reports a change

  @MGFLM-B08 @unit-level
  Scenario: A map or a configuration without a version gets one after its comment header
    Given a map starting with two comment lines and no version line, and a configuration "# config" then "project: x"
    When each is migrated to format 2
    Then the map has "version: 2" and the text does not start with it, and the configuration reads "# config", "version: 2", "project: x" and is format 2

  @MGFLM-B09 @unit-level
  Scenario: The result counts each rename by its old form
    Given a format 1 configuration with two "trinca_opcional:" keys and one "- name: regra-cumprida"
    When it is migrated to format 2
    Then the count for "trinca_opcional" is 2 and for "name: regra-cumprida" is 1, and the result reports a change

  @MGFLM-B10 @unit-level
  Scenario: A dry run reports without writing
    Given a format 1 map with "gerado_por: dev"
    When it is migrated to format 2 as a dry run
    Then the result reports a change and counts "gerado_por" once
    And the file on disk is unchanged

  @MGFLM-B11 @unit-level
  Scenario: A version line with a trailing comment is read and raised keeping the comment
    Given a map with "version: 1  # the first", and a configuration with "version: 3 # c"
    When their formats are read and the map is migrated to format 2
    Then they are formats 1 and 3, and the map's line reads "version: 2  # the first"

  @MGFLM-I01 @unit-level
  Scenario: Migrating twice changes nothing the second time
    Given a map with "version: 1"
    When it is migrated to format 2 twice
    Then the second run reports no change and the text is the same as after the first

  @MGFLM-X01 @unit-level
  Scenario: A file the YAML parser would refuse is still migrated
    Given a format 1 map with "gerado_por: dev" and a line with an unclosed bracket
    When it is migrated to format 2
    Then it has "generated_by: dev" and "version: 2" and keeps the broken line

  @MGFLM-E01 @unit-level
  Scenario: A missing file is an error and nothing is written
    Given a path where no file exists
    When it is migrated
    Then the answer is an error and no file is created

  @MGFLM-E02 @unit-level
  Scenario: A version that is not a number is an error and nothing is written
    Given a configuration with "version: 99999999999999999999999", and a map with "version: abc"
    When their formats are read, and the map is migrated to format 2
    Then each answer is an error saying the version is not a number, and the map's text is unchanged

  @MGFLM-E03 @unit-level
  Scenario: A hole in the chain leaves the file untouched
    Given only a step producing format 2 is registered
    And a format 1 map with "gerado_por: dev"
    When it is migrated to format 3
    Then the answer is an error naming format 3 and the file is unchanged

  @MGFLM-B12 @unit-level
  Scenario: A CRLF file is read and migrated as it is
    Given a current map with CRLF line endings, an older one, and one with no version line
    When their format is read and they are migrated
    Then the current one reads its version and is untouched, and the others get one version line ending in CRLF
