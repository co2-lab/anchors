// @anchors
//   code: VRFRC
//   ref: VRFRA

package gate

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/shell"
)

// versionRepo commits the sources with what they hold at HEAD, and returns the root and a
// configuration whose product app reads them.
func versionRepo(t *testing.T, atHead map[string]string, sources ...config.VersionSource) (string, *config.Config) {
	t.Helper()
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		c := exec.Command("git", append([]string{"-C", root}, args...)...)
		c.Env = append(shell.WithoutRepoEnv(os.Environ()), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null")
		if out, err := c.CombinedOutput(); err != nil {
			t.Skipf("git: %v %s", err, out)
		}
	}
	git("init", "-q")
	writeFile(t, root, "README.md", "x\n")
	for f, body := range atHead {
		writeFile(t, root, f, body)
	}
	git("add", ".")
	git("commit", "-q", "-m", "init")
	return root, &config.Config{Versioning: config.Versioning{"app": {Scheme: "semver", Source: sources}}}
}

var jsonSource = config.VersionSource{File: "app/version.json", Path: "version"}

func versionJSON(v string) string { return `{"version": "` + v + `"}` }

func TestVersionForward_passes(t *testing.T) {
	t.Run("VRFRA-B01: A version kept or raised passes", func(t *testing.T) {
		root, cfg := versionRepo(t, map[string]string{"app/version.json": versionJSON("1.2.0")}, jsonSource)
		if v, msg := checkVersionForward(config.Gate{}, root, nil, cfg); v != Pass {
			t.Errorf("kept: %s %s", v, msg)
		}
		writeFile(t, root, "app/version.json", versionJSON("1.3.0"))
		if v, msg := checkVersionForward(config.Gate{}, root, nil, cfg); v != Pass {
			t.Errorf("raised: %s %s", v, msg)
		}
	})
	t.Run("VRFRA-B03: A version HEAD did not hold passes", func(t *testing.T) {
		root, cfg := versionRepo(t, nil, jsonSource)
		writeFile(t, root, "app/version.json", versionJSON("0.1.0"))
		if v, msg := checkVersionForward(config.Gate{}, root, nil, cfg); v != Pass {
			t.Errorf("born now: %s %s", v, msg)
		}
	})
	t.Run("VRFRA-B04: No product with a source skips the gate", func(t *testing.T) {
		cfg := &config.Config{Versioning: config.Versioning{"cli": {Scheme: "semver"}}}
		if v, msg := checkVersionForward(config.Gate{}, t.TempDir(), nil, cfg); v != Skip || !strings.Contains(msg, "no version written") {
			t.Errorf("got %s %s", v, msg)
		}
	})
}

func TestVersionForward_fails(t *testing.T) {
	t.Run("VRFRA-B02: A version that goes back, moves sideways, cannot be read or disagrees fails naming the product", func(t *testing.T) {
		txt := config.VersionSource{File: "app/build.txt", Regex: `VERSION=(\S+)`}
		for _, c := range []struct {
			head, now, txtNow string
			says              []string
		}{
			{"1.2.0", "1.1.0", "1.1.0", []string{"app", "1.2.0", "1.1.0"}},
			{"1.2.0+a", "1.2.0+b", "1.2.0+b", []string{"app", "1.2.0+a", "1.2.0+b"}},
			{"1.2.0", "1.3", "1.3", []string{"app", `"1.3"`}},
			{"1.2.0", "1.3.0", "1.2.0", []string{"app", "app/version.json 1.3.0", "app/build.txt 1.2.0"}},
		} {
			root, cfg := versionRepo(t, map[string]string{"app/version.json": versionJSON(c.head), "app/build.txt": "VERSION=" + c.head + "\n"}, jsonSource, txt)
			writeFile(t, root, "app/version.json", versionJSON(c.now))
			writeFile(t, root, "app/build.txt", "VERSION="+c.txtNow+"\n")
			v, msg := checkVersionForward(config.Gate{}, root, nil, cfg)
			if v != Fail {
				t.Errorf("%s → %s: %s %s", c.head, c.now, v, msg)
				continue
			}
			for _, s := range c.says {
				if !strings.Contains(msg, s) {
					t.Errorf("%s → %s: %q does not say %q", c.head, c.now, msg, s)
				}
			}
		}
	})
}
