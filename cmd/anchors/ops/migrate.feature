# language: en
# @anchors
#   code: MGFTM
#   ref: MGCMM
#   updated_at: 2026-10-03
#   layer: feature

@MGCMM
Feature: MigrateCommand — the command the format error promises, bringing the map and the config up to this binary's format

  @MGCMM-B01 @unit-level
  Scenario: Both the map and the config reach the current format
    Given a map and an anchors.yaml both in format 1
    When migrate runs
    Then both are reported as migrated from format 1 to the current format
    And anchors.yaml carries the renamed key and the current version

  @MGCMM-B02 @unit-level
  Scenario: A missing file is reported and the other still migrates
    Given a project with a map in format 1 and no anchors.yaml
    When migrate runs
    Then the output has a line "· anchors.yaml:" with the reason
    And the map is migrated

  @MGCMM-B03 @unit-level
  Scenario: A file already current is reported as such
    Given a map already migrated
    When migrate runs again
    Then it reports "anchors.graph.yaml is already in format" the current format

  @MGCMM-B04 @unit-level
  Scenario: The renamed keys are listed in alphabetical order with their counts
    Given a map in format 1 with julgamentos, gerado_por and code_declarado in that order
    When migrate runs
    Then the keys are listed as code_declarado, gerado_por, julgamentos, each with "(1 occurrence(s))"

  @MGCMM-B05 @unit-level
  Scenario: A dry run reports without writing
    Given a map in format 1
    When migrate runs with --dry-run
    Then it says the map "would be migrated" and lists "gerado_por  (1 occurrence(s))"
    And it says "nothing was written"
    And the map on disk is unchanged

  @MGCMM-B06 @unit-level
  Scenario: A real migration asks for the commit
    Given a map in format 1
    When migrate runs
    Then the output tells the user to commit the migration

  @MGCMM-I01 @unit-level
  Scenario: A second run changes nothing
    Given a map migrated once
    When migrate runs again
    Then it asks for no commit
    And the map is byte for byte what the first run wrote

  @MGCMM-X01 @unit-level
  Scenario: The keys renamed are the migration package's
    Given a map in format 1 holding gerado_por
    When migrate runs
    Then the map holds generated_by, the name the migration step declares

  @MGCMM-X02 @unit-level
  Scenario: The command reminds of the commit and does not make it
    Given a map in format 1
    When migrate runs
    Then the output asks the user to commit the migration

  @MGCMM-B07 @unit-level
  Scenario: Crossing format 5 rewrites the letters of plans, flows and actions
    Given a format 4 project with a plan, a flow, an action, and a spec citing a phase beside a permission and a revision of its own
    When migrate runs with --dry-run, then for real, then again
    Then the dry run lists the rewrites and writes nothing, the real run rewrites the phase, step and result codes wherever cited and leaves the spec's own codes, and the second run rewrites nothing

  @MGCMM-B08 @unit-level
  Scenario: Crossing format 7 widens the four-character codes and renames the files named by them
    Given a format 6 project whose spec owns a four-character code cited by its test and naming a baseline file
    When migrate runs
    Then the code becomes five characters with the old one as prefix in every file, the baseline file is renamed, and the config reads code_lengths [5]

  @MGCMM-B09 @unit-level
  Scenario: Crossing format 7 gives every governed file a code of its own
    Given a format 6 project whose code and test files only ref their unit
    When migrate runs
    Then each of them gets a code of its own in its header and keeps its ref

  @MGCMM-B10 @unit-level
  Scenario: Crossing format 7 carries what each file was measured at to its new revision
    Given a format 6 map with every file measured at its current content
    When migrate runs
    Then the map has each rewritten file at its new revision

  @MGCMM-B11 @unit-level
  Scenario: Crossing format 7 records each renamed code in anchors.renames.yaml
    Given a format 6 project with a four-character code
    When migrate runs, then again
    Then anchors.renames.yaml records the old code and the new one, and the second run changes nothing

  @MGCMM-B12 @unit-level
  Scenario: Crossing format 7 turns a file carrying its spec's code into a ref with a code of its own
    Given a format 6 project whose feature carries its spec's code as its own
    When migrate runs
    Then the feature refs the unit's widened code and gets a code of its own

  @MGCMM-E02 @unit-level
  Scenario: Crossing format 7, a file that cannot be written fails the command, and a second run finishes it
    Given a format 6 project with a governed file that cannot be written
    When migrate runs, then runs again after the file is made writable
    Then the first run fails with the write error, and the second gives the file its code
