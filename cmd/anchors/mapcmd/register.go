package mapcmd

import "github.com/spf13/cobra"

// Register adds map domain commands to root.
func Register(root *cobra.Command) {
	root.AddCommand(newMapCmd())
	root.AddCommand(newImpactCmd())
	root.AddCommand(newIngestCmd())
	root.AddCommand(newJudgeCmd())
	root.AddCommand(newRecodeCmd())
	root.AddCommand(newFlowCmd())
}
