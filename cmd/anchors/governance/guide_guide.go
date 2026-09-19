package governance

// guideGuide é a régua de como escrever um GUIDE — o meta-guide. Um guide é a régua
// de um tipo de artefato; para ser CONFRONTÁVEL (e não só lido), ele precisa destilar
// suas regras em PONTOS DE CONFORMIDADE verificáveis. Sem eles, o gate de julgamento
// recai em heurística vaga ("respeita o espírito?"); com eles, a IA julga item a item.
const guideGuide = `# Guide guide (the ruler of how to write a ruler)

A guide is an ANCHOR that governs a type of artifact — it says how the artifact must be,
and it is the ruler against which the artifact is confronted. A guide nobody confronts is
just a document; a guide that governs but does not say WHAT TO VERIFY can only be confronted
by vague heuristic. This guide exists so that your guides are CONFRONTABLE.

## The central rule: distil into CONFORMANCE POINTS

The guide's prose explains and teaches. But judgment (the AI review gate) needs
verifiable TARGETS, not prose. That is why every governance guide has a section
of conformance points: the list of things that, if an artifact violates them, put it
OUTSIDE the ruler. Each point is objective enough for a pass/fail verdict.

This makes judgment LESS heuristic: instead of "does this screen respect atomic
design?" (vague, the AI guesses the criterion), it becomes "verify CK1, CK2, CK3" (each one a
verdict, with the criterion fixed in the guide). The checklist is the distilled prose — it cannot
diverge from it, because it is born from it.

## The mandatory section

Every guide that governs something (has a 'governs' rule) MUST have:

  ## Pontos de conformidade

  - CK1: <a verifiable, affirmative criterion> — <how to recognize the violation>
  - CK2: ...

Rules for the points:
- EACH point has a code (CK1, CK2, …) — the identity, so the report can reference the
  specific item (like the scenarios of a spec).
- AFFIRMATIVE and objective: "atoms come from the design system", not "good organization".
  Someone must be able to say pass/fail looking at the target, without guessing what you meant.
- INDEPENDENT: one point tests one thing. If you wrote "and" in the middle, it is probably
  two points.
- ANCHORED in the prose: each point distils a rule the guide's body already explains. If
  a point has no explanation above, either the explanation is missing, or the point is invention.
- HONEST about the non-computable: a point may be subjective ("is the separation of
  responsibilities clear?") — that is fine, it is exactly what AI judgment measures.
  But prefer the objective when you can; the more objective the point, the less the AI guesses.

## Points BY TARGET (when the guide governs more than one thing)

A guide usually governs MORE THAN ONE kind of target — the frontend guide governs screens AND
components; the backend one governs handlers AND models. A point that holds for screens
("atoms come from the design system") does not hold for a model. So do NOT make a flat
checklist that mixes everything: GROUP the points by target, with a subheading that names the
layer/tag they apply to. The codes gain a group prefix:

  ## Pontos de conformidade

  ### Para telas (tag: screen)
  - SCR-CK1: the screen injects components, it does not assemble inline layout blocks
  - SCR-CK2: atoms and molecules come from the design system, they are not redefined here

  ### Para componentes (tag: component)
  - CMP-CK1: the component does not fetch data — it receives everything by props
  - CMP-CK2: an atom does not compose other domain atoms

  ### Para todos (any governed target)
  - GEN-CK1: file names match what they export

How judgment uses this: when confronting a target, the AI looks at the target's LAYER/TAG
(the map tells it) and applies ONLY the points of the corresponding group + those of the "Para
todos" group. A screen point is never demanded of a model. That way each target is measured by the
right ruler, and the report cites the code of the matching point (SCR-CK1, CMP-CK2…).

Why this matters — it makes the analysis FOCUSED and OBJECTIVE:
- FOCUS: the AI verifies only the points of the target's layer, without weighing irrelevant rules.
- OBJECTIVITY: each verdict points at a named CK, not a general impression —
  reproducible (the same target against the same CKs converges between runs).
- ECONOMY: the AI consults the relevant group, not the guide's whole prose — which
  helps with batch adoption (see 'anchors guide', section ADOPTION).

If the guide governs a single kind of target, the flat checklist is enough — do not invent groups.

## How judgment uses this

When a 'measures: judgment' gate confronts a target against the guide, the AI:
1. reads the guide's Pontos de conformidade section;
2. evaluates the target against EACH point;
3. issues a report item by item — for each CK: compliant, or the violation (what, where,
   why, how to fix). See 'anchors guide' (section JUDGE).

A guide WITHOUT conformance points forces judgment to become heuristic —
'anchors doctor' flags that as debt.

## Structure of a guide

1. Purpose — what this guide governs and why (one or two sentences).
2. The rules, in prose — the teaching: what is right, examples, the reasoning. It is where
   someone learns to do it, not only to be verified.
3. Pontos de conformidade — the mandatory section above; the prose distilled into targets.
4. (Optional) Anti-patterns — the common mistakes and the antidote.

## What this guide governs

A guide must declare (via a 'governs' rule in anchors.yaml) WHICH tag it governs —
otherwise it confronts nobody and is not a guide, it is a doc. When creating a guide, add the
corresponding governs rule. If the guide is cross-cutting (it does not govern a type of artifact),
it is probably a doc, not a guide.

## Anti-patterns (refuse them)

- A guide without conformance points → it can only be confronted by heuristic; distil it.
- A vague point ("well structured") → nobody judges pass/fail; make it objective.
- A point that joins several checks with "and" → split it into independent points.
- A checklist that contradicts the prose → it must BE BORN from the prose, not compete with it.
- A guide that governs nobody → declare the governs, or it is a doc.
`
