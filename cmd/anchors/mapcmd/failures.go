package mapcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// newFailuresCmd — the third layer: what the observed failures have not explained yet.
//
// Anchors does NOT analyse the log. Root-cause analysis is the work of whoever understands
// the domain (or of an AI, against the context), and inventing a "root cause detector"
// here would promise what no heuristic delivers.
//
// What it does is what is missing today, and it is what closes the circuit: it gives the
// TARGET (which failures happened), the MATERIAL (how many times, in what window, in which
// spec), and the PLACE to record the conclusion — the spec itself, where whoever touches
// the code tomorrow will read it.
//
// It is the same design as `judge`: the CLI does not compute the verdict, it asks for it
// and keeps it.
func newFailuresCmd() *cobra.Command {
	var root, mapPath string
	var all bool
	cmd := &cobra.Command{
		Use:   "failures",
		Short: "The observed failures that the spec has not yet explained",
		Long: `Lists the failures that HAPPENED (ingested with 'anchors ingest --logs') and
that carry no conclusion in the spec yet.

A spec declares a failure as POSSIBLE without knowing how it will happen — an
` + "`if x == null`" + ` foresees the failure without knowing where the null comes from.
Once the application runs, the log has the context that answers it.

Three conclusions, and all three are progress:

  cause found    write it beside the rule — the handling can stop being generic
  resilient      ` + "`@resilient: <reason>`" + ` — it leaves the radar without leaving the record
  still unknown  ` + "`@observing: <what was ruled out>`" + ` — different from nobody having looked

The third is what no observability tool records, and it is knowledge: the next person to
look starts from what was already ruled out, instead of from zero.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return err
			}
			return reportFailures(absRoot, g, all)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	cmd.Flags().BoolVar(&all, "all", false,
		"also list the failures that already carry a conclusion (resilient or under observation)")
	return cmd
}

type observedFailure struct {
	mapx.FailureSignal
	Spec       string
	Conclusion gate.FailureConclusion
	Stale      bool
}

func reportFailures(root string, g *mapx.Graph, all bool) error {
	var open, concluded []observedFailure
	for _, n := range g.Nodes {
		if len(n.Failures) == 0 {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, n.ID))
		if err != nil {
			continue
		}
		conclusions := gate.FailureConclusions(string(b))
		for _, f := range n.Failures {
			o := observedFailure{
				FailureSignal: f,
				Spec:          n.ID,
				Conclusion:    conclusions[f.Rule],
				// The occurrence was measured against a VERSION of the spec. If the spec
				// changed, the rule may be another one — and concluding about the old
				// failure would assert about what was not observed.
				Stale: f.Rev != "" && n.Rev != "" && f.Rev != n.Rev,
			}
			if o.Conclusion.Resilient != "" || o.Conclusion.Observing != "" {
				concluded = append(concluded, o)
				continue
			}
			open = append(open, o)
		}
	}
	// The most frequent first: it is the one that costs the most, and the one that offers
	// the most material to whoever investigates.
	sort.Slice(open, func(i, j int) bool { return open[i].Count > open[j].Count })
	sort.Slice(concluded, func(i, j int) bool { return concluded[i].Rule < concluded[j].Rule })

	if len(open) == 0 && (!all || len(concluded) == 0) {
		fmt.Println(i18n.T("failure.review.none"))
		return nil
	}
	if len(open) > 0 {
		fmt.Printf(i18n.T("failure.review.header")+"\n", len(open))
		for _, o := range open {
			printFailure(o)
		}
		fmt.Println(i18n.T("failure.review.howto"))
	}
	if all && len(concluded) > 0 {
		fmt.Println()
		for _, o := range concluded {
			printFailure(o)
			if o.Conclusion.Resilient != "" {
				fmt.Printf(i18n.T("failure.review.resilient")+"\n", o.Conclusion.Resilient)
			}
			if o.Conclusion.Observing != "" {
				fmt.Printf(i18n.T("failure.review.observing")+"\n", o.Conclusion.Observing)
			}
		}
	}
	return nil
}

func printFailure(o observedFailure) {
	fmt.Printf(i18n.T("failure.review.item")+"\n", o.Rule, o.Count, o.First, o.Last)
	fmt.Printf(i18n.T("failure.review.where")+"\n", o.Spec)
	if o.Stale {
		fmt.Println(i18n.T("failure.review.stale"))
	}
}
