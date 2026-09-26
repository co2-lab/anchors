# language: en
# @anchors
#   ref: OBHNB
#   updated_at: 2026-09-26
#   layer: feature

@OBHNB
Feature: ObligationHonored — the cross-cutting duty that lives OUTSIDE the unit

  @OBHNB-B01 @unit-level
  Scenario: A node that carries the trigger and is absent from the demanded file fails
    Given a project declaring an obligation triggered by the personal-data attribute
    And a node whose header carries that attribute and whose token never appears in the purge script
    When the gate confronts it
    Then it returns Fail carrying the declared reason for the duty, so the reader learns
      what the absence costs

  @OBHNB-B02 @unit-level
  Scenario: A node that carries the trigger and does appear passes
    Given the same obligation and a node whose header carries the trigger
    And the purge script naming the token derived from that node
    When the gate confronts it
    Then it returns Pass, because the duty is fulfilled

  @OBHNB-B03 @unit-level
  Scenario: A node without the trigger contracts no obligation
    Given a node whose header does not declare the trigger attribute
    And a purge script that never names it
    When the gate confronts it
    Then it returns Pass, because the duty is charged by what the node declares about
      itself, not by what it might resemble

  @OBHNB-B04 @unit-level
  Scenario: A waiver exempts only when it carries a written reason
    Given a node that carries the trigger and is absent from the purge script
    When the gate confronts it once with the waiver followed by a reason and once with the
      waiver alone
    Then the first returns Pass and the second returns Fail, because the reason is what
      separates the honest exception from silence

  @OBHNB-B05 @unit-level
  Scenario: A project with no declared obligation is skipped
    Given a project whose Structure declares no cross-cutting obligation
    When the gate confronts a node that carries a trigger-looking attribute
    Then it returns Skip, because inventing duties would charge what nobody committed to

  @OBHNB-B06 @unit-level
  Scenario: An acknowledged debt with a written when yields Pending
    Given a node that carries the trigger and is absent from the purge script
    And a debt declaration naming the obligation and the phase in which it will be paid
    When the gate confronts it
    Then it returns Pending carrying that commitment, because the duty still holds and the
      record must stay visible in the report

  @OBHNB-B07 @unit-level
  Scenario: A bare debt marker keeps failing
    Given the same unfulfilled node with a debt marker naming the obligation and nothing more
    When the gate confronts it
    Then it returns Fail, because a marker with no when assumes no debt — it only hides better

  @OBHNB-B08 @unit-level
  Scenario: Waiver and debt stay distinct
    Given one unfulfilled node waiving the obligation with a reason and another acknowledging
      the debt with a when
    When the gate confronts both
    Then only the waiver passes, because the debt is still owed and cannot be stamped as fulfilled

  @OBHNB-B09 @unit-level
  Scenario: The failing verdict offers the three ways out
    Given a node that carries the trigger, is absent from the purge script and declares nothing
    When the gate confronts it
    Then the verdict names fulfilling, waiving with a reason and acknowledging the debt with
      a when, so nobody has to guess what the gate will accept

  @OBHNB-I01 @unit-level
  Scenario: The token is derived through the declared form
    Given a node named MetadataEntry
    When each declared identifier form is applied to it
    Then the raw form yields MetadataEntry, the screaming form METADATA_ENTRY, the snake form
      metadata_entry, the kebab form metadata-entry and a free template the composed token,
      because guessing the shape in the engine would put one project's mess inside the framework

  @OBHNB-I02 @unit-level
  Scenario: A glob that matches no file produces no violation
    Given an obligation whose destination glob matches no file in the project
    And a node that carries the trigger
    When the gate confronts it
    Then it returns Pass, because accusing where there was nothing to read would stamp what
      was never measured

  @OBHNB-I03 @unit-level
  Scenario: The node's own identified_as wins over the automatic form
    Given a node named MetadataEntry whose header declares it is referenced as a plural env var
    And an obligation whose automatic form would derive the singular one
    And the purge script naming only the plural declared by the node
    When the gate confronts it
    Then it returns Pass, because only the node knows the project's real irregularity —
      inverting this order accuses 28 correct models

  @OBHNB-X01 @unit-level
  Scenario: The gate does not decide which obligations exist
    Given a project whose Structure declares no obligation about personal data
    And a node a reviewer would consider obviously purgeable
    When the gate confronts it
    Then it returns Skip, because the duties are the project's decision and a gate that
      invented them would be turned off

  @OBHNB-X02 @unit-level
  Scenario: The gate does not understand what the destination does with the token
    Given a purge script that names the node's token in a dead branch and erases nothing
    When the gate confronts the node that carries the trigger
    Then it returns Pass, because the ruler is presence — separating forgotten from
      remembered is the defect this gate was built for, and judging the implementation is
      another ruler

  @OBHNB-X03 @unit-level
  Scenario: A declaration written in the body is not read
    Given a node that carries the trigger and is absent from the purge script
    And a waiver with a full reason written far down in the document's prose instead of the header
    When the gate confronts it
    Then it returns Fail, because without that cut a quotation in the prose would waive an
      obligation nobody meant to waive

  @OBHNB-E01 @unit-level
  Scenario: An unreadable destination file is not proof, and the others are still searched
    Given a node that carries the trigger
    And the only destination file naming its token cannot be read
    When the gate confronts the node
    Then it returns Fail naming the demanded glob
    And once a readable destination file names the token, it returns Pass

  @OBHNB-E02 @unit-level
  Scenario: A must_appear_in glob that does not parse fails naming it
    Given an obligation whose must_appear_in glob is malformed
    And a node that carries the trigger
    When the gate confronts the node
    Then it returns Fail naming the obligation and the glob
