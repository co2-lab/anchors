# language: en
# @anchors
#   code: LDQUB
#   ref: LDQUS
#   updated_at: 2026-11-19
#   layer: feature

@LDQUS
Feature: LoadQueue — the machine's runs admitted by what they cost

  @LDQUS-B01 @unit-level
  Scenario: A run starts when its cost fits in CPU and memory beside the leases; nothing running admits any cost
    Given 8 CPU and 16 GB, and a lease of 4 CPU and 6 GB
    When runs of 4 CPU/6 GB, 5 CPU/1 GB and 1 CPU/11 GB ask, and a run of 32 CPU with nothing running
    Then only the first fits, the costlier one runs alone, and an unread memory bounds nothing

  @LDQUS-B02 @unit-level
  Scenario: A run that waited past its limit stops the admission of every run that came after it, until it fits
    Given a heavy run waiting 10 minutes with a limit of 5, and a light run after it
    When the light run asks, then with the heavy within its limit, then the heavy alone
    Then the light waits, then passes, and the heavy starts once it fits

  @LDQUS-B03 @unit-level
  Scenario: The queue admits a project that declares load:, unless $ANCHORS_LOAD turns it off or the run is in CI and load.ci is not set
    Given no load:, a load:, ANCHORS_LOAD off, and GITHUB_ACTIONS true with and without ci
    When the queue is asked whether it admits
    Then only the declared project outside CI, or in CI with ci true, is admitted through it

  @LDQUS-B04 @unit-level
  Scenario: The capacity is the machine's CPUs and memory minus a reserve, or what $ANCHORS_LOAD_CAPACITY says
    Given the machine as it is, then ANCHORS_LOAD_CAPACITY cpu=3,memory=4GB, then cores=3
    When the capacity is read
    Then it is below the machine's memory, then 3 CPU and 4 GB, and an unknown part is refused

  @LDQUS-B05 @unit-level
  Scenario: Runs share the machine's leases: a run that does not fit waits, saying for what, until a lease is released
    Given 4 CPU and 8 GB, and app's test unit holding 3 CPU and 6 GB
    When web's mutation unit asks for 2 CPU and 2 GB, and then the first lease is released
    Then it waits saying it is held by app test unit, is recorded as waiting, and starts once released

  @LDQUS-B06 @unit-level
  Scenario: A lease whose process died is reclaimed
    Given a lease of the whole machine held by a process that does not exist
    When a run asks for the whole machine
    Then it starts
