# language: en
# @anchors
#   ref: RLTYR
#   updated_at: 2026-09-19
#   layer: feature

@RLTYR
Feature: RuleTypes — the rule vocabulary is extensible, but it must be declared

  @RLTYR-B01 @unit-level
  Scenario: A letter that is not declared in the vocabulary fails
    Given a vocabulary declaring the letters S, B and E
    And a spec cataloguing a rule under the letter P
    When the gate confronts it
    Then it returns Fail naming the letter P, because a letter the traceability cannot
      see makes the rule look covered when it is not

  @RLTYR-B02 @unit-level
  Scenario: A declared letter under a claimed section passes
    Given a vocabulary declaring the letters S, B and E with their sections
    And a spec cataloguing rules only under those letters and sections
    When the gate confronts it
    Then it returns Pass

  @RLTYR-B03 @unit-level
  Scenario: A section cataloguing rules under a title no letter claims fails
    Given a vocabulary whose declared sections do not include "Regras Inventadas"
    And a spec cataloguing a declared letter under that title
    When the gate confronts it
    Then it returns Fail naming the section, because the letter is claimed by a title and
      the catalogue has to sit where the vocabulary says it does

  @RLTYR-B04 @unit-level
  Scenario: The same letter claimed by two terms is a conflict in the vocabulary
    Given a vocabulary where the letter E is claimed by the terms Error and Estado
    When the gate confronts any spec
    Then it returns Fail signalling the CONFLICT, because each letter belongs to ONE term

  @RLTYR-B05 @unit-level
  Scenario: With no vocabulary declared the gate confronts the canonical letters
    Given a project that declares no vocabulary
    And a spec cataloguing a rule under the letter P, which is outside the canonical set
    When the gate confronts it
    Then it returns Fail saying it confronts the canonical vocabulary, because a gate that
      Skips forever gives the impression of a defence that does not exist

  @RLTYR-B06 @unit-level
  Scenario: A heading that is the rule code itself is not a category section
    Given a spec whose heading is the rule code followed by its title
    When the gate confronts it
    Then it returns Pass, because that heading is the rule's own header and not a
      category that has to claim a letter

  @RLTYR-B07 @unit-level
  Scenario: A section that only cites other sections' codes claims no letter
    Given a spec whose test-id section references codes of other sections in an inner column
    And that section title is claimed by no letter
    When the gate confronts it
    Then it returns Pass, because citing is not cataloguing

  @RLTYR-B08 @unit-level
  Scenario: A section that defines a code in the first table cell is charged
    Given a spec whose unclaimed section carries a table row opening with a rule code
    When the gate confronts it
    Then it returns Fail, because the first cell is where a definition lives

  @RLTYR-B09 @unit-level
  Scenario: A section declared as rule-cataloguing and filled without a code is Pending
    Given a vocabulary declaring "Eventos / Callbacks" as requiring a code
    And a spec whose section of that name carries a filled table and no code at all
    When the gate confronts it
    Then it returns Pending naming the section, because a row that asserts something
      verifiable and carries no code leaves the scenario with nothing to cite

  @RLTYR-B10 @unit-level
  Scenario: A section whose table already carries the code is not charged
    Given a vocabulary declaring "Eventos / Callbacks" as requiring a code
    And a spec whose section of that name carries the rule code in its table
    When the gate confronts it
    Then it does not return Pending, because there is nothing left to charge

  @RLTYR-B11 @unit-level
  Scenario: A declared section outside sections_require_code is not charged
    Given a vocabulary declaring "Variantes" under a letter but not as requiring a code
    And a spec whose section of that name merely enumerates values
    When the gate confronts it
    Then it does not return Pending, because demanding a rule of an index would invent a duty

  @RLTYR-B12 @unit-level
  Scenario: A project that does not use sections_require_code changes no behaviour
    Given a vocabulary declaring "Eventos / Callbacks" with no requires-code marking
    And a spec whose section of that name carries a filled table and no code
    When the gate confronts it
    Then it does not return Pending, because the ruler is born opt-in and would otherwise
      accuse an entire existing base at once

  @RLTYR-I01 @unit-level
  Scenario: A spec with no rule code at all is not this gate's problem
    Given a project with no vocabulary declared
    And a spec of pure prose that catalogues no rule
    When the gate confronts it
    Then it returns Pass, because charging the existence of a catalogued rule belongs to
      the spec-complete gate and doing it here would duplicate the ruler

  @RLTYR-I02 @unit-level
  Scenario: The verdict names the letter and where to declare it
    Given a project with no vocabulary declared
    And a spec using a letter outside the canonical set
    When the gate confronts it
    Then the verdict carries both the letter and the vocabulary key to declare it in,
      because failing without saying what transfers the diagnosis to whoever reads it

  @RLTYR-I03 @unit-level
  Scenario: A filled section with no code is the gap where the scenario loses its anchor
    Given a vocabulary declaring an events section as requiring a code
    And a spec whose events section is filled and carries no code
    When the gate confronts it
    Then it reports the finding, because without a code the scenario borrows the
      neighbouring section's and starts governing what is not its own

  @RLTYR-X01 @unit-level
  Scenario: The gate does not decide which letters exist
    Given a vocabulary declaring a letter the canonical set does not contain
    And a spec cataloguing rules under that letter and its claimed section
    When the gate confronts it
    Then it returns Pass, because the vocabulary is extensible by design and the canonical
      letters are the fallback, not a ceiling

  @RLTYR-X02 @unit-level
  Scenario: Without a declared vocabulary only the letter is charged
    Given a project that declares no vocabulary
    And a spec cataloguing a canonical letter under a title no vocabulary claims
    When the gate confronts it
    Then it returns Pass, because sections and terms only exist once the project declares
      them, and charging them against an implicit vocabulary would invent a rule nobody wrote

  @RLTYR-X03 @unit-level
  Scenario: The gate does not judge whether the letter suits the rule
    Given a vocabulary declaring both S for state and B for behaviour
    And a spec cataloguing a plainly behavioural rule under the state letter and section
    When the gate confronts it
    Then it returns Pass, because which letter a rule deserves is editorial judgment —
      the ruler here is that the traceability can see it

  @RLTYR-X04 @unit-level
  Scenario: The gate charges traceability, not format
    Given a spec whose unclaimed section carries prose and no rule code at all
    When the gate confronts it
    Then it returns Pass, because a section that catalogues nothing has no traceability
      to defend, whatever its shape
