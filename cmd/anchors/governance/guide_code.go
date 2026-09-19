package governance

// codeGuide é a régua universal do artefato CÓDIGO: como implementar guiado pela
// spec, com arquitetura que não apodrece. Agnóstico de linguagem/framework — a
// doutrina vem da experiência real (organização por eixo, dependência unidirecional,
// promoção sob demanda). O projeto especializa no seu guide de arquitetura.
const codeGuide = `# Code guide (implementing guided by the spec)

The code implements the behaviour the spec described. It is confronted against the
spec (the spec is the ruler; the code, the governed). Here the doctrine is not about syntax —
it is about WHERE each thing lives, so the project does not turn into a tangle over time.

## Two axes of organization, combined

All code organizes itself in two dimensions at the same time:
- BY TYPE/COMPOSITION (horizontal) — what is reusable and does NOT know the domain
  (generic blocks, utilities).
- BY DOMAIN/FEATURE (vertical) — slices isolated by business area, where the logic
  of that domain lives.
Business logic has no natural place on the by-type axis — if you only organize by
type, the logic spreads out. Hence: the domain organizes what has logic; the type organizes
what is reusable without a domain.

## The single classification rule

Reusable and domain-free → lives in the by-type (shared) layer.
Coupled to a domain → lives in that domain's slice.
Decisive test: "would this piece make sense in ANOTHER app, without this domain?"
Yes → shared. No → feature.

## Unidirectional dependency

- A module only imports from layers BELOW it. Never sideways, never upwards.
- A feature NEVER imports from another feature.
- Infra/shared NEVER import from a feature. The shared layer "rises"; it does not
  know who consumes it.

## Promotion on demand (DRY at the right time)

Something is born LOCAL to a feature. It only RISES to the shared layer when 2+ features
consume it — not before. Premature abstraction (rising "just in case") is debt, not
economy. Wait for the second consumer.

## Data boundaries

- Only the EDGES fetch data (the entries of the flow). The rest RECEIVES it by parameter.
  A block that receives its data by parameter is testable in isolation and reusable.
- Separate the PURE LOGIC (framework-free) into its own place. Pure domain functions,
  with no UI/IO dependency, are the easiest to test — and the most reused.

## The co-located triad

Every code artifact is born with its spec, its feature and its test beside it. Writing
code without all three leaves the unit with no ruler, no coverage and no proof. (It is
the core of Anchors — the map expects the triad.)

## Real violation vs. sanctioned exception

Not every broken rule is a bug. Before blindly "fixing" a dependency that
looks crooked, ask whether it is a sanctioned EXCEPTION (documented, with a reason) or a
real VIOLATION (a symptom of wrong coupling). Fix the violation; respect the
exception — and if it is not documented, document why.

## Conventions (in spirit, agnostic)

- The file name matches what it exports.
- Avoid re-exports/barrels that create import cycles.
- Zero magic values — use named constants/tokens.
- Import through a stable alias, not through fragile relative paths.

## Anti-patterns (refuse them)

- Business logic in the by-type layer → move it to the domain slice.
- A feature importing from another feature → promote the common part to the shared layer.
- Abstracting on first use → wait for the second consumer.
- A reusable block that fetches its own data → pass the data by parameter.
- Code without the triad (spec/feature/test) → the unit is orphaned in the map.

## Project specialization

The layer names, the stack, the tokens and the concrete folder structure belong to your
project — see the architecture/code guide in the 'guide' layer of anchors.yaml. This
guide is the universal doctrine; follow the project's dialect when it exists, and warn if it
does not.
`
