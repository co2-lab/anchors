// @anchors
//   code: RGCMB
//   ref: MPRGM

package mapcmd

import "github.com/spf13/cobra"

// Register adds map domain commands to root.
func Register(root *cobra.Command) {
	root.AddCommand(newMapCmd())
	root.AddCommand(newImpactCmd())
	root.AddCommand(newIngestCmd())
	root.AddCommand(newJudgeCmd())
	root.AddCommand(newReviewCmd())
	root.AddCommand(newRecodeCmd())
	root.AddCommand(newRenumberCmd())
	root.AddCommand(newFlowCmd())
	root.AddCommand(newFailuresCmd())
}
