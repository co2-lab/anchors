package flow

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/change"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/queue"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

// `anchors work <artefato> --for <alvo>` emite o PROMPT DE TRABALHO de uma etapa.
//
// Por que existe: o `anchors guide` ensina a DOUTRINA (o que é uma spec, como se escreve
// um teste) — é a régua, permanente e agnóstica de alvo. Mas quem vai executar UMA etapa
// sobre UM arquivo precisa de outra coisa: o que ler AGORA, nesta ordem; qual é a camada
// deste alvo e o que ela exige; onde nascem as peças da trinca; o que NÃO é escopo desta
// etapa; e como verificar no fim.
//
// Sem isso, cada orquestrador reescreve esse prompt à mão a cada vez — e sai diferente
// toda vez (o observado: um agente reinventou o molde de `.feature` de backend lendo os
// vizinhos, embora o guide do projeto tivesse a seção certa; outro gastou uma rodada
// descobrindo que uma camada declarativa não tem spec).
//
// O comando NÃO inventa: ele COMPÕE o que já está declarado na Estrutura —
//   - `governs:` → quais guides regem a camada do alvo (a régua a ler)
//   - `layers:`  → o que é a camada, e se ela é regida ou reconhecida
//   - `derived:` → onde nascem spec/feature/teste desse alvo
//   - `work:`    → passos extra que só o projeto sabe (override opcional por camada)
func newWorkCmd() *cobra.Command {
	var root, target string

	cmd := &cobra.Command{
		Use:   "work <artifact>",
		Short: "Emit the work prompt of a stage (spec|code|feature|test|review) for a target",
		Long: `Composes the WORK PROMPT of a stage over a concrete target — what an
agent (or subagent) needs to execute that piece without reinventing the script.

  anchors work spec   --for packages/backend/repositories/metadata.ts
  anchors work test   --for apps/mobile/src/hooks/useAuth.ts
  anchors work review --for packages/backend/business-logic/pricing.ts

The ` + "`review`" + ` is the last stage of the cycle, and the different one: it produces no artifact —
it CONFRONTS what was delivered. It exists because the gate verifies what is DECLARABLE,
and what is left over only shows up to whoever attacks from the outside (mutating the
rule, running with edge input).
Run it after ` + "`anchors deliver`" + `, which is what gives it scope.

Difference from ` + "`anchors guide`" + `: the guide teaches the DOCTRINE (permanent, without a target);
work delivers the TASK (what to read now, which layer this target is in, where the
pieces are born, what is not your scope, how to verify).

The content is COMPOSED from anchors.yaml — nothing is invented here.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			artifact := strings.ToLower(args[0])
			if !queue.ValidWorkArtifact(artifact) {
				return fmt.Errorf("unknown artifact %q — use: %s", artifact,
					strings.Join(queue.ArtefatosDeTrabalho, ", "))
			}
			if target == "" {
				return fmt.Errorf("provide the target with --for <path> " +
					"(e.g.: --for packages/backend/repositories/metadata.ts)")
			}
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			rel := common.RelTo(absRoot, target)
			g, _ := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath)) // sem mapa: cai na convenção

			// O alvo é a UNIDADE (o arquivo de código), não uma peça derivada dela. Apontar
			// para a spec é o engano previsível — é o artefato que já existe, então é o que
			// vem à mão. O resultado era silenciosamente absurdo: os caminhos saíam
			// `x.spec.spec.md`, `x.spec.feature`, e a spec aparecia como peça a produzir de
			// uma spec que já estava no disco. Nada avisava.
			//
			// Redirecionar (em vez de recusar) porque a intenção é inequívoca: quem pediu
			// `work feature --for x.spec.md` quer a feature da unidade que x.spec.md
			// descreve. Dizer o que se fez mantém o usuário no controle.
			if alvo, achou := derivedPieceUnit(absRoot, rel, cfg, g); achou {
				fmt.Fprintf(os.Stderr, "note: `%s` is a derived piece, not the unit. "+
					"Using `%s` as the target.\n\n", rel, alvo)
				rel = alvo
			}
			out, err := composeWorkPrompt(absRoot, rel, artifact, cfg, g)
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&target, "for", "", "the target file of this stage (required)")
	return cmd
}

// composeWorkPrompt monta o prompt. A ordem das seções é deliberada: PAPEL (quem você
// é) → ALVO (sobre o quê) → RÉGUA (o que ler antes) → FRONTEIRA (o que não é seu) →
// PROCEDIMENTO (como) → VERIFICAÇÃO (como saber que acabou).
func composeWorkPrompt(root, rel, artifact string, cfg *config.Config, g *mapx.Graph) (string, error) {
	layer, _ := scan.Classify(rel, cfg)
	l, hasLayer := cfg.Layers[layer]

	var b strings.Builder
	// O REVIEW não produz artefato — ele CONFRONTA. Chamar de "trabalho de review" e
	// "você vai produzir review" empurra o revisor para o modo autor, que é o oposto do
	// que se quer: quem escreve tende a explicar o que vê, quem confronta tenta quebrar.
	if artifact == "review-plan" {
		fmt.Fprintf(&b, "# WHOLE review — `%s`\n\n", rel)
	} else if artifact == "review" {
		fmt.Fprintf(&b, "# Review of `%s`\n\n", rel)
	} else {
		fmt.Fprintf(&b, "# Work: %s of `%s`\n\n", artifact, rel)
	}

	// PAPEL + a recusa antecipada. Uma camada RECONHECIDA não tem spec — dizer isso
	// AQUI evita a rodada de descoberta (e o arquivo errado nascendo).
	// A recusa vale para TODA peça da trinca, não só a spec. Antes cobria apenas `spec`, e
	// `anchors work feature --for <arquivo declarativo>` abria com "Você vai produzir
	// feature" para, duas seções abaixo, dizer "não tem spec, feature nem teste próprios"
	// — e prescrever a verificação de um arquivo que o mesmo prompt proíbe criar. Medido
	// num E2E real: a task teve de ser descartada à mão.
	// PARE também quando a PEÇA desta etapa é dispensada pela camada (`trinca_opcional`),
	// ainda que a camada seja REGIDA. Antes o prompt abria com "Você vai produzir
	// **feature**" e o roteiro completo de produção, e a dispensa aparecia quatro linhas
	// abaixo, como item de uma lista. Um worker que segue a manchete cria o arquivo
	// proibido — e a régua tinha dito as duas coisas.
	if hasLayer && waivedPieces(layer, cfg)[artifact] {
		fmt.Fprintf(&b, "## STOP\n\n`%s` — the layer **%s** WAIVES the piece `%s` "+
			"(`trinca_opcional` in anchors.yaml).\n\nThe waiver is declared, not an "+
			"oversight: this layer does not prove behavior with this piece. Creating it "+
			"would produce the empty artifact the declaration exists to avoid.\n\n"+
			"Do not create the %s. If the queue gave you this task, it is noise — discard it "+
			"(`anchors drop`) and report.\n", rel, layer, artifact, artifact)
		return b.String(), nil
	}
	if hasLayer && l.Regime == "declarativo" && (artifact == "spec" || artifact == "feature" || artifact == "test") {
		peca := map[string]string{"spec": "spec", "feature": "feature", "test": "test"}[artifact]
		fmt.Fprintf(&b, "## STOP\n\n`%s` belongs to the layer **%s**, declared as RECOGNIZED "+
			"(`regime: declarativo`) in anchors.yaml.\n\nSuch layers **have no spec, feature "+
			"or test of their own**: they do not originate rules (they only translate/configure), so there is "+
			"nothing to specify and nothing to prove in a scenario. What proves the behavior is the layer "+
			"that DECIDES, with its own triad.\n\nDo not create the %s. If there is a DECISION to document, "+
			"it belongs to the layer that decides — report the contradiction to whoever asked for this stage.\n",
			rel, layer, peca)
		return b.String(), nil
	}

	if artifact == "review-plan" {
		b.WriteString("You are the REVIEWER OF THE WHOLE. Each unit of this plan already went through " +
			"its own review, and each one is correct **on its own** — do not repeat that work.\n\n" +
			"Your target is the SEAM: what only appears when the pieces meet. Measured " +
			"in a real E2E: of 6 findings of an adversarial review, **4 crossed units**, " +
			"and in all of them each piece passed in isolation. It is the class no gate catches and that " +
			"the per-unit review cannot reach by definition of scope.\n\n" +
			"Your output is a REPORT with executed evidence, not a correction.\n\n")
	} else if artifact == "review" {
		b.WriteString("You are the REVIEWER of this unit. The work was already delivered and **passed " +
			"every gate** — that is the normal state of everything that arrives here, and it is " +
			"exactly why you exist: the gate confronts what is DECLARABLE, and what " +
			"is left over only appears to whoever attacks from the outside.\n\n" +
			"Your output is a REPORT with executed evidence, not a correction.\n\n")
	} else {
		fmt.Fprintf(&b, "You are going to produce **%s** for the target below, following the project ruler.\n\n", artifact)
	}
	fmt.Fprintf(&b, "## Target\n\n- File: `%s`\n", rel)
	if hasLayer {
		fmt.Fprintf(&b, "- Layer: **%s**", layer)
		if l.Regime != "" {
			fmt.Fprintf(&b, " (regime: %s)", l.Regime)
		}
		b.WriteString("\n")
		if len(l.Tags) > 0 {
			fmt.Fprintf(&b, "- Tags: %s\n", strings.Join(l.Tags, ", "))
		}
	} else {
		b.WriteString("- Layer: **unclassified** — confirm whether the path belongs to some " +
			"layer of `layers:` in anchors.yaml before proceeding.\n")
	}

	// RÉGUA: os guides que regem esta camada, lidos do `governs:`.
	if guides := guidesFor(l, cfg, artifact); len(guides) > 0 {
		b.WriteString("\n## Read first (in this order)\n\n")
		fmt.Fprintf(&b, "1. `anchors guide %s` — the doctrine of the artifact (what it is, what it is not)\n", artifact)
		for i, g := range guides {
			fmt.Fprintf(&b, "%d. `%s` — this project's ruler for this layer\n", i+2, g)
		}
		fmt.Fprintf(&b, "%d. The target (`%s`) and its layer neighbors, to follow the local dialect\n",
			len(guides)+2, rel)
	} else {
		b.WriteString("\n## Read first\n\n")
		fmt.Fprintf(&b, "1. `anchors guide %s` — the doctrine of the artifact\n", artifact)
		fmt.Fprintf(&b, "2. The target (`%s`) and its layer neighbors\n", rel)
		b.WriteString("\n> No project guide governs this layer (`governs:` in anchors.yaml). " +
			"Follow the dialect of the neighbors and record the gap.\n")
	}

	// TRINCA: onde nascem as peças, derivado do `derived:`.
	//
	// Camada RECONHECIDA (`regime: declarativo`) não tem trinca: listar spec/feature/test
	// ali contradizia a própria linha de cima do prompt, que declara o regime. Um agente
	// que confiasse nesta seção criaria a spec que a camada proíbe — e o `anchors new`
	// depois a recusaria, sem que nada explicasse a contradição.
	if artifact == "review" || artifact == "review-plan" {
		// O revisor não vê "onde as peças nascem" — ele vê O QUE EXISTE, que é o material
		// do confronto. E vê o registro de entrega, que dá escopo e traz a intenção
		// declarada pelo autor para ser confrontada contra o disco.
		b.WriteString("\n## What you are going to confront\n\n")
		writeTriadPaths(&b, rel, artifact, layer, cfg, g)
		writeDeliveryRecord(&b, root, rel, cfg)
	} else if hasLayer && l.Regime == "declarativo" {
		b.WriteString("\n## The pieces and where they are born\n\n")
		fmt.Fprintf(&b, "**Only the file itself.** `%s` is in a RECOGNIZED layer "+
			"(`regime: declarativo`): it does not originate rules — it translates, transports or "+
			"declares. That is why it **has no spec, feature or test of its own**.\n\n", rel)
		b.WriteString("The rule this file serves lives in the layer that DECIDES (the one that " +
			"consumes it). If you miss a spec here, the most likely thing is that " +
			"the decision is missing there — not that this layer needs one.\n")
	} else {
		b.WriteString("\n## The pieces and where they are born\n\n")
		writeTriadPaths(&b, rel, artifact, layer, cfg, g)
	}

	// REGIMES: as tags de nível que os cenários da feature DEVEM declarar. Estão no
	// anchors.yaml (`derived.regimes`) e nenhum comando as listava — um agente escreveu
	// `@integracao` (português, coerente com o resto do Gherkin) quando o certo
	// era `@integration-level`, e só descobriu garimpando o YAML.
	if artifact == "feature" || artifact == "test" {
		writeRegimes(&b, cfg)
	}

	// FRONTEIRA: o que NÃO é desta etapa. É o que o polvo chama de "belongs to the
	// layer agents" — sem isso o agente extrapola o escopo.
	b.WriteString("\n## What is NOT your scope\n\n")
	for _, s := range outOfScope(artifact) {
		fmt.Fprintf(&b, "- %s\n", s)
	}

	// PROCEDIMENTO: universal + override do projeto (híbrido).
	b.WriteString("\n## Procedure\n\n")
	// Em camada RECONHECIDA o procedimento padrão não serve: ele manda "leia a spec
	// inteira; ela é a régua", e ali não existe spec. Mandar ler o que não existe deixa o
	// agente sem chão — ou pior, o convence a criar a spec proibida.
	passos := procedureFor(artifact, cfg)
	if hasLayer && l.Regime == "declarativo" {
		passos = declarativeProcedure(rel)
	}
	for i, s := range passos {
		fmt.Fprintf(&b, "%d. %s\n", i+1, s)
	}
	if extra := l.Work[artifact]; len(extra) > 0 {
		b.WriteString("\n**Steps specific to this layer** (declared in anchors.yaml):\n\n")
		for _, s := range extra {
			fmt.Fprintf(&b, "- %s\n", s)
		}
	}

	// O QUE OS GATES VÃO COBRAR — dito ANTES, não depois.
	//
	// O relato foi este: "o guide diz 'rode `anchors check` para saber o formato'. Isso é
	// feedback pós-fato: eu escrevo, rodo, e só então descubro." Deu certo por copiar os
	// vizinhos, não porque a régua tenha dito. Um requisito que o autor só conhece depois
	// de reprovar é um requisito mal comunicado — o gate SABE o que exige, então dizer
	// antes custa nada e economiza uma rodada.
	if regras := gateRequirements(artifact, cfg); len(regras) > 0 {
		b.WriteString("\n## What the gates will demand (read BEFORE writing)\n\n")
		for _, r := range regras {
			fmt.Fprintf(&b, "- %s\n", r)
		}
	}

	// REGRA NOVA / REGRA VIOLADA PEDE RÉGUA. O buraco simétrico ao do §5.1 do QUALITY: lá
	// se caça a âncora sem gate; aqui, a REGRA sem gate — um dever declarado em prosa que
	// nenhum medidor confronta. Enquanto vive só no guide, vale enquanto alguém lembra.
	//
	// O caso que originou isto: um projeto tinha o combinado "todo sheet usa a lib de
	// bottom-sheet, nunca o modal nativo". Não era gate. Uma varredura achou 20 usos do
	// modal nativo em 11 arquivos, invisíveis por meses — e o custo não era estético: as
	// duas famílias aninhadas produzem falha SILENCIOSA (o sheet não abre, sem erro).
	// Virou uma entrada declarativa de fronteira em minutos. A regra existia; a régua não.
	b.WriteString("\n## If you ADD a rule — or CATCH a rule failing\n\n")
	b.WriteString("In both cases, ask: **does this fit in a gate?** Fixing the violation is " +
		"half the work; the other half is the ruler that prevents the relapse. A rule " +
		"that was violated once proves, by construction, that human memory does not sustain it.\n\n")
	b.WriteString("Order of preference of the measurer:\n\n")
	b.WriteString("1. **Declarative deterministic** — the rule is recognizable by pattern in the " +
		"artifact (forbidden import, required field, name out of convention): it fits in a " +
		"configuration entry of the framework, with no new code. Almost zero cost — try this first.\n")
	b.WriteString("2. **Deterministic with code** — the rule requires traversing structure " +
		"(resolving an edge, comparing two artifacts): it asks for a gate of its own.\n")
	b.WriteString("3. **AI judgment** — the rule is semantic and without fixed form " +
		"(\"the comment describes what the code does\").\n\n")
	b.WriteString("Only if none of the three serve does the rule stay in prose — and then it is declared as " +
		"explicit debt, not forgotten in silence. When you finish fixing a violation: " +
		"*\"what stops this from coming back tomorrow?\"*. If the answer is \"someone will remember\", " +
		"the work is not finished. (See QUALITY.md §5.1, \"The rule without a ruler\".)\n")

	// ISSUES ABERTAS DESTA UNIDADE. O ciclo escreve o achado do revisor numa issue e depois
	// nunca mais o lê: medido num E2E, os cinco achados mais valiosos de uma spec vieram de
	// uma issue que o próprio Anchors tinha escrito, e sem alguém trazê-la à mão a etapa
	// seguinte sairia com "✗0 bloqueantes" exatamente igual.
	//
	// É a falha de memória do framework: o defeito foi encontrado, registrado, e some do
	// caminho de quem poderia corrigi-lo. Trazê-lo para o prompt é o que fecha o laço entre
	// quem acha e quem conserta.
	writeOpenIssues(&b, root, rel)

	// SINAIS DE EXECUÇÃO. Os gates que medem o que só a EXECUÇÃO revela (o teste passa? o
	// teste PROVA a linha, ou só a executa?) leem de sinais ingeridos — e ficam `~` para
	// sempre se ninguém os ingerir. Medido num E2E real: `mutation-score` saiu `~828` no
	// repositório inteiro, e era exatamente o gate que teria pego 3 cenários decorativos
	// (o teste citava a restrição e não a provava). O instrumento existia, configurado, e
	// nunca tinha sido rodado.
	if artifact == "test" || artifact == "review" {
		b.WriteString("\n## Execution signals (what the gate does not see on its own)\n\n")
		b.WriteString("A green gate does not prove the test PROVES the rule — it proves that it runs. " +
			"What separates the two is mutation: change the line and see if the test falls.\n\n")
		b.WriteString("```sh\n")
		b.WriteString("# 1) the suite, and the result for the map\n")
		b.WriteString("#    (the gate `testes-passam` reads from here: without it it stays ~, and the gate that\n")
		b.WriteString("#     counts `it()` in the source gives ✓ even in a file that does not COMPILE)\n")
		b.WriteString("anchors ingest --junit <JUnit report of the suite>\n\n")
		fmt.Fprintf(&b, "# 2) mutation OF THIS unit (the whole project is expensive; scope it to the target)\n"+
			"#    run the project mutator over `%s`, then:\nanchors ingest --mutation <report.json>\n```\n\n", rel)
		b.WriteString("> A SURVIVING mutant is a change in the code that no test " +
			"noticed: there the rule is not proved. It is the finding the deterministic gates " +
			"do not reach.\n")
	}

	// COMO REGISTRAR A ENTREGA. Ausente de todo prompt até aqui: a seção "Como verificar
	// que terminou" acabava no `check`, e o `anchors deliver` — que é o que ENFILEIRA a
	// próxima etapa — não era citado em lugar nenhum. Medido em três execuções: o worker
	// terminava, nada acontecia, e o orquestrador tinha de descobrir o comando e suas
	// quatro flags obrigatórias por tentativa e erro (o CLI as revela uma por execução).
	//
	// Pior: o `deliver` repreende quem omite `--decision`/`--uncovered` ("você declarou
	// ZERO decisões livres… isso AFIRMA que a régua decidiu tudo") — a régua penalizando
	// uma omissão que ela própria induziu ao não pedir.
	if artifact != "review" && artifact != "review-plan" && artifact != "review-plan-draft" {
		b.WriteString("\n## How to RECORD the delivery (it is this that queues the next stage)\n\n")
		b.WriteString("```sh\n")
		fmt.Fprintf(&b, "anchors deliver --stage %s --unit %s \\\n", artifact, rel)
		b.WriteString("  --date <YYYY-MM-DD> \\\n")
		b.WriteString("  --intent \"what you set out to do, in one sentence\" \\\n")
		b.WriteString("  --decision \"every choice the ruler did NOT decide for you\" \\\n")
		b.WriteString("  --uncovered \"what was left without proof, and why\"\n```\n\n")
		b.WriteString("`--decision` and `--uncovered` accept the flag repeated (one per item). " +
			"Declaring ZERO decisions ASSERTS that the ruler decided everything — if you chose " +
			"anything on your own (a name, a limit, an order), it is a decision. It is what " +
			"the reviewer will confront against the disk.\n")
	}

	// FECHAMENTO DO REVIEW. O review é a única etapa que não produz artefato — ela produz
	// ACHADOS —, e por isso era a única que terminava sem dizer o que fazer com o
	// resultado. Medido num E2E real: o prompt repetia três vezes "NÃO corrija — quem
	// acha não conserta" e não emitia comando de registro. O ciclo `deliver → review → ?`
	// não fechava: o achado ficava no relatório do subagente, e sumia com ele.
	//
	// Pior era o ACHADO CRUZADO: o revisor de uma unidade encontrou o bug crítico numa
	// unidade IRMÃ e, corretamente, não o corrigiu. Nenhuma task foi enfileirada, porque
	// nada no framework recebia um achado sobre outro alvo. Sem alguém ler o relatório à
	// mão e rotear, o bug seguiria verde e entregue.
	if artifact == "review" || artifact == "review-plan" || artifact == "review-plan-draft" {
		b.WriteString("\n## How to CLOSE the review (REQUIRED)\n\n")
		b.WriteString("The review produces no artifact — it produces FINDINGS, and a finding that stays in your " +
			"report dies with the session. Record EACH one, on the target it belongs to:\n\n")
		b.WriteString("```sh\n# no finding:\n")
		fmt.Fprintf(&b, "anchors judge %s --gate review --verdict pass --reason \"what you confronted, and why it passed\"\n\n", rel)
		b.WriteString("# with findings (the --reason IS the body of the issue: a full report, not a verdict sentence):\n")
		fmt.Fprintf(&b, "anchors judge %s --gate review --verdict fail --reason \"$(cat <<'EOF'\n", rel)
		b.WriteString("## Report\n### 1. <what>  (file:line)\n- **why:** the rule this breaks\n" +
			"- **how to fix:** the path\nEOF\n)\"\n```\n\n")
		b.WriteString("**A finding in a unit DIFFERENT from the one you reviewed?** Record it on ITS " +
			"target, not on yours — `anchors judge <the-other-unit> --gate review --verdict fail`. " +
			"It is the only way the finding becomes work: whoever reviewed unit A does not fix " +
			"unit B, but needs to leave the defect where someone will find it.\n\n")
		b.WriteString("> You still do **not correct** — whoever finds does not fix. Recording is not " +
			"correcting: it is making the finding survive you.\n")
	}

	// VERIFICAÇÃO: sempre o mesmo comando, sempre citado.
	b.WriteString("\n## How to verify you are done\n\n")
	b.WriteString("```sh\nanchors map build\n")
	// O `--changed` recebe a peça que ESTA etapa produz, não o alvo da unidade: na etapa
	// `spec`, o `.ts` do alvo em geral ainda não existe, e o `check` aborta sem conferir
	// nada. Prescrever um comando que a própria etapa impede de funcionar é a régua se
	// contradizendo dentro do mesmo prompt.
	alvo := verificationTarget(rel, artifact, layer, cfg)
	fmt.Fprintf(&b, "anchors check --changed %s --no-record --deterministic\n```\n\n", alvo)
	b.WriteString("Every **blocking** gate must come out with `✗0`. A `~` is not a failure — it is the " +
		"gate saying it had nothing to confront (the reason comes written). If a gate fails, " +
		"fix it before considering the stage finished.\n\n")
	// O `--no-record` acima é para ITERAR: confere sem efeito colateral, quantas vezes
	// precisar. Mas a etapa que FECHA tem de registrar — o carimbo é a memória de que
	// aquela aresta foi validada, e é dele que o `anchors stale` vive.
	//
	// Medido: com o `work` prescrevendo só a forma `--no-record`, o carimbo nunca
	// acontecia durante o ciclo, e `anchors stale` reportava 9.975 de 9.975 arestas
	// "nunca validada" — um comando inutilizável, que parecia bug do framework e era o
	// prompt nunca mandando registrar. Um `check --all` bastou para derrubar o número a
	// 19.
	b.WriteString("When the gates are green, run one last time **without** " +
		"`--no-record`:\n\n```sh\n")
	fmt.Fprintf(&b, "anchors check --changed %s --deterministic\n```\n\n", alvo)
	b.WriteString("It is what STAMPS the edges of this unit as validated. Without that step the " +
		"map does not keep that the confrontation happened: the work is finished and `anchors stale` " +
		"keeps saying nobody ever looked.\n")

	// OPT-OUTS: a saída LEGÍTIMA quando um gate reprova algo que é decisão consciente.
	// Sem esta seção, o agente descobre a existência da dispensa pela mensagem de erro do
	// gate — foi o que aconteceu num teste real: `@no-scenario` não estava em guide algum,
	// e o autor teve de adivinhar a sintaxe por tentativa e erro.
	b.WriteString("\n## If a gate fails something that is a conscious decision\n\n")
	b.WriteString("There is a legitimate door, and it **requires the written reason** — a bare marker " +
		"waives nothing. The waiver stays on the line of what it waives (or in the comment " +
		"right above), dated by git and visible to whoever reads later:\n\n")
	switch artifact {
	case "spec":
		b.WriteString("- `@no-mark: <reason>` on the line of a rule that will NOT have a mark in the code " +
			"(a restriction satisfied by absence, an invariant proved by test). " +
			"**If while implementing the reason does not hold** — the rule after all needs code, " +
			"or should not exist —, that is not a formatting detail: it is the spec being wrong " +
			"about its own unit. Open an issue (`anchors judge <target> --gate review " +
			"--verdict fail --reason \"…\"`) instead of erasing the rule or inventing a " +
			"comment to silence the gate.\n")
		b.WriteString("- `@no-scenario: <reason>` on the line of a requirement that will have no scenario " +
			"— for what is truly not observable by scenario, not for what is hard to test.\n")
	case "code":
		b.WriteString("- `@no-paginate: <reason>` on a function that promises the whole set and does not paginate " +
			"(the limit is deliberate and you know the ceiling).\n" +
			"- `@allow-boundary: <reason>` on a line that crosses a layer boundary — acknowledged " +
			"debt, visible in the code instead of in a distant list.\n")
	}
	b.WriteString("See the complete list in `anchors guide header`. Waiving is a recorded " +
		"decision; erasing the problem is not.\n")

	// REGRA DE OURO: o que fazer quando a régua não decide.
	b.WriteString("\n## If the ruler does not decide\n\n")
	b.WriteString("**Do not guess.** If the spec/guide admits two readings with different " +
		"consequences, or does not decide something the code will need, the record has a PLACE — " +
		"it is not a comment on the PR, which dies at the merge:\n\n")
	if artifact == "spec" {
		// O TÍTULO da seção e o valor "nenhuma" vêm do mesmo catálogo que o
		// `anchors new spec` usa: mandar o agente escrever numa seção com outro nome que
		// o do artefato gerado é mandá-lo criar uma seção que gate nenhum encontra.
		fmt.Fprintf(&b, "- **Write the question in `## %s`**, in the spec you ", openSectionTitle())
		b.WriteString(
			"are producing. While there is an item there, the spec does not pass as finished (gate " +
				"`open-questions-resolved`) — which is the desired effect: implementing with an open " +
				"question is guessing, and guessing is confronted by no gate, because " +
				"every piece exists and references the others.\n" +
				"- **Do not invent the answer to empty the section.** When it comes, PROMOTE it " +
				"to a rule (with a code) and mark the item as resolved — the question stays in the trace.\n" +
				"- **If there is no doubt, write `" + noneValue() + "`.** Asserting that you looked is different from " +
				"omitting the section.\n")
	} else {
		fmt.Fprintf(&b, "- **Stop and read `## %s` in the spec.** ", openSectionTitle())
		b.WriteString("If there is an item there, " +
			"the decision has not been made yet: do not make it alone in the code.\n" +
			"- **If the gap is new** (you discovered it while implementing), it belongs to the SPEC. " +
			"Record it there and report — do not resolve it in the code, where nobody will find it.\n" +
			"- Choose the most defensible reading so as not to stall, but leave the choice " +
			"VISIBLE where it is decided. A silent decision becomes a bug nobody traces.\n")
	}
	return b.String(), nil
}

// guidesFor devolve os guides a ler para produzir `artifact` sobre um alvo da camada
// `l`. São DUAS origens, e ambas importam:
//   - o guide do ARTEFATO (quem rege a tag `spec`/`feature`/`test`) — como se escreve
//     uma spec neste projeto, qualquer que seja a camada do alvo;
//   - os guides da CAMADA do alvo (quem rege as tags dela) — o que esta camada exige.
//
// Considerar só a camada esconderia o SPEC_GUIDE de quem vai escrever a spec de um
// repository (o alvo tem a tag `repository`, não `spec`) — que é justamente o guide
// mais relevante da tarefa.
func guidesFor(l config.Layer, cfg *config.Config, artifact string) []string {
	tags := map[string]bool{artifact: true}
	for _, t := range l.Tags {
		tags[t] = true
	}
	var out []string
	seen := map[string]bool{}
	for _, g := range cfg.Governs {
		if tags[g.Governs] && !seen[g.From] {
			seen[g.From] = true
			out = append(out, g.From)
		}
	}
	sort.Strings(out)
	return out
}

// writeTriadPaths mostra onde cada peça da trinca nasce para este alvo, usando o
// `derived:` do projeto (co-location por padrão, overrides por camada).
func writeTriadPaths(b *strings.Builder, rel, artifact, layer string, cfg *config.Config, g *mapx.Graph) {
	if cfg.Derived == nil {
		b.WriteString("> The project does not declare `derived:` — confirm where the pieces live " +
			"by looking at the layer neighbors.\n")
		return
	}
	files, overridden := derivedPaths(rel, layer, cfg)
	dispensadas := waivedPieces(layer, cfg)
	order := []string{"spec", "feature", "test"}
	for _, k := range order {
		tpl, ok := files[k]
		if !ok {
			continue
		}
		p := tpl
		mark := " "
		if k == artifact {
			mark = "→" // a peça desta etapa
		}
		exists := ""
		if _, err := os.Stat(p); err == nil {
			exists = "  (already exists)"
		}
		note := ""
		if overridden[k] {
			note = "  ← layer override (not co-located)"
		}
		// A camada pode DISPENSAR uma peça (`trinca_opcional`), e o prompt tem de dizer
		// isso: listá-la como "onde nasce" faz o executor obediente criar exatamente o
		// artefato que a régua declarou não querer. Aconteceu num E2E real — a camada
		// dizia "opt-out honesto, em vez de 49 features/testes vazios", e o mesmo prompt
		// mandava criar os dois.
		if dispensadas[k] {
			fmt.Fprintf(b, "  `%s` — %s  ← WAIVED by this layer (`trinca_opcional`): "+
				"do NOT create\n", p, k)
			continue
		}
		fmt.Fprintf(b, "%s `%s` — %s%s%s\n", mark, p, k, exists, note)
	}
	b.WriteString("\n> Paths already resolved for THIS layer (co-location + overrides of " +
		"`derived:`). Do not create the piece anywhere else.\n")
}

// outOfScope é a FRONTEIRA do papel: o que pertence a outra etapa. Universal por
// artefato — a doutrina do Anchors, não do projeto.
func outOfScope(artifact string) []string {
	switch artifact {
	case "review-plan-draft":
		return []string{
			"**Do not review the CODE** — it does not exist yet. Your target is the plan as a script: " +
				"is it executable by someone who was not in the conversation that generated it?",
			"**Do not confuse it with the whole review** (`review-plan`), which runs at the END and " +
				"confronts the seam of what was delivered. This one runs BEFORE, and confronts the script.",
			"**Do not accept \"to be decided later\" inside an item.** Outside the checklist, in " +
				"prose, with whoever decides — that is the correct behavior and is not a finding.",
		}
	case "review-plan":
		return []string{
			"**Do not review the isolated piece** — that was already done. Your target is what only appears " +
				"when the pieces meet.",
			"**Do not correct anything.** Report; whoever finds does not fix.",
			"**Do not accept \"each part passed\" as proof of the whole.** It is the premise the 4 " +
				"findings that cross units refute.",
		}
	case "review":
		return []string{
			"**Do not correct the code** — not even the \"obviously wrong\". Report; whoever finds does not fix.",
			"**Do not rewrite the spec.** If it does not decide enough, that IS the finding.",
			"**Do not trust the green.** Every gate passing is the NORMAL state of work " +
				"that gets here — that is exactly why the review exists.",
			"**Do not report what you did not execute.** \"It seems that\" and \"probably\" are not findings; " +
				"command + output are.",
		}
	case "spec":
		return []string{
			"**Do not write code, feature or test** — this stage produces only the spec.",
			"**Do not describe implementation** (internal structures, algorithm, variable names): the spec describes observable BEHAVIOR and contract.",
			"**Do not decide for another layer**: if the decision belongs to the schema/the rule/the interface, point there instead of replicating.",
		}
	case "code":
		return []string{
			"**Do not rewrite the spec to fit the code** — if the spec is wrong, report it; do not adjust it in silence.",
			"**Do not invent a rule the spec does not have**: new behavior requires a spec first.",
			"**Do not mix layers**: respect what the target layer may import (see `layers:`).",
		}
	case "feature":
		return []string{
			"**Do not create a new scenario code**: every `@CODE-X##` of the feature must exist in the spec.",
			"**Do not describe implementation** in the steps — Gherkin speaks of observable behavior.",
			"**Do not write the test here** — the feature declares the scenarios; the test proves them.",
		}
	case "test":
		return []string{
			"**Do not test what the feature does not declare**: each test traces to a scenario with a code.",
			"**Do not mock the logic under test** (a mock that replicates the rule proves itself, not the code).",
			"**Do not change the code to make the test pass** without understanding whether the defect is in the code or in the test.",
		}
	}
	return nil
}

// procedureFor é o procedimento UNIVERSAL do artefato. O projeto acrescenta os seus
// passos via `layers.<camada>.work.<artefato>` (híbrido).
func procedureFor(artifact string, cfg *config.Config) []string {
	switch artifact {
	case "review-plan-draft":
		// O review do PLANO RECÉM-ESCRITO, antes de alguém executá-lo.
		//
		// É o único ponto do ciclo em que corrigir custa uma edição de texto. Depois, o
		// defeito do plano já se espalhou pelas specs que ele semeou — e o custo é
		// assimétrico: um erro na spec afeta uma unidade; um erro no plano afeta todas.
		//
		// A régua vem de defeitos medidos em planos reais: item cuja existência depende de
		// decisão não tomada (o executor decidiu sozinho), ambiguidade que três rodadas
		// seguidas tropeçaram na mesma pedra, e caminho que a estrutura já mudou.
		return []string{
			"**Hunt the item whose EXISTENCE depends on a decision not made.** Distinguish: " +
				"delegating WHERE it is decided is legitimate (\"the indexes stay in the schema spec, " +
				"decided when looking at the current schema\" — the item exists, the content lives in another " +
				"artifact). Postponing WHETHER the item exists is the defect (\"it is born if mobile calls " +
				"for X; otherwise mark it out of scope\") — whoever executes will decide alone or " +
				"stop the flow. It happened: the agent decided, and nobody had asked.",
			"**Confront every cited path against the DISK.** The plan cites files and " +
				"directories; the structure changes. A plan that tells you to edit a file that was " +
				"split guides the executor to the wrong place — and three consecutive rounds " +
				"stumbled on the same ambiguity (`models/` in two places with opposite regimes).",
			"**Confront the plan against the declared Structure.** Does it seed a spec in a layer " +
				"that HAS a spec? Does it tell you to create a piece in a declarative layer? The `anchors work spec " +
				"--for <target>` answers per target — use it on the items you doubt.",
			"**Read each item asking \"what does whoever executes this need to know and is not " +
				"here?\"** — not \"is it well written?\". The plan is executable or it is not.",
			"**Check the order.** An item that depends on another must come after it — and the " +
				"plan must SAY the dependency, not leave it implicit in the numbering.",
			"**Is the Definition of Done verifiable?** Each criterion must be confrontable " +
				"by someone who did not write the plan. \"It works well\" is not a criterion.",
			"**Do NOT rewrite the plan.** Report; whoever finds does not fix — and a plan " +
				"rewritten by the reviewer loses its owner.",
		}
	case "review-plan":
		// O review de CONJUNTO. Derivado da medição: dos 6 achados de um review adversarial
		// real, apenas 2 cabiam no escopo de uma unidade. Os outros 4 atravessavam — e em
		// todos eles CADA PEÇA estava correta sozinha, que é exatamente por que o review
		// por unidade não os alcança e nenhum gate os pega.
		return []string{
			"**Read the whole plan and the records in `changes/`.** Your scope is the SEAM, " +
				"not the piece: the per-unit review already guaranteed each one. Start by listing the " +
				"delivered units and how they connect.",
			"**Follow the DATA end to end.** Who writes, who reads, who deletes. The case that " +
				"motivated it: the model spec said to plant the ownership field, the DAO did not " +
				"plant it and the handler did not fill it — three correct pieces, and the record created " +
				"by Lambda stayed invisible to its owner on the screen.",
			"**Confront the PROMISE between layers.** Is what the lower layer delivers what the " +
				"upper one assumes? The case: the DAO truncated the query without paginating, and the rule that " +
				"consumes it picks \"the version with the greatest month\" — truncating gives no error, it gives the " +
				"WRONG answer in silence.",
			"**Confront RULE against RULE**, between specs of different layers and within the " +
				"same spec. No gate reads two rules together. An invariant declared in one " +
				"spec can be violated by the code another spec authorizes.",
			"**Prove the cross-cutting obligation by EFFECT, not by structure.** A duty can " +
				"pass in every piece (the table appears in the handler) and not work as a whole. " +
				"Mutate the block that fulfills it and run the suite: if nothing falls, the duty is not proved.",
			"**Confront the delivery records against the disk.** Each `changes/*.md` asserts " +
				"what was done, what was decided alone and what is not proved. " +
				"Divergence between the declared and the real is a finding — and it has already appeared: a record " +
				"claimed to have updated counts of a file that was not touched.",
			"**Check the plan's Definition of Done**, item by item, with evidence of " +
				"execution. What was left undone is a finding; what was left undone and was marked " +
				"as done is a serious finding.",
			"**Classify by severity and do NOT correct.** Your output is the report; whoever finds " +
				"does not fix.",
		}
	case "review":
		// Esta régua é DERIVADA de defeitos reais. Em três rodadas de um E2E, 7 defeitos
		// graves passaram com todos os gates verdes; nenhum foi achado lendo código. Os
		// que apareceram vieram de ATAQUE — mutar a regra, rodar com entrada de borda,
		// confrontar campo do modelo contra a regra que o interpreta. Por isso o
		// procedimento manda executar, não revisar no sentido usual.
		return []string{
			// O LUGAR do registro depende do modo, e o prompt não o repete aqui: a
			// seção de registros acima já diz onde ele está (arquivo no modo local,
			// comentário da issue no modo github). Nomear `changes/` neste passo mandava
			// o revisor procurar no disco um registro que vive na issue — medido na
			// revisão do card #321.
			"**Read the delivery record** (the section above says where it is) and confront the " +
				"declared INTENT against what is on the disk. Divergence between what the author " +
				"says he did and what he did is, by itself, a finding.",
			"**Confront spec ↔ code, rule by rule.** Is each catalogued rule " +
				"implemented? Was any implemented backwards? Is there behavior in the code " +
				"that no rule governs? (a model field nobody reads, a branch the spec " +
				"does not foresee).",
			"**MUTATE the rule and run the suite.** Delete the line that implements each rule; if the " +
				"tests stay green, that test does NOT prove that line. It is the only " +
				"instrument that measures the power of the test — and it has already caught a rule with 12 green tests " +
				"and zero proof. Restore the file after each mutation.",
			"**ATTACK by execution, with adversarial input.** Write a script that calls the " +
				"unit with: empty, limit, a value the spec allows but the code does not handle, " +
				"strange order, a special character in the key. Report with the command and the real " +
				"OUTPUT — not with what you deduced by reading.",
			"**Confront the rules against each other.** Two rules of the same spec can contradict each other " +
				"(one says absence, another says value, for the same case). A rule below can " +
				"contradict the one above. No gate confronts rule↔rule — this is the place.",
			"**Confront the PROSE against what it asserts.** Every sentence that describes ANOTHER " +
				"artifact — \"the contract emits three states\", \"the plan puts this out of scope\", " +
				"\"gate X demands this\" — is a verifiable assertion, and no gate verifies it. " +
				"Open the cited artifact and read it. Measured in ONE session: a rule demanded what the " +
				"plan put out of scope; three specs promised an obligation that was not " +
				"declared; a rule counted three values of a set the contract closes " +
				"at four. All three passed every gate — `ref-resolves` confirms that " +
				"the cited code EXISTS, not that the sentence about it is true. Distrust it " +
				"more when the prose CITES the source: citing gives authority without giving proof.",
			"**And confront the prose against the CODE, not only against other specs.** It is the same " +
				"attack on a worse target: the code is where the value actually lives. Measured: a spec " +
				"declared the closed set as `crítico` and the code exported `'critico'` — " +
				"DIFFERENT keys in a value that crosses a boundary as JSON and becomes a `case` of a " +
				"switch. The `as const` protected the contract side (it does not compile) and NOT the " +
				"consumer side, where the `case` would never match and the color would fall to the default in silence. " +
				"Open the file the spec says realizes the rule and compare the LITERALS, not the " +
				"meaning.",
			"**And distrust prose that describes the FUTURE.** \"Plan X will use this\", " +
				"\"phase Y will need it\" — there is no artifact to open and compare, so it is " +
				"UNVERIFIABLE, and that is worse than false. Measured: a spec justified two " +
				"architecture decisions with what two future plans would do, and neither of the two " +
				"mentioned them. A justification rests on what EXISTS today, plus a hypothesis " +
				"declared as a hypothesis — if the plan does not deliver, the decision is left with a reason " +
				"that does not check out, and whoever reviews it later does not know whether the rule still holds.",
			"**The exit is to REQUIRE instead of describe.** Not every mention of what does not yet " +
				"exist is a prediction: referencing a spec the PLAN seeds is legitimate — the plan " +
				"declares it, and there is a gate confronting the path. What does not hold is assuming what it " +
				"WILL DO. Measured: a rule said \"the consumers are named\" and listed two, " +
				"with only one existing; if the second implemented it another way, the argument " +
				"lost the half that sustained it. Rewritten as a requirement (\"every session " +
				"output goes through this operation\"), it holds even if the consumer does not exist — " +
				"and it finds it ready instead of deciding alone. A requirement DECIDES; a " +
				"description depends on a confirmation that cannot come yet.",
			"**And prose AGES: check whether a pending decision has already been made.** The " +
				"cases above are born wrong; this one is born RIGHT and rots. A spec that says \"unit " +
				"X decides this\" or \"assumes the worst case\" is correct while X " +
				"has not decided — and stays in the file after it decided. Measured: a rule " +
				"treated token rotation as a defensive hypothesis; the pool spec made it " +
				"fact, and nobody went back. The reading changes: \"we serialize as a precaution\" is what " +
				"someone removes in a refactor; \"we serialize because it rotates\" is not. No gate " +
				"flags it — `plano-alterado-justificado` demands from whoever ALTERS, and here nobody " +
				"altered: a new spec made another outdated, and the two remain " +
				"internally coherent.",
			"**And search by the CONTENT, not by the name you imagine.** The reviewer also " +
				"assumes. Measured: I searched for \"business indicators lambda\" to check whether " +
				"the origin existed, did not find it, and wrote in the report that it was a gap — the plan " +
				"seeded it under another name, and it was enough to search for the INDICATORS instead of the label " +
				"I had in my head. A `grep` for the name you expect confirms what you " +
				"assume; one for what the spec SAYS finds what is there. It is the same bias the steps " +
				"above fight, committed on the reviewer's side.",
			"**Check what the piece promises to the WORLD.** Does a declared resource exist in the infra? " +
				"Is a cited env var provided? Was the index the query needs created? Does the typecheck " +
				"pass by accident (`!`, object index) hiding something that only breaks in " +
				"production?",
			"**Classify by severity** (CRITICAL/SERIOUS/MINOR) and write the evidence of each " +
				"finding: command + output, or file:line. A finding without evidence is not actionable.",
			"**Do NOT correct.** Your output is the report and, if there is a critical or serious finding, an " +
				"issue. Whoever finds and whoever fixes are different roles — it is what preserves the " +
				"independence of the next review.",
		}
	case "spec":
		return []string{
			"Confirm the layer of the target and whether it is GOVERNED (has a spec) or RECOGNIZED (`regime: declarativo` → has none).",
			"Read 2–3 NEIGHBORING specs of the same layer: they show the real dialect of the project (which may diverge from the template).",
			"Generate the frame with `anchors new spec <Name> --preset <preset> --out <path>` (see the presets in `anchors new spec --list-sections`).",
			"Fill the sections cataloguing each rule with its code (`{CODE}-<letter><NN>`); the valid letters are in `rule_types` in anchors.yaml.",
			"In the Dependency Table, use `backticks` only on symbols the code WILL use — that becomes a verifiable contract. Free description stays in prose.",
		}
	case "code":
		return []string{
			"Read the whole spec before writing the first line; it is the ruler.",
			"Read 1–2 neighboring files of the same layer to follow the local pattern (imports, error, style).",
			ruleMarking(cfg),
			"If the spec promises a symbol in the Dependency Table, USE that symbol — the gate `dependency-honored` confronts it.",
		}
	case "feature":
		return []string{
			"List the scenario codes the spec declares — the feature covers those, and only those.",
			"Generate the frame with `anchors new feature <Name> --out <path>`.",
			"Write one scenario per observable behavior, with the code tag and the regime tag (`@unit-level` etc., see `derived.regimes`).",
			"The TITLE of the scenario must describe the behavior — the test will mirror it (the gate confronts code AND description).",
		}
	case "test":
		return []string{
			"Read the feature: each of its scenarios becomes a test, with the code in the title.",
			"An INVARIANT (`-I##`) is proved in a CLOSED CYCLE: apply the producer result and verify through the CONSUMER. Building the state by hand and verifying only the read tests the consumer, not the composition — and that is where the bugs hide.",
			"After writing, MUTATE the rule: delete the line it implements and run. If the test stays green, it does not prove that line — fix the test, not the code.",
			"Find out WHERE the test of this layer lives (co-location or `derived.overrides`) before creating the file.",
			"Mirror the scenario title in the `it(...)` — a divergent description is drift the gate flags.",
			"Run the tests and confirm they pass because they assert the right behavior, not by accident.",
		}
	}
	return nil
}

// writeRegimes imprime o vocabulário de tags de regime do projeto — a tradução da tag
// que o cenário declara para o regime canônico, e a superfície onde ele é provado.
func writeRegimes(b *strings.Builder, cfg *config.Config) {
	if cfg.Derived == nil || len(cfg.Derived.Regimes) == 0 {
		return
	}
	b.WriteString("\n## Valid regime tags (use EXACTLY these)\n\n")
	tags := make([]string, 0, len(cfg.Derived.Regimes))
	for t := range cfg.Derived.Regimes {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	for _, t := range tags {
		canon := cfg.Derived.Regimes[t]
		surface := ""
		if cfg.Derived.Surfaces != nil {
			if sfc, ok := cfg.Derived.Surfaces[canon]; ok {
				surface = fmt.Sprintf(" → proved on the surface `%s`", sfc)
			}
		}
		fmt.Fprintf(b, "- `@%s` (regime %s)%s\n", t, canon, surface)
	}
	// O contra-exemplo usa a tag que o projeto REALMENTE declarou, listada logo acima —
	// cravar um par aqui ensinaria um vocabulario que pode nao ser o dele, que e' o defeito
	// que esta mensagem existe para evitar.
	b.WriteString("\nThe tag belongs to the PROJECT and is not translatable: any spelling other " +
		"than the ones listed above makes the scenario be confronted by no gate.\n")
}

// derivedPieceUnit resolve a UNIDADE quando o alvo dado é uma peça derivada
// (spec/feature/test). Prefere o mapa — a aresta `specifies` diz exatamente qual código a
// spec descreve; sem mapa, cai na convenção de nome (tronco + extensões usuais).
func derivedPieceUnit(root, rel string, cfg *config.Config, g *mapx.Graph) (string, bool) {
	layer, _ := scan.Classify(rel, cfg)
	l, ok := cfg.Layers[layer]
	if !ok {
		return "", false
	}
	// É peça derivada? O kind da camada decide (spec/feature/test), não a extensão.
	switch l.Kind {
	case "spec", "feature", "test":
	default:
		return "", false
	}

	// 1) pelo mapa: a spec APONTA o código (`specifies`); feature/test chegam via a spec.
	if g != nil {
		if alvo := targetByEdge(g, rel, "specifies"); alvo != "" {
			return alvo, true
		}
		// feature/test → sobe até a spec (`covered-by`/`tested-by` chegam NELES)
		for _, e := range g.Edges {
			if e.To != rel {
				continue
			}
			if e.Type == "covered-by" || e.Type == "tested-by" {
				if alvo := targetByEdge(g, e.From, "specifies"); alvo != "" {
					return alvo, true
				}
			}
		}
	}

	// 2) pela convenção: mesmo tronco, extensões de código usuais no mesmo diretório.
	dir := filepath.Dir(rel)
	base := filepath.Base(rel)
	for _, suf := range []string{".spec.md", ".feature", ".test.ts", ".test.tsx", ".spec.ts", "_test.go", "_test.py", "_spec.rb"} {
		base = strings.TrimSuffix(base, suf)
	}
	if i := strings.Index(base, "."); i > 0 {
		base = base[:i]
	}
	for _, ext := range []string{".ts", ".tsx", ".go", ".py", ".java", ".rs", ".rb", ".php", ".kt", ".cs"} {
		// ToSlash porque o que se devolve é o id de um nó, não um caminho de disco: no
		// Windows o Join daria "src\x.ts", e o id não casaria com o mapa (gravado com "/").
		cand := filepath.ToSlash(filepath.Join(dir, base+ext))
		if _, err := os.Stat(filepath.Join(root, cand)); err == nil {
			return cand, true
		}
	}
	return "", false
}

func targetByEdge(g *mapx.Graph, from, tipo string) string {
	for _, e := range g.Edges {
		if e.From == from && string(e.Type) == tipo {
			return e.To
		}
	}
	return ""
}

// declarativeProcedure é o procedimento de uma camada RECONHECIDA (`regime: declarativo`).
// Não há spec para ler — a régua é o CONTRATO da camada vizinha que este arquivo serve, e
// o dialeto dos irmãos. O risco característico aqui não é divergir de uma spec: é a camada
// declarativa DECIDIR alguma coisa, virando regra escondida onde ninguém procura.
// openSectionTitle e noneValue vêm do catálogo i18n — o mesmo que o `anchors new spec`
// usa para escrever a seção. Cravá-los aqui faria o prompt mandar escrever numa seção que
// o artefato gerado não tem.
func openSectionTitle() string {
	if t := i18n.TIn(i18n.Current(), "section.title.open"); t != "" {
		return t
	}
	return "Open Decisions"
}

func noneValue() string {
	if t := i18n.TIn(i18n.Current(), "spec_guide.none"); t != "" {
		return t
	}
	return "none"
}

func declarativeProcedure(rel string) []string {
	return []string{
		"This layer **has no spec**: the ruler is the contract of whoever consumes this file, " +
			"plus the dialect of its neighbours. Read 2–3 siblings of the same layer BEFORE writing.",
		"Translate, transport or declare — **do not decide**. If you catch yourself writing a " +
			"business `if`, a default that changes the result, or a rule validation, " +
			"stop: that belongs to the layer that decides, and here it becomes a rule hidden where " +
			"nobody will look.",
		"Keep the surface predictable: same naming, error and return conventions as the " +
			"siblings. Whoever consumes this layer counts on it.",
		"If a decision is missing for you to be able to write, it is missing in the spec of whoever " +
			"consumes it — report it there, do not resolve it here.",
	}
}

// gateRequirements traduz, em requisitos legíveis, o que os gates INTERNOS declarados
// para este artefato vão confrontar. Só descreve gate que o projeto realmente declarou —
// prometer cobrança que não existe é tão ruim quanto esconder a que existe.
func gateRequirements(artifact string, cfg *config.Config) []string {
	// o que cada checker interno exige, em uma frase acionável
	porChecker := map[string]string{
		"spec-sections": "**Every rule must be CATALOGUED** — code + structured place. " +
			"A heading counts (`### ABCDX-B01 — ...`), a table row (`| \\`ABCDX-B01\\` | ... |`) " +
			"or a bold bullet (`- **ABCDX-B01** ...`). A loose mention in prose does NOT count, and " +
			"an unfilled placeholder fails.",
		"has-code": "The file carries at least one **scenario code** (the identity).",
		"header-valid": "The **`@anchors` header** at the top, with the identity: `code:` if this " +
			"artifact owns it, `ref:` if it references the spec.",
		"rule-types": "Each **code letter** (the `B` of `-B01`) must be declared in the " +
			"project `rule_types` vocabulary — and the section that defines it, too.",
		"route-declared":     "A SCREEN spec declares the **route** in the header.",
		"dependency-honored": "Every symbol promised in the **Dependency Table** (in backticks) is used in the code.",
		"spec-feature-match": "Every declared requirement has a **scenario in the feature** (or `@no-scenario: <reason>`).",
		"rule-implemented":   "Every catalogued rule appears **in the code** (the excerpt that realizes it carries its code in a comment) — or is waived on its line with `@no-code: <reason>`, for what is satisfied by the ABSENCE of code. Declare rule by rule: it is what trades guessing for confrontation.",
		"open-questions-resolved": "The section **`## Decisões em aberto`** is REQUIRED and never " +
			"stays empty: either it lists what the spec does not decide, or it carries `nenhuma`. Writing " +
			"`nenhuma` is an ASSERTION — \"I looked and there is no doubt\" —, different from omitting the " +
			"section, which says nothing. Whatever is open there holds the spec until it is decided.",
		"code-reference-valid": "Every **cited code** must exist in the project (an orphan citation fails).",
		"scenario-asserts":     "The RESULT step asserts an **observable result** — not `\\\"effect X is verified\\\"`.",
		"feature-test-match": "Each scenario of the feature has a test that proves it, and the test " +
			"**declares which scenario it proves**: the code goes in the TITLE of the case (`it('[ABCDX-B01] " +
			"…')` or `describe('… [ABCDX-B01]')`), not in a comment — the gate removes comments " +
			"before searching. Marking is ASSERTING that that test proves that scenario; " +
			"a marker on the wrong test creates false traceability, and every relational gate starts " +
			"confronting the wrong pair with everything green.",
		"non-empty":          "The file **cannot be an empty skeleton**.",
		"triad-complete":     "The unit needs the **complete triad** (spec + feature + test).",
		"ref-resolves":       "The `ref:` must point to the **`code:` of the sibling spec** — not to another.",
		"pagination-honored": "A function that promises a set **does not return the first page** in silence.",
		"layer-boundary":     "Respect the **layer boundaries** declared in `boundaries:`.",
	}
	var out []string
	visto := map[string]bool{}
	for _, gt := range cfg.Gates {
		if !gateApplies(gt, artifact) || visto[gt.Check] {
			continue
		}
		if frase, ok := porChecker[gt.Check]; ok {
			visto[gt.Check] = true
			if !gt.IsBlocking() {
				frase += " *(informative)*"
			}
			out = append(out, frase)
		}
	}
	return out
}

// gateApplies: este gate se aplica ao artefato que está sendo produzido?
//
// O `on:` do gate diz sobre QUAL NÓ ele roda; esta função responde outra pergunta — quem
// precisa CONHECER a exigência ao escrever. Nem sempre é o mesmo.
//
// `feature-test-match` roda sobre a feature e cobra do TESTE: exige que o código do
// cenário apareça no título do caso. Enquanto ele só aparecia na etapa `feature`, quem
// escrevia o teste não era avisado — medido, um agente escreveu 36 features corretas e
// depois esbarrou num gate BLOQUEANTE vermelho sem entender por quê; a exigência só estava
// legível no código-fonte do gate.
func gateApplies(gt config.Gate, artifact string) bool {
	for _, k := range gt.On {
		if k == artifact {
			return true
		}
	}
	// Gates que rodam numa peça e cobram de OUTRA: quem produz a peça cobrada tem de saber.
	if artifact == "test" && gt.Check == "feature-test-match" {
		return true
	}
	return false
}

// writeDeliveryRecord anexa o REGISTRO DE ENTREGA ao prompt do revisor, quando existe.
//
// É a metade declarada do confronto: o autor escreveu o que acha que fez, quais decisões
// tomou sozinho e o que sabe não estar provado. O revisor confronta isso contra o disco —
// e a divergência entre as duas versões é, por si só, um achado que não existe forma de
// detectar sem ter as duas.
//
// Sem registro, o review ainda roda (a unidade está no disco), mas perde o escopo e a
// declaração de intenção. Por isso a ausência é reportada, não silenciada.
func writeDeliveryRecord(b *strings.Builder, root, rel string, cfg *config.Config) {
	pend, _ := change.Pending(root)
	unidade := filepath.ToSlash(rel)
	// TODOS os registros da unidade, não o primeiro. Uma unidade acumula uma entrega por
	// etapa (spec, feature, code, test), e servir só a primeira dá ao revisor a foto mais
	// VELHA: medido, o registro da etapa `code` dizia "NADA está provado: o .test.ts é da
	// próxima etapa" — verdade quando foi escrito, falso quando o review rodou, com o
	// teste já no disco. Um revisor obediente pularia a mutação por achar que não há suíte.
	var achados []string
	for _, p := range pend {
		raw, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if m := changeUnitRE.FindStringSubmatch(string(raw)); m == nil || m[1] != unidade {
			continue
		}
		relPath, _ := filepath.Rel(root, p)
		achados = append(achados, fmt.Sprintf("### Delivery record — `%s`\n\n"+
			"```markdown\n%s\n```\n", relPath, strings.TrimSpace(string(raw))))
	}
	if len(achados) > 0 {
		sort.Strings(achados)
		fmt.Fprintf(b, "\n## Delivery records of this unit (%d)\n\n", len(achados))
		b.WriteString("The author declared what follows, once per stage. **Confront each " +
			"assertion against the disk** — above all \"no free decision\" and \"nothing without " +
			"proof\", which are strong assertions and where the defects hide. Note that " +
			"a record describes the state AT THAT moment: what it says is missing may already " +
			"exist.\n\n")
		for _, a := range achados {
			b.WriteString(a + "\n")
		}
		return
	}
	// NO MODO `github` O REGISTRO NÃO ESTÁ NO DISCO, e dizer "sem registro" ali é um
	// falso negativo — o pior tipo, porque manda o revisor registrar no relatório uma
	// lacuna que não existe.
	//
	// Medido na revisão do card #321: o guia avisou "Sem registro de entrega" e o registro
	// estava nos comentários da issue, onde o `deliver` o põe naquele modo (v0.1.61). O
	// revisor teve de descobrir isso sozinho.
	if cfg != nil && cfg.GitHubMode() {
		b.WriteString("\n> **The delivery record is in the COMMENTS of the issue** of this " +
			"card — that is where `anchors deliver` puts it in `github` mode. Read it: it is the DECLARED " +
			"half of the confrontation, what the author says he did, against the disk.\n")
		return
	}
	b.WriteString("\n> **No delivery record** for this unit (`anchors deliver`). " +
		"You review what is on the disk, but without the author's declared intent — hence, without " +
		"a way to confront what he THINKS he did against what he did. Record that in the report.\n")
}

// derivedPaths resolve, para um alvo e sua camada, ONDE cada peça da trinca nasce —
// aplicando co-location e os overrides do `derived:`, com os placeholders já substituídos.
//
// Existe separada porque duas coisas precisam da mesma resposta: a seção que MOSTRA os
// caminhos e o comando de verificação que o prompt prescreve. Enquanto só a primeira
// resolvia, o comando prescrito apontava para o alvo (`--for <alvo>.ts`) mesmo na etapa
// `spec` — onde o `.ts` ainda não existe e o `check` aborta:
//
//	Error: "packages/.../metadataVersioning.ts" não existe no disco nem no mapa
//
// O mesmo prompt dizia "não escreva código nesta etapa" e prescrevia um comando que só
// funciona com o código escrito. Os dois agentes de spec de um E2E real bateram nisso,
// independentemente.
func derivedPaths(rel, layer string, cfg *config.Config) (map[string]string, map[string]bool) {
	files, overridden := map[string]string{}, map[string]bool{}
	if cfg.Derived == nil {
		return files, overridden
	}
	dir := filepath.Dir(rel)
	// O NOME vem do `mapx`, e não de um corte local.
	//
	// `filepath.Ext("X.spec.md")` é `.md`, então cortar por ela deixa `X.spec` — e o
	// prompt passava a mandar criar `X.spec.ts` e `X.spec.test.ts`, enquanto o mapa (que
	// usa `StemOfAnchor`) liga a trinca por `X`. Duas implementações do mesmo corte, e a
	// deste comando estava errada: quem seguisse o prompt criaria arquivo que nenhum gate
	// encontra.
	name, _ := mapx.StemOfAnchor(rel)
	// O módulo é o diretório-pai — usado por overrides que agrupam por Lambda/módulo
	// (ex.: `packages/backend/__tests__/unit/lambdas/{{module}}.test.ts`).
	module := filepath.Base(dir)
	bruto := map[string]string{}
	for k, v := range cfg.Derived.PadroesDe() {
		// O PRIMEIRO padrão: este prompt mostra ONDE escrever cada peça, e uma lista de
		// caminhos alternativos no lugar de um responderia a pergunta com outra.
		if len(v) > 0 {
			bruto[k] = v[0]
		}
	}
	// OVERRIDE da camada vence a co-location. Resolvê-lo aqui é o ponto do comando:
	// mandar o leitor "conferir os overrides" seria devolver a ele exatamente a
	// descoberta manual que este prompt existe para eliminar.
	for _, ov := range cfg.Derived.Overrides {
		if ov.When != layer {
			continue
		}
		for k, v := range ov.PadroesDe() {
			if len(v) == 0 {
				continue
			}
			bruto[k] = v[0]
			overridden[k] = true
		}
	}
	// `{{ext}}` é a extensão do ALVO, não "ts" fixo. Numa tela (`.tsx`), o valor fixo
	// fazia o prompt prescrever `Tela.test.ts` sob a frase "Não crie a peça em outro
	// lugar" — um arquivo que, criado, é erro de sintaxe (JSX em `.ts`). O arquivo real
	// ao lado era `.test.tsx`, e o mapa o resolvia certo: só o prompt mentia.
	ext := strings.TrimPrefix(filepath.Ext(rel), ".")
	if ext == "" || ext == "md" {
		ext = "ts" // alvo sem extensão de código (ou uma spec): o default do projeto
	}
	r := strings.NewReplacer("{{dir}}", dir, "{{name}}", name, "{{module}}", module, "{{ext}}", ext)
	for k, v := range bruto {
		files[k] = r.Replace(v)
	}
	return files, overridden
}

// verificationTarget devolve o caminho que o `anchors check --changed` deve receber nesta
// etapa: a peça que ELA produz, não o alvo da unidade. Cai no alvo quando a etapa não
// produz peça derivada (código, review).
func verificationTarget(rel, artifact, layer string, cfg *config.Config) string {
	switch artifact {
	case "spec", "feature", "test":
		if files, _ := derivedPaths(rel, layer, cfg); files[artifact] != "" {
			return files[artifact]
		}
	}
	return rel
}

// waivedPieces traduz o `trinca_opcional` da camada (declarado por ARESTA) para as
// PEÇAS que ele dispensa. `covered-by` é a aresta spec→feature, logo dispensa a feature;
// `tested-by` é feature→test, logo dispensa o teste.
func waivedPieces(layer string, cfg *config.Config) map[string]bool {
	out := map[string]bool{}
	if cfg == nil || layer == "" {
		return out
	}
	l, ok := cfg.Layers[layer]
	if !ok {
		return out
	}
	porAresta := map[string]string{"covered-by": "feature", "tested-by": "test", "specifies": "spec"}
	for _, aresta := range l.OptionalTriadEdges {
		if peca := porAresta[aresta]; peca != "" {
			out[peca] = true
		}
	}
	return out
}

// writeOpenIssues traz para o prompt os achados ainda em aberto sobre esta unidade —
// de `issues/todo/` e `issues/doing/`, que são os estados vivos.
func writeOpenIssues(b *strings.Builder, root, rel string) {
	base := strings.TrimSuffix(rel, filepath.Ext(rel))
	for _, suf := range []string{".spec.md", ".feature", ".test", ".spec"} {
		base = strings.TrimSuffix(base, suf)
	}
	slug := strings.ReplaceAll(base, "/", "-")

	var achadas []string
	for _, st := range []issue.State{issue.Todo, issue.Doing} {
		nomes, _ := issue.List(root, st)
		for _, nome := range nomes {
			if !strings.Contains(nome, slug) {
				continue
			}
			achadas = append(achadas, filepath.Join(issue.Dir, string(st), nome))
		}
	}
	if len(achadas) == 0 {
		return
	}
	sort.Strings(achadas)
	b.WriteString("\n## Findings ALREADY RECORDED about this unit (read BEFORE writing)\n\n")
	b.WriteString("Someone already confronted this unit and found what is below. It is not " +
		"history: it is pending work, and the gates do NOT repeat it — they confront what is " +
		"declarable, and these findings are exactly what was left outside that.\n\n")
	for _, a := range achadas {
		fmt.Fprintf(b, "- `%s`\n", a)
	}
	b.WriteString("\n> Read each one. What still holds, resolve in this stage; what no longer " +
		"holds, close (`anchors judge <target> --gate review --verdict pass --reason \"...\"`). " +
		"An issue nobody reads is a defect the pipeline already saw and let through.\n")
}

// ruleMarking emite o passo de ligar regra↔código conforme o projeto a EXIGE ou não.
//
// A ressalva "se o projeto usa esse padrão" existia para não impor a prática a quem não
// a adotou — mas ela também dava saída a quem a adotou: quem implementa lê "se", decide
// que não, e a marcação nunca acontece. O gate `regra-implementada` então cobra algo que
// o procedimento apresentou como opcional, e a dívida só aparece depois de reprovar.
//
// Com `derived.rule_marking: required` declarado, o passo vira obrigação — e o
// procedimento passa a ensinar ANTES o que o gate cobra DEPOIS.
func ruleMarking(cfg *config.Config) string {
	if cfg != nil && cfg.Derived != nil &&
		strings.EqualFold(strings.TrimSpace(cfg.Derived.RuleMarkingPolicy), "required") {
		return "Implement each rule of the spec and MARK in the code the excerpt that realizes it " +
			"(`// {CODE}-B01: …`) — this project REQUIRES the marking, and the gate " +
			"`regra-implementada` confronts it. Whatever has no code, waive on the line " +
			"of the rule (`@no-code: <reason>`)."
	}
	return "Implement each rule of the spec; cite the rule code (`{CODE}-B01`) in the " +
		"comment of the excerpt that realizes it, if the project uses that pattern."
}
