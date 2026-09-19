# language: en
# @anchors
#   ref: CDCTC
#   updated_at: 2026-09-19
#   layer: feature

@CDCTC
Feature: CodeCataloged — what the code exports must be in the spec, or waived in the code

  @CDCTC-B01 @unit-level
  Scenario: An artifact that is not a spec leaves without a verdict
    Given a node whose kind is code, test or feature
    When the gate confronts it
    Then it returns Skip, because the ruler starts from the spec that governs the code

  @CDCTC-B02 @unit-level
  Scenario: An exported symbol the spec never names fails, and the verdict names it
    Given a spec cataloguing one rule and a code file exporting three functions
    When the gate confronts it
    Then it returns Fail naming each orphan and its line, because the name alone would
      make the reader hunt for the symbol in the file

  @CDCTC-B03 @unit-level
  Scenario: What the spec already catalogues is never accused
    Given a spec cataloguing one of the exported functions by name
    When the gate confronts it
    Then the verdict does not name that function, because it is already catalogued

  @CDCTC-B04 @unit-level
  Scenario: A no-rule marker with a written reason waives the symbol
    Given an exported function carrying a no-rule marker with a written reason
    And the spec that catalogues every other symbol
    When the gate confronts it
    Then it returns Pass, because not every export deserves a rule and the waiver is what
      makes the noise manageable without lying

  @CDCTC-B05 @unit-level
  Scenario: A bare no-rule marker does not waive
    Given an exported function carrying a no-rule marker with nothing written after it
    When the gate confronts it
    Then it returns Fail, because a bare marker would be a silent way to quiet the gate
      and the trace that a decision was taken would vanish

  @CDCTC-B06 @unit-level
  Scenario: A spec cataloguing every exported symbol passes
    Given a code file exporting two constants and a spec naming both
    When the gate confronts it
    Then it returns Pass

  @CDCTC-B07 @unit-level
  Scenario: With no code linked the gate leaves without a verdict
    Given a spec with no code file linked to it in the map
    When the gate confronts it
    Then it returns Skip, because the absence belongs to the triad gate and accusing it
      in both places would duplicate the debt

  @CDCTC-B08 @unit-level
  Scenario: Without a declared export pattern the gate skips and says so
    Given a Go file and a project that declared no export pattern
    When the gate confronts the spec that governs it
    Then it does not return Pass, and the verdict names how to enable the pattern,
      because green over what was never read is worse than an honest red

  @CDCTC-B09 @unit-level
  Scenario: With the pattern declared the gate confronts for real in any language
    Given a project declaring an export pattern for Go
    And a Go file exporting a function the spec never names
    When the gate confronts the spec
    Then it returns Fail naming that function

  @CDCTC-B10 @unit-level
  Scenario: The declared dialect family also supplies the pattern
    Given a project that names its dialect family as Go and declares no pattern of its own
    And a Go file exporting a function the spec never names
    When the gate confronts the spec
    Then it returns Fail naming that function, because naming the family is enough

  @CDCTC-I01 @unit-level
  Scenario: The waiver holds in the comment block above the symbol
    Given the no-rule declaration written inline, one line above, in a two-line comment
      block, in a block with paragraphs and in a doc comment
    When the gate reads the context of the symbol in each case
    Then the declaration holds in all of them, because it is documentation and the
      explanation rarely fits on one line

  @CDCTC-I02 @unit-level
  Scenario: The waiver does not leak between symbols
    Given two exported functions where only the first carries a no-rule declaration
    When the gate reads the context of each symbol
    Then only the first carries the declaration, because inheriting would let one marker
      exempt the whole file, which is the opposite of what it is

  @CDCTC-I03 @unit-level
  Scenario: The gate never approves a language it cannot read
    Given a Go file whose exports match no TypeScript syntax
    And a project that declared no export pattern
    When the gate confronts the spec
    Then it does not return Pass, because stamping approval over what was never read is
      the worst possible failure in a measuring instrument

  @CDCTC-X01 @unit-level
  Scenario: The gate does not judge whether the rule describes the symbol well
    Given a spec whose rule names the exported function and describes it wrongly
    When the gate confronts it
    Then it returns Pass, because the ruler is whether the spec NAMES the symbol —
      judging what the rule says about it belongs to another gate

  @CDCTC-X02 @unit-level
  Scenario: The gate does not decide which symbols deserve a rule
    Given an exported function of pure formatting, which a reviewer would exempt
    And no no-rule declaration anywhere near it
    When the gate confronts the spec
    Then it returns Fail, because the exemption is the project's call and the waiver is
      where it records it — deciding here would remove the calibration

  @CDCTC-X03 @unit-level
  Scenario: The gate knows no language, the project declares what is public
    Given two projects whose export patterns recognise different syntaxes
    And the same file, public under one pattern and invisible under the other
    When the gate confronts each
    Then the verdicts differ, because recognising what is public depends on the language
      and Anchors does not presume

  @CDCTC-X04 @unit-level
  Scenario: The gate does not charge the absence of code
    Given a spec cataloguing rules with no code file linked to it
    When the gate confronts it
    Then it returns Skip rather than Fail, because accusing the same debt in two gates
      would duplicate the finding
