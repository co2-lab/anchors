package quality

import "github.com/spf13/cobra"

// Register adds quality domain commands to root.
func Register(root *cobra.Command) {
	root.AddCommand(newCheckCmd())
	root.AddCommand(newVerifyCmd())
	root.AddCommand(newDoctorCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newStaleCmd())
	root.AddCommand(newCoverageCmd())
	root.AddCommand(newReportCmd())
	root.AddCommand(newTestCmd())
	root.AddCommand(newMutationCmd())
	root.AddCommand(newStampCmd())
}
