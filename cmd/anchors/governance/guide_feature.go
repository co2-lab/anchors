package governance

// featureGuide é a régua universal do artefato FEATURE: os cenários de comportamento
// que traduzem a spec em histórias verificáveis. Agnóstico de Gherkin/ferramenta.
const featureGuide = `# Feature guide (the behaviour scenarios)

The feature translates the spec into SCENARIOS: concrete stories of what the unit does. It is
the bridge between the spec (what must be) and the test (the executable proof). Each scenario
reuses a CODE from the spec — it never invents a new one.

## The shape of a scenario

Every scenario has three tenses, without mixing them:
- CONTEXT (past) — the state before; all the preconditions, right here.
- TRIGGER (present) — the single event that fires the behaviour.
- RESULT (future) — what is expected to be observed afterwards.
Do not put an action in the context, nor a verification in the trigger, nor a new action in the result.

## Feature rules

- SELF-CONTAINED. Each scenario declares its own preconditions. "Shared context" at the
  top of the file is forbidden — the reader should not hunt for preconditions far from the
  scenario. A scenario is a complete story.
- EXPLICIT NAMES. "the signup screen", "the Save button" — never "the screen", "the button".
  Business language, not implementation language.
- COVERS THE SPEC'S CATALOG. For EACH item the spec catalogued — each state,
  validation, action, behaviour, permission, restriction, message, and each field of the
  data contract/state — there is at least one scenario. That is how the feature
  proves the spec was honoured.
- TRACEABLE CODE. Each scenario carries the code of the spec item it covers.
  A scenario without a code is invisible to the coverage gate.
- LITERAL COPY. A message scenario verifies the EXACT TEXT from the spec's catalog
  (real accents and punctuation). For interpolated text, verify the stable part.
- DATES AND DEFAULTS HAVE A SCENARIO. A date/timezone gets a dedicated scenario ("I must not see
  the previous day"); a default gets an absence scenario; a conditional field gets a
  present AND an absent scenario.

## One scenario per catalogued item

The spec catalogues items, and each item carries a LETTER that states its nature — state,
validation, action, restriction, message, navigation, behaviour. That letter is the bridge:
the scenario reuses the item's code, so the letter says what KIND of scenario it is. There
is no separate list of scenario types to memorise — the types ARE the spec's sections.

So the count is not a matter of taste: for each catalogued item there is at least one
scenario, and a section with five rules produces at least five. A feature with three
scenarios for a spec with twelve items is not concise — it is incomplete, and the coverage
gate says which ones are missing.

## Classify each scenario by TEST LEVEL

Where the behaviour will be proven in the cheapest, most isolated way. Ask in this order
and stop at the first yes:

1. Is it a pure rule or function, with no environment to stand up? → UNIT.
2. Is it the internal behaviour of ONE component or module — internal state, local
   interaction, its edges mocked — that crosses no boundary? → INTEGRATION.
3. Is it a flow that crosses the system's boundaries — navigation, persistence, a real
   external dependency? → END-TO-END.

The levels are not exclusive: a rule proved by a cheap test and also exercised inside a
real flow declares both. But the NATURAL level is always declared.

Anti-end-to-end signal: if the observable is internal state, an animation or a scroll, it
is NOT end-to-end — forcing it there produces a fragile test, or an empty one that passes
without proving anything.

The tag that names each level is the PROJECT's, declared in the map of regimes. Reuse
exactly the tags it declares — a spelling the project did not declare makes the scenario
confronted by no gate.

## The unit with no interface (traceability, not BDD)

A backend rule, a gate, a repository, a pure function: the feature of a unit with no
interface is a TRACEABILITY document. It does not run as BDD — the executable proof is the
test beside it — and it exists so the catalogued item, the scenario and the test are bound
by one code.

What changes for this profile:

- The context is the INPUT — the received value, the state of the data, the config in
  force. Not a screen state, not navigation.
- It is never end-to-end: a unit with no interface crosses no screen. Forcing it there
  produces the empty test the signal above describes.
- One scenario per catalogued item, same rule as always; what does not exist here is the
  visual-regression scenario, because there is nothing to capture.

Everything else holds: self-contained, explicit names, the item's code, and the three
tenses without mixing them.

## The feature↔test contract

The feature and the executable test are ONE contract. Each step of the scenario appears
mirrored in the test; if you edit a feature step without updating the test, the gate
must break. Do not let the two diverge silently.

## Anti-patterns (refuse them)

- A precondition at the top of the file (shared context) → move it inside the
  scenario.
- A scenario without a code → invisible to the gate; reuse the spec's code.
- "the screen", "the button" → name what, explicitly.
- Forcing E2E over internal state → fragile test; use the natural level.
- Paraphrased copy → verify the literal text from the catalog.

## Project specialization

The concrete format (the scenario syntax, the names of the level/priority/
category tags, the execution tool) belongs to your project — see the feature guide in the
'guide' layer of anchors.yaml. This guide is the doctrine; follow the project's dialect if
it exists, and warn if it does not.
`
