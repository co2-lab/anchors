package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

// flagsInit guarda as respostas vindas da linha de comando. Ponteiros de propósito: é o
// que distingue "não respondi" (vale o default inferido do disco) de "respondi vazio"
// (`--artifacts=""`, nenhum artefato) — dois casos com resultados opostos.
type flagsInit struct {
	preset     string
	header     bool
	artifacts  []string
	gates      bool
	colocation bool
	layers     []string
	governs    []string
	workflow   string
	repo       string
	labels     []string
}

// respostasDeFlags converte as flags em Respostas, consultando quais foram REALMENTE
// passadas. `cmd.Flags().Changed` é o que permite não confundir o zero-value de um bool
// com uma escolha deliberada de `false`.
func respostasDeFlags(cmd *cobra.Command, f *flagsInit) (initx.Respostas, error) {
	var r initx.Respostas
	if cmd.Flags().Changed("preset") {
		r.Preset = &f.preset
	}
	if cmd.Flags().Changed("header") {
		r.Header = &f.header
	}
	if cmd.Flags().Changed("gates") {
		r.Gates = &f.gates
	}
	if cmd.Flags().Changed("colocation") {
		r.Colocation = &f.colocation
	}
	if cmd.Flags().Changed("artifacts") {
		r.Artifacts = &f.artifacts
	}
	if cmd.Flags().Changed("layers") {
		r.Layers = &f.layers
	}
	if cmd.Flags().Changed("workflow") {
		r.Workflow = &f.workflow
	}
	if cmd.Flags().Changed("repo") {
		r.Repo = &f.repo
	}
	if cmd.Flags().Changed("labels") {
		r.Labels = &f.labels
	}
	if len(f.governs) > 0 {
		r.Governs = map[string][]string{}
		for _, g := range f.governs {
			guide, tags, ok := strings.Cut(g, "=")
			if !ok {
				return r, fmt.Errorf("--governs espera GUIDE=tag1,tag2 (recebi %q)", g)
			}
			r.Governs[guide] = strings.Split(tags, ",")
		}
	}
	return r, nil
}

// runInitNaoInterativo é o `init` para quem não tem terminal: um agente operando o CLI.
//
// Sem ele, o fluxo em que o usuário pede a uma IA para iniciar o projeto (BOOTSTRAP.md
// §5) trava no comando central — a guarda de TTY aborta, e o agente fica sem como
// prosseguir.
//
// São DUAS chamadas, e a separação é o que dá ao agente o que perguntar ao usuário:
//
//  1. `--questions` devolve as perguntas em JSON, com opções, default inferido do disco,
//     e o que cada resposta MUDA no projeto;
//  2. as respostas voltam em flags, e a saída traz o veredito de CADA uma.
func runInitNaoInterativo(cmd *cobra.Command, root string, f *flagsInit, aceitarDefaults bool) error {
	p, err := initx.Infer(root)
	if err != nil {
		return fmt.Errorf("inferência: %w", err)
	}
	qs := initx.Perguntas(p, initx.PresetNames())

	r, err := respostasDeFlags(cmd, f)
	if err != nil {
		return err
	}

	// SEM respostas, o comando PERGUNTA — não decide. Escrever aceitando todos os
	// defaults produziria, num projeto novo, um anchors.yaml com zero camadas e zero
	// artefatos: carrega sem erro e não governa nada, que é o mesmo arquivo inútil que
	// a guarda de TTY existe para impedir.
	//
	// Aceitar os defaults continua possível, mas tem de ser DITO (`--defaults`):
	// assim "não respondi" e "aceito tudo" nunca são a mesma coisa.
	if !respondeuAlgo(r) && !aceitarDefaults {
		return emiteJSON(map[string]any{
			"projeto":           root,
			"precisa_descobrir": initx.PrecisaDescobrir(root, p),
			"escrito":           false,
			"perguntas":         qs,
			"como_responder": "passe as respostas em flags (ex.: --artifacts=spec,feature,test " +
				"--colocation) e chame de novo; ou --defaults para aceitar os defaults " +
				"inferidos do disco",
		})
	}
	status := initx.ValidaRespostas(qs, r)

	// Uma resposta inválida recusa o CONJUNTO. Escrever as válidas produziria um
	// anchors.yaml que ninguém decidiu por completo — e um arquivo assim carrega sem
	// erro, governa errado, e não acusa a causa.
	if !initx.TudoAceito(status) {
		_ = emiteJSON(map[string]any{
			"escrito":   false,
			"respostas": status,
			"erro":      "há respostas inválidas — nada foi escrito",
		})
		return fmt.Errorf("respostas inválidas; nada foi escrito")
	}

	if err := aplicaRespostas(root, p, status); err != nil {
		return err
	}
	return emiteJSON(map[string]any{
		"escrito":       true,
		"arquivo":       filepath.Join(root, config.DefaultFile),
		"respostas":     status,
		"proximo_passo": proximoPassoApos(root, p),
	})
}

// respondeuAlgo diz se veio ao menos uma resposta. É o que separa "quero as perguntas"
// de "aqui estão as respostas" — sem precisar de uma flag para cada intenção.
func respondeuAlgo(r initx.Respostas) bool {
	return r.Preset != nil || r.Header != nil || r.Artifacts != nil || r.Gates != nil ||
		r.Colocation != nil || r.Layers != nil || len(r.Governs) > 0 ||
		r.Workflow != nil || r.Repo != nil || r.Labels != nil
}

// aplicaRespostas monta o anchors.yaml a partir dos status já validados, na mesma ordem
// da TUI — cada decisão restringe a seguinte.
func aplicaRespostas(root string, p *initx.Proposal, status []initx.StatusResposta) error {
	cfg := p.Config
	valor := func(id string) any {
		for _, s := range status {
			if s.ID == id {
				return s.Valor
			}
		}
		return nil
	}

	if nome, _ := valor("preset").(string); nome != "" && nome != "nenhum" {
		for _, pr := range initx.Presets {
			if pr.Name == nome {
				preset := pr
				initx.ApplyPreset(cfg, preset, detectModules(root, preset))
				break
			}
		}
	}

	artefatos := map[string]bool{}
	for _, a := range comoLista(valor("artifacts")) {
		artefatos[a] = true
	}
	initx.ApplyArtifactChoice(cfg, artefatos, map[string]string{
		"guide": p.GuideDir,
		"plan":  p.PlanDir,
	})

	if sim, _ := valor("gates").(bool); sim {
		// Projeto novo (sem código nem artefato no disco) nasce com os gates
		// BLOQUEANTES — ver DefaultGates para o porquê da distinção.
		novo := len(p.CodeDirs) == 0 && !p.HasSpecMD && !p.HasFeature && !p.HasTest
		cfg.Gates = initx.DefaultGates(artefatos, novo)
	}
	colocado, _ := valor("colocation").(bool)
	initx.ApplyColocation(cfg, colocado, artefatos)

	if l := comoLista(valor("layers")); len(l) > 0 {
		keep := map[string]bool{}
		for _, n := range l {
			keep[n] = true
		}
		initx.PruneCodeLayers(cfg, keep)
	}

	// O modo é decisão HUMANA (onde a fila mora), e os modos são EXCLUDENTES: só se
	// declara o bloco quando a escolha é `github`. Escrever `mode: local` explícito seria
	// ruído — a ausência do bloco já significa local (WORKFLOW.md §2).
	if modo, _ := valor("workflow").(string); modo == "github" {
		repo, _ := valor("repo").(string)
		cfg.Workflow = &config.Workflow{
			Mode:   config.ModeGitHub,
			Repo:   repo,
			Labels: comoLista(valor("labels")),
		}
	}

	if err := config.Save(cfg, filepath.Join(root, config.DefaultFile)); err != nil {
		return fmt.Errorf("salvar: %w", err)
	}

	// O header guide é semeado se pedido: é o padrão mandatório do bloco @anchors, e
	// sem ele o gate de header não tem régua para confrontar.
	if sim, _ := valor("header").(bool); sim {
		guideDir := p.GuideDir
		if guideDir == "" {
			guideDir = "guides"
		}
		dest := filepath.Join(root, guideDir, "HEADER_GUIDE.md")
		body := initx.RenderHeaderGuide(initx.Preset{}, nil)
		if os.MkdirAll(filepath.Dir(dest), 0o755) == nil {
			_ = os.WriteFile(dest, []byte(body), 0o644)
		}
	}
	return nil
}

// proximoPassoApos diz ao agente o que fazer em seguida. Um comando que termina sem
// dizer isso obriga quem o chamou a adivinhar a ordem do ciclo.
func proximoPassoApos(root string, p *initx.Proposal) string {
	if initx.PrecisaDescobrir(root, p) {
		return "a fase DESCOBRIR não aconteceu: rode `anchors guide project` e conduza a " +
			"entrevista com o usuário antes de escrever qualquer código"
	}
	return "`anchors map build` — sem o mapa, nenhum arquivo existe para os gates"
}

func comoLista(v any) []string {
	switch t := v.(type) {
	case []string:
		return t
	case []any:
		var out []string
		for _, x := range t {
			if s, ok := x.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// emiteJSON escreve na saída padrão. Indentado: quem lê isto é um agente, mas um humano
// depurando o fluxo lê a mesma saída.
func emiteJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
