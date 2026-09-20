# language: en
# @anchors
#   ref: RLUEX
#   updated_at: 2026-09-19
#   layer: feature

@RLUEX
Feature: Rule — the identity of a verification inside a gate, and the waiver that names it

  @RLUEX-B01 @unit-level
  Scenario: Building an identifier joins the gate and the rule
    Given a gate name and a rule name
    When the identifier is built
    Then it carries both halves, because short rule names would collide between gates

  @RLUEX-B02 @unit-level
  Scenario: A gate with a single verification gains no separator
    Given a gate name and an empty rule name
    When the identifier is built
    Then it is the gate name alone, because there is no second half to distinguish

  @RLUEX-B03 @unit-level
  Scenario: The identifier decomposes into gate and rule
    Given an identifier carrying both halves
    When each half is asked for
    Then the gate half and the rule half come back separately

  @RLUEX-B04 @unit-level
  Scenario: The rule half is empty when the identifier carries only a gate
    Given an identifier that is just a gate name
    When the rule half is asked for
    Then it is empty, because the verdict belongs to the gate as a whole

  @RLUEX-B05 @unit-level
  Scenario: A waiver with no reason is refused
    Given the textual form of a waiver with no reason, and one whose reason is blank
    When both are parsed
    Then both are refused, because the reason is the only thing separating a deliberate
      waiver from an ignored gate

  @RLUEX-B06 @unit-level
  Scenario: A waiver with no rule name is refused
    Given a textual entry that carries a reason and no rule before it
    When it is parsed
    Then it is refused, because there is nothing to waive

  @RLUEX-B07 @unit-level
  Scenario: A path is refused as the waiver target
    Given waiver entries whose target is a file path, a folder pattern and a file name
    When they are parsed
    Then each is refused, and the error names the artifact code as what to use instead

  @RLUEX-B08 @unit-level
  Scenario: A target marker with nothing after it is refused
    Given a waiver entry that opens a target and names none
    When it is parsed
    Then it is refused, because accepting it would produce a waiver that waives nothing

  @RLUEX-B09 @unit-level
  Scenario: The waiver accepts both granularities
    Given one waiver naming a single rule and another naming a whole gate
    When each is asked about a rule of that gate
    Then the rule waiver covers only its own rule and the gate waiver covers them all

  @RLUEX-B10 @unit-level
  Scenario: A waiver that declares targets does not hold for the whole gate
    Given a waiver restricted to two artifact codes
    When it is asked without naming any target
    Then it answers not waived, because the gate must still run to confront the others

  @RLUEX-B11 @unit-level
  Scenario: A waiver by target spares the named codes and confronts the rest
    Given a waiver naming two artifact codes with a written reason
    When each code and a third, unnamed one are asked about
    Then the two named are waived with their reason and the third is still confronted

  @RLUEX-B12 @unit-level
  Scenario: A waiver with no declared target holds for every code
    Given a waiver naming a rule with a reason and no target
    When any artifact code is asked about
    Then it is waived, because a freshly declared gate still needs the coarse exit

  @RLUEX-B13 @unit-level
  Scenario: An artifact with no code is not reached by a waiver restricted to targets
    Given a waiver restricted to one artifact code
    When an artifact with no code at all is asked about
    Then it is not waived, because an artifact with no identity is an earlier problem

  @RLUEX-B14 @unit-level
  Scenario: Each target carries its own reason
    Given a commit message waiving the same rule for two codes, for different reasons
    When each code is asked about
    Then each brings back its own reason, and neither overwrites the other

  @RLUEX-B15 @unit-level
  Scenario: The commit message declares waivers that survive in the history
    Given a commit message carrying two waiver markers with codes and reasons
    When the message is read
    Then the named codes are waived and an unnamed code is still confronted

  @RLUEX-B16 @unit-level
  Scenario: A commit marker whose reason is blank is refused
    Given a commit message whose marker ends with nothing after the colon
    When the message is read
    Then it is refused, by the same guarantee the textual form holds

  @RLUEX-B17 @unit-level
  Scenario: Two waivers merge instead of forcing a choice
    Given one waiver from the commit message and another from the environment
    When they are merged
    Then the result honours both, because a pipeline hook and an author can coexist

  @RLUEX-I01 @unit-level
  Scenario: A waiver never reaches what nobody waived
    Given a waiver naming one gate, and an empty waiver
    When a different gate is asked about
    Then neither waives it, because a waiver leaking onto an unnamed gate is the error
      that would cost most here

  @RLUEX-I02 @unit-level
  Scenario: Every accepted waiver carries a written reason
    Given the textual form and the commit-message form, each with a blank reason
    When both are read
    Then both are refused, because the report would otherwise show a waiver without
      saying why

  @RLUEX-I03 @unit-level
  Scenario: A refusal always reaches the caller as an error
    Given several malformed waiver entries
    When they are parsed
    Then an error is collected for each, because a waiver that does not waive would fail
      the next commit with no visible explanation

  @RLUEX-X01 @unit-level
  Scenario: The unit does not accept a path as the waiver target
    Given a waiver whose target is a path a reviewer would consider obvious
    When it is parsed
    Then it is refused, because a path changes when somebody reorganises folders and the
      waiver would stop holding in silence

  @RLUEX-X02 @unit-level
  Scenario: The unit does not decide whether a rule passes
    Given a rule identifier that no waiver mentions
    When the waiver is asked about it
    Then it answers only not waived, and no verdict, because the verdict belongs to the
      checker

  @RLUEX-X03 @unit-level
  Scenario: The unit reads neither files nor the map
    Given two waivers built from the same text in different working trees
    When each is asked the same question
    Then both answer the same, because a unit that looked at disk would make one waiver
      mean two things
