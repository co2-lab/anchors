package governance

// headerGuide é a régua UNIVERSAL e MANDATÓRIA do bloco de cabeçalho de arquivo — o
// padrão de comentário no topo de todo artefato do projeto, que carrega as marcações
// que o Anchors lê: identidade, carimbo de alteração, tags de agrupamento e opt-outs.
// É o guide transversal: rege TODOS os arquivos, não um tipo. O `anchors init` semeia
// um HEADER_GUIDE.md concreto no projeto a partir desta régua.
const headerGuide = `# Header guide (the block of markers at the top of each file)

Every file that takes part in the graph carries, at the top, a standardized HEADER
BLOCK — a comment that Anchors reads to know the file's identity,
when it changed, which groups it belongs to, and which opt-outs it declares. It is MANDATORY:
a file without a conforming header is invisible to what Anchors does best.

This guide is CROSS-CUTTING — it governs every file, of any layer. It is the only
ruler that does not ask "what kind of artifact?"; it asks "does every file have its
label?".

## The form: key-value for data, @tag for flags

The block is a comment (in the dialect of the file's language: '//' in TS/Go, '#' in
Python/shell, '<!-- -->' in markdown). Inside it:

- 'key: value' LINES — the DATA (identity, dates, grouping with a value).
- '@flag' LINES — the boolean FLAGS and opt-outs (presence = on).

Example (TS/Go) — a screen:

  // @anchors
  //   ref: LGNN            (the SCREEN realizes the LGNN spec — it references, it does not own)
  //   updated_at: 2026-08-08
  //   layer: screen
  //   @feature: auth
  //   @noPropagation

Example (markdown, in the owning SPEC) — the spec OWNS the code:

  <!-- @anchors
    code: LGNN            (the spec IS the owner of the identity)
    updated_at: 2026-08-08
    layer: screen
    @feature: auth
  -->

The block opens with '@anchors' and the following lines are the markers. What does not
apply, omit — but every file in the graph needs IDENTITY: 'code:' if it is the owner
(the spec), 'ref:' if it references (the rest of the triad), OR 'layer:' if it belongs to a
RECOGNIZED layer with no spec (infra/dao/presentation/domain vocabulary — see below).

## The markers

### Identity: 'code:' (ownership) vs 'ref:' (reference)

This is the distinction that avoids the most common confusion. A scenario code belongs to
ONE unit; the files that revolve around it either OWN it or REFERENCE it:

- 'code: <CODE>' — OWNERSHIP. The file IS the canonical owner of the identity. Only ONE value.
  Who the owner is: the SPEC (the source of truth, SPEC.md). The spec DEFINES the requirement and
  its scenarios; it owns the code. Generate it with 'anchors code <name>' (uniqueness).
- 'ref: <CODE> [, <CODE>...]' — REFERENCE. The file is NOT the owner; it realizes,
  covers or proves the owning unit(s). It can be MULTIPLE. Thus:
    - the CODE (Divider.tsx) that realizes the spec  → 'ref: DIVI'
    - the FEATURE that covers the spec's scenarios   → 'ref: DIVI'
    - the TEST that proves the scenarios             → 'ref: DIVI'
  And a file may reference SEVERAL units: a util tested by scenarios of two
  screens → 'ref: TXDT, MNDT'; a screen that composes components → 'ref' with their codes.

Why it matters: putting 'code:' in a test would say the test OWNS the code — but it only
references it (it proves the spec's unit). Swapping ownership for reference crosses the
traceability and produces a false identity collision. When in doubt: the spec has 'code:'; the
rest of the triad has 'ref:'.

### RECOGNIZED layers: identity by 'layer:' (and 'dep:' for dependencies)

Not every layer has a spec. The Structure may declare RECOGNIZED layers — infra
(pure helpers), data access (dao), presentation (domain→UI map), domain
vocabulary (enums/catalogs). They exist to LEAVE THE SCRUTINY of a spec (there is no rule
of their own to document), but they still need IDENTITY. Since they have no owning spec nor
a sibling to reference, their minimal honest identity is the layer itself:

- 'layer: <layer>' — for a file of a recognized layer, 'layer:' IS the identity
  (it satisfies the gate). Do not invent a 'code:' (it owns no spec) nor a forced 'ref:'
  (it realizes no spec). It honestly declares "I belong to this layer".
- 'dep: <file> [, <file>...]' — since they have no spec, they have no Dependency
  Table (SPEC_TYPES §5). So they declare in their OWN header the FILES they
  depend on (the file path, not a code — the target may be another recognized layer with no
  code). Each one becomes a 'depends-on' edge. That is how one recognized layer references
  another, or points at the governed file that consumes it. E.g.: a presentation map that takes colours from
  'theme/tokens.ts' → 'dep: theme/tokens.ts'.

GOVERNED layers (business-logic, validation, hook, store, repository, service, and the spec)
still require 'code:'/'ref:' — they have a rule to document, so they have a spec and a
code identity. 'layer:' alone is NOT enough for a governed one.

### Other data (key: value)

- 'updated_at: <YYYY-MM-DD>' — the CHANGE STAMP: when the content last
  changed. Do NOT maintain it by hand (it will lie) — Anchors fills/validates it against
  the real history (git). Declaring it here is optional; the true value comes from the graph.

### Grouping tags (key: value OR @tag)

Cross-cutting labels that categorize the file, for queries and gate scope:

- 'layer: <layer>' — the Structure layer it belongs to (screen, component,
  service, model…). Normally Anchors INFERS this from the path (the layer pattern);
  declare it only if you want to override.
- '@feature: <name>' — a free GROUPING label for the vertical module (auth, dashboard,
  budgets…), alongside '@experimental' and '@legacy'. It says what the file belongs to,
  and nothing reads it: no gate, no edge, no confrontation.
  The vertical axis that IS confronted is product doctrine — 'product/<name>.doctrine.md',
  which the spec points at with '@realizes'. If what you want is for a rule spanning
  several units to have one place and be verified, that is the axis: run
  'anchors guide product'. The tag groups; the doctrine decides.
  What the tag is and is not is written down, with the reasons: see the vertical-grouping
  doctrine ('VGRUP'). The mechanism grows toward FILTER and VIEW ('--feature'), never
  toward a gate that demands the label — a label that becomes an obligation gets filled in
  to silence the charge, and the grouping fills with files nobody classified.
- '@<tag>' — any other free grouping tag (e.g. @experimental, @legacy).

### Opt-outs (only @flag — presence = on)

Honest waivers of an Anchors rule, recorded in the file itself (dated
by git, localized). Anchors ACTUALLY reads these today:

- '@noPropagation' — this child does NOT depend on the parent: the propagation wave does not pass
  through it. See PROPAGATION.
- '@anchors-shared-code' — the scenario codes here belong to ANOTHER unit on
  purpose (a util/handler test that proves the scenario of the unit it serves); they do not
  count as an identity collision. See TRACEABILITY.

### Opt-outs WITH A REASON ('@flag: <reason>')

These are not booleans: they require the reason written after the colon, on the SAME line.
A bare marker ('@no-scenario:' with no text) waives NOTHING — it keeps failing. The
difference exists because these waive a rule about a specific line, and without the
why the waiver becomes a hole with a pretty name: whoever reads it later cannot tell whether it was a decision
or an oversight.

They go on the line of what they waive (or in the comment immediately above, where the line does not
fit a readable comment):

- '@no-scenario: <reason>' — on the line of a spec REQUIREMENT: this requirement will have no
  scenario in the feature. For what is truly non-observable by scenario (a
  structural restriction proven by another instrument). Gate 'spec-feature-match'.
- '@no-paginate: <reason>' — on a function that promises a set and does not paginate: the limit is
  deliberate and known (e.g. a table with dozens of rows fixed in a migration). Gate
  'pagination-honored'.
- '@allow-boundary: <reason>' — on a line that crosses a layer boundary declared in
  'boundaries:'. For acknowledged debt, which stays visible and dated in the code instead of
  in a distant exception list. Gate 'layer-boundary'.

An opt-out WAIVES the rule, never the record — it is the legitimate door, the opposite of the
silent coverage hole.

## Header rules

- ALWAYS at the top of the file (before the code), so it is the first thing read.
- code is the only essential marker; the rest is as needed.
- updated_at is NOT maintained by hand — let Anchors/git take care of it; writing it wrong is
  worse than omitting it (an anchor lying about itself).
- An opt-out always with a WHY beside it (a comment on the same line or below): whoever
  reads it later needs to know why the rule was waived.
- The comment dialect is the file's — Anchors only reads TEXT and recognizes the
  markers in any comment; do not invent syntax outside the comment.

## Anti-patterns (refuse them)

- A file without a code → an orphan invisible to Traceability; give it identity.
- updated_at maintained by hand → it will go stale and lie; let git be the truth.
- An opt-out without a why → nobody knows whether it still holds; explain it beside.
- A marker outside a comment (in executable code) → Anchors reads comments; and
  it would pollute the code.

## The project specializes

'anchors init' seeds a concrete HEADER_GUIDE.md for your project — with your
stack's comment dialect and the real tags/features. This is the universal floor;
the project's adapts it. If your project has no header guide, generate it via init or
copy this pattern — and warn that the standardization is a loose end.
`
