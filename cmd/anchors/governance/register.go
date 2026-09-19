package governance

import "github.com/spf13/cobra"

// Register adds governance domain commands to root.
func Register(root *cobra.Command) {
	root.AddCommand(newGuideCmd())
	root.AddCommand(newAuditCmd())
	root.AddCommand(newGovernsCmd())
	root.AddCommand(newComplianceCmd())
}
