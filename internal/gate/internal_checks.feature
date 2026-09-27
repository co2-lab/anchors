# language: en
# @anchors
#   ref: INCHN
#   updated_at: 2026-09-27
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

  @INCHN-B17 @unit-level
  Scenario: With a tests source a scenario is written only when a title cites it
    Given a spec whose code appears in its test's fixture but in no title
    When scenario-coverage runs with the project's tests declared, and without
    Then with the declaration the scenario has no test, and without it the scenario is written
    And a file the source lists no test in, such as a YAML flow under a Jest pattern, is read as without a source

  @INCHN-E02 @unit-level
  Scenario: A failing tests source fails scenario-coverage naming the error
    Given a project whose tests script exits with an error
    When scenario-coverage runs
    Then it fails naming the script's error

  @INCHN-B18 @unit-level
  Scenario: A support file is neither run nor counted as naming a scenario
    Given a support file with no execution, and a support file that is the only one citing a spec's code
    When tests-pass judges the first and scenario-coverage reads the second
    Then tests-pass skips saying it is support, and the scenario has no test

  @INCHN-B19 @unit-level
  Scenario: A feature scenario is recognised in any Gherkin language
    Given features whose scenarios open with Portuguese, English, Spanish and French keywords,
      outlines, the Example synonym and unaccented spellings
    And a feature holding only an examples table, and one holding only its title
    When the emptiness ruler confronts each
    Then every feature with a scenario passes and the other two fail

  @INCHN-B20 @unit-level
  Scenario: Scenario coverage charges what the spec defines, not what it cites
    Given a spec defining two proven requirements whose prose cites other units' codes
    And a spec defining two requirements of which only one was proven
    When scenario-coverage confronts each
    Then the first passes naming none of the cited codes
    And the second fails naming the unproven requirement

  @INCHN-B21 @unit-level
  Scenario: Scenario coverage tells a missing test apart from a test never run
    Given a spec with two requirements and no ingested execution
    When scenario-coverage runs with no test naming them, and with a test naming only the first
    Then the first run fails naming both requirements
    And the second fails asking to ingest the run of the named one while still naming the other
    And once an ingested execution proves both, it passes

  @INCHN-B22 @unit-level
  Scenario: Scenario coverage honours a layer that dispenses tested-by
    Given a spec of a layer dispensing tested-by and a spec of a layer that does not, neither tested
    When scenario-coverage confronts each
    Then the first is skipped saying tested-by is dispensed, and the second fails

  @INCHN-B23 @unit-level
  Scenario: Mutation score passes at the threshold and fails below it naming the survivors
    Given fresh mutation signals scoring above the threshold, exactly at it, below it, and with every mutant ignored
    When mutation-score confronts each
    Then the one below fails naming how many mutants survived and the threshold
    And the others pass

  @INCHN-B24 @unit-level
  Scenario: A missing or stale mutation signal is pending
    Given a file with no mutation signal, and a file whose perfect score was measured at an older revision
    When mutation-score confronts each
    Then both are pending, the first saying what to ingest and the second that the signal is stale

  @INCHN-B25 @unit-level
  Scenario: A score between acceptable and desirable is pending, not failed
    Given an acceptable threshold of 70 and a desirable one of 90
    When mutation-score confronts scores of 75 and 92
    Then 75 is pending, naming both ranges and the 15 points left, and 92 passes clean
    And with no desirable threshold, or one below the acceptable, 75 passes

  @INCHN-B26 @unit-level
  Scenario: The mutation thresholds come from the report
    Given a score of 65
    When the report declares an acceptable threshold of 60, and when it declares none
    Then it passes against 60 and fails against the default of 70

  @INCHN-B27 @unit-level
  Scenario: The verdict follows the isolated scope and the report reads the delta
    Given a file scoring 8% isolated and 77% full, and one scoring 30% isolated and 100% full
    When mutation-score confronts each
    Then both fail, and the first report names both scores, their delta and the dependents
    And a low delta at a low score is reported as a missing assertion, never as coupling at the same time
    And an isolated score above the threshold passes, as does a total score with no scopes

  @INCHN-B28 @unit-level
  Scenario: A scope measured at an older revision decides nothing
    Given an isolated score of 30% measured at the previous revision and a full one of 95% at the current
    When mutation-score confronts the file
    Then it does not fail and the old score is not in the report
    And the same isolated score measured at the current revision, or with no revision stamp, fails

  @INCHN-B29 @unit-level
  Scenario: Without a repository updated-at skips instead of blaming the file
    Given a file with an identity date in a directory outside any git repository
    When updated-at-atual confronts it
    Then it skips naming the missing repository, not saying the file is uncommitted
    And in a fresh repository the same file dated today passes and dated 2020 fails

  @INCHN-B30 @unit-level
  Scenario: A section title in another language fails unless the gate waives it
    Given a Portuguese project and a spec whose section titles are in English
    When spec-sections confronts it
    Then it fails naming the expected Portuguese title
    And it passes with enforce_section_language false, with no configuration, or when the titles are in Portuguese
    And a section title outside the catalogue is not charged

  @INCHN-B31 @unit-level
  Scenario: The skeletons anchors new emits are born conforming
    Given the spec, feature and test that anchors new emits for a unit
    When the header ruler confronts each and spec-sections confronts the spec
    Then every one of them passes
