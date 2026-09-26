# language: en
# @anchors
#   ref: INCHN
#   updated_at: 2026-09-26
#   layer: feature

@INCHN
Feature: InternalChecks — the registry that routes a declared check name to a function

  @INCHN-B01 @unit-level
  Scenario: A declared name routes to the function registered under it
    Given a gate declaring a check name that the registry knows
    When the routing resolves the name
    Then the registered function answers, and its verdict is the gate's

  @INCHN-B02 @unit-level
  Scenario: A name that does not resolve answers undetermined
    Given a gate declaring a check name no registry holds
    When the routing resolves the name
    Then it answers Pending and never Pass, because a checker that never ran has
      measured nothing

  @INCHN-B03 @unit-level
  Scenario: The routing tries the relational registry before the simpler ones
    Given a name registered among the relational checkers
    When the routing resolves it
    Then the relational function answers, so a checker that grows a dependency changes
      registry without changing name

  @INCHN-B04 @unit-level
  Scenario: The per-node path reads the target file, and a failed read is a failure
    Given a node whose file does not exist on disk
    When a per-node checker is routed to it
    Then it returns Fail, because a checker of content with no content has nothing
      to answer

  @INCHN-B05 @unit-level
  Scenario: The aggregate path reads no file at all
    Given a gate of aggregate scope whose checker is relational
    When the routing resolves it
    Then the checker receives an empty node, because trying to read a file whose scope
      is the SET would fail over a file that never existed

  @INCHN-B06 @unit-level
  Scenario: An unresolved name in the aggregate path is undetermined too
    Given a gate of aggregate scope declaring a check name no registry holds
    When the routing resolves it
    Then it answers Pending, and the report names the check that failed to route

  @INCHN-B07 @unit-level
  Scenario: An aggregate checker that needs the declaring gate receives it
    Given a parameterised gate whose checker reads its own configuration
    When the routing resolves it
    Then the checker receives that gate, because otherwise every instance would see
      the same configuration or none

  @INCHN-B08 @unit-level
  Scenario: Setting the rule letters reconfigures every dependent pattern together
    Given a project declaring a rule-type letter outside the canonical vocabulary
    When the grammar is reconfigured
    Then both the scenario-code ruler and the feature ruler recognise that letter,
      because a pattern left behind reports green over what it never looked at

  @INCHN-B09 @unit-level
  Scenario: A file that is empty or only whitespace fails the emptiness ruler
    Given a file holding nothing but blanks, and another holding text
    When each is confronted
    Then the blank one fails and the other passes

  @INCHN-B10 @unit-level
  Scenario: The identity ruler charges the presence of a scenario code
    Given a file carrying a scenario code and another carrying none
    When each is confronted
    Then the first passes and the second fails, because a piece with no code is an
      invisible orphan

  @INCHN-B11 @unit-level
  Scenario: A governed file with no identity block fails the header ruler
    Given a source file carrying no identity block at all
    When the header ruler confronts it
    Then it returns Fail

  @INCHN-B12 @unit-level
  Scenario: A governed file passes with ownership or with reference, never with layer alone
    Given three governed files: one declaring ownership, one declaring reference, and
      one declaring only its layer
    When the header ruler confronts each
    Then the first two pass and the third fails, because a governed layer demands
      ownership or reference

  @INCHN-B13 @unit-level
  Scenario: A file of a recognised layer passes with the layer alone
    Given a file whose declared layer is a recognised one, carrying no code and no reference
    When the header ruler confronts it
    Then it returns Pass, because it has neither an owning spec nor a sibling to reference

  @INCHN-B14 @unit-level
  Scenario: A binary file steps aside from the header ruler
    Given a node whose content carries the bytes of an image
    When the header ruler confronts it
    Then it returns Skip, because there is no comment syntax in an image and charging
      one would bar every visual baseline commit

  @INCHN-B15 @unit-level
  Scenario: An executable test script steps aside by a different path
    Given a test node whose file is a data document consumed by an external runner
    When the header ruler confronts it
    Then it returns Skip, because the format belongs to the runner and the identity is
      in the file name

  @INCHN-B16 @unit-level
  Scenario: A guide without compliance points fails
    Given a guide with no compliance-points section, and another with the section and
      no item inside it
    When each is confronted
    Then both fail, because the AI judgment gate would otherwise fall back on vague
      heuristics

  @INCHN-I01 @unit-level
  Scenario: An unresolved name never approves, on either path
    Given a check name that no registry holds
    When it is routed through the per-node path and through the aggregate one
    Then neither returns Pass, because approving would stamp green over a verification
      that never ran

  @INCHN-I02 @unit-level
  Scenario: Every registered name is reachable through exactly one routing path
    Given every name the registries hold
    When each is routed
    Then each resolves, because a name reachable from nowhere is a gate accepted in
      silence that measures nothing

  @INCHN-I03 @unit-level
  Scenario: The compliance ruler is recognised in every language of the catalogue
    Given a guide whose compliance-points heading is written in a language other than
      the engine's own
    When the guide is confronted
    Then the section is recognised, because a project seeded in another language would
      otherwise be born failing a guide its own tooling had just written

  @INCHN-X01 @unit-level
  Scenario: The registry does not decide which checks a project runs
    Given a registry holding many checkers and a project declaring one gate
    When the project is confronted
    Then only the declared check is routed, because charging the rest would bill a
      project for a ruler it never adopted

  @INCHN-X02 @unit-level
  Scenario: The registry does not invoke external tooling
    Given a check routed with no executable of any kind available to it
    When it answers
    Then it answers from the text alone, because a registry lookup must not depend on
      what is installed on the machine

  @INCHN-X03 @unit-level
  Scenario: The registry does not judge whether the text is good
    Given a file whose header is present and whose content a reviewer would call wrong
    When the header ruler confronts it
    Then it returns Pass, because the ruler here is presence and shape — judging the
      content belongs to another gate

  @INCHN-E01 @unit-level
  Scenario: A test missing from disk does not hide the codes the other tests name
    Given a map listing a test that is gone from disk and a test that names a scenario code
    When scenario-coverage checks which codes are written
    Then the code the present test names is reported as written, not as missing a test
