package ops

import (
	"fmt"
	"github.com/co2-lab/anchors/internal/i18n"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/charmbracelet/huh"
	"github.com/co2-lab/anchors/cmd/anchors/governance"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var root string
	var naoInterativo, aceitarDefaults bool
	var f flagsInit
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Configure the project (anchors.yaml) through questions and answers",
		Long: `Scans the project deterministically (without AI), proposes a
Structure, and confirms/adjusts it with you through questions — arriving at a correct
anchors.yaml. The bulk is inferred; the questions cover only the human decisions
(co-location, layer granularity, and which guides govern which tags).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			// O modo não-interativo existe para o agente: sem ele, a guarda de TTY
			// aborta e o fluxo em que o usuário pede a uma IA para iniciar o projeto
			// (BOOTSTRAP.md §5) trava no comando central.
			if naoInterativo {
				return runInitNonInteractive(cmd, absRoot, &f, aceitarDefaults)
			}
			return runInit(absRoot)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&naoInterativo, "non-interactive", false, "no TUI: with no answers, emit the questions as JSON; with answers in flags, apply them")
	cmd.Flags().BoolVar(&aceitarDefaults, "defaults", false, "with --non-interactive and no answers: accept the defaults inferred from disk")
	cmd.Flags().StringVar(&f.preset, "preset", "", "stack preset (see --questions)")
	cmd.Flags().BoolVar(&f.header, "header", true, "seed guides/HEADER_GUIDE.md")
	cmd.Flags().StringSliceVar(&f.artifacts, "artifacts", nil, "anchor kinds of the project (spec,feature,test,…)")
	cmd.Flags().BoolVar(&f.gates, "gates", true, "seed the default gates (informational)")
	cmd.Flags().BoolVar(&f.colocation, "colocation", false, "derived files next to the code")
	cmd.Flags().StringSliceVar(&f.layers, "layers", nil, "code directories to treat as layers")
	cmd.Flags().StringArrayVar(&f.governs, "governs", nil, "governs rule: GUIDE=tag1,tag2 (repeatable)")
	cmd.Flags().StringVar(&f.workflow, "workflow", "", "where the work queue lives: local|manual|github (manual: like local, and no command writes an issue on its own)")
	cmd.Flags().StringVar(&f.repo, "repo", "", "repository owner/name (required in github mode)")
	cmd.Flags().StringSliceVar(&f.labels, "labels", nil, "labels that mark the Anchors cards (required in github mode)")
	return cmd
}

// errNoTTY é a mensagem de quando o `init` interativo não tem terminal. Nomeia a
// SAÍDA, e não só o problema: existe um modo não-interativo, e quem cai aqui (um agente,
// um pipe, o CI) tem como prosseguir — sem esta indicação, o comando parece um beco.
//
// `contexto` diz o que já aconteceu antes da falha, porque isso muda o que o leitor
// precisa saber: se o git foi tocado, se algo foi escrito.
func errNoTTY(contexto string) error {
	return fmt.Errorf("`anchors init` is interactive and there is no terminal available "+
		"(no TTY: pipe, CI or agent without an interactive shell).\n"+
		"%s\n\n"+
		"Use the NON-INTERACTIVE mode, in two calls of the SAME flag:\n"+
		"  1. anchors init --non-interactive\n"+
		"       with no answers, returns the decisions in JSON (options, the default\n"+
		"       inferred from disk, and what each answer changes in the project)\n"+
		"  2. anchors init --non-interactive --artifacts=spec,feature,test --colocation\n"+
		"       with answers in flags, applies them and reports the verdict of each one\n\n"+
		"Or run it in an interactive terminal.", contexto)
}

// runInit orquestra: infere (puro) → pergunta (TUI) → aplica decisões (puro) →
// salva. A lógica de decisão vive em initx (funções puras, testadas); esta função
// só faz o I/O e a coleta de respostas.
func runInit(root string) error {
	outPath := filepath.Join(root, config.DefaultFile)

	if _, err := os.Stat(outPath); err == nil {
		var overwrite bool
		if err := huh.NewConfirm().
			Title(config.DefaultFile + " already exists. Overwrite?").
			Value(&overwrite).Run(); err != nil {
			erroDePrompt = true
		}
		if !overwrite {
			fmt.Println(i18n.T("init.aborted"))
			return nil
		}
	}

	// 0) GIT — antes de escanear. É o substrato: sem versionamento o carimbo de
	// alteração, a cobertura de diff e o pre-commit ficam desligados, e nenhum deles
	// falha ruidosamente. Vem primeiro para que tudo que o init escrever daqui em
	// diante já nasça sob versionamento.
	if !gitStep(root) {
		return errNoTTY("Nothing was written, and git was not touched.")
	}

	fmt.Println(i18n.T("init.scanning"))
	p, err := initx.Infer(root)
	if err != nil {
		return fmt.Errorf("inference: %w", err)
	}
	printFindings(p)

	// 0.4) A FASE ANTERIOR — antes de perguntar qualquer coisa, reconhecer se a fase
	// DESCOBRIR ainda não aconteceu. Sem PROJECT.md e sem código, as perguntas abaixo
	// saem sem resposta boa; o que muda é só QUEM recebe a instrução (pessoa ou IA).
	if !discoverStep(root, p) {
		return errNoTTY("Nothing was written.")
	}

	cfg := p.Config
	empty := len(p.CodeDirs) == 0 && !p.HasSpecMD && !p.HasFeature && !p.HasTest
	if empty {
		fmt.Println(i18n.T("init.empty_project"))
	}

	// 0.5) PRESET DE STACK — oferece uma estrutura consagrada. Opcional: "nenhum"
	// mantém a inferência/edição manual. Se escolhido, preenche as layers de código
	// e, para presets modulares, deduz os prefixos de módulo (liga na identidade).
	chosenPreset, presetOK := askPreset()
	var presetModules []string
	if presetOK {
		presetModules = detectModules(root, chosenPreset)
		prefixes := initx.ApplyPreset(cfg, chosenPreset, presetModules)
		fmt.Printf("Preset '%s' applied (%s, %d layers).\n", chosenPreset.Title, chosenPreset.Pattern, len(chosenPreset.Layers))
		if chosenPreset.Modular && len(prefixes) > 0 {
			fmt.Printf("Detected modules and their identity prefixes (Layer 1):\n")
			for m, pfx := range prefixes {
				fmt.Printf("  %-20s → %s\n", m, pfx)
			}
		} else if chosenPreset.Modular {
			fmt.Printf("(modular preset: the module prefixes will be deduced when there are modules in %s)\n", chosenPreset.ModuleGlob)
		}
		if chosenPreset.CoverageHint != "" {
			fmt.Printf("\nTest signals (for `anchors ingest` to measure real coverage):\n  %s\n", chosenPreset.CoverageHint)
		}
	}

	// 0.6) GUIA DE CABEÇALHO — semeia um HEADER_GUIDE.md no projeto (a régua do bloco
	// @anchors, instanciada com o dialeto de comentário da stack e as features reais).
	// É o padrão MANDATÓRIO de cabeçalho; o init o materializa para o projeto seguir.
	if askConfirmDefault("Seed guides/HEADER_GUIDE.md (the @anchors header standard of the files)?", true) {
		var moduleBasenames []string
		for _, m := range presetModules {
			moduleBasenames = append(moduleBasenames, filepath.Base(m))
		}
		guideDir := p.GuideDir
		if guideDir == "" {
			guideDir = "guides"
		}
		dest := filepath.Join(root, guideDir, "HEADER_GUIDE.md")
		body := initx.RenderHeaderGuide(chosenPreset, moduleBasenames)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err == nil && os.WriteFile(dest, []byte(body), 0o644) == nil {
			fmt.Printf("Header guide seeded: %s\n", filepath.Join(guideDir, "HEADER_GUIDE.md"))
		}
	}

	// 1) ARTEFATOS — SEMPRE perguntado. Pré-marca os detectados; num projeto vazio,
	// o usuário marca o que PRETENDE usar. As layers de artefato são (re)construídas
	// a partir da escolha, não da detecção.
	chosenArtifacts := askMultiSelectPre(
		"Which anchor kinds does your project use (or will use)?",
		initx.ArtifactNames(), p.DetectedArtifacts())
	initx.ApplyArtifactChoice(cfg, chosenArtifacts, map[string]string{
		"guide": p.GuideDir,
		"plan":  p.PlanDir,
	})

	// 1.5) GATES PADRÃO — o ciclo nasce com os gates dos artefatos escolhidos
	// (spec-completa, tests-green, scenario-coverage, …), todos INFORMATIVOS. É o que
	// amarra os sinais de teste ao ciclo de vida sem o usuário escrever tudo à mão.
	// `empty` é o projeto SEM código, spec, feature nem teste — e é ele que decide o
	// estado de maturação: num projeto novo os gates nascem BLOQUEANTES, porque não há
	// débito a acomodar e o primeiro desvio é o mais barato de impedir (ver DefaultGates).
	if gates := initx.DefaultGates(chosenArtifacts, empty); len(gates) > 0 {
		names := make([]string, len(gates))
		for i, g := range gates {
			names[i] = g.Name
		}
		if askConfirmDefault(
			fmt.Sprintf("Seed %d default gate(s) (informative) for the chosen artifacts? [%s]",
				len(gates), strings.Join(names, ", ")), true) {
			cfg.Gates = gates
			if chosenArtifacts["test"] {
				fmt.Println(i18n.T("init.test_gates_note"))
			}
		}
	}

	// 2) CO-LOCATION — SEMPRE perguntado (texto adaptado a detectado vs. a decidir).
	title := "How are the derivatives (spec/feature/test) organized?"
	if p.Colocated {
		title = "I detected the derivatives next to the code (co-location). Use that convention?"
	}
	useColocation := askConfirmDefault(title, p.Colocated)
	initx.ApplyColocation(cfg, useColocation, chosenArtifacts)

	// 3) CAMADAS DE CÓDIGO — sempre perguntado se há dirs detectados; se vazio, avisa.
	if names := initx.CodeLayerNames(cfg); len(names) > 0 {
		keep := askMultiSelect("Which code directories to treat as layers?", names)
		initx.PruneCodeLayers(cfg, keep)
	} else if empty {
		fmt.Println(i18n.T("init.no_code_yet"))
	}

	// 3.5) MODO DE TRABALHO — onde a fila mora. É decisão HUMANA e EXCLUDENTE
	// (WORKFLOW.md §2): não se infere do remote, e não há fallback de um modo para o
	// outro. Perguntada aqui, e não no começo, porque só faz sentido depois de o projeto
	// ter forma — mas ANTES de salvar, porque muda o arquivo.
	if askConfirmDefault("Will the work queue live in GitHub issues (instead of local)?", false) {
		repo := askText("Which repository? (owner/name — never inferred from the remote)")
		if repo != "" {
			cfg.Workflow = &config.Workflow{
				Mode:   config.ModeGitHub,
				Repo:   repo,
				Labels: []string{"anchors"},
			}
			fmt.Println(i18n.T("init.github_mode"))
		} else {
			fmt.Println(i18n.T("init.local_mode"))
		}
	} else if !askConfirmDefault("Should each `anchors check` write its findings as files in issues/? (no = manual mode: the check only reports)", true) {
		cfg.Workflow = &config.Workflow{Mode: config.ModeManual}
	}

	// 4) GOVERNS — só se houver guides (a régua precisa existir para reger).
	if len(p.GuideFiles) > 0 {
		answers := askGovernAnswers(p.GuideFiles, initx.Tags(cfg))
		cfg.Governs = initx.BuildGovernRules(answers)
	}

	// Nada é escrito se algum prompt não pôde rodar. Ver `erroDePrompt`: sem TTY, cada
	// pergunta devolve vazio e o resultado é um config que carrega e não governa nada —
	// escrito com um "✓" na frente. Falhar aqui, ANTES de gravar, é o que impede o
	// arquivo inútil de existir e ser tomado por configuração válida.
	// O guide de SPEC é semeado AQUI, não junto do header: naquele ponto as escolhas do
	// usuário ainda não foram aplicadas ao cfg, e o guide sairia sem o dialeto do projeto
	// (as letras de `rule_types`, o comprimento de `code_lengths`) — genérico, que é
	// exatamente o que ele existe para não ser.
	//
	// Vai antes da guarda de TTY de propósito? NÃO: depois. Um init abortado não deve
	// deixar arquivo, e o header já é um resíduo conhecido — não vamos criar um segundo.
	if erroDePrompt {
		return errNoTTY("Nothing was written: an `anchors.yaml` generated without the answers " +
			"would come out with 0 layers and 0 gates, would load without error and would govern nothing.")
	}
	if err := config.Save(cfg, outPath); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	if p.GuideDir != "" {
		destSpec := filepath.Join(root, p.GuideDir, "SPEC_GUIDE.md")
		if os.MkdirAll(filepath.Dir(destSpec), 0o755) == nil &&
			os.WriteFile(destSpec, []byte(governance.RenderSpecGuide(cfg, "")), 0o644) == nil {
			fmt.Printf("Spec guide seeded: %s\n", filepath.Join(p.GuideDir, "SPEC_GUIDE.md"))
		}
	}
	fmt.Printf("\n✓ %s written (%d layers, %d governs rules)\n", config.DefaultFile, len(cfg.Layers), len(cfg.Governs))
	fmt.Println("  review it and run: anchors map build")
	return nil
}

// --- casca de prompt: só coleta respostas do usuário (huh) ---

// erroDePrompt registra que algum prompt não conseguiu rodar (tipicamente: sem TTY —
// pipe, CI, agente executando `anchors init` sem terminal interativo).
//
// Existe porque descartar esses erros produzia o pior resultado possível: cada prompt
// devolvia a resposta vazia, o `init` seguia até o fim e imprimia "✓ anchors.yaml escrito
// (0 camadas, 0 regras de governs)" — sucesso anunciado, arquivo inútil. E o efeito é em
// cascata, porque as escolhas alimentam umas às outras: sem os artefatos escolhidos,
// `ApplyArtifactChoice` não cria camada e `DefaultGates` devolve zero gate. O usuário
// terminava com um config que carrega, não governa nada, e faz o `anchors check` morrer
// com "nenhum gate declarado".
//
// Medido ao adotar o Anchors no próprio Anchors: o `init` detectou "código em:
// [cli/internal cli/cmd]" na tela, e escreveu `layers: {}` — descartou a própria detecção.
var erroDePrompt bool

// askText coleta uma resposta livre. Usado pelo `repo` do modo github, que não tem
// conjunto de opções para escolher.
func askText(title string) string {
	var v string
	if err := huh.NewInput().Title(title).Value(&v).Run(); err != nil {
		erroDePrompt = true
	}
	return strings.TrimSpace(v)
}

func askConfirmDefault(title string, def bool) bool {
	v := def
	if err := huh.NewConfirm().Title(title).Value(&v).Run(); err != nil {
		erroDePrompt = true
	}
	return v
}

// askMultiSelectPre pré-marca só os itens presentes em `pre` (os demais entram
// desmarcados). Usado quando a inferência sugere um subconjunto, não tudo.
func askMultiSelectPre(title string, items []string, pre map[string]bool) map[string]bool {
	opts := make([]huh.Option[string], 0, len(items))
	for _, it := range items {
		opts = append(opts, huh.NewOption(it, it).Selected(pre[it]))
	}
	var picked []string
	if err := huh.NewMultiSelect[string]().Title(title).Options(opts...).Value(&picked).Run(); err != nil {
		erroDePrompt = true
	}
	set := map[string]bool{}
	for _, p := range picked {
		set[p] = true
	}
	return set
}

// askMultiSelect devolve um set dos itens escolhidos (todos pré-selecionados).
func askMultiSelect(title string, items []string) map[string]bool {
	opts := make([]huh.Option[string], 0, len(items))
	for _, it := range items {
		opts = append(opts, huh.NewOption(it, it).Selected(true))
	}
	var picked []string
	if err := huh.NewMultiSelect[string]().Title(title).Options(opts...).Value(&picked).Run(); err != nil {
		erroDePrompt = true
	}
	set := map[string]bool{}
	for _, p := range picked {
		set[p] = true
	}
	return set
}

// askGovernAnswers pergunta, para cada guide, qual tag ele rege. Devolve o mapa
// guide→tag que BuildGovernRules (puro) transforma em regras.
func askGovernAnswers(guides, tags []string) map[string]string {
	options := append(append([]string{}, tags...), initx.NoneTag)
	answers := map[string]string{}
	for _, guide := range guides {
		var tag string
		if err := huh.NewSelect[string]().
			Title(fmt.Sprintf("The guide %s governs which tag?", filepath.Base(guide))).
			Options(huh.NewOptions(options...)...).
			Value(&tag).Run(); err != nil {
			erroDePrompt = true
		}
		answers[guide] = tag
	}
	return answers
}

// askPreset oferece o menu de presets de stack (+ "nenhum"). Devolve o preset
// escolhido e ok=false se o usuário optou por não usar nenhum.
func askPreset() (initx.Preset, bool) {
	const none = "none (infer / edit by hand)"
	opts := []string{none}
	for _, p := range initx.Presets {
		opts = append(opts, p.Title)
	}
	var choice string
	if err := huh.NewSelect[string]().
		Title("Use a project structure preset (established stack)?").
		Options(huh.NewOptions(opts...)...).
		Value(&choice).Run(); err != nil {
		erroDePrompt = true
	}
	if choice == none || choice == "" {
		return initx.Preset{}, false
	}
	for _, p := range initx.Presets {
		if p.Title == choice {
			return p, true
		}
	}
	return initx.Preset{}, false
}

// detectModules lista os diretórios de módulo existentes sob o ModuleGlob de um
// preset modular (ex.: os subdirs de src/features/). Vazio para presets não-modulares
// ou projeto novo — nesse caso os prefixos são deduzidos depois, quando os módulos
// nascerem. Usa doublestar para casar o glob de diretório.
func detectModules(root string, p initx.Preset) []string {
	if !p.Modular || p.ModuleGlob == "" {
		return nil
	}
	glob := strings.TrimRight(p.ModuleGlob, "/")
	matches, err := doublestar.Glob(os.DirFS(root), glob)
	if err != nil {
		return nil
	}
	var mods []string
	for _, m := range matches {
		if fi, err := os.Stat(filepath.Join(root, m)); err == nil && fi.IsDir() {
			mods = append(mods, m)
		}
	}
	return mods
}

func printFindings(p *initx.Proposal) {
	fmt.Println("\nFound:")
	if p.HasSpecMD {
		fmt.Println("  • specs (*.spec.md)")
	}
	if p.HasFeature {
		fmt.Println("  • features (*.feature)")
	}
	if p.HasTest {
		fmt.Println("  • tests (*.test.*)")
	}
	if p.GuideDir != "" {
		fmt.Printf("  • %d guides in %s/\n", len(p.GuideFiles), p.GuideDir)
	}
	if len(p.CodeDirs) > 0 {
		fmt.Printf("  • code in: %v (extensions: %v)\n", p.CodeDirs, p.CodeExts)
	}
	if p.Colocated {
		fmt.Println("  • co-location: derivatives next to the code")
	}
	fmt.Println()
}
