# language: en
# @anchors
#   ref: BRSRB
#   updated_at: 2026-09-26
#   layer: feature

@BRSRB
Feature: BoardServe — the board page served locally with live state, read from the host only when something changed

  @BRSRB-B01 @unit-level
  Scenario: The live payload says it is live and stamps the read time
    Given one item read at 2026-09-24 12:00 in UTC-3
    When the board data is assembled
    Then it holds live true, takenAt 2026-09-24T15:00:00Z and the one item

  @BRSRB-B02 @unit-level
  Scenario: A read inside the floor does not call the host
    Given a board read once with a floor of one hour
    When it is read again
    Then the same bytes come back and no call is made

  @BRSRB-B03 @unit-level
  Scenario: Past the floor the board sweeps only when something newer exists
    Given a board whose newest card was updated at 2026-09-21T10:00:00Z
    When it is read past the floor with the host answering that same time, then with an older baseline
    Then the first read asks once and does not sweep
    And the second sweeps again

  @BRSRB-B04 @unit-level
  Scenario: A failed sweep after a good read serves the good read
    Given a board with a previous read and a host that refuses everything
    When it is read
    Then the previous read comes back with no error

  @BRSRB-B05 @unit-level
  Scenario: The full sweep stitches pages, drops pull requests and keeps the owner
    Given two pages of issues where one item is a pull request, and a comment "anchors-owner: alice" on card 1
    When the full sweep runs
    Then the board has cards 1 and 2 only
    And card 1 is owned by alice
    And the issues were read through the REST route

  @BRSRB-B06 @unit-level
  Scenario: The comments of open cards follow the cursor and failures yield none
    Given a host answering two pages of comments linked by a cursor
    When the comments are read
    Then both pages are read in two calls
    And and a refusal or a repository name without a slash yields no comments

  @BRSRB-B07 @unit-level
  Scenario: The repository comes from the clone and the banner says what is served
    Given a clone of acme/app
    When board serve runs with --repo acme/app on an invalid port
    Then it prints "repository: acme/app" and the address before failing to listen

  @BRSRB-B08 @unit-level
  Scenario: The server hands out the published page and the live JSON
    Given board serve running for acme/app
    When the root and /board.json are requested
    Then the root returns the published board page as text/html
    And /board.json returns status 200, application/json, no-store, live and the two cards

  @BRSRB-I01 @unit-level
  Scenario: The incremental baseline is the newest update of the served board
    Given items updated 2026-01-02, 2026-03-01 and 2026-02-01
    When the baseline is taken
    Then it is 2026-03-01
    And and an unreadable board has no baseline

  @BRSRB-X01 @unit-level
  Scenario: The page is the pipeline's own board page
    Given board serve running
    When the root is requested
    Then the body is byte for byte the published board page

  @BRSRB-E01 @unit-level
  Scenario: A refused sweep carries the host's message
    Given a host answering "API rate limit exceeded"
    When the full sweep runs
    Then it fails with "collect from GitHub" and the host's message

  @BRSRB-E02 @unit-level
  Scenario: With nothing to serve the data route answers 502
    Given board serve running with a host that refuses everything and no previous read
    When /board.json is requested
    Then the status is 502 and the body is a JSON error holding the host's message

  @BRSRB-E03 @unit-level
  Scenario: Without a repository the command points at --repo
    Given a directory where the clone cannot name its repository
    When board serve runs with no --repo
    Then it fails pointing at --repo <owner/name>
