# language: en
# @anchors
#   ref: RCDLR
#   updated_at: 2026-09-26
#   layer: feature

@RCDLR
Feature: RecodeDialect — the project's own surfaces of a code: testID prefixes and file names

  @RCDLR-B01 @unit-level
  Scenario: The lower convention derives the lower-case prefix
    Given the code "TCDTX"
    When its testID prefix is derived with the convention "lower" and with no convention
    Then the prefixes are "tcdtx" and empty

  @RCDLR-B02 @unit-level
  Scenario: Only quoted, hyphenated testID tokens are rewritten
    Given the text holding "tcdt-amount" in double quotes, 'tcdt-row-1' in single quotes, and 'tcdtsomething'
    When the prefix "tcdt" is rewritten to "tctx"
    Then 2 rewrites are counted, both testIDs read "tctx-", and "tcdtsomething" is unchanged

  @RCDLR-B03 @unit-level
  Scenario: An empty or unchanged prefix rewrites nothing
    Given the text holding the testID "tcdt-x"
    When the empty prefix is rewritten to "tctx", and "tcdt" is rewritten to "tcdt"
    Then both count zero rewrites

  @RCDLR-B04 @unit-level
  Scenario: The occurrences of a prefix are counted
    Given the text holding "txdt-a", 'txdt-b' and "outro-c"
    When the prefixes "txdt", "tcdt" and the empty prefix are counted
    Then the counts are 2, 0 and 0

  @RCDLR-B05 @unit-level
  Scenario: The testID attributes are counted whatever their prefix
    Given the text holding testID="zzzz-root", testID = 'a-b' and a "testID" word with no "=" before its value
    When the testID attributes are counted
    Then the count is 2

  @RCDLR-B06 @unit-level
  Scenario: A path matches the file patterns with the exact code
    Given the patterns "**/{{code}}-*.yaml", "**/{{code}}-suite.yaml" and "**/*.{{code}}-*.png"
    When the paths of "TCDTX" flows, a suite, a snapshot, another code's flow, a component and a lower-case "tcdtx" flow are matched
    Then only the flow, the suite and the snapshot of "TCDTX" match

  @RCDLR-B07 @unit-level
  Scenario: The code is renamed in the file name only
    Given the paths "a/b/TCDTX-B01.yaml", "a/b/TCDTX-suite.yaml", "a/Foo.TCDTX-VR-loaded.png" and "a/b/ActionLink.tsx"
    When "TCDTX" is renamed to "TCTXX"
    Then they become "a/b/TCTXX-B01.yaml", "a/b/TCTXX-suite.yaml", "a/Foo.TCTXX-VR-loaded.png" and stay "a/b/ActionLink.tsx"

  @RCDLR-X01 @unit-level
  Scenario: A longer code in a file name is not renamed
    Given the path "a/b/TCDTXX-B01.yaml"
    When "TCDTX" is renamed to "TCTXX"
    Then the path is unchanged
