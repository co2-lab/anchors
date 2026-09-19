# language: en
# @anchors
#   ref: DMDCD
#   updated_at: 2026-09-19
#   layer: feature

@DMDCD
Feature: DomainDeclared — the spec declares what the unit ACCEPTS, and who rejects the invalid

  @DMDCD-B01 @unit-level
  Scenario: An artifact that is not a spec leaves without a verdict
    Given a node whose kind is code, test or feature
    When the gate confronts it
    Then it returns Skip, because only a spec has a domain to declare

  @DMDCD-B02 @unit-level
  Scenario: A spec without the domain section is failed
    Given a spec that catalogues rules and never opens the domain section
    And no declared waiver anywhere in it
    When the gate confronts it
    Then it returns Fail
    And the verdict names the waiver marker, so whoever reads it learns the declared way out

  @DMDCD-B03 @unit-level
  Scenario: A section opened and left empty is failed
    Given a spec whose domain section has only the table header, with no data row
    When the gate confronts it
    Then it returns Fail, because opening the title without declaring anything is the same
      gap wearing the appearance of compliance

  @DMDCD-B06 @unit-level
  Scenario: A row filled only with placeholders is not a declaration
    Given a spec whose domain section carries a single row reading "TODO" in every column
    When the gate confronts it
    Then it returns Fail, because the untouched template asserts nothing

  @DMDCD-B04 @unit-level
  Scenario: An entry with no owner is failed, and the verdict names it
    Given a spec declaring the entry "chave" with its accepted values
    And the column that says who guarantees it is left blank
    When the gate confronts it
    Then it returns Fail
    And the verdict names "chave", so the reader does not have to hunt for the orphan

  @DMDCD-B07 @unit-level
  Scenario: An entry whose owner is named passes
    Given a spec declaring the entry "chave" with its accepted values
    And the column that says who guarantees it reads "the interface, before calling"
    When the gate confronts it
    Then it returns Pass

  @DMDCD-B05 @unit-level
  Scenario: A waiver with a written reason silences the gate
    Given a spec with no domain section
    And a waiver marker followed by the reason "receives only typed values from its own code"
    When the gate confronts it
    Then it returns Skip, and the reason stays in the spec as the record that someone looked

  @DMDCD-I01 @unit-level
  Scenario: A bare waiver, with no reason, does not waive
    Given a spec with no domain section
    And a waiver marker with nothing written after it
    When the gate confronts it
    Then it returns Fail, because a waiver with no why is the silence the gate exists to end

  @DMDCD-I02 @unit-level
  Scenario: A non-answer in the owner column is not an owner
    Given a spec declaring the entry "chave"
    And the column that says who guarantees it reads "I do not validate (MTVRX-X04)"
    When the gate confronts it
    Then it returns Fail, because carrying a restriction into the owner column names nobody —
      it is the sentence that creates the orphan

  @DMDCD-X01 @unit-level
  Scenario: The gate does not judge whether the declared domain is correct
    Given a spec declaring the entry "month" as accepting any text, which is wider than the real domain
    And the column that says who guarantees it names the caller
    When the gate confronts it
    Then it returns Pass, because the ruler here is the PRESENCE of the declaration —
      whether the accepted set matches reality is judgment, and judgment belongs to another gate

  @DMDCD-X02 @unit-level
  Scenario: The gate does not read the code to check the validation exists
    Given a spec whose domain section is complete and every entry has a named owner
    And the code that realises it performs no validation at all
    When the gate confronts it
    Then it returns Pass, because this layer reads TEXT — crossing the declaration with the
      implementation belongs to the relational gate, which has the map
