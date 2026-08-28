package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/charmbracelet/huh"
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
		Short: "Configura o projeto (anchors.yaml) por perguntas e respostas",
		Long: `Escaneia o projeto de forma determinística (sem IA), propõe uma
Estrutura, e confirma/ajusta com você via perguntas — chegando a um anchors.yaml
correto. O grosso é inferido; as perguntas cobrem só as decisões humanas
(co-location, granularidade das camadas, e quais guides regem quais tags).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			// O modo não-interativo existe para o agente: sem ele, a guarda de TTY
			// aborta e o fluxo em que o usuário pede a uma IA para iniciar o projeto
			// (BOOTSTRAP.md §5) trava no comando central.
			if naoInterativo {
				return runInitNaoInterativo(cmd, absRoot, &f, aceitarDefaults)
			}
			return runInit(absRoot)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().BoolVar(&naoInterativo, "non-interactive", false, "sem TUI: sem respostas, emite as perguntas em JSON; com respostas em flags, aplica-as")
	cmd.Flags().BoolVar(&aceitarDefaults, "defaults", false, "com --non-interactive e sem respostas: aceita os defaults inferidos do disco")
	cmd.Flags().StringVar(&f.preset, "preset", "", "preset de stack (ver --questions)")
	cmd.Flags().BoolVar(&f.header, "header", true, "semear guides/HEADER_GUIDE.md")
	cmd.Flags().StringSliceVar(&f.artifacts, "artifacts", nil, "tipos de âncora do projeto (spec,feature,test,…)")
	cmd.Flags().BoolVar(&f.gates, "gates", true, "semear os gates padrão (informativos)")
	cmd.Flags().BoolVar(&f.colocation, "colocation", false, "derivados ao lado do código")
	cmd.Flags().StringSliceVar(&f.layers, "layers", nil, "diretórios de código a tratar como camadas")
	cmd.Flags().StringArrayVar(&f.governs, "governs", nil, "regra governs: GUIDE=tag1,tag2 (repetível)")
	cmd.Flags().StringVar(&f.workflow, "workflow", "", "onde a fila de trabalho mora: local|github")
	cmd.Flags().StringVar(&f.repo, "repo", "", "repositório owner/nome (obrigatório no modo github)")
	cmd.Flags().StringSliceVar(&f.labels, "labels", nil, "labels que marcam os cards do Anchors (obrigatório no modo github)")
	return cmd
}

// erroSemTTY é a mensagem de quando o `init` interativo não tem terminal. Nomeia a
// SAÍDA, e não só o problema: existe um modo não-interativo, e quem cai aqui (um agente,
// um pipe, o CI) tem como prosseguir — sem esta indicação, o comando parece um beco.
//
// `contexto` diz o que já aconteceu antes da falha, porque isso muda o que o leitor
// precisa saber: se o git foi tocado, se algo foi escrito.
func erroSemTTY(contexto string) error {
	return fmt.Errorf("`anchors init` e interativo e nao ha terminal disponivel "+
		"(sem TTY: pipe, CI ou agente sem shell interativo).\n"+
		"%s\n\n"+
		"Use o modo NAO-INTERATIVO, em duas chamadas da MESMA flag:\n"+
		"  1. anchors init --non-interactive\n"+
		"       sem respostas, devolve as decisoes em JSON (opcoes, o default\n"+
		"       inferido do disco, e o que cada resposta muda no projeto)\n"+
		"  2. anchors init --non-interactive --artifacts=spec,feature,test --colocation\n"+
		"       com respostas em flags, aplica-as e reporta o veredito de cada uma\n\n"+
		"Ou rode num terminal interativo.", contexto)
}

// runInit orquestra: infere (puro) → pergunta (TUI) → aplica decisões (puro) →
// salva. A lógica de decisão vive em initx (funções puras, testadas); esta função
// só faz o I/O e a coleta de respostas.
func runInit(root string) error {
	outPath := filepath.Join(root, config.DefaultFile)

	if _, err := os.Stat(outPath); err == nil {
		var overwrite bool
		if err := huh.NewConfirm().
			Title(config.DefaultFile + " já existe. Sobrescrever?").
			Value(&overwrite).Run(); err != nil {
			erroDePrompt = true
		}
		if !overwrite {
			fmt.Println("abortado — nada foi alterado.")
			return nil
		}
	}

	// 0) GIT — antes de escanear. É o substrato: sem versionamento o carimbo de
	// alteração, a cobertura de diff e o pre-commit ficam desligados, e nenhum deles
	// falha ruidosamente. Vem primeiro para que tudo que o init escrever daqui em
	// diante já nasça sob versionamento.
	if !etapaGit(root) {
		return erroSemTTY("Nada foi escrito, e o git nao foi tocado.")
	}

	fmt.Println("Escaneando o projeto…")
	p, err := initx.Infer(root)
	if err != nil {
		return fmt.Errorf("inferência: %w", err)
	}
	printFindings(p)

	// 0.4) A FASE ANTERIOR — antes de perguntar qualquer coisa, reconhecer se a fase
	// DESCOBRIR ainda não aconteceu. Sem PROJECT.md e sem código, as perguntas abaixo
	// saem sem resposta boa; o que muda é só QUEM recebe a instrução (pessoa ou IA).
	if !etapaDescobrir(root, p) {
		return erroSemTTY("Nada foi escrito.")
	}

	cfg := p.Config
	empty := len(p.CodeDirs) == 0 && !p.HasSpecMD && !p.HasFeature && !p.HasTest
	if empty {
		fmt.Println("Projeto novo/vazio — vou perguntar a estrutura que você pretende usar.")
	}

	// 0.5) PRESET DE STACK — oferece uma estrutura consagrada. Opcional: "nenhum"
	// mantém a inferência/edição manual. Se escolhido, preenche as layers de código
	// e, para presets modulares, deduz os prefixos de módulo (liga na identidade).
	chosenPreset, presetOK := askPreset()
	var presetModules []string
	if presetOK {
		presetModules = detectModules(root, chosenPreset)
		prefixes := initx.ApplyPreset(cfg, chosenPreset, presetModules)
		fmt.Printf("Preset '%s' aplicado (%s, %d camadas).\n", chosenPreset.Title, chosenPreset.Pattern, len(chosenPreset.Layers))
		if chosenPreset.Modular && len(prefixes) > 0 {
			fmt.Printf("Módulos detectados e seus prefixos de identidade (Camada 1):\n")
			for m, pfx := range prefixes {
				fmt.Printf("  %-20s → %s\n", m, pfx)
			}
		} else if chosenPreset.Modular {
			fmt.Printf("(preset modular: os prefixos de módulo serão deduzidos quando houver módulos em %s)\n", chosenPreset.ModuleGlob)
		}
		if chosenPreset.CoverageHint != "" {
			fmt.Printf("\nSinais de teste (para `anchors ingest` medir cobertura de verdade):\n  %s\n", chosenPreset.CoverageHint)
		}
	}

	// 0.6) GUIA DE CABEÇALHO — semeia um HEADER_GUIDE.md no projeto (a régua do bloco
	// @anchors, instanciada com o dialeto de comentário da stack e as features reais).
	// É o padrão MANDATÓRIO de cabeçalho; o init o materializa para o projeto seguir.
	if askConfirmDefault("Semear guides/HEADER_GUIDE.md (o padrão de cabeçalho @anchors dos arquivos)?", true) {
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
			fmt.Printf("Guia de cabeçalho semeado: %s\n", filepath.Join(guideDir, "HEADER_GUIDE.md"))
		}
	}

	// 1) ARTEFATOS — SEMPRE perguntado. Pré-marca os detectados; num projeto vazio,
	// o usuário marca o que PRETENDE usar. As layers de artefato são (re)construídas
	// a partir da escolha, não da detecção.
	chosenArtifacts := askMultiSelectPre(
		"Quais tipos de âncora seu projeto usa (ou vai usar)?",
		initx.ArtifactNames(), p.DetectedArtifacts())
	initx.ApplyArtifactChoice(cfg, chosenArtifacts, map[string]string{
		"guide": p.GuideDir,
		"plan":  p.PlanDir,
	})

	// 1.5) GATES PADRÃO — o ciclo nasce com os gates dos artefatos escolhidos
	// (spec-completa, tests-green, scenario-coverage, …), todos INFORMATIVOS. É o que
	// amarra os sinais de teste ao ciclo de vida sem o usuário escrever tudo à mão.
	if gates := initx.DefaultGates(chosenArtifacts); len(gates) > 0 {
		names := make([]string, len(gates))
		for i, g := range gates {
			names[i] = g.Name
		}
		if askConfirmDefault(
			fmt.Sprintf("Semear %d gate(s) padrão (informativos) para os artefatos escolhidos? [%s]",
				len(gates), strings.Join(names, ", ")), true) {
			cfg.Gates = gates
			if chosenArtifacts["test"] {
				fmt.Println("  (os gates de teste leem os sinais de `anchors ingest` — rode a suíte com coverage e ingira)")
			}
		}
	}

	// 2) CO-LOCATION — SEMPRE perguntado (texto adaptado a detectado vs. a decidir).
	title := "Como os derivados (spec/feature/teste) se organizam?"
	if p.Colocated {
		title = "Detectei os derivados ao lado do código (co-location). Usar essa convenção?"
	}
	useColocation := askConfirmDefault(title, p.Colocated)
	initx.ApplyColocation(cfg, useColocation, chosenArtifacts)

	// 3) CAMADAS DE CÓDIGO — sempre perguntado se há dirs detectados; se vazio, avisa.
	if names := initx.CodeLayerNames(cfg); len(names) > 0 {
		keep := askMultiSelect("Quais diretórios de código tratar como camadas?", names)
		initx.PruneCodeLayers(cfg, keep)
	} else if empty {
		fmt.Println("  (sem código ainda — declare as camadas de código no anchors.yaml quando existirem)")
	}

	// 3.5) MODO DE TRABALHO — onde a fila mora. É decisão HUMANA e EXCLUDENTE
	// (WORKFLOW.md §2): não se infere do remote, e não há fallback de um modo para o
	// outro. Perguntada aqui, e não no começo, porque só faz sentido depois de o projeto
	// ter forma — mas ANTES de salvar, porque muda o arquivo.
	if askConfirmDefault("A fila de trabalho vai morar nas issues do GitHub (em vez de local)?", false) {
		repo := askTexto("Qual repositório? (owner/nome — nunca inferido do remote)")
		if repo != "" {
			cfg.Workflow = &config.Workflow{
				Mode:   config.ModeGitHub,
				Repo:   repo,
				Labels: []string{"anchors"},
			}
			fmt.Println("  modo `github`: a label `anchors` marca os cards do fluxo no board compartilhado.")
			fmt.Println("  rode `anchors doctor --fix` depois, para semear os pipelines.")
		} else {
			fmt.Println("  sem repositório declarado — seguindo no modo local.")
		}
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
		return erroSemTTY("Nada foi escrito: um `anchors.yaml` gerado sem as respostas " +
			"sairia com 0 camadas e 0 gates, carregaria sem erro e nao governaria nada.")
	}
	if err := config.Save(cfg, outPath); err != nil {
		return fmt.Errorf("salvar: %w", err)
	}
	if p.GuideDir != "" {
		destSpec := filepath.Join(root, p.GuideDir, "SPEC_GUIDE.md")
		if os.MkdirAll(filepath.Dir(destSpec), 0o755) == nil &&
			os.WriteFile(destSpec, []byte(renderSpecGuide(cfg, "")), 0o644) == nil {
			fmt.Printf("Guia de spec semeado: %s\n", filepath.Join(p.GuideDir, "SPEC_GUIDE.md"))
		}
	}
	fmt.Printf("\n✓ %s escrito (%d camadas, %d regras de governs)\n", config.DefaultFile, len(cfg.Layers), len(cfg.Governs))
	fmt.Println("  revise e rode: anchors map build")
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

// askTexto coleta uma resposta livre. Usado pelo `repo` do modo github, que não tem
// conjunto de opções para escolher.
func askTexto(title string) string {
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
			Title(fmt.Sprintf("O guide %s rege qual tag?", filepath.Base(guide))).
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
	const none = "nenhum (inferir / editar à mão)"
	opts := []string{none}
	for _, p := range initx.Presets {
		opts = append(opts, p.Title)
	}
	var choice string
	if err := huh.NewSelect[string]().
		Title("Usar um preset de estrutura de projeto (stack consagrada)?").
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
	fmt.Println("\nEncontrei:")
	if p.HasSpecMD {
		fmt.Println("  • specs (*.spec.md)")
	}
	if p.HasFeature {
		fmt.Println("  • features (*.feature)")
	}
	if p.HasTest {
		fmt.Println("  • testes (*.test.*)")
	}
	if p.GuideDir != "" {
		fmt.Printf("  • %d guides em %s/\n", len(p.GuideFiles), p.GuideDir)
	}
	if len(p.CodeDirs) > 0 {
		fmt.Printf("  • código em: %v (extensões: %v)\n", p.CodeDirs, p.CodeExts)
	}
	if p.Colocated {
		fmt.Println("  • co-location: derivados ao lado do código")
	}
	fmt.Println()
}
