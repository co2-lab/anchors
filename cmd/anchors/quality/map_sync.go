// @anchors
//   code: MSCMA
//   ref: MPSYN

package quality

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/mapcmd"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// syncMapForCommit writes the map as it will be once the commit exists, and stages it.
//
// The map committed with a change was the map of the moment before: the hook had just dated
// the staged files (`updated_at`), which changed their revision, and a node's date comes from
// the last commit that touched the file — the one not made yet. So the CI's `map build` of
// the commit differed from the committed map on nearly every commit, and the next rebuild
// dropped the proofs of the dated files over a date.
//
// Here: the files the hook dated keep what was measured at their revision (`RebaseRev` —
// only the date changed), the staged files take today as their date, and the map is built
// from the INDEX (`WalkStaged`), not the tree, so an unstaged edit or an untracked file does
// not enter a map the commit will not match. Stamps, judgments, signals and the flow come
// from the map on disk, re-read under its lock. A map git does not track is left alone —
// nothing compares it.
func syncMapForCommit(absRoot string, cfg *config.Config, bumped []touchDecision) (string, error) {
	mapPath := filepath.Join(absRoot, mapx.DefaultPath)
	if _, err := os.Stat(mapPath); err != nil || cfg == nil {
		return "", nil
	}
	if err := exec.Command("git", "-C", absRoot, "ls-files", "--error-unmatch", mapx.DefaultPath).Run(); err != nil {
		return "", nil
	}
	files, err := scan.WalkStaged(absRoot, cfg)
	if err != nil {
		return "", err
	}
	out, err := exec.Command("git", "-C", absRoot, "diff", "--cached", "--name-only", "--relative", "-z").Output()
	if err != nil {
		return "", err
	}
	dates := gitmeta.AllCommitDates(absRoot)
	staged := 0
	for _, p := range strings.Split(string(out), "\x00") {
		if p != "" {
			dates[filepath.ToSlash(p)] = gitmeta.Today()
			staged++
		}
	}
	fresh := mapx.Build(files, cfg, dates)
	ahead, err := scan.GovernedTreeChanges(absRoot, cfg)
	if err != nil {
		return "", err
	}
	head := headMap(absRoot)
	apart := len(ahead) > 0
	if err := mapx.WithLock(mapPath, func() error {
		before, err := mapx.Load(mapPath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("read the map: %w", err)
		}
		// A file the commit leaves as HEAD has it, and the tree has edited since, takes HEAD's
		// measurement of that content: carried back from the edit, its proofs could have lost
		// the scenarios of a rule the edit changed.
		asHead := map[string]*mapx.Node{}
		if err == nil && head != nil {
			for _, n := range fresh.Nodes {
				h, b := head.Node(n.ID), before.Node(n.ID)
				if h != nil && h.Rev == n.Rev && h.Signal != nil && b != nil && b.Rev != n.Rev {
					asHead[n.ID] = h
				}
			}
		}
		if err == nil {
			for _, d := range bumped {
				before.RebaseRev(d.File, scan.ShortHash([]byte(d.Old)), scan.ShortHash([]byte(d.Content)))
			}
			mapx.FillOldEvidence(fresh, before, mapcmd.AtHeadReader(absRoot), cfg)
			mapx.PreserveStamps(fresh, before)
			fresh.Flow = before.Flow
		}
		for i := range fresh.Nodes {
			if h, ok := asHead[fresh.Nodes[i].ID]; ok {
				fresh.Nodes[i].Signal, fresh.Nodes[i].EvidenceKept = h.Signal, h.EvidenceKept
			}
		}
		// A file the commit leaves as HEAD has it keeps what was measured of that content:
		// the map on disk may hold another revision of it — an unstaged edit, measured or
		// not —, and the committed map must not lose the proofs HEAD already had.
		mapx.CarryUnchangedEvidence(fresh, head)
		mapx.FillSignals(fresh, head)
		if !apart {
			return mapx.Save(fresh, mapPath)
		}
		// The tree is ahead of the commit (another session's unstaged edits): the commit
		// gets its own map, straight into the index, and the map on disk stays the tree's,
		// with what was measured of those edits.
		if err := stageMapBlob(absRoot, fresh); err != nil {
			return err
		}
		if before != nil {
			return mapx.Save(before, mapPath)
		}
		return nil
	}); err != nil {
		return "", err
	}
	if apart {
		return fmt.Sprintf("· map: the commit's map staged apart — %d governed file(s) differ in the tree, and the map on disk keeps them (%d staged file(s) dated today, %d re-dated kept their proofs)", len(ahead), staged, len(bumped)), nil
	}
	if err := exec.Command("git", "-C", absRoot, "add", "--", mapx.DefaultPath).Run(); err != nil {
		return "", fmt.Errorf("stage the map: %w", err)
	}
	return fmt.Sprintf("· map: written as the commit will have it and staged (%d staged file(s) dated today, %d re-dated kept their proofs)", staged, len(bumped)), nil
}

// stageMapBlob writes the map into the index without touching the file on disk.
func stageMapBlob(absRoot string, g *mapx.Graph) error {
	dir, err := os.MkdirTemp("", "anchors-map-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	tmp := filepath.Join(dir, "map.yaml")
	if err := mapx.Save(g, tmp); err != nil {
		return err
	}
	b, err := os.ReadFile(tmp)
	if err != nil {
		return err
	}
	return stageBlob(absRoot, mapx.DefaultPath, b)
}
