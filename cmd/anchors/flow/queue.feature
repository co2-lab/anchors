# language: en
# @anchors
#   ref: WRQUW
#   updated_at: 2026-09-26
#   layer: feature

@WRQUW
Feature: WorkQueue — list, pull, close and discard the work, from the local queue or from the board

  @WRQUW-B01 @unit-level
  Scenario: queue lists the live tasks with the hygiene hints
    Given an empty queue, and then a queue with a task claimed by "w1" and a task waiting for triage
    When `anchors queue` runs
    Then the first says "empty queue"
    And the second lists "2 task(s) in the queue", "◐ [claimed]", "by:         w1", "○ [pending]", "1 claimed" and "1 in 'triage'"

  @WRQUW-B02 @unit-level
  Scenario: next claims the next task in local mode
    Given a local queue with the task "a-code" and one informative gate declared
    When `anchors next --worker w9` runs
    Then the output says "task claimed: a-code", "suggestion: feature" and "anchors done a-code"
    And it says "1 informative gate(s) declared"
    And the task is claimed by "w9"

  @WRQUW-B03 @unit-level
  Scenario: A cold start seeds a plan that still has work
    Given an empty queue and a plan citing a spec that does not exist
    When `anchors next` runs
    Then the output says "seeded with 1 plan(s)" and "task claimed:"
    And a project whose every plan is fulfilled says "empty queue — nothing to do"

  @WRQUW-B04 @unit-level
  Scenario: A cited spec exists by path or by a unique name
    Given the citation "Pricing.spec.md"
    When it is checked with no file, with one file of that name, and with two
    Then it exists only with exactly one file of that name

  @WRQUW-B05 @unit-level
  Scenario: The seed count is recomputed when printed
    Given a seeded plan citing three specs
    When its task is printed with none on disk, and then with two delivered
    Then it says "3 of 3", and then "1 of 3"
    And a task whose reason already has a count, and a task that is not a plan seed, get no count

  @WRQUW-B06 @unit-level
  Scenario: done closes by id and in batch
    Given five tasks on files a, a, b, c and d
    When `anchors done` runs with no argument, with the id "t4", with `--file src/a.ts`, with `--kind spec`, with `--kind code`, with `--all`, and with the id "nope"
    Then the first fails with "provide the <id>", "t4" is done, the file closes two, the spec kind matches nothing, the code kind closes one, all empties the queue, and "nope" is an error

  @WRQUW-B07 @unit-level
  Scenario: drop deletes a task without archiving it
    Given the task "junk" in the queue
    When `anchors drop junk` runs twice
    Then the first says "task discarded: junk" and the task is gone
    And the second is an error

  @WRQUW-B08 @unit-level
  Scenario: reclaim respects a live worker unless forced
    Given a task claimed a moment ago by "w1"
    When `anchors reclaim` runs, and then `anchors reclaim --force`
    Then the first says "0 task(s) returned" and "1 task(s) claimed RECENTLY stayed"
    And the second says "1 task(s) returned" and the task is pending again

  @WRQUW-B09 @unit-level
  Scenario: The default worker is pid at host
    Given no worker name
    When the default worker is made
    Then it is this process's pid, "@", and the host

  @WRQUW-B10 @unit-level
  Scenario: The board is not claimed without a session
    Given a project in github mode and no ANCHORS_SESSION
    When `anchors next` runs
    Then it fails with "ANCHORS_SESSION is not set" and says to "export ANCHORS_SESSION="

  @WRQUW-B11 @unit-level
  Scenario: The board identity falls back to the OS user, and says so
    Given no ANCHORS_SESSION and the OS user "alice"
    When the agent's identity is made, and again with the session "devA"
    Then the first ends in "/alice" and standard error warns "ANCHORS_SESSION is not set"
    And the second ends in "/devA" with no warning

  @WRQUW-B12 @unit-level
  Scenario: The agent's own card is resumed before asking the pipeline
    Given a board where card 42 is owned by this host's "dev3", and the session "dev3"
    When `anchors next` runs
    Then the output says "resuming your card" and "card claimed: #42"
    And nothing is asked of the claim pipeline

  @WRQUW-B13 @unit-level
  Scenario: A claim without a card names the run to follow
    Given a claim that timed out, and a claim run 77 that succeeded with no card
    When each is reported
    Then the first says "did not finish within" and the runs list to follow, without error
    And the second says "no free card on the board" and "gh run view 77 --log", without error

  @WRQUW-B14 @unit-level
  Scenario: The claimed card is printed with its state and owner
    Given card 42 in progress owned by this host's "dev3"
    When `anchors next` resumes it
    Then the output says "card claimed: #42", "state:    in-progress" and the owner

  @WRQUW-B15 @unit-level
  Scenario: The deliverable follows the card's title
    Given a map where PRICX is "src/pricing.ts"
    When a plan card, a spec card whose body names PRICX, and another card whose code is gone are printed
    Then the plan card asks for "every spec the plan seeds"
    And the spec card lists "3. anchors work test    --for src/pricing.ts"
    And the other card points to its body and falls back to the folder "packages/infra"

  @WRQUW-B16 @unit-level
  Scenario: The documentation duties follow the unit
    Given a project requiring "docs/openapi.yaml" for changes in the layer "logic", and one requiring nothing
    When the duties of "src/pricing.ts" are printed in each
    Then the first says "Changing `src/pricing.ts` REQUIRES touching:" and names "docs/openapi.yaml"
    And the second prints nothing

  @WRQUW-B17 @unit-level
  Scenario: Only an agent that does not decide the product is told to escalate
    Given an agent with no role, and then one whose role is product owner
    When a card is printed for each
    Then the first says "YOU DO NOT DECIDE THE DIRECTION" and the second does not

  @WRQUW-B18 @unit-level
  Scenario: A card under review asks for a verdict
    Given a card under review whose body names PRICX, and review card 931 naming no unit
    When each is printed for the agent "host/dev3"
    Then the first asks for "anchors work review --for src/pricing.ts"
    And the second points to "#931 in:body"
    And both end with "anchors-review: approved by host/dev3" and "anchors-review: rejected by host/dev3", without "anchors pr-body"

  @WRQUW-B19 @unit-level
  Scenario: Other cards end naming the pull request body command
    Given card 42 in progress
    When `anchors next` resumes it
    Then the output says "anchors pr-body --cards 42"

  @WRQUW-B20 @unit-level
  Scenario: The role is asked at most once, and never without a terminal
    Given a declared developer role, and then standard input from a pipe and from the null device
    When the local decision is ensured, and the terminal is checked
    Then the declared role is kept silently
    And neither the pipe nor the null device is a terminal

  @WRQUW-I01 @unit-level
  Scenario: A cold start seeds one plan
    Given two plans with specs still missing
    When `anchors next` runs on an empty queue
    Then exactly one plan task is seeded

  @WRQUW-X01 @unit-level
  Scenario: queue claims nothing
    Given a queue with a pending task
    When `anchors queue` runs
    Then the task is still pending

  @WRQUW-X02 @unit-level
  Scenario: The board is not claimed without a repository
    Given a github configuration whose repository is empty
    When the board is asked for the next card
    Then it fails with "workflow.repo is empty"

  @WRQUW-E01 @unit-level
  Scenario: A failed claim run is an error
    Given a claim run 78 that ended in failure with no card
    When it is reported
    Then it fails with "#78 ended as \"failure\""
