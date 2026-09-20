# language: en
# @anchors
#   ref: SCLTR
#   updated_at: 2026-09-20
#   layer: feature

@SCLTR
Feature: ScenarioLetterDeclared — the letter of a scenario code exists in the vocabulary

  @SCLTR-B01 @unit-level
  Scenario: An artifact that is not a feature leaves without a verdict
    Given a node whose kind is spec, code or test
    When the gate confronts it
    Then it returns Skip, because scenario codes live in features

  @SCLTR-B02 @unit-level
  Scenario: With no declared vocabulary the gate leaves without a verdict
    Given a project that declares no rule types
    When the gate confronts a feature carrying scenario codes
    Then it returns Skip, because with nothing to compare against every letter would be
      either all valid or all invented, and both answers are noise

  @SCLTR-B03 @unit-level
  Scenario: A feature carrying no scenario code leaves without a verdict
    Given a feature whose scenarios carry no code at all
    When the gate confronts it
    Then it returns Skip, because there is nothing to judge

  @SCLTR-B04 @unit-level
  Scenario: Every letter inside the vocabulary passes
    Given a feature whose scenario codes all carry declared letters
    When the gate confronts it
    Then it returns Pass

  @SCLTR-B05 @unit-level
  Scenario: A letter outside the vocabulary is undetermined, not a failure
    Given a feature whose scenario code carries a letter the project never declared
    When the gate confronts it
    Then it returns Pending, because the nature may deserve declaring and that decision
      is not the gate's

  @SCLTR-B06 @unit-level
  Scenario: The verdict names the letters that are outside and the codes carrying them
    Given a feature whose scenario code carries an undeclared letter
    When the gate confronts it
    Then the verdict names that letter and that code, because that is what the reader
      needs in order to choose between declaring and remapping

  @SCLTR-B07 @unit-level
  Scenario: Codes sharing one unknown letter are grouped into a single line
    Given a feature where several scenarios carry the same undeclared letter
    When the gate confronts it
    Then the verdict reports one letter and lists its codes together, because the reader
      needs to know which letters are outside, not to reread the same accusation

  @SCLTR-I01 @unit-level
  Scenario: The scan is over the shape of a code, never over the vocabulary
    Given a feature whose scenario code carries a letter outside the vocabulary
    When the gate confronts it
    Then the code is seen at all, because a scan built from the declared letters would
      discard precisely what the gate exists to find

  @SCLTR-I02 @unit-level
  Scenario: The code-length pattern is read at every call
    Given a project declaring a code length other than the default
    When the gate confronts a feature whose codes use that length
    Then those codes are recognised, because a pattern frozen at load time would ignore
      the project's declaration

  @SCLTR-I03 @unit-level
  Scenario: A valid letter is never named in the verdict
    Given a feature mixing codes with declared letters and codes with invented ones
    When the gate confronts it
    Then only the invented letters appear in the verdict, so the reader does not hunt for
      the real finding

  @SCLTR-X01 @unit-level
  Scenario: The gate does not choose between declaring and remapping
    Given a feature carrying an undeclared letter
    When the gate confronts it
    Then the verdict shows what is outside without prescribing a repair, because both
      repairs are legitimate and the choice needs domain knowledge

  @SCLTR-X02 @unit-level
  Scenario: The tags accompanying a code are not judged
    Given a feature whose scenario carries a declared letter alongside a free-form tag
    When the gate confronts it
    Then it returns Pass, because a tag is free vocabulary by design and charging it
      would turn a precise instrument into a style opinion
