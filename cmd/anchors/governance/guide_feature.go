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

## Classify each scenario on two axes

1. TEST LEVEL — where it will be proven in the cheapest, most isolated way:
   • unit (pure logic, no environment)
   • integration (one isolated component/module, with its edges mocked)
   • end-to-end (the real flow crossing the system boundaries)
   The levels are not exclusive; choose the NATURAL one for that behaviour.
   Anti-E2E signal: if the observable is internal state, animation or scroll, it is NOT E2E —
   forcing it produces a fragile or empty test.
2. PRIORITY — exactly ONE per scenario (from critical to low), by the rule of the LARGEST
   axis: business impact, kind of behaviour, frequency of use, regression
   risk. Hard rule: a critical scenario ALSO REQUIRES an end-to-end test, besides
   the cheap test.

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
