# language: en
# @anchors
#   code: LCFTL
#   ref: MPLCK
#   updated_at: 2026-10-03
#   layer: feature

@MPLCK
Feature: MapLock — the map changed by one writer at a time, each applying only what it changes

  @MPLCK-B01 @unit-level
  Scenario: The lock is a file beside the map with its owner
    Given a map path
    When the lock is taken and released
    Then the lock file held this process's pid and host, and is gone after the release

  @MPLCK-B02 @unit-level
  Scenario: A second writer waits for the first
    Given a writer holding the lock
    When another writer asks for it and the first releases it
    Then the second takes it only after the release

  @MPLCK-B03 @unit-level
  Scenario: An abandoned lock is taken over, a live one is not
    Given a lock of a dead process of this host, an old lock of another host, and a fresh lock of another host
    When a writer asks for the lock
    Then the first two are taken over and the third is waited for

  @MPLCK-B04 @unit-level
  Scenario: Parallel processes each changing their part all reach the map
    Given four processes, each recording twenty signals of its own on the same map
    When they run at the same time
    Then the map holds all eighty

  @MPLCK-B05 @unit-level
  Scenario: The lock is released when the function returns
    Given a function run under the lock that fails
    When it returns
    Then the error comes back and the lock is free

  @MPLCK-E01 @unit-level
  Scenario: A lock held past the timeout
    Given a lock held by a live process
    When a writer waits past the timeout
    Then it fails saying another process is writing the map, and the lock stays

  @MPLCK-E02 @unit-level
  Scenario: No map, or a refused change, writes nothing
    Given no map on disk, and a map with a change that refuses
    When each is updated
    Then the error comes back, nothing is written, and the lock is free

  @MPLCK-E03 @unit-level
  Scenario: A lock that cannot be created names the file
    Given a map path in a directory that does not exist
    When the lock is taken
    Then the error names the lock file
