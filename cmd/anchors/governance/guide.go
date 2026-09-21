package governance

import (
	"fmt"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/spf13/cobra"
)

// agentGuide é o playbook que ensina uma IA a OPERAR o Anchors. A inversão: o
// Anchors não invoca IA; a IA (no CLI dela) invoca o Anchors como ferramenta. Ela
// roda `anchors guide`, aprende o fluxo, e opera com os comandos existentes.
const agentGuide = `# Operating Anchors (guide for AI agents)

You are an AI agent and Anchors is a command-line TOOL you use to develop with
rigor. You read and write the files (specs, code, features, tests) — Anchors does
NOT generate content; it says what to do, keeps the dependency map, queues the
work, and verifies what you did.

## The model

- An ANCHOR is a document that guides development and confronts what was done
  (spec, feature, guide, doc). They live in the repository.
- The MAP (anchors.graph.yaml) links the anchors to the files: who depends on whom.
- A GUIDE is the ruler for producing each artifact (spec, code, feature, test,
  plan). ALWAYS read the relevant guide before writing.
- The QUEUE (.anchors/tasks/) is the pending work. The WATCHER feeds it: every time
  a file changes, it queues a task with the suggested next step. Whoever works
  PULLS the task — the watcher never calls you.

## Two roles (read carefully — this is what keeps you free)

There are TWO modes in which you act, and do not mix them:

- CONVERSATION: the session in which you talk to the user. It must NEVER get stuck
  running a work queue. Here you plan, dispatch, and report.
- WORKER: whoever actually executes a task from the queue (writes the spec,
  implements, tests). A worker takes ONE task, does ONE step, closes it, and ends.

### How to run a worker WITHOUT holding the conversation hostage

When there is work in the queue and it is time to execute it:

  • IF you can delegate to a subagent/process in the BACKGROUND (e.g. a subagent
    tool): delegate the worker and GO BACK to talking to the user. They stay free.
  • IF you do NOT have background (many AI clients do not): do NOT hijack the
    conversation silently. Tell the user, for example:
      "There are N tasks in the queue (anchors queue). I can process them now, but
       that will occupy me for a while — do you want me to go ahead, or would you
       rather open another session/terminal to run 'anchors next' in parallel?"
    The decision is the user's. The queue is a file on disk and the claim is
    atomic, so several sessions can run 'anchors next' at once without collision.

Golden rule: queue work runs in background if it can; otherwise, only with the
user's explicit approval. The conversation is theirs, not yours to monopolize.

## The development flow

### 0. PREPARE THE MAP AND START THE WATCHER (first thing)
Run:  anchors map build     (the watcher and impact/check need the map to exist)
Run:  anchors watch start
From here on, every file change becomes a task in the queue — including YOURS.
Yes, it is redundant when you were the one who changed the file (you already knew).
But it is the SAME path as when the change comes from outside (the user edited in
the editor, a 'git pull' brought something). One mechanism only, for every origin.
Trust the queue, not your memory.

### 0.5. DISCOVER — ONLY if the project does NOT YET EXIST (in the CONVERSATION)
If the directory is empty (or nearly — no code, no anchors.yaml), there is a pass
BEFORE the plan: discover what the project is technically. Run 'anchors guide project'
and follow the ruler. In short: YOU interview the user in 5 stages (purpose and form →
language → architecture and paradigm → macro structure and conventions → tooling and
formatting), ONE question at a time, waiting for each answer. At the end, and only at
the end, write two files at the root:
  • PROJECT.md   the TECHNICAL summary decided (stack, paradigm, structure, indentation,
                 extensions, editors) — this is what you read before writing each file
  • INSIGHTS.md  the transcript: each question, each answer, and what was DISCARDED
                 and why — this is what answers "why this choice?" later
Skip this stage in a project that already has code: there 'anchors init' infers from disk.
(Even so it is worth writing the PROJECT.md of what you OBSERVED — see the guide.)

### 1. PLAN (in the CONVERSATION)
Read the plan guide:  anchors guide plan
With the user, draft IN TEXT in the conversation a PLAN that decides WHICH specs need
to be born or change (it seeds specs, never code directly). Iterate until they APPROVE
the scope. ONLY THEN write the plan file — saving it already makes the watcher queue
the "specify" task (the belt starts). Writing the file before the "yes" starts the
machine without approval. → approved? write the plan.

### 2+. WORK BY PULLING FROM THE QUEUE
From the plan on, the work flows through the queue. A WORKER's cycle is always the same:

  a) anchors next               pulls and claims the next task (atomic)
  b) anchors map build          incorporate into the map the file the task cites (if it
                                is NEW, it is not in the map yet — without this step,
                                'impact' and 'check' answer "is not in the map")
  c) anchors impact <file>      now yes: the fine detail of what to propagate/validate
  d) execute the step, reading the right GUIDE before writing:
       • task "specify"   → write the .spec.md (spec guide). Give each requirement
                            its scenario code (identity) and its regime.
       • task "implement" → write code + feature (code/feature guides).
       • task "test"      → write the tests (test guide).
       • task "verify"    → only confront (there is no new artifact to write).
       • task "verify-tests" → RUN the test suite producing the coverage
                            artifacts (see your stack's CoverageHint: jest
                            --coverage + jest-junit, go test -coverprofile + go-junit-
                            report, pytest --cov --junitxml…), then:
                              anchors ingest --junit <r.xml> --lcov <cov.info>
                            so Anchors knows which tests passed and which spec
                            requirements got proven. Then ask the two things that
                            catch NEW bugs:
                              anchors coverage --diff <base> --lcov <cov>  (is what you
                                changed covered? — catches the new line with no test)
                              anchors coverage --delta                     (did coverage
                                drop vs. before? — catches coverage regression)
                            ONLY THEN confront.
  e) anchors map build          incorporate the files YOU just wrote
  f) anchors check --changed <file>   confront (see CONFRONT below)
  g) anchors done <id>          close the task (it goes to .anchors/done/)
  h) go back to (a) until 'anchors next' says "queue empty"

RULE: 'impact' and 'check' read the MAP, not the disk. Every new or moved file only
exists for them after an 'anchors map build'. When in doubt, run 'map build' first.

Each file you save in (d)/(e) makes the watcher queue the NEXT task
(spec→implement, feature→test, ...). That is why the cycle sustains itself: you do not
need to remember what comes next; the queue tells you.

### CONFRONT (step (f), in detail)
  anchors check --changed <file>      incremental, what the change touched
  anchors check --all                 the whole project
  • A BLOCKING gate that failed prevents promotion (exit 1); an informative one only records.
  • REPORT to the user the gates that failed and why. Fix them and run again.
  • NEVER close the task (done) with a blocking gate still red.

### JUDGE (the gate a script does not compute — YOU are the meter)
Some gates measure what no script knows: "does this screen break down in atomic design?",
"does the spec describe behaviour and not implementation?", "does every public export of
this unit trace to a spec rule — or is it dead code/scope-creep?" (the inverse of
feature-test-match: here a greppable SYMBOL is not enough, because the spec describes
behaviour and does not cite impl names — judging requires understanding that rule X
justifies export Y). These are 'measures: judgment' gates.
'anchors check' does not compute them — it marks the targets as ⏳ and queues them. You:
  a) anchors judge --pending            see the targets and, in each task, the guide + the question
  b) READ the indicated guide. Go straight to the "## Pontos de conformidade" section — the list of
     CK items is what you must verify. Do NOT judge by the whole prose "by eye";
     judge the target against EACH point. If the checklist is GROUPED BY TARGET (subsections
     "### Para <camada>"), apply only the points of THIS target's layer group + those of the
     "Para todos" group — a screen point does not apply to a model. That makes the
     verdict focused, objective and reproducible. (If the guide lacks the section, it is debt —
     the deterministic gate 'guide-has-checklist' already warns; judge by what is there.)
  c) anchors judge <target> --gate <g> --verdict pass|fail --reason "<REPORT>"
You evaluated the target against each CK; do NOT return one sentence — return the REPORT per item.
For each CK: compliant, or the non-compliance with what, WHERE (file:line), which CK
it violated, and HOW to fix it. That text becomes the issue body — so nobody reprocesses the
target later just to find out what to fix.
'fail' opens the issue with the report; 'pass' resolves the previous issue (if any). The
verdict is stamped with the target's rev, so it AGES (stale) if the target changes.

### HEALTH (when you want the panorama)
  anchors doctor    SYSTEMIC loose ends that check does not see: orphans, coverage
                    holes, loose layers, dead edges. Does not block.

## ADOPTION — bringing Anchors into a project that ALREADY EXISTS

A project BORN with Anchors does not have this problem: each file is confronted
and stamped when it is created, so it never accumulates debt. The heavy cold start
belongs only to WHOEVER ADOPTS — a large, old project whose map is born with thousands
of (guide, file) pairs, all without a stamp. If that is your case, do NOT try to audit
everything at once (it can cost millions of tokens). Do it in BATCHES, in this order:

1. STRUCTURAL CUT first (free, no AI). Run 'anchors governs' and 'anchors
   doctor'. If several guides govern the SAME set (redundancy), narrow the tags in
   anchors.yaml so each guide governs only its own scope. That cuts the work at the root
   before spending a token. Also complete the missing 'governs' rules (a guide
   without governance is not a guide — the doctor warns).

2. BATCH BY GUIDE, not by target. The guide is expensive to read (hundreds of lines) and
   the target is cheap. So do NOT reread the guide for each file: read the guide ONCE and
   judge many targets of that guide in sequence. Use 'anchors governs <guide>' to
   get the list of targets, read the guide, and judge the batch. That amortizes the guide's
   fixed cost across the targets (a ~4× cut). If you delegate to subagents to
   parallelize, give ONE WHOLE GUIDE to each subagent (never one target per subagent —
   that makes each one reread the guide, the opposite of the saving). Judging is reading: it
   does not need an isolated worktree/workspace.

3. THE STAMP IS THE CURSOR. The first pass does NOT have to finish in one go. Each
   'anchors judge' stamps the target; re-confronting skips what already has a valid stamp at
   the current rev. So do one batch today, another tomorrow — Anchors only re-queues what
   has not been judged yet (or what changed since). Report to the user what you covered
   and what is left; never fake total coverage that did not happen.

## BATCH FIXING — sweeping the project file by file (parallel, without rework)

When there are MANY pending items spread around (e.g. applying the @anchors header to the
whole project, completing specs), the efficient way is to sweep by FILE, with several
workers in parallel — but there are two rules that avoid waste:

1. ONE FILE AT A TIME, ALL ITS PENDING ITEMS. If you are going to open a file, fix
   EVERYTHING pending in it at once — not just the header. Use:
     anchors audit <file>            ITS pending items (gates + doctor)
     anchors audit <file> --impact   the whole UNIT (the triad spec↔code↔feature↔
                                     test) — fixes the unit in a single worker.
   It is a waste for one worker to touch the .spec header and another the sibling .tsx;
   take the unit (--impact) and resolve the whole triad at once.

2. TOP DOWN THE TREE. Fixing a CHILD and then the PARENT (which governs it)
   forces redoing the child. Process in topological order — rulers/specs (parents)
   BEFORE code/tests (governed). Anchors gives you the order ready-made:
     anchors map show --worklist --pending   the files with pending items, ALREADY ordered
                                             top down. Process in that order.
   When parallelizing: distribute SLICES of the worklist keeping the order (one worker per
   feature/subtree), never a parent and its child in concurrent workers.

## What to report to the user

- When starting the flow: "watcher on; I will plan with you and then process the queue
  in background / or ask your approval to run, depending on my capability."
- After 'check': which gates failed, which block, the debt left behind.
- After 'doctor': systemic loose ends.
- Divergences you do NOT resolve on your own → present the options: fix the target,
  update the anchor (spec/guide), or turn it into a new plan.
- NEVER "fix" by making the anchor lie about the code. If the code diverges from the
  spec, either the code is wrong (fix it) or the spec aged (update the spec, which
  may generate new work — report that).

## Command reference

  anchors init                      configures the project (anchors.yaml) — once
  anchors install-hooks             installs the git pre-commit that runs the gates on staged files
  anchors new <kind> <name>         emits the skeleton of an artifact (spec|feature|test)
  anchors new <kind> --list-sections  the sections (default/optional) of the kind
  anchors recode <old> <new>        renames a code and propagates (dry-run; --apply writes)
  anchors watch start               starts the watcher (queues tasks in background)
  anchors watch status|stop|pause|resume|logs   controls the watcher
  anchors queue                     lists the live tasks (read-only)
  anchors next                      pulls+claims the next task (the worker)
  anchors done <id>                 closes a finished task
  anchors map build                 (re)builds the dependency map
  anchors map show <file>           a node's neighbourhood (↑governed / ↓propagates)
  anchors impact <file>             what a change reaches (↓propagates / ↑validates)
  anchors check --changed <file>    runs the gates over the impact path
  anchors check --all               runs the gates over everything
  anchors ingest --junit/--lcov     ingests test signals from the runner (execution/coverage)
  anchors coverage [<spec>]         coverage by scenario (requirement proven?) and by line
  anchors coverage --diff <base> --lcov <c>   is what I CHANGED covered? (patch coverage)
  anchors coverage --delta          did coverage DROP since the last ingestion?
  anchors ingest --junit r --layer <unit|integration|e2e>   ingests by LAYER (merges)
  anchors report tests|quality|structure|config|issues|inconsistencies   reports in docs/
  anchors report all                generates every perspective + index in docs/anchors/
  anchors doctor                    ecosystem health (systemic loose ends)
  anchors guide                     this guide
  anchors guide project             the ruler of the DISCOVER phase (new project → PROJECT.md)
  anchors guide plan                the ruler of the PLAN phase

## Where the project's guides are

The guides (the rulers of each artifact) live in the repository — see the anchors.yaml
section for the 'guide' layer (typically in guides/). Read the right guide
BEFORE writing each type of artifact. If the project has no guide for something
you are going to produce, warn the user — it is a loose end.
`

func newGuideCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guide",
		Short: "Print the Anchors guides for AI agents",
		Long: `With no arguments, prints the operating playbook: the development flow
(plan → specify → map → implement → test → confront), the commands,
and what to report to the user. The AI runs this first and operates Anchors as a tool.

Subcommands print the guides for the specific rulers:
  anchors guide project  how to discover a project that does not yet exist (PROJECT.md)
  anchors guide plan     how to structure a plan (the origin of the movement)
  anchors guide spec     how to write a spec (the source of truth)
  anchors guide code     how to implement the code guided by the spec
  anchors guide feature  how to write the feature (the behaviour scenarios)
  anchors guide test     how to write the tests (the executable ruler)
  anchors guide guide    how to write a guide (the ruler of a ruler)
  anchors guide header   the header block of every file (cross-cutting, mandatory)`,
		// sem subcomando → o playbook de operação
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(agentGuide)
			return nil
		},
	}
	cmd.AddCommand(
		// Os DOIS guias que dizem o que fazer diante do que não se sabe ganham a seção de
		// AUTONOMIA, que muda conforme a declaração local (`.anchors/settings.yaml`).
		//
		// Ela não é um aviso no fim do texto: quem não decide o produto lê uma instrução
		// diferente, no lugar onde ela importa. O `settings user-issues` fecha a porta do
		// claim; esta seção fecha a que mais se usa — perguntar a quem está rodando o
		// agente, e receber uma resposta razoável de quem não tinha autoridade para dá-la.
		newGuideComAutonomia("review", "how to review a PR: what is yours and what check already measured", reviewGuide),
		newGuideComAutonomia("work", "how to work a card: the order, what to do with a finding that is not its own", WorkGuide),
		newGuideSubCmd("project", "how to discover a project that does not yet exist (PROJECT.md + INSIGHTS.md)", projectGuide),
		newGuidePlanCmd(),
		newGuideProductCmd(),
		newGuideFlowCmd(),
		newGuideSubCmd("spec", "how to write a spec (the source of truth)", specGuide),
		newGuideSubCmd("code", "how to implement the code guided by the spec", codeGuide),
		newGuideSubCmd("feature", "how to write the feature (behaviour scenarios)", featureGuide),
		newGuideSubCmd("test", "how to write the tests (the executable ruler)", testGuide),
		newGuideSubCmd("guide", "how to write a guide (the ruler of a ruler)", guideGuide),
		newGuideSubCmd("header", "the header block of every file (cross-cutting, mandatory)", headerGuide),
	)
	return cmd
}

// newGuideSubCmd fabrica um subcomando de guia que só imprime um texto embutido.
// As quatro réguas (spec/code/feature/test) compartilham essa casca fina.
func newGuideSubCmd(use, short, body string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(body)
			return nil
		},
	}
}

func newGuideFlowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "flow",
		Short: "Print the flow guide (work driven by shape, not by memory)",
		Long: `The flow guide is the ruler for ACTIONS (the puzzle pieces, each declaring
its results) and FLOWS (the assembly that says where each result goes).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(flowGuide)
			return nil
		},
	}
}

func newGuideProductCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "product",
		Short: "Print the product doctrine guide (the rule that cuts across targets)",
		Long: `The product doctrine guide is the ruler for the rule that belongs to no
single unit: it lives in product/, and the specs point at it with ` + "`@realizes`" + `.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(productGuide)
			return nil
		},
	}
}

func newGuidePlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Print the plan guide (how a plan is structured and seeds specs)",
		Long: `The plan guide is the ruler of the PLAN phase: how the AI and the user
write a plan that decides WHICH specs are born or change — never code directly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(planGuide)
			return nil
		},
	}
}

// newGuideComAutonomia fabrica um guia que anexa a seção de autonomia.
//
// A seção depende do `.anchors/settings.yaml`, e por isso este subcomando aceita `--root`:
// sem ele, um agente que rodasse o guia de outro diretório leria a régua errada — e a
// régua errada aqui é a que autoriza perguntar.
func newGuideComAutonomia(use, short, body string) *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				// Fora de um projeto, o guia ainda serve — só não tem a declaração local
				// para consultar. Melhor imprimir a régua geral que recusar a ajuda.
				fmt.Print(body)
				return nil
			}
			fmt.Print(body)
			printAutonomy(absRoot)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	return cmd
}
