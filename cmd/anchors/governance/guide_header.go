// @anchors
//   code: GHCGD
//   ref: GVGDG

package governance

// headerGuide é a régua UNIVERSAL e MANDATÓRIA do bloco de cabeçalho de arquivo — o
// padrão de comentário no topo de todo artefato do projeto, que carrega as marcações
// que o Anchors lê: identidade, carimbo de alteração, tags de agrupamento e opt-outs.
// É o guide transversal: rege TODOS os arquivos, não um tipo. O `anchors init` semeia
// um HEADER_GUIDE.md concreto no projeto a partir desta régua.
// flagAt is the flags' at sign, spliced into the guide's examples: written whole in this
// source, each example would be read as a real flag of this file.
const flagAt = "@"

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
  //   code: LGNSC          (the file's OWN code — its address)
  //   ref: LGNNA           (the SCREEN realizes the LGNNA spec — the unit it belongs to)
  //   updated_at: 2026-08-08
  //   layer: screen
  //   @feature: auth
  //   @noPropagation

Example (markdown, in the owning SPEC) — the spec OWNS the code:

  <!-- @anchors
    code: LGNNA           (the spec's own code IS the unit's code)
    updated_at: 2026-08-08
    layer: screen
    @feature: auth
  -->

The block opens with '@anchors' and the following lines are the markers. What does not
apply, omit — but every file in the graph needs a CODE OF ITS OWN ('code:'), and the
identity of where it belongs: the spec's own code is its unit's; the rest of the unit adds
'ref:' to it; a file of a RECOGNIZED layer with no spec (infra/dao/presentation/domain
vocabulary — see below) adds 'layer:'.

## The markers

### Identity: 'code:' (the file's own) and 'ref:' (its unit)

- 'code: <CODE>' — the FILE'S OWN code, one value, unique in the project: the address by
  which anything else names this file (the dependency chain does). Every file has one. The
  SPEC's own code is also its UNIT's code: the rules and scenarios are numbered from it
  (DIVID-B01). 'anchors check --fix' generates a missing one from the file's name and type;
  'anchors code <name>' gives one by hand.
- 'ref: <CODE> [, <CODE>...]' — the UNIT(S) the file realizes, covers or proves. It can be
  MULTIPLE. Thus, beside its own code:
    - the CODE (Divider.tsx) that realizes the spec  → 'ref: DIVID'
    - the FEATURE that covers the spec's scenarios   → 'ref: DIVID'
    - the TEST that proves the scenarios             → 'ref: DIVID'
  And a file may reference SEVERAL units: a util tested by scenarios of two
  screens → 'ref: TXDTS, MNDTS'; a screen that composes components → 'ref' with their codes.

Why it matters: a test whose 'code:' were the spec's would be a second file at the spec's
address — gate 'header-valid' fails a code two files carry. The unit is said by 'ref:'; the
file is named by its own 'code:'.

### RECOGNIZED layers: identity by 'layer:' (and 'dep:' for dependencies)

Not every layer has a spec. The Structure may declare RECOGNIZED layers — infra
(pure helpers), data access (dao), presentation (domain→UI map), domain
vocabulary (enums/catalogs). They exist to LEAVE THE SCRUTINY of a spec (there is no rule
of their own to document), but they still need IDENTITY. Since they have no owning spec nor
a sibling to reference, their minimal honest identity is the layer itself:

- 'layer: <layer>' — for a file of a recognized layer, 'layer:' is where it belongs,
  beside its own 'code:'. Do not force a 'ref:' (it realizes no spec). It honestly declares
  "I belong to this layer".
- 'dep: <file> [, <file>...]' — since they have no spec, they have no Dependency
  Table (SPEC_TYPES §5). So they declare in their OWN header the FILES they
  depend on (the file path, not a code — the target may be another recognized layer with no
  code). Each one becomes a 'depends-on' edge. That is how one recognized layer references
  another, or points at the governed file that consumes it. E.g.: a presentation map that takes colours from
  'theme/tokens.ts' → 'dep: theme/tokens.ts'.

GOVERNED layers (business-logic, validation, hook, store, repository, service, and the spec)
still require 'ref:' to their unit (the spec, its own 'code:') — they have a rule to
document, so they have a spec. 'layer:' alone is NOT enough for a governed one.

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

## Flags beside the code: the dependency and navigation chain

Outside the header, on the line they speak of, six flags tie files to files by their own
'code:' — never by path, so a rename breaks nothing. Each is a comment in the file's dialect;
'anchors check --fix' writes the ones the code proves, and the gates of the chain confront
them with the code (they are named in each line below).

- '@dep: <CODE>' — at the end of each import line of a governed file: this import ties this
  file to that one; the symbols are the import's own. 'dep-declared' asks one on every
  import that resolves to a file of the project; 'dep-honored' asks it to name the code of
  that file. A package of the ecosystem takes none.
    import { PALETTE } from './tokens' // ` + flagAt + `dep: TOKNS
- '@dep[type]: <CODE>' — instead of '@dep:' on an import that brings types only: nothing of
  that file runs with this one, so a change to it stales no test's proof ('anchors check
  --fix' writes it on 'import type …'; the project's dialect says how a type import reads
  in other languages, 'type_import_pattern').
    import type { RootStackParamList } from '@/navigation/RootNavigator' // ` + flagAt + `dep[type]: RNIRT
- '@dep[<kind>]: <name>' — what the code reaches that is no import (a table of the database,
  an external API, a queue), on the line that calls it. The kinds are the project's
  ('dependency_kinds:' in anchors.yaml); 'anchors map deps --kind <kind>' lists who reaches each.
    await db.query(sql) // ` + flagAt + `dep[db]: transactions
- '` + flagAt + `no-dep: <reason>' — on the import line instead: this import stays out of the chain
  (a test double…).
- What has no side effect on the tests ties no test's proof: a line that is only a comment
  (a rule citation, a note), and the flags '@navigates' and '@dep[type]' by default —
  'evidence.no_side_effect.flags' in anchors.yaml lists them, and
  'evidence.no_side_effect.sections' the spec sections whose change proves nothing.
- '@used-by: <CODE>[, <CODE>...]' — in the comment right above each symbol another file
  imports: exactly who imports it ('used-by-declared'). A default import is the module's
  'default', flagged above 'export default' or 'module.exports ='. Above an export list,
  which declares several names on one line, the flag names its symbol after the codes:
    // ` + flagAt + `used-by: ARNAA, WLLTW
    export const PALETTE = {}
    // ` + flagAt + `used-by: DRSSD (DeletionStatus)
    export type { DeletionScope, DeletionStatus }
- '@navigates: <CODE>[, <CODE>...] [<RULE>]' — on a navigation call's line, or the line
  before: this call leads to that screen (several, for a back navigation that returns to
  more than one), through that rule of the screen's Out table. Every call counts — navigate,
  push, replace, goBack, reset ('nav-annotated', 'nav-matches-spec', 'nav-symmetric',
  'nav-reachable'; 'anchors guide navigation' has the whole of it).
    navigation.navigate('GoalDetail') // ` + flagAt + `navigates: GLDTG [HOMEH-A01]
- '` + flagAt + `no-nav: <reason>' — on the call's line: this call is no screen edge (it closes a modal of
  the same screen).

The header 'dep:' above remains for a dependency no import line carries; where one does, the
flag on the import line is the declaration — removing or changing the import shows at once
which flag goes with it. 'anchors map deps <CODE>' walks the chain as a tree; 'anchors map
nav' prints the screens' edges.

## Header rules

- ALWAYS at the top of the file (before the code), so it is the first thing read. Only
  blank lines, comments and a shebang may come before it: a block below the code is read as
  text, not as the header (a comment may precede a directive such as 'use client', so the
  header goes above it). When the language truly demands something before it, declare why
  inside the block: '@fixed-header: <why>'.
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
