# language: en
# @anchors
#   ref: GTENG
#   updated_at: 2026-09-26
#   layer: feature

@GTENG
Feature: GateEngine — which gates reach which node, and what the run concludes

  @GTENG-B01 @unit-level
  Scenario: A gate reaches only the kinds it declares
    Given a gate declaring two kinds of node
    When a node of a declared kind and one of another kind are routed
    Then the first is reached and the second is not

  @GTENG-B02 @unit-level
  Scenario: A gate that names labels reaches only the nodes carrying one
    Given a gate naming a label
    When a node carrying that label and one carrying another are routed
    Then only the one carrying the named label is reached

  @GTENG-B03 @unit-level
  Scenario: One excluded label is enough to keep a node out
    Given a gate excluding two labels
    When nodes carrying each of them, and one carrying neither, are routed
    Then the two carrying an excluded label stay out and the third is still charged

  @GTENG-B04 @unit-level
  Scenario: Exclusion wins over the positive label filter
    Given a gate that both names a label and excludes another
    When a node carrying both is routed
    Then it stays out, because layers carry transversal labels and declaring the
      exception must have effect exactly where it matters

  @GTENG-B05 @unit-level
  Scenario: A gate demanding a mark reaches only the targets that carry it
    Given a gate demanding a mark in the target's text
    When a target carrying the mark and one without it are routed
    Then only the first is reached, because otherwise the pending counter would measure
      the size of the project instead of the work

  @GTENG-B06 @unit-level
  Scenario: An unreadable target does not apply
    Given a gate demanding a mark in the target's text
    And a target whose file cannot be read
    When it is routed
    Then it is not reached, because charging blind against unknown content is worse
      than not charging

  @GTENG-B07 @unit-level
  Scenario: A gate with no applicable target does not run at all
    Given a whole-project gate and a slice holding no node it declares
    When the run confronts the slice
    Then the gate produces no verdict, so a commit touching one document does not fire
      a whole-project check

  @GTENG-B08 @unit-level
  Scenario: A gate whose required binary is absent steps aside
    Given a gate requiring a binary that is not installed
    When the run confronts an applicable node
    Then it returns Skip naming the missing binary, because the gate did not measure
      and what is missing is the tool

  @GTENG-B09 @unit-level
  Scenario: A waiver by target spares one node and confronts the rest
    Given a waiver naming one artifact code with a written reason
    And two applicable nodes, one of them that code
    When the run confronts both
    Then the named one returns Skip carrying the reason and the other is still confronted

  @GTENG-B10 @unit-level
  Scenario: An aggregate gate runs once and reports against the scope
    Given a gate of whole-project scope and several applicable nodes
    When the run confronts them
    Then exactly one verdict comes back, and its target names the scope rather than
      one of the files

  @GTENG-B11 @unit-level
  Scenario: A judgment gate asks for judgment when nothing has answered it
    Given a gate declared as judged by an intelligence, and no stamp in the map
    When the run confronts an applicable node
    Then it emits a request for judgment rather than computing a verdict

  @GTENG-B12 @unit-level
  Scenario: A judgment gate reads the stamp an earlier judgement left
    Given a map carrying a recorded judgement for this gate and this node
    When the run confronts it
    Then the recorded answer becomes the verdict, so work already done stops being
      invisible and the pending counter can come down

  @GTENG-B13 @unit-level
  Scenario: A judgement recorded as waived becomes Skip and never Pass
    Given a map carrying a judgement recorded as waived for this node
    When the run confronts it
    Then it returns Skip, because the gate did not measure and Pass would assert an
      approval nobody gave

  @GTENG-B14 @unit-level
  Scenario: A gate declaring neither a command nor a check is undetermined
    Given a gate that declares no way of answering
    When the run confronts an applicable node
    Then it returns Pending naming the omission, rather than passing in silence

  @GTENG-B15 @unit-level
  Scenario: Only the pending item that says a decision is still to take bars promotion
    Given one pending item marked as a decision still to take and another that merely
      had nothing to confront
    When each is produced
    Then only the first is marked as barring promotion, because treating them alike
      failed four hundred and eleven nodes at once

  @GTENG-B16 @unit-level
  Scenario: Only the obligations gate produces assumed debt
    Given a pending verdict from the obligations gate and one from another gate
    When each is produced
    Then only the first is marked as debt and carries its deadline into the record

  @GTENG-B17 @unit-level
  Scenario: The engine reconfigures the code grammar before running anything
    Given a project declaring a rule-type letter outside the canonical vocabulary
    When the run starts
    Then the grammar already recognises that letter, because a gate reading the old
      vocabulary reports green over what it never looked at

  @GTENG-B18 @unit-level
  Scenario: The plain entry point runs with the map and no Structure
    Given gates, nodes and a map, and no Structure at all
    When the plain entry point is called
    Then the gates run, and the absence of a Structure is read as no mapping declared

  @GTENG-B19 @unit-level
  Scenario: The entry point that carries the Structure hands it to the checkers
    Given gates, nodes, a map and a loaded Structure
    When the entry point carrying the Structure is called
    Then the relational checkers receive it, because the regimes and the surfaces of
      the triad are declared there

  @GTENG-B20 @unit-level
  Scenario: The entry point that knows the sweep kind honours the full-sweep scope
    Given a gate declaring one scope for a cut and another for the whole project
    When the entry point is called once for a cut and once for a full sweep
    Then the scope reported differs between the two, because a gate able to sweep on
      its own should run once instead of receiving thousands of targets in batches

  @GTENG-B21 @unit-level
  Scenario: The entry point that honours a waiver by target keeps the gate running
    Given a waiver restricted to one artifact code
    When the entry point that honours it is called over several nodes
    Then the gate produces a verdict for every node, one of them spared

  @GTENG-B22 @unit-level
  Scenario: A vendored file is out of every internal ruler and still reached by an external command
    Given a pipeline Anchors seeded, still carrying its template marker
    When an internal ruler and an external command are routed
    Then the internal ruler does not reach it
    And the external command does

  @GTENG-B23 @unit-level
  Scenario: A gate declared with run executes the custom runner even when canonical specifies check
    Given a canonical gate declaring check
    When a project gate declares run
    Then the custom runner executes rather than the canonical check


  @GTENG-I01 @unit-level
  Scenario: A waiver by target never removes the gate from the list
    Given a waiver naming one of two applicable nodes
    When the run confronts both
    Then the unnamed one still receives a real verdict, because removing the gate would
      erase the ruler for the whole repository and let a defect elsewhere pass along

  @GTENG-I02 @unit-level
  Scenario: Stepping aside, not measuring and failing are three different answers
    Given a gate whose tool is missing and a gate that measured and found a defect
    When both run
    Then their verdicts differ, because each says something the others do not: the gate
      does not apply, the gate could not measure, the target failed

  @GTENG-I03 @unit-level
  Scenario: A waived target leaves the failure tally without leaving the report
    Given a waiver naming a target with a written reason
    When the run confronts it
    Then the verdict carries that reason in its detail, because silence is the
      difference between waiving and hiding

  @GTENG-I04 @unit-level
  Scenario: The reported target of an aggregate gate is the scope
    Given an aggregate gate and several applicable files
    When the run confronts them
    Then the verdict names the scope, because blaming one of many files for an answer
      about the set would be a statement the engine cannot support

  @GTENG-X01 @unit-level
  Scenario: The engine does not decide whether a target is correct
    Given a gate whose checker always approves and a target a reviewer would fail
    When the run confronts it
    Then the verdict is the checker's, because an engine that judged would give the
      same defect two owners

  @GTENG-X02 @unit-level
  Scenario: The engine does not compute the verdict of a judgment gate
    Given a judgment gate and a node whose content plainly answers the question
    When the run confronts it with no stamp recorded
    Then it still asks for judgment, because all the engine knows is whether somebody
      already answered

  @GTENG-X03 @unit-level
  Scenario: The engine invents neither a map nor a Structure nor a waiver
    Given gates and nodes, and no map, no Structure and no waiver
    When the run confronts the nodes
    Then it runs on what it received, because fabricating any of them would answer
      about a project state that does not exist
