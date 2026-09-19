package governance

// testGuide é a régua universal do artefato TEST: a prova executável dos cenários.
// Agnóstico de framework (Jest, pytest, go test, …).
const testGuide = `# Test guide (the executable ruler)

The test is the executable proof of the feature's scenarios. It closes the
traceability chain: the same code crosses spec → feature → test. A test that passes
by mistake is worse than no test — it lies about the system's health.

## The pyramid

Organize the tests by cost and scope, from cheap to expensive:
- UNIT — pure logic, no environment. Fast, runs anywhere.
- INTEGRATION — one module with its edges (infrastructure) mocked; the domain runs
  for real.
- END-TO-END — the real flow crossing the system boundaries.
Each feature scenario already arrived classified at its natural level — prove it there.

## The golden rules (what separates a real test from theatre)

- COVERS THE HAPPY PATH + AT LEAST ONE ERROR. Every tested behaviour proves the main
  path AND at least one failure case.
- TEST THE REAL PRODUCTION CODE. An inline mock that REPLICATES the logic inside
  the test is forbidden — it does not exercise the real code and produces zero coverage. Writing and
  reading through the same low-level path in the same test is forbidden — that tests the library, not
  your code.
- MOCK THE INFRASTRUCTURE AT THE EDGE, not the internal units. Let the domain execute
  for real; mock only what leaves the process (network, disk, external service). Mocking the
  whole internal unit zeroes the coverage that matters.
- ISOLATION AND CLEANUP. Test data carries an identifiable prefix; clean up at the end (and
  keep a global safety net). Never leave permanent test data.
- CENTRALIZE THE SETUP. Do not create data inside each case when there is shared
  state — centralize the seed; it avoids races and contention.
- ASYNCHRONY/EVENTUAL CONSISTENCY MADE EXPLICIT. If the system is eventually
  consistent, wait on a predicate (retry/poll), do not read right after writing.

## Environment safety (when the test touches external state)

- DISCOVER THE ENVIRONMENT BY IDENTITY, NOT BY NAME. Validate that the resources belong
  to the test environment through an unambiguous mark (tag/attribute), never inferring from the
  name. Fail BEFORE writing if there is ambiguity — so it never runs against production.
- DO NOT EDIT THE ENVIRONMENT FILE BY HAND. Generate it by script, reproducibly.

## Honesty about green (a cultural rule)

Distinguish the suites that run without credentials from those that require an external environment. NEVER
claim "N/N green" without that caveat. "What could pass here passed" is an honest
sentence; "all green" when half of it did not even run is a lie that costs dearly
later. (The project's golden rule holds: absence of proof is not proof of absence.)

## Anti-patterns (refuse them)

- A mock that replicates the target's logic → you tested the mock, not the code.
- Writing and reading through the same low-level path → you tested the library.
- Mocking the whole internal unit → zero coverage where it matters.
- Reading right after writing in an eventual system → intermittent (flaky) test.
- Inferring the environment by name → risk of running against production.
- "All green" without saying what did not run → a dishonest report.

## Project specialization

The concrete framework, the helpers, the folder structure and the PR gates belong to your
project — see the test guide in the 'guide' layer of anchors.yaml. This guide is the
universal doctrine; follow the project's dialect when it exists, and warn if it does not.
`
