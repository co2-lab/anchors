package ops

import "github.com/spf13/cobra"

// Register adds ops domain commands to root.
func Register(root *cobra.Command) {
	root.AddCommand(newInitCmd())
	root.AddCommand(newNewCmd())
	root.AddCommand(newMigrateCmd())
	root.AddCommand(newFreezeCmd())
	root.AddCommand(newThawCmd())
	root.AddCommand(newSettingsCmd())
	root.AddCommand(newInstallHooksCmd())
	root.AddCommand(newCommitMsgCmd())
	root.AddCommand(newBoardCmd())
	root.AddCommand(newDocsCmd())
	root.AddCommand(newSuggestCmd())
	root.AddCommand(newSynthesizeCmd())
	root.AddCommand(newCodeCmd())
	root.AddCommand(newGeneratedPathsCmd())
}
