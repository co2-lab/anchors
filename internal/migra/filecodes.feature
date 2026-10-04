# language: en
# @anchors
#   code: FLFTD
#   ref: MGFCD
#   updated_at: 2026-10-04
#   layer: feature

@MGFCD
Feature: FileCodes — every file a code of its own, of five characters

  @MGFCD-B01 @unit-level
  Scenario: A four-character code is widened, keeping it as the prefix
    Given the code ARNA of ArenaScreen, with ARNAR taken
    When it is widened
    Then the new code starts with ARNA, has five characters and is not ARNAR

  @MGFCD-B02 @unit-level
  Scenario: The files of one unit get different code names
    Given ArenaScreen.tsx in layer screen, ArenaScreen.feature in layer feature and ArenaScreen.test.tsx in layer test
    When their code names are read
    Then they are "ArenaScreen screen", "ArenaScreen feature" and "ArenaScreen test"

  @MGFCD-B03 @unit-level
  Scenario: A file's code is new and of five characters
    Given the code the file would get is taken
    When the file's code is generated
    Then it has five characters and is not the taken one

  @MGFCD-B04 @unit-level
  Scenario: Only a text file with a comment syntax carries the line
    Given a TypeScript file, a JSON file, and a binary
    When each is asked whether it can carry a code
    Then only the TypeScript file can

  @MGFCD-B05 @unit-level
  Scenario: The line goes below the header's opener, or in a new header
    Given a TypeScript file with a header, a Markdown file with one, a script with a shebang and no header
    When the code is written in each
    Then it sits below the opener in the file's syntax, and the script keeps the shebang first

  @MGFCD-B06 @unit-level
  Scenario: A header carrying its unit's code turns it into a ref beside a code of its own
    Given a feature header and a Markdown header carrying the unit's code, and a file with no header
    When the code is turned into a ref with an own code
    Then each header has its own code followed by the unit's ref, and the file with no header is unchanged

  @MGFCD-B07 @unit-level
  Scenario: The renamed codes are read old to current, a code renamed twice to the last one
    Given a project with no renames file, and then one recording ARNA → ARNAA, LOGI → LOGIN and later ARNAA → ARENA
    When the renames are read
    Then none come from the missing file, and ARNA and ARNAA both resolve to ARENA, LOGI to LOGIN

  @MGFCD-B08 @unit-level
  Scenario: The header's updated_at is set to the day given, and only in the header
    Given a file whose header and body both write an updated_at, and a header without one
    When the date is set
    Then only the header's changes, and the header without one is left as it was
