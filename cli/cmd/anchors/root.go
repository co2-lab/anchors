package main

import (
	"fmt"

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
	}
	root.AddCommand(newGuideCmd())
	root.AddCommand(newWorkCmd())
	root.AddCommand(newInitCmd())
	root.AddCommand(newInstallHooksCmd())
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
