// @anchors
//   code: LDCMV
//   ref: LDCMS

package ops

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/co2-lab/anchors/internal/load"
	"github.com/spf13/cobra"
)

// A run that waits looks like a hang unless something says what it waits for. `anchors load`
// shows the machine's capacity, the runs holding it and the runs waiting, from every worktree
// and project on the machine.

func newLoadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "load",
		Short: "Show what the machine lends its runs, the runs holding it and the runs waiting (`load:`)",
		Long: `The heavy runs of a project that declares ` + "`load:`" + ` — its test and mutation suites, its
project and batch tools — wait in a queue shared by every worktree and project on the machine:
each states what it costs, in CPUs and memory, and starts when that fits in what the machine
has free. A run that waited past its limit stops the admission of new runs until it fits.
The commit hook jumps the queue, and CI skips it unless ` + "`load.ci`" + ` asks for it.

  capacity   the machine's CPUs and memory, minus a reserve left to whoever uses it
             ($ANCHORS_LOAD_CAPACITY=cpu=8,memory=16GB sets it)
  running    each run holding the machine: project, task, level, cost, since when
  waiting    each run waiting, how long, and whether it waited past its limit`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := load.Read()
			if err != nil {
				return err
			}
			printLoad(cmd.OutOrStdout(), load.Sorted(st), load.MachineCapacity(), time.Now())
			return nil
		},
	}
}

func printLoad(w io.Writer, st load.State, capacity load.Capacity, now time.Time) {
	used := load.Used(st)
	fmt.Fprintf(w, "capacity   %.0f CPU, %s\n", capacity.CPU, load.Size(capacity.Memory))
	fmt.Fprintf(w, "in use     %.0f CPU, %s\n", used.CPU, load.Size(used.Memory))
	if len(st.Leases) == 0 {
		fmt.Fprintln(w, "running    —")
	}
	for i, l := range st.Leases {
		label := "running "
		if i > 0 {
			label = "        "
		}
		fmt.Fprintf(w, "%s   %s %s (%s: %.0f CPU, %s), for %s, pid %d\n", label,
			filepath.Base(l.Project), l.Task, l.Level, l.Cost.CPU, load.Size(l.Cost.Memory), now.Sub(l.Since).Round(time.Second), l.PID)
	}
	if len(st.Waiters) == 0 {
		fmt.Fprintln(w, "waiting    —")
	}
	for i, l := range st.Waiters {
		label := "waiting "
		if i > 0 {
			label = "        "
		}
		waited := now.Sub(l.Since)
		past := ""
		if waited.Seconds() > l.Limit {
			past = " — past its limit: no new run starts until it fits"
		}
		fmt.Fprintf(w, "%s   %s %s (%s: %.0f CPU, %s), for %s, pid %d%s\n", label,
			filepath.Base(l.Project), l.Task, l.Level, l.Cost.CPU, load.Size(l.Cost.Memory), waited.Round(time.Second), l.PID, past)
	}
}
