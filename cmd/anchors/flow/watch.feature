# language: en
# @anchors
#   code: WTFTW
#   ref: WTCHA
#   updated_at: 2026-10-03
#   layer: feature

@WTCHA
Feature: Watch — the background watcher that turns "a file changed" into "there is work in the queue"

  @WTCHA-B01 @unit-level
  Scenario: status tells stopped, running with its metadata, and paused
    Given a project with no watcher, then one whose pid is a live process with metadata, then paused
    When `anchors watch status` runs in each state
    Then it says "watcher: stopped", then "watcher: running (pid N)" with the root, then "watcher: paused"

  @WTCHA-B02 @unit-level
  Scenario: A second start is refused
    Given a project whose watcher is running
    When `anchors watch start` runs
    Then the command fails with "already running"

  @WTCHA-B03 @unit-level
  Scenario: pause needs a running watcher, and resume undoes it
    Given a project with no watcher
    When `anchors watch pause` runs
    Then it fails with "not running"
    And with a running watcher, pause leaves the pause flag and resume removes it

  @WTCHA-B04 @unit-level
  Scenario: stop terminates the watcher once
    Given a project whose watcher is a live process
    When `anchors watch stop` runs, and then again
    Then the process was terminated and the pid file is gone
    And the second stop fails

  @WTCHA-B05 @unit-level
  Scenario: logs prints the log as is
    Given a project with no log, and then the log "● src/a.ts [code] → task queued"
    When `anchors watch logs` runs
    Then the first fails with "no log" and the second prints the log unchanged

  @WTCHA-B06 @unit-level
  Scenario: run names the configuration or the map it could not load
    Given a project without anchors.yaml, and then one with a broken map
    When `anchors watch run` runs
    Then it fails with "load config", and then with "load map"

  @WTCHA-B07 @unit-level
  Scenario: Files created after the start become tasks until the signal
    Given the loop running over a project
    When "src/pricing.ts" and "src/fresh/handler.ts" are created, and the process receives SIGTERM
    Then both are queued for their feature
    And the output says "watching" and "shutting down", and the pid file is gone

  @WTCHA-B08 @unit-level
  Scenario: The tree is watched without ignored folders, and new folders are swept
    Given a project with "src/a/x.ts", "src/b.ts" and "node_modules/lib/y.js"
    When the tree is added to the watcher
    Then three directories are watched and none under node_modules
    And the files born in "src" are only "b.ts"

  @WTCHA-B09 @unit-level
  Scenario: A governed change queues the next missing piece
    Given the unit "src/pricing.ts" whose feature already exists
    When "src/pricing.ts" changes
    Then the task "src/pricing.ts→test" is queued and reported with "queue: 1 task(s)"

  @WTCHA-B10 @unit-level
  Scenario: A waived piece is skipped
    Given the layer "model" waiving the feature and the test, and the unit "models/user.ts"
    When "models/user.ts" changes
    Then the task "models/user.ts→review" is queued
    And the layer of "models/user.spec.md" is "model", while a spec with no code and a file outside the structure have none

  @WTCHA-B11 @unit-level
  Scenario: What is not work is not queued
    Given "README.md" outside the structure, "src/gone.ts" missing, "src/probes/probe1.ts" ignored and the editor temporary "src/.!21662!pricing.spec.md"
    When each changes
    Then nothing is queued

  @WTCHA-B12 @unit-level
  Scenario: A delivery record triggers the review when the unit closes
    Given the record "changes/a.md" for the unit "src/pricing.ts"
    When it is handled with no code, with code only, and with code and test
    Then the first says "the review waits for the unit to close", the second queues nothing, and the third queues "changes/a.md→review"
    And the plan record "changes/plan.md" queues "review-plan", and "changes/reviewed/old.md" queues nothing

  @WTCHA-B13 @unit-level
  Scenario: A waived test, a Go test, and a record with no unit do not hold the review
    Given "models/user.ts" in a layer that waives the test, and "svc/handler.go" beside "svc/handler_test.go"
    When the missing pieces are checked
    Then neither holds the review, and neither does a record with no unit

  @WTCHA-B14 @unit-level
  Scenario: Any markdown directly under changes is a delivery record
    Given "changes/a.md" with the unit "src/x.ts" and "changes/b.md" with no unit line
    When each is read as a delivery
    Then the first names "src/x.ts" and the second is a delivery with no unit
    And "changes/missing.md", "changes/a.txt" and "docs/a.md" are not deliveries

  @WTCHA-B15 @unit-level
  Scenario: The task id is stable and path-safe
    Given the path "src/x/pricing.ts" and the step "test"
    When its task id is made twice, and once for the step "review"
    Then the two ids are equal, start with "src-x-pricing-test-", and hold no "/"
    And the review id differs

  @WTCHA-I01 @unit-level
  Scenario: A task in the queue is not queued twice
    Given the task "src/pricing.ts→test" already queued
    When "src/pricing.ts" changes again
    Then the output says "already in the queue (test)" and the queue keeps one task

  @WTCHA-I02 @unit-level
  Scenario: A change pending at the signal is handled before the loop returns, and none after
    Given a running loop with a 400ms debounce and the node "src/pricing.ts" in the map
    When "src/pricing.ts" is written and the loop receives SIGTERM inside the debounce window
    Then the task "src/pricing.ts→feature" is queued and the node's rev updated when the loop returns, and the node does not change after it

  @WTCHA-X01 @unit-level
  Scenario: Handling a change writes nothing outside the queue
    Given the unit "src/pricing.ts"
    When "src/pricing.ts" changes
    Then every file written is under ".anchors/"

  @WTCHA-B16 @unit-level
  Scenario: A file the watcher sees for the first time enters the map
    Given a project with a map, and a spec created while the watcher runs
    When the watcher handles the change
    Then the map on disk and the watcher's copy have the spec
