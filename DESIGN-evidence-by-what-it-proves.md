<!-- @anchors
  code: EBWIP
  layer: desenho
-->

# Evidence by what it proves — a proof goes stale when what it proves changes, not when the file does

> IMPLEMENTED on 2026-10-08 (W01–W03, one delivery). A test's proof is stamped with the revision of the whole file it
> proves, so any edit stales it — a date, a changelog line, a navigation row, a flag comment.
> The proof should follow what it proves: a file's behaviour, and, at the end, each scenario's
> own rule and text.

## The problem

Every proof Anchors keeps is compared with the revision of the whole file (`Rev`, the hash of
its content): a scenario's proof with its spec's, a capture's closure with each file's, a
mutation with the code's. Any edit stales it, whatever it changed:

- **MIF, after adopting navigation:** 51 specs edited only in their Navigation section (In/Out
  rows), their `updated_at` and their change history lost their Maestro and VR proofs — only
  the junit suite, rerun at the new revision, was left: a simulator rerun for edits that
  changed no behaviour (`scenario-coverage` ✗51).
- **Flag comments:** a code file that gains a `@dep:` or `@used-by:` written by hand (not by
  `check --fix`, which carries the evidence) moves the closure of every test reaching it, and
  every flow through it goes stale for a comment.
- **Before:** an `@no-vr` on a state line, a renamed VR tag — each stale every e2e scenario of
  the spec (MIF, 2026-09).

The answer so far is `anchors keep-evidence`: whoever edits declares the edit proves nothing new.
It works when someone remembers it; on MIF nobody ran it on those 51 specs, and the proofs were
gone by the next rebuild. A tool that depends on remembering is the hole the gates exist to close.

Part of the answer exists already: since DNDDP-W05 each Out row has a revision of its own, and a
flow asserting a navigation goes stale when its row changes (EVFRA-B11).

## The design

- **An evidence revision beside the file's revision.** `Rev` stays what it is — the file's
  content, which says whether the map matches the tree. Each node gains an `EvidenceRev`: the hash
  of the content with what changes no behaviour left out (`scan.EvidenceOf`). When a rebuild finds
  a file whose `Rev` moved and whose `EvidenceRev` held, it carries the file's evidence to the new
  revision — what held at the old one only, as `check --fix` does (`CarryUnchangedEvidence`).
  Every comparison of a proof keeps reading `Rev`, and nobody has to remember anything.
- **What a spec's evidence leaves out:** the `@anchors` header (its `updated_at` above all), the
  change history section, the navigation sections and an Out table under any title (each Out
  row is confronted by its own revision already), and blank-line and trailing-space differences.
- **What a code file's evidence leaves out:** the header — the hook dates code files too —, and
  the chain's flags — `@dep` (kinded too), `@used-by`, `@navigates`, `@no-dep`, `@no-nav`, the line
  that only carries one —, read the way the mock stamps read them (MCSTM-B21, B22, one reading
  now). Coverage and mutation, which record line numbers, follow a `LineRev` that keeps every line
  in place: a flag appended to a line keeps their proof; a flag line inserted moves their lines,
  and they wait for the next run, as before.
- **Per rule.** A spec's evidence also has a revision per rule definition (`RuleRevs`) and one
  for the rest (`RestRev`). A spec whose rules alone changed is carried with the proven scenarios
  of the rules that changed marked stale (`StaleCodes`): `scenario-coverage` reads them as stale
  — Pending — and the others as proven. A run that proves them again, or `keep-evidence`, makes
  them fresh. The feature's text is not part of it: a feature's change is the feature node's.
  This is the per-scenario freshness MIF proposed on 2026-09, parked then for `keep-evidence`.
- **Moving to it keeps what is fresh.** A map written before evidence revisions has none to
  compare with: the first build reads each moved file as the commit (or the index) has it and,
  when its hash is the map's revision, takes its evidence from there (`FillOldEvidence`). Nothing
  fresh goes stale on the upgrade, and nothing stale becomes fresh — only what held at the old
  revision is carried.
- **Everywhere a revision moves.** The build, the commit hook's map (which prefers HEAD's
  measurement of a file the commit leaves as HEAD has it), `anchors test`'s reading of the tree,
  and the watcher.

## The phases

All three delivered together, at the user's request (2026-10-08).

| Phase | What | Proof |
| --- | --- | --- |
| `EBWIP-W01` | **The spec's evidence.** `EvidenceRev` on spec nodes, leaving out header, history and navigation; a rebuild carries what it shows unchanged; the upgrade reads the old revision's evidence from the commit. | **Measured** on a clone of MIF at HEAD with the 89 specs its tree edits for navigation: the released 0.1.310 left all 89 without a fresh proof; the new build keeps 79 fresh, with no `keep-evidence`. Of the other 10, 8 were already stale at HEAD and stay so, and 2 changed their route line — outside navigation —, rightly stale. |
| `EBWIP-W02` | **The code's evidence.** `EvidenceRev` and `LineRev` on code and test nodes, leaving out the header and the chain's flags; closures carried with the file. | **Measured** on a clone of MIF: its commit adopting the chain (959 files flagged, headers re-dated) built over the previous map. The released 0.1.310 staled 393 of 1673 tests; the new build keeps the 1673 fresh, and the coverage of files that gained flag lines waits for its next run. |
| `EBWIP-W03` | **Per rule.** `RuleRevs` and `RestRev` on specs; a rule that changed alone marks its scenarios stale (`StaleCodes`), read apart by `scenario-coverage`; a run or `keep-evidence` freshens them. | Unit proofs (EDSTD-B20, B21, INCHN-B42, KPEVD-B06): editing one rule stales that rule's scenarios, each variant, and no other; a change outside the rules stales the whole spec. |

## What NOT to do

- Leave out of the evidence anything that can change behaviour. A rule row, a validation, a data
  contract, a state, a message, a test's own text — all stay in.
- Make `keep-evidence` redundant by guessing. It stays for what no rule can read: an edit the
  author knows proves nothing new.
- Let an upgrade stale every proof. The first build reads the old revision's evidence from the
  commit (W01).

## Decisions

| Code | Question | Decided |
| --- | --- | --- |
| `EBWIP-D01` | Should a navigation-only spec edit, or a flag-only code edit, keep the proofs without `keep-evidence`? | Yes — the evidence follows what it proves (the user, 2026-10-08, on MIF's report). |
| `EBWIP-D02` | Is W03 — per-scenario freshness — part of this plan? | Yes, delivered with W01 and W02 (the user, 2026-10-08). |

## Open Decisions

None.

## Lessons

- **The comparison sites stayed as they were.** The plan proposed every proof reading
  `EvidenceRev`; carrying the evidence when the revision moves — the mechanism `check --fix`
  already used — gives the same verdicts with no comparison site touched.
- **Code headers carry the hook's date too.** The first measurement on MIF kept nothing: each
  flagged file had also been re-dated. The header left the code's evidence as it had left the
  spec's.
- **The upgrade needed the old content.** A map with no evidence revisions has nothing to
  compare with; reading the commit's copy, checked against the map's revision, is what made the
  first build count.
