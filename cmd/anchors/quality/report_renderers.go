package quality

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/health"
	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/queue"
)

// Os renderers das perspectivas de relatório. Cada um RECORTA fontes que o Anchors já
// mede — nenhum inventa dado. Todos devolvem markdown.

// --- QUALITY: o veredito dos gates + o débito ---
func renderQuality(ctx reportCtx) string {
	var b strings.Builder
	b.WriteString(reportHeader("Quality report (gates)", ctx.when))
	if ctx.cfg == nil || len(ctx.cfg.Gates) == 0 {
		b.WriteString("_No gate declared in anchors.yaml._\n")
		return b.String()
	}
	prof := gate.Aggregate(gate.RunWithConfig(ctx.cfg.Gates, ctx.g.Nodes, ctx.root, ctx.g, ctx.cfg))

	b.WriteString("## Verdicts by gate\n\n")
	b.WriteString("| gate | force | ✓ | ✗ | ⏳/~ |\n|---|---|---:|---:|---:|\n")
	for _, name := range prof.GateNames() {
		s := prof.ByGate[name]
		force := "informative"
		if s.Blocking {
			force = "**blocking**"
		}
		fmt.Fprintf(&b, "| %s | %s | %d | %d | %d |\n", name, force, s.Pass, s.Fail, s.Skip+s.Pending+s.Judge)
	}
	b.WriteString("\n")
	if prof.Passed {
		b.WriteString("✓ **Promotable** — every blocking gate passes.\n\n")
	} else {
		fmt.Fprintf(&b, "✗ **Barred** — %d blocking gate(s) failed.\n\n", len(prof.Blocked))
	}

	if len(prof.Failures) > 0 {
		b.WriteString("## Divergences (they become issues)\n\n")
		shown := 0
		for _, r := range prof.Failures {
			if shown >= 40 {
				fmt.Fprintf(&b, "- … and %d more\n", len(prof.Failures)-shown)
				break
			}
			mark := "informative"
			if r.Blocking {
				mark = "BLOCKS"
			}
			fmt.Fprintf(&b, "- [%s] `%s` @ `%s` — %s\n", mark, r.Gate, r.Target, firstLineOf(r.Detail))
			shown++
		}
	}
	if len(prof.Judged) > 0 {
		fmt.Fprintf(&b, "\n## Awaiting AI judgment (%d)\n\nPending judgment gates — run `anchors judge --pending`.\n", len(prof.Judged))
	}
	b.WriteString(reportFooter())
	return b.String()
}

// --- STRUCTURE: camadas, governança, colisões, órfãos ---
func renderStructure(ctx reportCtx) string {
	var b strings.Builder
	b.WriteString(reportHeader("Structure report", ctx.when))
	st := ctx.g.Statistics()

	b.WriteString("## Nodes by type\n\n| kind | nodes |\n|---|---:|\n")
	for _, k := range sortedKindKeys(st.NodesByKind) {
		fmt.Fprintf(&b, "| %s | %d |\n", k, st.NodesByKind[k])
	}
	fmt.Fprintf(&b, "\n%d nodes, %d edges.\n\n", st.Nodes, st.Edges)

	b.WriteString("## Governance (who governs whom)\n\n")
	gov := ctx.g.GovernanceSummary()
	if len(gov) == 0 {
		b.WriteString("_No guide governs anything (no governs rules)._\n\n")
	} else {
		b.WriteString("| guide | governs (direct) |\n|---|---:|\n")
		type row struct {
			g string
			n int
		}
		var rows []row
		for g, n := range gov {
			rows = append(rows, row{g, n})
		}
		sort.Slice(rows, func(i, j int) bool { return rows[i].n > rows[j].n })
		for _, r := range rows {
			fmt.Fprintf(&b, "| `%s` | %d |\n", r.g, r.n)
		}
		b.WriteString("\n")
	}

	// colisões e órfãos vêm do health (recorte por check)
	rep := health.Diagnose(ctx.g, ctx.cfg, ctx.root)
	b.WriteString(findingSection(rep, "identidade-duplicada", "Identity collisions"))
	b.WriteString(findingSection(rep, "orfao", "Orphans (code with no spec)"))
	b.WriteString(findingSection(rep, "identidade-ausente", "Missing identity"))
	b.WriteString(reportFooter())
	return b.String()
}

// --- CONFIG: estado do anchors.yaml e o que falta ---
func renderConfig(ctx reportCtx) string {
	var b strings.Builder
	b.WriteString(reportHeader("Configuration report", ctx.when))
	if ctx.cfg == nil {
		b.WriteString("_anchors.yaml not found._ Run `anchors init`.\n")
		return b.String()
	}
	fmt.Fprintf(&b, "## Declared layers (%d)\n\n| layer | kind | tags | code_prefix |\n|---|---|---|---|\n", len(ctx.cfg.Layers))
	for _, name := range sortedKeys(ctx.cfg.Layers) {
		l := ctx.cfg.Layers[name]
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", name, l.Kind, strings.Join(l.Tags, ","), l.CodePrefix)
	}
	b.WriteString("\n## Declared governance\n\n")
	if len(ctx.cfg.Governs) == 0 {
		b.WriteString("_No governs rule._\n")
	} else {
		for _, gr := range ctx.cfg.Governs {
			fmt.Fprintf(&b, "- `%s` governs the tag `%s`\n", gr.From, gr.Governs)
		}
	}
	fmt.Fprintf(&b, "\n## Declared gates (%d)\n\n", len(ctx.cfg.Gates))
	if len(ctx.cfg.Gates) == 0 {
		b.WriteString("_No gate._ Consider `anchors init` to seed the defaults.\n")
	} else {
		for _, gt := range ctx.cfg.Gates {
			force := "informative"
			if gt.IsBlocking() {
				force = "blocking"
			}
			kind := gt.Check
			if gt.IsJudgment() {
				kind = "AI-judgment"
			} else if gt.Run != "" {
				kind = "external: " + gt.Run
			}
			fmt.Fprintf(&b, "- `%s` (%s, %s) over %v\n", gt.Name, kind, force, gt.On)
		}
	}
	b.WriteString("\n## Derived (co-location): ")
	if ctx.cfg.Derived == nil {
		b.WriteString("not configured\n")
	} else {
		fmt.Fprintf(&b, "anchor `%s`, %d template(s)\n", ctx.cfg.Derived.Anchor, len(ctx.cfg.Derived.PadroesDe()))
	}
	// o que falta — recorte do health por checks estruturais de config
	rep := health.Diagnose(ctx.g, ctx.cfg, ctx.root)
	b.WriteString("\n" + findingSection(rep, "guide-sem-governo", "Guides with no governance (declare governs or turn it into a doc)"))
	b.WriteString(findingSection(rep, "kind-sem-gate", "Kinds with no gate (coverage hole)"))
	b.WriteString(reportFooter())
	return b.String()
}

// --- ISSUES: o trabalho rastreado (issues + tasks) ---
func renderIssues(ctx reportCtx) string {
	var b strings.Builder
	b.WriteString(reportHeader("Issues and tasks report", ctx.when))

	b.WriteString("## Issues (the adopted debt)\n\n")
	future, _ := issue.List(ctx.root, issue.Future)
	todo, _ := issue.List(ctx.root, issue.Todo)
	doing, _ := issue.List(ctx.root, issue.Doing)
	done, _ := issue.List(ctx.root, issue.Done)
	fmt.Fprintf(&b, "| state | count |\n|---|---:|\n| future | %d |\n| todo | %d |\n| doing | %d |\n| done | %d |\n\n",
		len(future), len(todo), len(doing), len(done))
	// SEPARADAS POR DONO, pelo mesmo motivo que `future/` fica à parte: o que espera uma
	// PESSOA e o que o agente resolve sozinho são listas com leitores diferentes, e
	// misturá-las faz as duas pararem de ser lidas. A do usuário vem primeiro — é a que
	// trava o resto, porque ninguém além dele pode destravá-la.
	abertas := append(append([]string{}, doing...), todo...)
	var deQuemDecide, doAgente []string
	for _, id := range abertas {
		st := issue.Todo
		if contémNome(doing, id) {
			st = issue.Doing
		}
		if issue.FileOwner(filepath.Join(ctx.root, issue.Dir, string(st), id)) == issue.DonoUsuário {
			deQuemDecide = append(deQuemDecide, id)
			continue
		}
		doAgente = append(doAgente, id)
	}
	if len(deQuemDecide) > 0 {
		b.WriteString("### Waiting on YOU (the agent cannot resolve it)\n\n")
		b.WriteString("A product decision, an answer that is not in the code. " +
			"While it lasts, whoever implements it will be guessing.\n\n")
		for _, id := range deQuemDecide {
			fmt.Fprintf(&b, "- %s\n", strings.TrimSuffix(id, ".md"))
		}
		b.WriteString("\n")
	}
	if len(doAgente) > 0 {
		b.WriteString("### Open (the work of now)\n\n")
		for _, id := range doAgente {
			fmt.Fprintf(&b, "- %s\n", strings.TrimSuffix(id, ".md"))
		}
		b.WriteString("\n")
	}
	// `future/` fica numa lista SEPARADA, e não somada às abertas: a dívida assumida não
	// é trabalho pendente, é trabalho adiado com um momento declarado. Misturá-la com o
	// que é para agora afogaria a fila — e é justamente por afogar que as listas param de
	// ser lidas.
	if len(future) > 0 {
		b.WriteString("### Deferred (adopted debt — it comes due, it does not vanish)\n\n")
		for _, id := range future {
			fmt.Fprintf(&b, "- %s\n", strings.TrimSuffix(id, ".md"))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Tasks (the work queue)\n\n")
	tasks, _ := queue.List(ctx.root)
	if len(tasks) == 0 {
		b.WriteString("_Empty queue._\n")
	} else {
		fmt.Fprintf(&b, "%d live task(s):\n\n", len(tasks))
		for _, t := range tasks {
			fmt.Fprintf(&b, "- [%s] `%s` — %s (%s)\n", t.State, t.Changed, t.SuggestedNext, t.Kind)
		}
	}
	b.WriteString(reportFooter())
	return b.String()
}

// --- INCONSISTENCIES: a lista (potencialmente infinita) de coisas a arrumar ---
func renderInconsistencies(ctx reportCtx) string {
	var b strings.Builder
	b.WriteString(reportHeader("Inconsistencies to fix", ctx.when))
	b.WriteString("> Everything the validators detect that has NOT yet become an issue — the raw\n")
	b.WriteString("> debt backlog. Triage and promote to an issue whatever will be handled (detection → issue → plan).\n\n")

	// 1. saúde do grafo (health) — todos os findings, agrupados por check
	rep := health.Diagnose(ctx.g, ctx.cfg, ctx.root)
	byCheck := map[string][]health.Finding{}
	var order []string
	for _, f := range rep.Findings {
		if _, ok := byCheck[f.Check]; !ok {
			order = append(order, f.Check)
		}
		byCheck[f.Check] = append(byCheck[f.Check], f)
	}
	sort.Strings(order)
	total := len(rep.Findings)
	fmt.Fprintf(&b, "## From the health validator (%d)\n\n", total)
	for _, check := range order {
		fs := byCheck[check]
		fmt.Fprintf(&b, "### %s (%d)\n\n", check, len(fs))
		for i, f := range fs {
			if i >= 30 {
				fmt.Fprintf(&b, "- … and %d more\n", len(fs)-i)
				break
			}
			sev := " "
			if f.Severity == health.Warn {
				sev = "⚠"
			}
			fmt.Fprintf(&b, "- %s `%s` — %s\n", sev, f.Subject, f.Detail)
		}
		b.WriteString("\n")
	}

	// 2. gates que reprovam (recorte do quality, sem repetir a tabela)
	if ctx.cfg != nil && len(ctx.cfg.Gates) > 0 {
		prof := gate.Aggregate(gate.RunWithConfig(ctx.cfg.Gates, ctx.g.Nodes, ctx.root, ctx.g, ctx.cfg))
		if len(prof.Failures) > 0 {
			fmt.Fprintf(&b, "## From the quality gates (%d failures)\n\n", len(prof.Failures))
			shown := 0
			for _, r := range prof.Failures {
				if shown >= 50 {
					fmt.Fprintf(&b, "- … and %d more\n", len(prof.Failures)-shown)
					break
				}
				fmt.Fprintf(&b, "- `%s` @ `%s` — %s\n", r.Gate, r.Target, firstLineOf(r.Detail))
				shown++
			}
		}
	}
	fmt.Fprintf(&b, "\n**Total inconsistencies: %d** (health) — triage and convert into issues whatever will be handled.\n", total)
	b.WriteString(reportFooter())
	return b.String()
}

// --- helpers ---

func reportFooter() string {
	return "\n---\n_Anchors — report generated from what the project measures._\n"
}

func firstLineOf(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}

// findingSection recorta os findings de um check específico do health numa seção.
func findingSection(rep health.Report, check, title string) string {
	var fs []health.Finding
	for _, f := range rep.Findings {
		if strings.HasPrefix(f.Check, check) {
			fs = append(fs, f)
		}
	}
	if len(fs) == 0 {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "## %s (%d)\n\n", title, len(fs))
	for i, f := range fs {
		if i >= 25 {
			fmt.Fprintf(&b, "- … and %d more\n", len(fs)-i)
			break
		}
		fmt.Fprintf(&b, "- `%s` — %s\n", f.Subject, f.Detail)
	}
	b.WriteString("\n")
	return b.String()
}

func sortedKindKeys(m map[mapx.Kind]int) []mapx.Kind {
	ks := make([]mapx.Kind, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Slice(ks, func(i, j int) bool { return ks[i] < ks[j] })
	return ks
}

// contémNome diz se o nome está na lista — para saber de qual estado a issue veio, já que
// `doing` e `todo` são concatenadas antes de separar por dono.
func contémNome(lista []string, nome string) bool {
	for _, n := range lista {
		if n == nome {
			return true
		}
	}
	return false
}
