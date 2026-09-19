# language: en
# @anchors
#   ref: LYBNL
#   updated_at: 2026-09-19
#   layer: feature

@LYBNL
Feature: LayerBoundary — a layer does not reach what is not its own

  @LYBNL-B01 @unit-level
  Scenario: An artifact that is not code leaves without a verdict
    Given a node whose kind is spec, test or feature
    When the gate confronts it
    Then it returns Skip, because there is no import to forbid outside code

  @LYBNL-B02 @unit-level
  Scenario: Content matching a forbidden pattern fails, naming line and reason
    Given a boundary forbidding the screens layer to import from the repositories
    And a screen whose second line imports straight from a repository
    When the gate confronts it
    Then it returns Fail naming line 2 and the declared reason, because a prohibition
      with no motive turns into ritual

  @LYBNL-B03 @unit-level
  Scenario: A rule scoped to a layer charges only that layer
    Given the same boundary declared for the screens layer
    And a hook that imports from a repository
    When the gate confronts it
    Then it does not fail, because the hook is exactly who is allowed to reach the data

  @LYBNL-B04 @unit-level
  Scenario: A rule with no layer holds for all code
    Given a boundary with no layer forbidding the raw clock
    And files of the screens, hooks and repositories layers each reading the raw clock
    When the gate confronts each of them
    Then every one fails, because a rule without a layer is how a global prohibition
      is declared

  @LYBNL-B05 @unit-level
  Scenario: Severity warn records without failing, and the default is error
    Given a boundary marked severity warn and a file that violates it
    When the gate confronts it
    Then it returns Pending, and the same boundary with no severity returns Fail,
      because the default is error and warn is the per-rule maturation

  @LYBNL-B06 @unit-level
  Scenario: A waiver with a written reason on the line waives that line
    Given a forbidden import carrying an allow-boundary marker with a written reason
    When the gate confronts it
    Then it returns Pass, and the acknowledged debt stays visible and dated in the code

  @LYBNL-B07 @unit-level
  Scenario: The waiver also holds in the comment on the line above
    Given a forbidden import whose allow-boundary marker with reason sits on the line above
    When the gate confronts it
    Then it returns Pass, because an import has nowhere to carry a readable end-of-line
      comment and demanding it inline would push the author not to declare at all

  @LYBNL-B08 @unit-level
  Scenario: A bare marker with no reason does not waive
    Given a forbidden import carrying an allow-boundary marker with nothing written after it
    When the gate confronts it
    Then it returns Fail, because a bare marker is a silent way to quiet the gate

  @LYBNL-B09 @unit-level
  Scenario: With no boundary declared the verdict is Pending, never Pass
    Given a project that declares no boundary at all
    When the gate confronts a code file
    Then it returns Pending naming what to declare, because pretending it checked is
      worse than saying what is missing

  @LYBNL-B10 @unit-level
  Scenario: An invalid forbid pattern fails visibly
    Given a boundary whose forbid pattern does not compile
    When the gate confronts a code file
    Then it returns Fail explaining the config problem, because swallowed in silence it
      would switch the rule off with nobody knowing

  @LYBNL-B11 @unit-level
  Scenario: The pattern is matched against the whole file, catching a multi-line import
    Given a boundary forbidding a named import from a package
    And a file whose import of that name is wrapped over several lines
    When the gate confronts it
    Then it returns Fail pointing at the line where the match starts

  @LYBNL-I01 @unit-level
  Scenario: The same rule is expressible in six language dialects
    Given the same architectural rule written in the import dialect of TypeScript,
      Python, Go, Java, Rust and Ruby
    When the gate confronts a violation and a legitimate import in each dialect
    Then the violation fails and the legitimate import passes in all six, because the one
      who writes the pattern is the project and the engine knows no language

  @LYBNL-I02 @unit-level
  Scenario: A single-line import of the same shape is still caught
    Given a boundary forbidding a named import from a package
    And a file whose import of that name fits on one line
    When the gate confronts it
    Then it returns Fail, because accusing only the wrapped form would accuse formatting
      rather than the violation

  @LYBNL-I03 @unit-level
  Scenario: The waiver holds on any line of the matched stretch
    Given a multi-line import whose allow-boundary marker sits on the from line
    When the gate confronts it
    Then it returns Pass, because demanding the marker on the first line of the match
      would require the author to know where the regex started matching

  @LYBNL-I04 @unit-level
  Scenario: A line anchor keeps holding per line
    Given a boundary whose pattern anchors the forbidden text to a whole line
    And a file carrying that text inside a string in the middle of a line
    When the gate confronts it
    Then it returns Pass, because the whole-file match changes the dot, not the meaning
      of the anchors

  @LYBNL-X01 @unit-level
  Scenario: The gate does not decide which boundaries exist
    Given a project whose Structure declares no boundary
    And a screen importing straight from a repository, which a reviewer would forbid
    When the gate confronts it
    Then it does not fail, because inventing boundaries would charge what nobody
      committed to

  @LYBNL-X02 @unit-level
  Scenario: The gate does not parse the language, it matches text
    Given a boundary whose pattern is plain text with no notion of imports
    And a file where the forbidden text appears outside any import statement
    When the gate confronts it
    Then it returns Fail, because the ruler is TEXT — understanding the import graph of
      every language would tie the engine to a set of ecosystems

  @LYBNL-X03 @unit-level
  Scenario: The gate does not judge whether the boundary is the right one to draw
    Given a boundary forbidding something a reviewer would consider harmless
    And a file that matches it
    When the gate confronts it
    Then it returns Fail, because the ruler is the DECLARATION — whether the boundary is
      worth drawing is design judgment, and that belongs to whoever writes the Structure
