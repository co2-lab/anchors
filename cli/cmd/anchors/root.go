package main

import (
	"fmt"
	"path/filepath"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "anchors",
		Short: "Anchors — framework de continuidade para desenvolvimento assistido por IA",
		Long: `Anchors mantém um projeto coerente ao longo do tempo através de âncoras:
documentos que guiam o desenvolvimento e confrontam o que foi feito.

Este CLI exercita o ciclo de vida do Anchors: constrói o mapa de dependências,
propaga alterações, roda os gates de qualidade e reporta a saúde do projeto.`,
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
			return refuseIfFrozen(cmd)
		},
	}
	root.AddCommand(newGuideCmd())
	root.AddCommand(newWorkCmd())
	root.AddCommand(newInitCmd())
	root.AddCommand(newInstallHooksCmd())
	// O botão de pânico: congela e descongela o projeto inteiro.
	root.AddCommand(newFreezeCmd())
	root.AddCommand(newThawCmd())
	root.AddCommand(newNewCmd())
	root.AddCommand(newRecodeCmd())
	root.AddCommand(newMapCmd())
	root.AddCommand(newImpactCmd())
	root.AddCommand(newCheckCmd())
	root.AddCommand(newVerifyCmd())
	root.AddCommand(newAuditCmd())
	root.AddCommand(newDoctorCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newWatchCmd())
	root.AddCommand(newQueueCmd())
	root.AddCommand(newNextCmd())
	root.AddCommand(newDoneCmd())
	root.AddCommand(newDropCmd())
	root.AddCommand(newEscalateCmd())
	root.AddCommand(newCommitMsgCmd())
	root.AddCommand(newPRBodyCmd())
	root.AddCommand(newReclaimCmd())
	root.AddCommand(newStaleCmd())
	root.AddCommand(newCodeCmd())
	root.AddCommand(newJudgeCmd())
	root.AddCommand(newSuggestCmd())
	root.AddCommand(newGovernsCmd())
	root.AddCommand(newIngestCmd())
	root.AddCommand(newTestCmd())
	root.AddCommand(newMutationCmd())
	root.AddCommand(newDeliverCmd())
	root.AddCommand(newComplianceCmd())
	root.AddCommand(newCoverageCmd())
	root.AddCommand(newReportCmd())
	return root
}

// comandosQueRodamCongelado são os que continuam valendo com o projeto parado.
//
// A régua: um comando pode rodar congelado se ele NÃO PRODUZ estado do projeto. Ler é
// permitido — quem está investigando o problema precisa do `status`, do `doctor`, dos
// guias. Escrever no mapa, julgar, ingerir sinal, abrir card: não, porque é isso que o
// congelamento existe para parar.
//
// `thaw` está aqui pelo motivo mais óbvio e mais fácil de esquecer: sem ele, o
// congelamento seria irreversível pelo próprio Anchors.
var comandosQueRodamCongelado = map[string]bool{
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
		if comandosQueRodamCongelado[c.Name()] {
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
	return fmt.Errorf("🛑 o projeto está CONGELADO — `%s` não roda.\n\n"+
		"   Motivo: %s\n\n"+
		"   O trabalho que você já fez continua no seu branch local.\n"+
		"   Para liberar (quem congelou): `anchors thaw`\n"+
		"   Para investigar: `anchors status`, `anchors doctor` e `anchors guide` continuam valendo",
		cmd.Name(), cfg.FreezeReasonText())
}
