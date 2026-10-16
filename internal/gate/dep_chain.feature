# language: en
# @anchors
#   code: DCFDP
#   ref: DCGDP
#   layer: feature

@DCGDP
Feature: DependencyChain — every import flagged with the code it uses, and every symbol with who uses it

  @DCGDP-B01 @unit-level
  Scenario: The imports are read by the dialect's pattern and resolved to files of the map
    Given a TS screen importing its tokens with names across lines and an alias, a hook through a project alias, a package, and a waived import
    When its imports are read
    Then there are four statements, the tokens import ends on its fourth line with its names, the alias resolves into the project and the package stays outside

  @DCGDP-B02 @unit-level
  Scenario: dep-declared names each import of a governed file with no flag, and the code it would carry
    Given the screen, whose tokens import has no flag
    When dep-declared runs
    Then it names the tokens import by line, path and code, and neither the package nor the waived import

  @DCGDP-B03 @unit-level
  Scenario: dep-honored names each flag whose code is not the one its import resolves to
    Given the screen, whose hook import is flagged with a wrong code
    When dep-honored runs
    Then it names the wrong code with the right one

  @DCGDP-B04 @unit-level
  Scenario: used-by-declared names each imported symbol whose flag is missing or wrong, and each flag nobody imports
    Given the tokens file, one symbol flagged with another code, one imported with no flag, and one flagged that nobody imports
    When used-by-declared runs
    Then it names the three

  @DCGDP-B05 @unit-level
  Scenario: The fixers write the dependency and used-by flags, correct a wrong code, and remove a stale flag
    Given the screen and the tokens file as above
    When check --fix runs the chain's fixers
    Then the flags are written and corrected, the stale one removed, and the three gates pass

  @DCGDP-B06 @unit-level
  Scenario: A re-export declares the names it lists, and an inline import brings the member it reads
    Given a barrel re-exporting one name, a type list of two and a list spanning lines, and a file reading a type through an inline import
    When check --fix runs the chain's fixers
    Then each re-exported name gets its flag above its list, naming its symbol, the inline import carries its flag, and used-by-declared passes

  @DCGDP-B07 @unit-level
  Scenario: Tests, test support and flows take part in the chain
    Given a test importing the screen and a fixture, and a flow composing a login flow inline and by file
    When dep-declared runs on the test, the flow's imports are read, and check --fix runs the chain's fixers
    Then the test's imports are charged and flagged, the flow's compositions are flagged with its own comment, and the fixture's symbol lists the test

  @DCGDP-B08 @unit-level
  Scenario: A kinded dependency needs a declared kind, and no import
    Given a file flagging `@dep[db]: transactions` and `@dep[queue]: jobs` on its calls, in a project declaring the kind db
    When dep-honored runs
    Then only the queue flag is named, for its undeclared kind

  @DCGDP-B09 @unit-level
  Scenario: An import of types only is flagged as such
    Given a TypeScript file with an unflagged type import, a type import with a plain flag, a runtime import flagged as types, and an inline type import the author flagged as types
    When the fixer runs
    Then the type imports carry `@dep[type]:`, the runtime one `@dep:`, the author's inline one stays, and dep-honored passes

  @DCGDP-B10 @unit-level
  Scenario: Imports by name are flagged per file
    Given an OCaml module that opens a lexer, calls an evaluator by its qualified name, names a module in a comment, a standard-library module and a constructor, and flags a module it does not use
    When the chain is confronted and fixed
    Then the lexer and the evaluator are its imports, the missing flags are asked and the unused one refused, and the fixer writes one flag line per module below the header and removes the unused one

  @DCGDP-B11 @unit-level
  Scenario: The fixer finds a declaration whose name ends in a letter outside ASCII
    Given a file declaring caféForte and then café
    When the declaration of café is looked up
    Then it is the second line, not the first

  @DCGDP-B12 @unit-level
  Scenario: Without import patterns the chain is pending, and a commented import is none
    Given a ts file with an import, one commented out on its own line and one after code, in a project with the ts example and in one declaring no imports
    When dep-declared runs
    Then the commented imports are not read, and with no imports declared the gate is pending naming dialect.imports

  @DCGDP-B13 @unit-level
  Scenario: A Go import block is read import by import
    Given a Go file with a single import, a block of two and a slice of strings, in a project with the go example
    When its imports are read
    Then there are three imports, each on its own line, and the slice is none

  @DCGDP-B14 @unit-level
  Scenario: An import path divided by another separator resolves by it
    Given a Python module importing app.core.money through an alias, .helpers beside it, ..core one level up, and requests
    When its imports are read with the python example and the separator
    Then each resolves to its file, the package to its __init__.py, and requests to nothing

  @DCGDP-B15 @unit-level
  Scenario: The symbol of a listed used-by flag is written inside the comment
    Given a listed symbol and the line comment, the OCaml block comment and the markup comment
    When the used-by flag is written
    Then the symbol stands before each closer, and an unlisted symbol is not named

