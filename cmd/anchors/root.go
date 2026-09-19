package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/cmd/anchors/flow"
	"github.com/co2-lab/anchors/cmd/anchors/governance"
	"github.com/co2-lab/anchors/cmd/anchors/mapcmd"
	"github.com/co2-lab/anchors/cmd/anchors/ops"
	"github.com/co2-lab/anchors/cmd/anchors/quality"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "anchors",
		Short: "Anchors — a continuity framework for AI-assisted development",
		Long: `Anchors keeps a project coherent over time through anchors:
documents that guide development and confront what was done.

This CLI exercises the Anchors lifecycle: it builds the dependency map,
propagates changes, runs the quality gates and reports the project's health.`,
		SilenceUsage: true,
		// Sem isto o erro sai DUAS vezes: o cobra imprime "Error: x" ao voltar do
		// Execute, e o main imprime "erro: x" logo em seguida. Quem trata a saída é o
		// main — é ele que também decide o código de saída (ver ExitNaoRegido) —,
		// então o cobra silencia.
		SilenceErrors: true,
		Version:       fmt.Sprintf("%s (commit %s, built %s)", version, commit, date),
		// O CONGELAMENTO é conferido AQUI, e não em cada comando.
		//
		// Espalhar a checagem por comando é garantir que o próximo nasça sem ela — e um
		// comando que escapa do freio o torna decorativo. O `PersistentPreRunE` roda antes
		// de todo subcomando, inclusive dos que ainda não existem.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// O IDIOMA vem primeiro: a própria recusa por congelamento tem de sair no
			// idioma do projeto.
			//
			// Definir aqui e não só no `config.Load` é o que faz TODA mensagem sair
			// traduzida, inclusive as que um comando imprime ANTES de carregar a
			// configuração. Medido: o `anchors status` confere o git antes de carregar
			// o `anchors.yaml`, e um projeto com `lang: es` recebia essas linhas em
			// inglês — a chave existia, o idioma é que ainda não valia.
			applyProjectLang(cmd)
			// O AVISO DE TELEMETRIA, e aqui pelo mesmo motivo do congelamento acima: é o
			// único ponto por onde TODO comando passa.
			//
			// Os candidatos descartados e por quê:
			//
			//	init    → só alcança projeto NOVO. Quem instala o binário num projeto que
			//	          outra pessoa configurou nunca o roda — e é o caso comum num time.
			//	next    → só alcança quem PEDE CARD. Um revisor que roda `check` não vê.
			//	doctor  → é opcional, e muita gente nunca o roda.
			//
			// E há a propriedade que decide: o `PersistentPreRunE` roda ANTES do comando,
			// então o aviso aparece antes de o primeiro evento sair. Com `next` ou
			// `doctor`, o dado já teria ido.
			common.NoticeTelemetry(cmd)
			return refuseIfFrozen(cmd)
		},
	}

	governance.Register(root)
	mapcmd.Register(root)
	quality.Register(root)
	flow.Register(root)
	ops.Register(root)

	return root
}

// commandsAllowedWhileFrozen são os que continuam valendo com o projeto parado.
//
// A régua: um comando pode rodar congelado se ele NÃO PRODUZ estado do projeto. Ler é
// permitido — quem está investigando o problema precisa do `status`, do `doctor`, dos
// guias. Escrever no mapa, julgar, ingerir sinal, abrir card: não, porque é isso que o
// congelamento existe para parar.
//
// `thaw` está aqui pelo motivo mais óbvio e mais fácil de esquecer: sem ele, o
// congelamento seria irreversível pelo próprio Anchors.
var commandsAllowedWhileFrozen = map[string]bool{
	"thaw": true, "freeze": true,
	"status": true, "doctor": true, "guide": true, "version": true,
	"help": true, "completion": true, "coverage": true, "impact": true,
}

// refuseIfFrozen barra o comando quando o projeto declara `enabled: false`.
//
// A config é lida do disco em vez de vir do comando: cada um a carrega do seu jeito (uns
// com `--root`, outros do cwd), e depender disso deixaria buracos. Se ela não carregar,
// o comando segue — a ausência de config é outro problema, e responder "congelado" ali
// mandaria quem investiga para o lado errado.
func refuseIfFrozen(cmd *cobra.Command) error {
	// A CADEIA inteira, não só o nome do comando invocado.
	//
	// `cmd.Name()` de `anchors guide work` devolve "work", não "guide" — e o guia é
	// justamente o que a pessoa precisa ler durante o congelamento. Medido: `guide work`
	// era recusado enquanto `guide` sozinho passava.
	//
	// Subir até a raiz cobre qualquer profundidade de subcomando, inclusive as que ainda
	// não existem.
	for c := cmd; c != nil; c = c.Parent() {
		if commandsAllowedWhileFrozen[c.Name()] {
			return nil
		}
	}
	root := "."
	if f := cmd.Flags().Lookup("root"); f != nil && f.Value.String() != "" {
		root = f.Value.String()
	}
	absRoot, err := config.AbsRoot(root)
	if err != nil {
		return nil
	}
	cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
	if err != nil || !cfg.Frozen() {
		return nil
	}
	cmd.SilenceUsage = true
	return fmt.Errorf("%s\n\n%s\n\n%s\n%s\n%s",
		i18n.T("freeze.blocked.cli", cmd.Name()),
		i18n.T("freeze.blocked.reason", cfg.FreezeReasonText()),
		i18n.T("freeze.blocked.work_safe"),
		i18n.T("freeze.blocked.how_to_thaw"),
		i18n.T("freeze.blocked.what_still_works"))
}

// applyProjectLang lê o `lang:` do anchors.yaml e o define, ANTES de qualquer saída.
//
// Falha em silêncio de propósito: um projeto sem configuração, ou com ela quebrada, tem
// outro problema — e recusar aqui impediria o `anchors init` de rodar justamente onde
// ainda não há o que ler. O idioma cai no padrão, e o comando segue.
//
// O `config.Load` também define o idioma, e isso não é redundância: ele valida o valor e
// recusa um `lang:` que o Anchors não suporta. Aqui a leitura é frouxa, porque o objetivo
// é só não imprimir em inglês antes de saber.
func applyProjectLang(cmd *cobra.Command) {
	root := "."
	if f := cmd.Flags().Lookup("root"); f != nil && f.Value.String() != "" {
		root = f.Value.String()
	}
	absRoot, err := config.AbsRoot(root)
	if err != nil {
		return
	}
	b, err := os.ReadFile(filepath.Join(absRoot, config.DefaultFile))
	if err != nil {
		return
	}
	if m := langNoYAML.FindSubmatch(b); m != nil {
		_ = i18n.Set(strings.TrimSpace(string(m[1])))
	}
}

// langNoYAML pega o `lang:` de topo do anchors.yaml.
//
// Por regex e não por parse: este ponto roda antes de tudo, e um YAML inválido não pode
// impedir o comando de rodar — quem reclama do YAML é o `config.Load`, com a linha e a
// chave. O `^` exige coluna zero, então um `lang:` aninhado noutro bloco não conta.
var langNoYAML = regexp.MustCompile(`(?m)^lang:[ \t]*["']?([A-Za-z-]+)["']?[ \t]*$`)
