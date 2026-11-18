# language: en
# @anchors
#   code: CMFTA
#   ref: CMMRC
#   updated_at: 2026-10-16
#   layer: feature

@CMMRC
Feature: Comments — how the project's files write a comment, declared in the dialect

  @CMMRC-B01 @unit-level
  Scenario: A file's syntax is the first entry that covers it, or Anchors' own for its formats
    Given a dialect declaring the ts and python examples
    When the syntax of a TS file, a Python file at the root, a spec, a feature and a JSON file is asked
    Then each gets its entry, the spec and the feature Anchors' own, and the JSON none

  @CMMRC-B02 @unit-level
  Scenario: Blanking keeps strings and lines, and turns comments into spaces
    Given TS code with a line comment, a block across lines, a marker inside a string, an escaped quote and a template across lines
    When it is blanked
    Then only the comments are spaces, every line and offset kept, and with no syntax nothing changes

  @CMMRC-B03 @unit-level
  Scenario: Lines that are only a comment, the code of a content, a line cut
    Given a content with a block comment's inner lines, a line of code with a comment at its end, and a blank line
    When its comment-only lines, its code and a line cut are read
    Then the block's lines are only a comment, the code drops them and the end comment, and the blank line stays

  @CMMRC-B04 @unit-level
  Scenario: The writer's marker and where a comment cannot go
    Given the ts example, OCaml's entry and Markdown's
    When the writer's marker, its closer and a line after a JSX tag are asked
    Then ts writes `//` and takes no comment after `>`, OCaml writes `(*` closed by ` *)`, Markdown `<!--` by ` -->`

  @CMMRC-B05 @unit-level
  Scenario: A malformed entry refuses the config
    Given entries with no marker, an empty marker, an invalid glob, a broken no_line_after and a block of one delimiter
    When the config is loaded
    Then each refuses it, naming the entry

  @CMMRC-B06 @unit-level
  Scenario: The families have an example of their comments, never applied by itself
    Given the ts, python and ocaml examples and every family's together
    When they read their files and a project that declared none is asked
    Then ts reads JSX's comment and its template literals, python its triple-quoted strings, ocaml its blocks and dune's files, and the project has no syntax for its code

  @CMMRC-I01 @unit-level
  Scenario: Blanking never moves a line
    Given contents of every example's language
    When each is blanked
    Then its length and the positions of its line breaks are the original's

  @CMMRC-X01 @unit-level
  Scenario: No language's comment is built in
    Given a project with no dialect
    When the syntax of a Go, a Python and an OCaml file is asked
    Then none has one

  @CMMRC-B07 @unit-level
  Scenario: The examples' entries covering a set of files
    Given the files a.go, web/b.tsx, run.sh, c.zig, Makefile and a.spec.md
    When the entries covering them are asked for
    Then the TS, Go and shell entries come in the examples' order, and .zig and Makefile are named
