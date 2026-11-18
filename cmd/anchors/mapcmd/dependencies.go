// @anchors
//   code: DPNDC
//   ref: DPNDS

package mapcmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// dependsRelations say how an item is depended on by the node it was reached from, by edge type:
// the other end of the ties `impact` reads.
var dependsRelations = map[mapx.EdgeType]string{
	mapx.EdgeDependsOn: "is depended on", mapx.EdgeReadsFrom: "writes what it reads", mapx.EdgeCaptures: "is captured",
	mapx.EdgeRealizes: "is realized", mapx.EdgeGatedBy: "gates it", mapx.EdgeNeeds: "must finish first",
	mapx.EdgeGoverns: "governs it", mapx.EdgeSpecifies: "specifies it", mapx.EdgeCoveredBy: "is covered by it",
	mapx.EdgeTestedBy: "is tested by it", mapx.EdgeSeeds: "seeds it",
}

// dependenciesOut is the answer for the files asked about: what they depend on.
type dependenciesOut struct {
	Origins      []string     `json:"origins"`
	Dependencies []impactItem `json:"dependencies"`
}

// newDependenciesCmd is `impact` the other way: not who depends on the files, but everything
// they depend on — what they import or declare @dep of, read, capture, realize, and the specs,
// features and guides that rule them.
func newDependenciesCmd() *cobra.Command {
	var root, mapPath, propagate, format string
	var stopAt []string
	cmd := &cobra.Command{
		Use:   "dependencies <file>...",
		Short: "What these files depend on — impact the other way",
		Long: `Queries the map and shows, for one or more files, everything they depend on: what they
import or declare @dep of, the writers of what they read, what they capture or realize, and
the specs, features and guides that rule them. It is the other end of the ties 'impact'
follows.

  • Direct — what the files themselves depend on (the default).
  • --propagate indirect — one level further: what the direct ones depend on.
  • --propagate cascade  — the whole chain.

Each item past the first level carries the path that ties it to the file; an item that
depends on many others says how many it pulls (--stop-at <glob> lists it and goes no
further from it).

  --format text    grouped by level (default)
  --format files   the dependencies, one per line
  --format json    levels and paths

Query only — it opens no issue and changes nothing.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q, err := readTiesQuery(root, mapPath, propagate, format, stopAt, args)
			if err != nil {
				return err
			}
			ties := q.g.Dependencies(q.targets, q.maxLevel, q.stop)
			rep := dependenciesOut{Origins: q.targets, Dependencies: tieItems(q.g, q.targets, ties, dependsRelations, q.g.DependencyCount)}
			switch format {
			case "files":
				for _, it := range rep.Dependencies {
					fmt.Println(it.ID)
				}
			case "json":
				b, err := json.MarshalIndent(rep, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(b))
			default:
				fmt.Printf("dependencies of: %s\n\n", strings.Join(rep.Origins, ", "))
				frozen := printLevels(rep.Dependencies, "they depend on these", "(nothing — it depends on nothing)", propagate, "goes")
				if len(frozen) > 0 {
					fmt.Printf("\n❄ Frozen (%d) — kept as they are by decision:\n", len(frozen))
					for _, it := range frozen {
						fmt.Printf("   %s — %s (%s)\n", it.ID, it.Relation, it.Frozen)
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map (default: <root>/anchors.graph.yaml)")
	tiesFlags(cmd, &propagate, &format, &stopAt)
	return common.TakesFiles(cmd)
}
