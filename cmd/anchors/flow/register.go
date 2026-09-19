package flow

import "github.com/spf13/cobra"

// Register adds flow domain commands to root.
func Register(root *cobra.Command) {
	root.AddCommand(newWorkCmd())
	root.AddCommand(newWatchCmd())
	root.AddCommand(newQueueCmd())
	root.AddCommand(newNextCmd())
	root.AddCommand(newDoneCmd())
	root.AddCommand(newDropCmd())
	root.AddCommand(newEscalateCmd())
	root.AddCommand(newDecidedCmd())
	root.AddCommand(newTaskStatusCmd())
	root.AddCommand(newUnblockCmd())
	root.AddCommand(newBackfillLabelsCmd())
	root.AddCommand(newDiscardCmd())
	root.AddCommand(newReclaimCmd())
	root.AddCommand(newDeliverCmd())
	root.AddCommand(newProgressMergeCmd())
	root.AddCommand(newPRBodyCmd())
}
