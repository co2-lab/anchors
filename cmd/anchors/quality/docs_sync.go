package quality

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/doct"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// syncDocsForCommit compiles the pages of `docs/` the commit leaves out of date and stages
// them with it, as the map is.
//
// A commit that changed a spec was barred by `docs-fresh` until someone ran `anchors docs
// build` by hand; a project could not make the gate blocking without making every spec
// change a two-step commit. Here the pages are compiled from what the commit records — the
// staged specs and templates, and the commit's map —, so another session's unstaged edit
// does not enter them. Only the pages out of date are compiled, and a page written by hand
// is never touched.
//
// When the tree holds governed changes the index does not, the pages go straight into the
// index and the files on disk stay as they are, as the map does: the tree's pages are the
// tree's to build. Otherwise they are written and staged. A page git ignores is left alone.
func syncDocsForCommit(absRoot string, cfg *config.Config) (string, error) {
	if cfg == nil || !cfg.DocsOnPreCommit() {
		return "", nil
	}
	if fi, err := os.Stat(filepath.Join(absRoot, doct.Dir)); err != nil || !fi.IsDir() {
		return "", nil
	}
	mapPath := filepath.Join(absRoot, mapx.DefaultPath)
	g := stagedMap(absRoot, mapPath)
	if g == nil {
		loaded, err := mapx.Load(mapPath)
		if err != nil {
			return "", nil // no map to compile against: `docs-fresh` has nothing to confront either
		}
		g = loaded
	}
	read, err := scan.IndexReader(absRoot)
	if err != nil {
		return "", err
	}
	c, err := doct.NewWith(absRoot, g, read)
	if err != nil {
		return "", err
	}
	pages, err := c.Compiled()
	if err != nil {
		return "", err
	}
	ahead, err := scan.GovernedTreeChanges(absRoot, cfg)
	if err != nil {
		return "", err
	}
	apart := len(ahead) > 0
	var staged []string
	for _, p := range pages {
		if exec.Command("git", "-C", absRoot, "check-ignore", "-q", "--", p.Path).Run() == nil {
			continue
		}
		if apart {
			if err := stageBlob(absRoot, p.Path, p.Content); err != nil {
				return "", err
			}
		} else {
			dst := filepath.Join(absRoot, filepath.FromSlash(p.Path))
			if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
				return "", err
			}
			if err := os.WriteFile(dst, p.Content, 0o644); err != nil {
				return "", err
			}
			if out, err := exec.Command("git", "-C", absRoot, "add", "--", p.Path).CombinedOutput(); err != nil {
				return "", fmt.Errorf("stage %s: %v %s", p.Path, err, strings.TrimSpace(string(out)))
			}
		}
		staged = append(staged, p.Path)
	}
	if len(staged) == 0 {
		return "", nil
	}
	where := "written and staged"
	if apart {
		where = "staged apart — the tree holds changes the commit does not, and the files on disk keep them"
	}
	return fmt.Sprintf("· docs: %d page(s) compiled from the commit's specs, %s: %s", len(staged), where, strings.Join(staged, ", ")), nil
}

// stageBlob writes content into the index at rel (a path relative to the project root),
// without touching the file on disk.
func stageBlob(absRoot, rel string, content []byte) error {
	dir, err := os.MkdirTemp("", "anchors-blob-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	tmp := filepath.Join(dir, "blob")
	if err := os.WriteFile(tmp, content, 0o644); err != nil {
		return err
	}
	out, err := exec.Command("git", "-C", absRoot, "hash-object", "-w", tmp).Output()
	if err != nil {
		return fmt.Errorf("write the blob of %s: %w", rel, err)
	}
	sha := strings.TrimSpace(string(out))
	// `--cacheinfo` names the path from the repository's top, not from the directory the
	// command runs in: a project below the top needs its own prefix.
	prefix, err := exec.Command("git", "-C", absRoot, "rev-parse", "--show-prefix").Output()
	if err != nil {
		return fmt.Errorf("locate the project in the repository: %w", err)
	}
	path := strings.TrimSpace(string(prefix)) + filepath.ToSlash(rel)
	if out, err := exec.Command("git", "-C", absRoot, "update-index", "--add", "--cacheinfo", "100644,"+sha+","+path).CombinedOutput(); err != nil {
		return fmt.Errorf("stage %s: %v %s", rel, err, strings.TrimSpace(string(out)))
	}
	return nil
}
