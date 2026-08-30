package main

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
	b.WriteString(reportHeader("Relatório de qualidade (gates)", ctx.when))
	if ctx.cfg == nil || len(ctx.cfg.Gates) == 0 {
		b.WriteString("_Nenhum gate declarado no anchors.yaml._\n")
		return b.String()
	}
	prof := gate.Aggregate(gate.RunWithConfig(ctx.cfg.Gates, ctx.g.Nodes, ctx.root, ctx.g, ctx.cfg))

	b.WriteString("## Vereditos por gate\n\n")
	b.WriteString("| gate | força | ✓ | ✗ | ⏳/~ |\n|---|---|---:|---:|---:|\n")
	for _, name := range prof.GateNames() {
		s := prof.ByGate[name]
		force := "informativo"
		if s.Blocking {
			force = "**bloqueante**"
		}
		fmt.Fprintf(&b, "| %s | %s | %d | %d | %d |\n", name, force, s.Pass, s.Fail, s.Skip+s.Pending+s.Judge)
	}
	b.WriteString("\n")
	if prof.Passed {
		b.WriteString("✓ **Promovível** — todos os gates bloqueantes passam.\n\n")
	} else {
		fmt.Fprintf(&b, "✗ **Barrado** — %d gate(s) bloqueante(s) reprovam.\n\n", len(prof.Blocked))
	}

	if len(prof.Failures) > 0 {
		b.WriteString("## Divergências (viram issue)\n\n")
		shown := 0
		for _, r := range prof.Failures {
			if shown >= 40 {
				fmt.Fprintf(&b, "- … e mais %d\n", len(prof.Failures)-shown)
				break
			}
			mark := "informativo"
			if r.Blocking {
				mark = "BLOQUEIA"
			}
			fmt.Fprintf(&b, "- [%s] `%s` @ `%s` — %s\n", mark, r.Gate, r.Target, firstLineOf(r.Detail))
			shown++
		}
	}
	if len(prof.Judged) > 0 {
		fmt.Fprintf(&b, "\n## Aguardando julgamento de IA (%d)\n\nGates de julgamento pendentes — rode `anchors judge --pending`.\n", len(prof.Judged))
	}
	b.WriteString(reportFooter())
	return b.String()
}

// --- STRUCTURE: camadas, governança, colisões, órfãos ---
func renderStructure(ctx reportCtx) string {
	var b strings.Builder
	b.WriteString(reportHeader("Relatório de estrutura", ctx.when))
	st := ctx.g.Statistics()

	b.WriteString("## Nós por tipo\n\n| kind | nós |\n|---|---:|\n")
	for _, k := range sortedKindKeys(st.NodesByKind) {
		fmt.Fprintf(&b, "| %s | %d |\n", k, st.NodesByKind[k])
	}
	fmt.Fprintf(&b, "\n%d nós, %d arestas.\n\n", st.Nodes, st.Edges)

	b.WriteString("## Governança (quem rege quem)\n\n")
	gov := ctx.g.GovernanceSummary()
	if len(gov) == 0 {
		b.WriteString("_Nenhum guide rege nada (sem regras governs)._\n\n")
	} else {
		b.WriteString("| guide | rege (direto) |\n|---|---:|\n")
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
	b.WriteString(findingSection(rep, "identidade-duplicada", "Colisões de identidade"))
	b.WriteString(findingSection(rep, "orfao", "Órfãos (código sem spec)"))
	b.WriteString(findingSection(rep, "identidade-ausente", "Identidade ausente"))
	b.WriteString(reportFooter())
	return b.String()
}

// --- CONFIG: estado do anchors.yaml e o que falta ---
func renderConfig(ctx reportCtx) string {
	var b strings.Builder
	b.WriteString(reportHeader("Relatório de configuração", ctx.when))
	if ctx.cfg == nil {
		b.WriteString("_anchors.yaml não encontrado._ Rode `anchors init`.\n")
		return b.String()
	}
	fmt.Fprintf(&b, "## Camadas declaradas (%d)\n\n| layer | kind | tags | code_prefix |\n|---|---|---|---|\n", len(ctx.cfg.Layers))
	for _, name := range sortedKeys(ctx.cfg.Layers) {
		l := ctx.cfg.Layers[name]
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", name, l.Kind, strings.Join(l.Tags, ","), l.CodePrefix)
	}
	b.WriteString("\n## Governança declarada\n\n")
	if len(ctx.cfg.Governs) == 0 {
		b.WriteString("_Nenhuma regra governs._\n")
	} else {
		for _, gr := range ctx.cfg.Governs {
			fmt.Fprintf(&b, "- `%s` rege a tag `%s`\n", gr.From, gr.Governs)
		}
	}
	fmt.Fprintf(&b, "\n## Gates declarados (%d)\n\n", len(ctx.cfg.Gates))
	if len(ctx.cfg.Gates) == 0 {
		b.WriteString("_Nenhum gate._ Considere `anchors init` para semear os padrões.\n")
	} else {
		for _, gt := range ctx.cfg.Gates {
			force := "informativo"
			if gt.IsBlocking() {
				force = "bloqueante"
			}
			kind := gt.Check
			if gt.IsJudgment() {
				kind = "julgamento-IA"
			} else if gt.Run != "" {
				kind = "externo: " + gt.Run
			}
			fmt.Fprintf(&b, "- `%s` (%s, %s) sobre %v\n", gt.Name, kind, force, gt.On)
		}
	}
	b.WriteString("\n## Derived (co-location): ")
	if ctx.cfg.Derived == nil {
		b.WriteString("não configurado\n")
	} else {
		fmt.Fprintf(&b, "âncora `%s`, %d template(s)\n", ctx.cfg.Derived.Anchor, len(ctx.cfg.Derived.PadroesDe()))
	}
	// o que falta — recorte do health por checks estruturais de config
	rep := health.Diagnose(ctx.g, ctx.cfg, ctx.root)
	b.WriteString("\n" + findingSection(rep, "guide-sem-governo", "Guides sem governança (declare governs ou vire doc)"))
	b.WriteString(findingSection(rep, "kind-sem-gate", "Kinds sem gate (buraco de cobertura)"))
	b.WriteString(reportFooter())
	return b.String()
}

// --- ISSUES: o trabalho rastreado (issues + tasks) ---
func renderIssues(ctx reportCtx) string {
	var b strings.Builder
	b.WriteString(reportHeader("Relatório de issues e tasks", ctx.when))

	b.WriteString("## Issues (o débito adotado)\n\n")
	future, _ := issue.List(ctx.root, issue.Future)
	todo, _ := issue.List(ctx.root, issue.Todo)
	doing, _ := issue.List(ctx.root, issue.Doing)
	done, _ := issue.List(ctx.root, issue.Done)
	fmt.Fprintf(&b, "| estado | qtd |\n|---|---:|\n| future | %d |\n| todo | %d |\n| doing | %d |\n| done | %d |\n\n",
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
		if issue.DonoDoArquivo(filepath.Join(ctx.root, issue.Dir, string(st), id)) == issue.DonoUsuário {
			deQuemDecide = append(deQuemDecide, id)
			continue
		}
		doAgente = append(doAgente, id)
	}
	if len(deQuemDecide) > 0 {
		b.WriteString("### Esperando VOCÊ (o agente não pode resolver)\n\n")
		b.WriteString("Uma decisão de produto, uma resposta que não está no código. " +
			"Enquanto durar, quem implementar vai adivinhar.\n\n")
		for _, id := range deQuemDecide {
			fmt.Fprintf(&b, "- %s\n", strings.TrimSuffix(id, ".md"))
		}
		b.WriteString("\n")
	}
	if len(doAgente) > 0 {
		b.WriteString("### Abertas (o trabalho de agora)\n\n")
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
		b.WriteString("### Adiadas (dívida assumida — vence, não some)\n\n")
		for _, id := range future {
			fmt.Fprintf(&b, "- %s\n", strings.TrimSuffix(id, ".md"))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Tasks (a fila de trabalho)\n\n")
	tasks, _ := queue.List(ctx.root)
	if len(tasks) == 0 {
		b.WriteString("_Fila vazia._\n")
	} else {
		fmt.Fprintf(&b, "%d task(s) viva(s):\n\n", len(tasks))
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
	b.WriteString(reportHeader("Inconsistências a corrigir", ctx.when))
	b.WriteString("> Tudo que os validators detectam e ainda NÃO virou issue — o backlog\n")
	b.WriteString("> de débito cru. Triar e promover a issue o que for tratar (detecção → issue → plano).\n\n")

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
	fmt.Fprintf(&b, "## Do validador de saúde (%d)\n\n", total)
	for _, check := range order {
		fs := byCheck[check]
		fmt.Fprintf(&b, "### %s (%d)\n\n", check, len(fs))
		for i, f := range fs {
			if i >= 30 {
				fmt.Fprintf(&b, "- … e mais %d\n", len(fs)-i)
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
			fmt.Fprintf(&b, "## Dos gates de qualidade (%d reprovações)\n\n", len(prof.Failures))
			shown := 0
			for _, r := range prof.Failures {
				if shown >= 50 {
					fmt.Fprintf(&b, "- … e mais %d\n", len(prof.Failures)-shown)
					break
				}
				fmt.Fprintf(&b, "- `%s` @ `%s` — %s\n", r.Gate, r.Target, firstLineOf(r.Detail))
				shown++
			}
		}
	}
	fmt.Fprintf(&b, "\n**Total de inconsistências: %d** (health) — triar e converter em issues o que for tratar.\n", total)
	b.WriteString(reportFooter())
	return b.String()
}

// --- helpers ---

func reportFooter() string {
	return "\n---\n_Anchors — relatório gerado do que o projeto mede._\n"
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
			fmt.Fprintf(&b, "- … e mais %d\n", len(fs)-i)
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
