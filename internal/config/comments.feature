# language: en
# @anchors
#   code: CMFTA
#   ref: CMMRC
#   updated_at: 2026-10-03
#   layer: feature

@CMMRC
Feature: CommentMarkers — which text opens a line comment in each kind of file

  @CMMRC-B01 @unit-level
  Scenario: A known extension answers with its line-comment prefixes
    Given the extension ".php"
    When the reader asks for its line-comment prefixes
    Then it gets the double slash and the hash, in that order

  @CMMRC-B02 @unit-level
  Scenario: An unknown extension answers with no prefix
    Given the extension ".unknown", which the table does not list
    When the reader asks for its line-comment prefixes
    Then it gets no prefix at all

  @CMMRC-B03 @unit-level
  Scenario: The line comment of a path follows its last extension, ignoring case
    Given the paths "web/a.test.ts" and "db/001.SQL"
    When the writer asks which comment to use for each
    Then "web/a.test.ts" gets the double slash and "db/001.SQL" gets the double dash

  @CMMRC-B04 @unit-level
  Scenario: A markup file gets the opening of a block comment
    Given the paths "docs/README.md" and "templates/index.html"
    When the writer asks which comment to use for each
    Then both get the opening of an HTML comment, with no closing

  @CMMRC-B05 @unit-level
  Scenario: A path with no extension or an unknown one gets the hash comment
    Given the paths "Makefile" and "assets/logo.svg"
    When the writer asks which comment to use for each
    Then both get the hash comment

  @CMMRC-B06 @unit-level
  Scenario: An extension with several prefixes gets the first one declared
    Given the path "web/index.php", whose extension accepts the double slash and the hash
    When the writer asks which comment to use
    Then it gets the double slash

  @CMMRC-I01 @unit-level
  Scenario: For every extension of the table, the line comment is the table's first prefix
    Given every extension the table lists
    When the writer is asked about a file with each extension
    Then its answer is always the first prefix the reader lists for that extension

  @CMMRC-X01 @unit-level
  Scenario: The prefix lookup takes the extension as given and does not normalise it
    Given the inputs ".GO", "go" and "main.go"
    When the reader asks for their line-comment prefixes
    Then none of them gets a prefix
