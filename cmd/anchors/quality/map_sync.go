package quality

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	if err := mapx.WithLock(mapPath, func() error {
		if before, err := mapx.Load(mapPath); err == nil {
			for _, d := range bumped {
				before.RebaseRev(d.File, scan.ShortHash([]byte(d.Old)), scan.ShortHash([]byte(d.Content)))
			}
			mapx.PreserveStamps(fresh, before)
			fresh.Flow = before.Flow
		}
		return mapx.Save(fresh, mapPath)
	}); err != nil {
		return "", err
	}
	if err := exec.Command("git", "-C", absRoot, "add", "--", mapx.DefaultPath).Run(); err != nil {
		return "", fmt.Errorf("stage the map: %w", err)
	}
	return fmt.Sprintf("· map: written as the commit will have it and staged (%d staged file(s) dated today, %d re-dated kept their proofs)", staged, len(bumped)), nil
}
