# language: en
# @anchors
#   ref: GVGDG
#   updated_at: 2026-10-03
#   layer: feature

@GVGDG
Feature: GovernanceGuides — the guides an agent reads to operate Anchors, and the contracts other code relies on

  @GVGDG-B01 @unit-level
  Scenario: Bare guide prints the operating playbook
    Given the governance commands registered on a root
    When anchors guide runs with no subcommand
    Then it prints "# Operating Anchors (guide for AI agents)" with the development flow and the command reference

  @GVGDG-B14 @unit-level
  Scenario: guide --help lists every subcommand
    Given the governance commands registered on a root
    When the help of anchors guide is read
    Then its list names each of the thirteen subcommands exactly once, and nothing else

  @GVGDG-B02 @unit-level
  Scenario: Each guide subcommand prints its own guide
    Given the governance commands registered on a root
    When each subcommand of anchors guide runs
    Then the subcommands are exactly changelog, code, feature, flag, flow, guide, header, plan, product, project, report-bug, review, spec, test and work, each with a short description
    And each prints a different text opening with its own title, such as "# Work guide" for work

  @GVGDG-B03 @unit-level
  Scenario: The review and work guides append the autonomy section of the root they are given
    Given a directory whose local settings declare the role dev, and a path that is not a project
    When anchors guide review and anchors guide work run with --root on each
    Then with the first they print their guide followed by "**Your role (Dev) does NOT decide"
    And with the second they print their guide followed by "**You did not declare a role**"

  @GVGDG-B04 @unit-level
  Scenario: Register adds exactly the four governance commands
    Given an empty root command
    When the governance package registers its commands
    Then the root holds exactly audit, compliance, governs and guide

  @GVGDG-B05 @unit-level
  Scenario: The review guide teaches who counts and which verdict wins
    Given the review guide as anchors guide review prints it
    When a reviewer reads the verdict section
    Then it says "Only the reviewer the claim assigned counts", "only a line posted after the assignment" and "the last one wins"
    And that a line inside a code block is an example, not a verdict

  @GVGDG-B06 @unit-level
  Scenario: The review guide says the reviewer does not move the card
    Given the review guide as anchors guide review prints it
    When a reviewer reads it
    Then it says "YOU DO NOT MOVE THE CARD", that the checks and the MERGE move the states
    And that the wrong state is "more expensive than the late state"

  @GVGDG-B07 @unit-level
  Scenario: The review guide carries continuous conformance points anchored in its prose
    Given the review guide as anchors guide review prints it
    When its conformance points are read
    Then the heading "## Pontos de conformidade" is one the guide-checklist gate recognises
    And the points run from REV-CK1 to REV-CK18 with no gap, each with its anchor in the prose above the list

  @GVGDG-B08 @unit-level
  Scenario: The work guide teaches the claim as the first step
    Given the work guide as anchors guide work prints it
    When an agent reads it from the top
    Then "anchors next" appears before "## The order", with ANCHORS_SESSION and the reason "find it empty"

  @GVGDG-B09 @unit-level
  Scenario: The work guide makes the agent wait for the CI verdict
    Given the work guide as anchors guide work prints it
    When an agent reads what follows the push
    Then it reads "is not a stopping point", "--watch", "work review" and "half the work"

  @GVGDG-B10 @unit-level
  Scenario: The work guide forbids closing the card by hand
    Given the work guide as anchors guide work prints it
    When an agent reads it
    Then it reads "YOU DO NOT CLOSE THE CARD" and "Who closes it is the MERGE", with "25 green PRs" as the measured cost

  @GVGDG-B11 @unit-level
  Scenario: The work guide sends a finding through escalate and reads the decision queue first
    Given the work guide as anchors guide work prints it
    When an agent looks for what to do with a finding that is not its card
    Then it reads "anchors escalate" with "--card"
    And "gh issue list --label anchors:needs-user" appears before "--for-user"

  @GVGDG-B12 @unit-level
  Scenario: The project guide covers the discover phase and the playbook points to it before planning
    Given the project guide and the playbook as anchors guide prints them
    When an agent starts a project that does not exist yet
    Then the project guide covers PROJECT.md, INSIGHTS.md and "Stage 1 — Purpose and form" through "Stage 5 — Tooling and formatting"
    And the playbook's "### 0.5. DISCOVER" comes before "### 1. PLAN"

  @GVGDG-B13 @unit-level
  Scenario: The test guide names the instrument per input shape and teaches the stamp refresh
    Given the test guide as anchors guide test prints it
    When a boundary rule is to be proven
    Then it names "SMALL AND CLOSED", "LARGE BUT STRUCTURED" and "OPEN" with their instruments
    And it teaches "anchors stamp --refresh" in the "same commit"

  @GVGDG-I01 @unit-level
  Scenario: Every command a guide cites exists in the command tree
    Given the whole command tree, built from every package's registration
    When every `anchors <command>` cited by the playbook and the thirteen guides is looked up
    Then each one is a registered command

  @GVGDG-I02 @unit-level
  Scenario: The work guide uses the real label names
    Given the work guide as anchors guide work prints it
    When its labels are compared with the workflow label constants
    Then it contains "anchors:under-" and "anchors:needs-user" and never "anchors:sob-"

  @GVGDG-I03 @unit-level
  Scenario: The verdict line the review guide teaches is the one the pipeline parses
    Given the verdict expression of the seeded anchors-pr-checks.yml
    When it reads "anchors-review: approved by agent-b" and "anchors-review: rejected by agent-b" as the guide teaches them
    Then it extracts approved and rejected, each by agent-b

  @GVGDG-X01 @unit-level
  Scenario: Only the review and work guides depend on the project
    Given the thirteen guide subcommands
    When their flags and output are inspected
    Then only review and work take --root
    And no other guide prints the autonomy section

  @GVGDG-B15 @unit-level
  Scenario: The guides tell a fix from a bug and ask for the marker
    Given the work, code and review guides
    When they are printed
    Then the work guide says a bug is a defect that shipped, shows the fix commit with its Bug footer and asks for the failing test first, and the code and review guides point to it

  @GVGDG-B16 @unit-level
  Scenario: The changelog guide tells the technical changelog from the product one
    Given the changelog guide
    When it is printed
    Then it says the generated changelog is technical and recommends a product changelog synthesized from it, with bugs but not fixes, and chores only when they matter to the product

  @GVGDG-B17 @unit-level
  Scenario: The spec guide asks for four passes and a review
    Given the spec guide
    When it is read
    Then it walks every section with its questions, derives the variations, generalizes into invariants, reviews intent against mechanism, and has a defect's rule written before its fix

  @GVGDG-B18 @unit-level
  Scenario: The report-bug guide says how to tell, report and go on
    When `anchors guide report-bug` runs, and `anchors guide`
    Then the guide tells Anchors' bugs from the project's, asks for a made-up minimal case and a dry run, and says what to do while the fix does not come
    And the playbook points to the guide and to anchors report-bug

  @GVGDG-B19 @unit-level
  Scenario: The guides recommend visual regression for every state of a visual unit
    When the test, spec and feature guides are printed
    Then the test guide asks for the VR scenario, a baseline per state and a capture test
    And the spec and feature guides point a visual unit's states to it

  @GVGDG-B20 @unit-level
  Scenario: The test guide recommends a contract test against the compiled OpenAPI
    When the test guide is printed
    Then it recommends a contract test named by {CODE}-CT that loads the OpenAPI document

  @GVGDG-B21 @unit-level
  Scenario: The guides tie validations to states and errors to messages
    When the spec and test guides are printed
    Then validations are triggers of State Flow transitions, errors cite their messages, and every message is captured

