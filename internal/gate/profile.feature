# language: en
# @anchors
#   code: PRFTC
#   ref: PRFLO
#   updated_at: 2026-10-03
#   layer: feature

@PRFLO
Feature: Profile — the verdicts of a run, gathered per gate and per node

  @PRFLO-B01 @unit-level
  Scenario: Each gate gets a summary counting its verdicts
    Given results of gate "a" with two passes and a failure, of gate "b" with two skips and two pending, and of gate "c" with two awaiting judgement
    When the results are aggregated
    Then gate "a" counts 2 passes and 1 failure, gate "b" 2 skips and 2 pending, gate "c" 2 awaiting judgement, and the gate names come back as a, b, c

  @PRFLO-B02 @unit-level
  Scenario: The summary carries the total time and the most expensive run
    Given three results of one gate taking 7, 7 and 1 milliseconds
    When the results are aggregated
    Then the gate's total time is 15 milliseconds and its most expensive run is 7 milliseconds

  @PRFLO-B03 @unit-level
  Scenario: Every failed result is listed as a failure
    Given a failed result on target "z" among passes, skips and pending verdicts
    When the results are aggregated
    Then the failures list holds exactly the result on "z"

  @PRFLO-B04 @unit-level
  Scenario: Only a failure of a blocking gate blocks promotion
    Given one failure of a non-blocking gate, and separately one failure of a blocking gate
    When each set of results is aggregated
    Then the non-blocking failure leaves promotion allowed, and the blocking failure refuses it

  @PRFLO-B05 @unit-level
  Scenario: A pending verdict blocks only when its gate blocks and it impedes
    Given a blocking pending that does not impede, a non-blocking pending that impedes, and a blocking pending that impedes
    When each result is aggregated
    Then only the blocking pending that impedes refuses promotion

  @PRFLO-B06 @unit-level
  Scenario: Results awaiting judgement are listed as awaiting judgement
    Given two results whose verdict awaits an AI judgement
    When the results are aggregated
    Then the list of results awaiting judgement holds both

  @PRFLO-B07 @unit-level
  Scenario: Per-node verdicts include only confronted nodes, sorted
    Given results on "c.go", "a.go", "b.go", "pend.go", a skip on "skip.go" and a verdict awaiting judgement on "judge.go"
    When the per-node verdicts are collapsed
    Then the nodes are a.go, b.go, c.go and pend.go, in that order, without skip.go or judge.go

  @PRFLO-B08 @unit-level
  Scenario: A node is failed when a blocking gate failed on it or left it an impeding pending
    Given "c.go" failed by a blocking gate, "impede.go" left an impeding pending by a blocking gate, and "a.go" failed by a non-blocking gate
    When the per-node verdicts are collapsed
    Then "c.go" and "impede.go" are marked failed and "a.go" is confronted but not failed

  @PRFLO-I01 @unit-level
  Scenario: Promotion is refused exactly when something blocks
    Given each case of the promotion table
    When the result is aggregated
    Then the pass flag is false exactly when the blocked list is not empty

  @PRFLO-X01 @unit-level
  Scenario: The profile does not decide whether a pending impedes
    Given a blocking pending result that its gate did not mark as impeding
    When the result is aggregated
    Then promotion stays allowed, because the profile does not reinterpret the verdict
