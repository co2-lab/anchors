// @anchors
//   code: RLSTV
//   ref: RLSTA

package ops

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/changelog"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/version"
)

// A stage that releases writes the release's record itself: the tag of its product, on the
// sha that went out, with the version that sha carries, and the GitHub Release of that tag
// with the release's changelog. A tag the workflow wrote on its own could name another
// version, or another commit, than the one the stage passed; written here, the record and
// the tag cannot disagree. Each tag is then a version really released, which is what lets
// the changelog go release by release.

// releaseAPI calls GitHub's API, its error carrying what gh said — a 404 is an answer.
var releaseAPI = func(root string, args ...string) ([]byte, error) {
	if _, err := exec.LookPath("gh"); err != nil {
		return nil, errors.New("the gh CLI is not in the PATH: the release's tag is written through GitHub's API")
	}
	c := exec.Command("gh", append([]string{"api"}, args...)...)
	c.Dir = root
	var out, errb bytes.Buffer
	c.Stdout, c.Stderr = &out, &errb
	if err := c.Run(); err != nil {
		return out.Bytes(), fmt.Errorf("gh api %s: %v: %s", args[len(args)-1], err, strings.TrimSpace(errb.String()+" "+out.String()))
	}
	return out.Bytes(), nil
}

func notFound(err error) bool { return err != nil && strings.Contains(err.Error(), "404") }

// releaseOf writes the tag of the stage's product on its sha, and its GitHub Release. A tag
// already on this sha is a rerun: only a missing Release is written. A tag on another sha is
// the same version released twice, and fails.
func releaseOf(t stageTarget) (string, error) {
	name := t.stage.Release
	st, err := stateAt(t.root, t.cfg.Versioning, name, "release:", t.branch, t.sha)
	if err != nil {
		return "", err
	}
	ver := or(st.Current, st.Next)
	if ver == "" {
		return "", fmt.Errorf("%s has no version to release on %s: %s", name, short(t.sha), st.Reason)
	}
	tag := st.product.TagFor(ver)
	out, err := releaseAPI(t.root, "repos/{owner}/{repo}/git/ref/tags/"+tag)
	switch {
	case err == nil:
		var ref struct {
			Object struct {
				SHA  string `json:"sha"`
				Type string `json:"type"`
			} `json:"object"`
		}
		if err := json.Unmarshal(out, &ref); err != nil {
			return "", fmt.Errorf("read the tag %s: %w", tag, err)
		}
		if ref.Object.SHA != t.sha {
			return "", fmt.Errorf("the tag %s is already on %s: %s %s was released from another commit — bump to a greater version", tag, short(ref.Object.SHA), name, ver)
		}
	case notFound(err):
		if _, err := releaseAPI(t.root, "-X", "POST", "repos/{owner}/{repo}/git/refs", "-f", "ref=refs/tags/"+tag, "-f", "sha="+t.sha); err != nil {
			return "", fmt.Errorf("write the tag %s: %w", tag, err)
		}
	default:
		return "", fmt.Errorf("read the tag %s: %w", tag, err)
	}
	if _, err := releaseAPI(t.root, "repos/{owner}/{repo}/releases/tags/"+tag); err == nil {
		return fmt.Sprintf("✓ %s %s is released as %s on %s", name, ver, tag, short(t.sha)), nil
	} else if !notFound(err) {
		return "", fmt.Errorf("read the release %s: %w", tag, err)
	}
	body, err := releaseBody(t, st, ver)
	if err != nil {
		return "", err
	}
	if _, err := releaseAPI(t.root, "-X", "POST", "repos/{owner}/{repo}/releases", "-f", "tag_name="+tag, "-f", "name="+tag, "-f", "body="+body); err != nil {
		return "", fmt.Errorf("write the release %s: %w", tag, err)
	}
	return fmt.Sprintf("✓ %s %s released as %s on %s, with its changelog", name, ver, tag, short(t.sha)), nil
}

// releaseBody is the release's changelog: the commits since the product's release before it,
// rendered with the template of the product's first output, else the `changelog:` block's.
func releaseBody(t stageTarget, st *versionStateT, ver string) (string, error) {
	g := changelog.Git{Root: t.root}
	prev, err := previousRelease(g, st.product, st.scheme, t.sha, ver)
	if err != nil {
		return "", err
	}
	commits, err := g.Commits(prev, t.sha)
	if err != nil {
		return "", err
	}
	r := changelog.Classify(ver, time.Now().UTC().Format("2006-01-02"), commits)
	r.Tag = st.product.TagFor(ver)
	cl, label := t.cfg.Changelog, "changelog.template"
	if pc := st.product.Changelog; pc != nil && len(pc.Outputs) > 0 {
		cl, label = &pc.Outputs[0].Changelog, "versioning."+t.stage.Release+".changelog.outputs[0].template"
	}
	blocks, err := renderBlocks(t.root, cl, label, []changelog.Release{r})
	if err != nil || len(blocks) == 0 {
		return "", err
	}
	return strings.TrimRight(blocks[0].Text, "\n") + "\n", nil
}

// previousRelease is the tag of the highest version of the product reachable from `to`
// other than ver: the release the one being written follows.
func previousRelease(g changelog.Git, p *config.Product, s version.Scheme, to, ver string) (string, error) {
	tags, err := g.Tags(to)
	if err != nil {
		return "", err
	}
	var best *version.V
	bestTag := ""
	for _, tg := range tags {
		raw, ok := p.VersionOfTag(tg)
		if !ok || raw == ver {
			continue
		}
		v, err := s.Parse(raw)
		if err != nil {
			return "", fmt.Errorf("tag %s: %w", tg, err)
		}
		if best == nil || version.Compare(v, *best) > 0 {
			best, bestTag = &v, tg
		}
	}
	return bestTag, nil
}
