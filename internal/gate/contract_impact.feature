# language: en
# @anchors
#   code: CIFCN
#   ref: CTRIM
#   updated_at: 2026-10-07
#   layer: feature

@CTRIM
Feature: ContractImpact — a changed field names the rules that use it, and their tests

  @CTRIM-B01 @unit-level
  Scenario: Changed and removed fields change, added ones do not
    Given a committed spec with amount, currency, fee, note and a state a rule uses, then amount, note and the state edited, fee removed and tip added
    When its impact is read
    Then amount, fee and the state are impacted, not currency, tip nor note, which no rule uses

  @CTRIM-B02 @unit-level
  Scenario: A changed field names its rules, here and in the dependents, and their tests
    Given rules of this spec and of a spec depending on its code that use amount, and a test citing one of them
    When amount changes
    Then the impact names every rule and the test file

  @CTRIM-B03 @unit-level
  Scenario: The gate is divergence with the impact, and passes without
    Given the spec with amount changed, the spec unchanged, and a feature node
    When contract-impact confronts each
    Then it is divergence naming the field, passes, and skips

  @CTRIM-B04 @unit-level
  Scenario: The impacted tests are listed for the selection
    Given the spec with uncommitted changes to amount
    When the impacted tests are asked for
    Then the test file of the affected rule is listed

  @CTRIM-B05 @unit-level
  Scenario: The rules the change's own revision revises or checks answer the impact
    Given a field changed that several rules use, and a new revision naming all of them, then one naming only one
    When the impacts are computed
    Then the first leaves no impact, and the second keeps the impact of the rules nobody answered

  @CTRIM-B06 @unit-level
  Scenario: A rule is answered by a revision added to its own spec, even one not yet in git, or by the changed spec's naming its full code; a short code answers only its own spec's rule
    Given the amount changed, read by rules of the pay spec, the checkout spec and a cart spec not yet in git
    When the pay spec's revision names B01 and V01, then each dependent answers its own rule, then the pay spec names CHKOT-B01
    Then the short codes leave the other specs' rules open, each spec's own revision answers its rule, and the full code answers the checkout's

