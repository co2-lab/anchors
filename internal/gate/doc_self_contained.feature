# language: en
# @anchors
#   ref: DSCDC
#   updated_at: 2026-09-19
#   layer: feature

@DSCDC
Feature: DocSelfContained — the spec has to stand on its own

  @DSCDC-B01 @unit-level
  Scenario: A reference that brings the passage it announces passes
    Given a spec naming the path of a plan node and quoting the passage just below it
    When the gate confronts it
    Then it returns Pass, because the reader has the argument in hand and does not leave the page

  @DSCDC-B02 @unit-level
  Scenario: A reference that only points is accused
    Given a spec whose line names a plan path, or its bare file name, or a revision code,
      and carries nothing of what it announces
    When the gate confronts it
    Then it returns Fail telling the author to bring the text, because the scaffold sends
      the reader to a file they do not have

  @DSCDC-B03 @unit-level
  Scenario: A quotation counts in any written tradition
    Given a spec whose reference is followed by the passage in straight quotes, typographic
      quotes, guillemets, German low quotes or CJK corner brackets
    When the gate confronts each one
    Then every one returns Pass, because demanding Latin quotes would accuse a French or a
      Japanese project unjustly

  @DSCDC-B04 @unit-level
  Scenario: A path inside a code fence is an example, not a reference
    Given a spec showing a command to run inside a fenced block, and that command names a plan path
    When the gate confronts it
    Then it returns Pass, because inside the fence the path is what the reader types, not
      where the reader is being sent

  @DSCDC-B05 @unit-level
  Scenario: The spec citing its own path is identifying itself
    Given a spec whose body names its own file path
    When the gate confronts it
    Then it returns Pass, because it is sending nobody anywhere

  @DSCDC-B06 @unit-level
  Scenario: A rule code is not a revision
    Given a spec whose line names a sibling rule code and explains what that rule proved
    When the gate confronts it
    Then it returns Pass, because the revision form is the one the doctrine reserves and
      accusing the other would charge the spec for naming its own subject

  @DSCDC-B07 @unit-level
  Scenario: Only the spec is charged
    Given the same empty reference written in a plan, a feature, a test and a code node
    When the gate confronts each of them
    Then every one returns Skip, because the plan references sibling plans by function and
      the feature does not become documentation prose

  @DSCDC-B08 @unit-level
  Scenario: A revision cited with an explanation on the same line passes
    Given a spec whose line names a revision code and then says, in substantive prose, what
      that revision changed
    When the gate confronts it
    Then it returns Pass, because the code is the label and the sentence is the content

  @DSCDC-B09 @unit-level
  Scenario: With no map the confrontation is skipped
    Given a spec whose body names what looks like a plan path
    And no graph built
    When the gate confronts it
    Then it returns Skip, because the list of what counts as a reference comes from the map
      and nowhere else

  @DSCDC-I01 @unit-level
  Scenario: The ruler matches structure, never vocabulary
    Given the same bare path reference written in English, Spanish, German and Japanese
    When the gate confronts each spec
    Then every one returns Fail, because a gate that matched words would pass in silence over
      the project written in the other language, and silence is worse than absence

  @DSCDC-I02 @unit-level
  Scenario: The on-line explanation escape belongs to the revision, not to the path
    Given a spec whose line names a plan path and surrounds it with long prose that never says
      what is in the file
    When the gate confronts it
    Then it returns Fail, because a path is the place the person would have to go and no
      amount of surrounding prose says what waits there

  @DSCDC-I03 @unit-level
  Scenario: The verdict names the line and shows what it says
    Given a spec whose fourth line carries a bare plan reference
    When the gate confronts it
    Then the finding carries that line number and an excerpt of the line, so the reader does
      not have to hunt for it

  @DSCDC-X01 @unit-level
  Scenario: The gate does not judge whether the accompanying content is faithful
    Given a spec whose reference is followed by a quotation that says something else entirely
    When the gate confronts it
    Then it returns Pass, because the ruler is whether the reader is left with something to
      read — judging fidelity is another class of gate

  @DSCDC-X02 @unit-level
  Scenario: The gate marks and does not block
    Given a spec with a bare reference
    When the gate confronts it
    Then the finding is recorded as informative and no command is offered to fix it, because
      rewriting a sentence is the work of whoever wrote it

  @DSCDC-X03 @unit-level
  Scenario: Measuring explanation errs on the permissive side
    Given a spec whose line names a revision code followed by prose just past the threshold
      and saying little
    When the gate confronts it
    Then it returns Pass, because the threshold is a coarse ruler and mass false positives
      are what make someone switch an informative gate off
