# language: en
# @anchors
#   code: DPFTA
#   ref: GTDPG
#   layer: feature

@GTDPG
Feature: Duplicates — each gate confronts the repeats of what it declares

  @GTDPG-B01 @unit-level
  Scenario: A key declared twice turns the gate's verdict into a failure naming the lines, unless switched off or skipped
    Given a spec defining a rule code twice, and a gate with the rule-types reader
    When the gate's verdict is confronted with the repeats, then with duplicates false, then on a skipped node, then with a reader that measures repeats as a divergence
    Then it fails naming the code and both lines beside the gate's own finding; switched off or skipped it stays; the divergence reader diverges, and a failure stays a failure

  @GTDPG-B02 @unit-level
  Scenario: rule-types counts the rule codes a file defines, not those it cites
    Given a spec defining B01 in two sections, and citing it in what rules use, an Out table, the state flow, an events table under the title the project declares, an open decision, an alias, a retired line and a prose bullet, with a heading followed by its own row
    When its occurrences are read
    Then B01 is counted at its two definitions only, and the heading with its own row once
