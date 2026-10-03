// @anchors
//   ref: GVGDG

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
- VISUAL REGRESSION — for screens and components: each state captured and compared with
  its baseline image (see below).
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

## Proving a boundary: the instrument follows the SHAPE of the input space

A boundary rule says what may NOT come out, or what alone is accepted: "the screen
asserts no state of X", "no credential crosses", "only these fields". The rule applies to
EVERY input, and a test that checks ONE representative leaves the neighbour open. The
instrument is chosen by the shape of the input space, not by taste:

- SMALL AND CLOSED (a handful of values: the verdicts, the states, the environments) →
  EXHAUSTIVE. ` + "`it.each`" + ` over the whole set, read from the source of truth (the exported
  list, the type's members), never a copy typed in the test. A value added later is
  tested without anyone remembering to.
- LARGE BUT STRUCTURED (numbers, strings, nested objects built from a few kinds of part)
  → A TABLE OF CLASSES. One case per class that can behave differently (empty, zero,
  null, at the limit, one past it, each variant of the union), each named for the class.
- OPEN (anything a source may send, any text) → CLOSE THE OUTPUT OR THE SOURCE, never a
  list of the forbidden. Assert the exact set of keys that comes out, or read the unit's
  own code for what it cannot reach (imports, module state). A blacklist of names only
  proves the names someone remembered.

Measured, in a real project: nine surviving mutants across four units, all the same
shape — the boundary closed for one case, the neighbour left open.

## Doubles of a module: the stamp, and who refreshes it

When a test doubles a module of the project (` + "`jest.mock`" + `, ` + "`vi.mock`" + `…), the double
carries a stamp — ` + "`// @contract: <file> | <anchor line> | <lines> | <hash>`" + ` — that the
` + "`mock-stamped`" + ` gate recomputes against the real file. Write missing stamps with
` + "`anchors stamp`" + `; it never rewrites an existing one.

WHEN YOU CHANGE A FUNCTION THAT OTHERS DOUBLE, you own the doubles too:

    anchors stamp --refresh <file you changed>

It lists every double stamped against the previous version — test, line, member, and how
the stamped block changed from HEAD — and updates those stamps. Each double listed
reproduces the OLD contract: adjust it in the same commit — after the refresh the stamp
says the double matches, and you are the one who checked it. If you skip the refresh, the
pre-commit (` + "`check --changed`" + ` of the file) fails on every stale stamp. A member you renamed or removed
is not refreshed: adjust the double, delete its stamp, and run ` + "`anchors stamp`" + `.

The refresh rewrites the tests it lists, and each one then needs its ` + "`updated_at`" + ` bumped:
` + "`anchors touch`" + ` does it for every changed file (a ` + "`@contract`" + ` stamp never covers the header,
so the order of the two does not matter).

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

## Visual regression (screens and components)

Every visual unit — a screen, a component — is strongly recommended to have a VISUAL
REGRESSION test for EVERY state its spec registers. Which units are visual is the
anchors.yaml's to say: the code layers tagged ` + "`screen`" + ` or ` + "`component`" + `. A state with no visual
value of its own — a transient loading, a state that looks like another — is exempted
where it is declared, with the reason: ` + "`@no-vr: <reason>`" + ` on its heading or its row. An assertion proves what a state shows;
only a capture proves what it LOOKS like, and a style change that breaks a state passes
every assertion. Each state is tied to its proof both ways, and all of it is part of the
unit:

- the state has a VR scenario in the feature, tagged with the project's visual regime and
  the state's code;
- the VR scenario has a VR TEST — a capture flow or a screenshot test — that names
  ` + "`{CODE}-VR-<state>`" + ` (` + "`BUTTN-VR-S01`" + `) in its file name or its text, and a BASELINE IMAGE
  beside the unit: ` + "`<Unit>.{CODE}-VR-<state>.<ext>`" + `, an optional variant after the state, in
  png, jpg, webp or svg;
- the other way round, every VR scenario is of a state the spec registers, and every VR
  test is of a VR scenario the feature declares — a capture that outlived its state is
  removed with it.

The tool is the project's — a browser screenshot assertion, a device flow, a component
story runner; the naming is what lets the gates find it. Regenerate a state's baseline in
the same change that alters the state on purpose, and review the image as you review code:
an unreviewed baseline approves whatever it captured. Declare the baselines as a layer of
the unit (proof, like a test) so a commit made only of them is governed.

A capture's evidence goes stale when the unit's own file or one of its baselines changes —
` + "`evidence-fresh`" + ` and ` + "`anchors stale`" + ` say so, and ` + "`anchors test`" + ` runs it again. One level only: a
component the screen uses has its own capture, and its change stales that one, not the
screen's.

Four gates ask the four questions: ` + "`vr-states-covered`" + ` (every state has a VR scenario),
` + "`vr-scenarios-tested`" + ` (every VR scenario has a VR test and an image),
` + "`vr-scenarios-of-states`" + ` (every VR scenario is of a state) and ` + "`vr-tests-of-scenarios`" + `
(every VR test is of a scenario).

## Contract tests (APIs)

Every API unit is strongly recommended to have a CONTRACT TEST: one that runs requests
against the API and validates each answer against the project's OpenAPI document — the
one ` + "`anchors docs build`" + ` compiles from the specs, so the test proves the implementation keeps
what the specs promise. Every language has the tool: Schemathesis or Dredd for any stack,
kin-openapi in Go, jest-openapi or openapi-response-validator in JavaScript, openapi-core in
Python, swagger-request-validator with RestAssured in Java. The tool is the project's.

What makes it traceable, and what ` + "`contract-tested`" + ` asks of each API unit:

- the scenario ` + "`{CODE}-CT`" + ` in the unit's feature, tagged with the project's contract regime;
- a test of the unit that names ` + "`{CODE}-CT`" + ` (its case or its file) and LOADS the OpenAPI
  document — a schema copied into the test is a second contract that drifts;
- each refusal of ` + "`Error Responses`" + ` asserted with its status, its error code and its
  message, which an OpenAPI validator does not check by itself.

A contract test's evidence goes stale when the unit's code, its spec or the OpenAPI
document changes — ` + "`evidence-fresh`" + ` says so, and ` + "`anchors test`" + ` runs it again.

## Anti-patterns (refuse them)

- A mock that replicates the target's logic → you tested the mock, not the code.
- Writing and reading through the same low-level path → you tested the library.
- Mocking the whole internal unit → zero coverage where it matters.
- Reading right after writing in an eventual system → intermittent (flaky) test.
- Inferring the environment by name → risk of running against production.
- "All green" without saying what did not run → a dishonest report.
- A boundary proven by one representative, or by a list of forbidden names → the
  neighbour case is still open.

## Project specialization

The concrete framework, the helpers, the folder structure and the PR gates belong to your
project — see the test guide in the 'guide' layer of anchors.yaml. This guide is the
universal doctrine; follow the project's dialect when it exists, and warn if it does not.
`
