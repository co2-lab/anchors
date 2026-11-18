# language: en
# @anchors
#   code: SYMXF
#   ref: SYMXS
#   updated_at: 2026-11-12
#   layer: feature

@SYMXS
Feature: SymbolExists — a name marked as a symbol is in the project's code

  @SYMXS-B01 @unit-level
  Scenario: A marked name the code no longer has fails, naming where it is cited
    Given a code comment marking `#isEvalSupported` and `$cfg.readerOptions`, and a project whose code has readerOptions and cfg but not isEvalSupported
    When the gate confronts the file
    Then it fails naming the file, the line and isEvalSupported, and passes once the mark names a name the code has

  @SYMXS-B02 @unit-level
  Scenario: Only real marks in comments and in the spec's prose are confronted
    Given a code file whose code holds a marked-looking string, whose comments hold a link anchor, a path fragment, a CSS id, a folding directive, a color and a comment opened by `#`, and a spec whose fence holds a shell variable
    When the gate confronts each
    Then none of them is confronted, and a `@@name` in the spec's prose is

  @SYMXS-B03 @unit-level
  Scenario: A line that cites a name on purpose is waived
    Given a comment citing a removed option with `@no-symbol-check:` and a reason on its line
    When the gate confronts the file
    Then the name is not confronted

  @SYMXS-B04 @unit-level
  Scenario: No mark, or another kind of artifact, is skipped
    Given a code file whose comments mark nothing, and a feature
    When the gate confronts each
    Then each is skipped

  @SYMXS-E01 @unit-level
  Scenario: A code file whose comments the project does not declare is pending
    Given a code file of an extension no comment syntax matches
    When the gate confronts it
    Then it is pending, naming the file

  @SYMXS-E02 @unit-level
  Scenario: Without a map the gate is pending
    Given a code comment with a mark and no built map
    When the gate confronts it
    Then it is pending
