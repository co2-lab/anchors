# language: en
# @anchors
#   code: LDQRB
#   ref: LDQRS
#   updated_at: 2026-11-19
#   layer: feature

@LDQRS
Feature: SuiteLoad — a suite's run waits for the load queue, and the commit hook jumps it

  @LDQRS-B01 @unit-level
  Scenario: A suite's run holds a lease at the level its files' scope costs, released after; with no load: or with the queue off it holds none
    Given a suite of the workspace packages/backend, with no load:, then with mutation scoped heavy there
    When a whole mutation run, a run over apps/web/a.ts, and a run with the queue off are admitted
    Then no load: holds nothing, the whole run holds heavy, the other medium, and with the queue off nothing
