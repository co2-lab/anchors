// @anchors
//   ref: PRRPR

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// rootFixture builds <tmp>/outer/anchors.yaml, <tmp>/outer/inner/anchors.yaml and the
// plain directories <tmp>/outer/pkg/deep and <tmp>/outer/inner/src, and returns <tmp>.
func rootFixture(t *testing.T) string {
	t.Helper()
	tmp, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"outer/pkg/deep", "outer/inner/src", "loose/dir"} {
		if err := os.MkdirAll(filepath.Join(tmp, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{"outer/" + DefaultFile, "outer/inner/" + DefaultFile} {
		if err := os.WriteFile(filepath.Join(tmp, f), []byte("lang: en\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return tmp
}

func TestProjectRoot_walksUpToTheNearestProject(t *testing.T) {
	t.Run("PRRPR-B01: The project root is the nearest directory above the start that holds the config", func(t *testing.T) {})
	tmp := rootFixture(t)
	for start, want := range map[string]string{
		"outer/pkg/deep":  "outer",
		"outer":           "outer",
		"outer/inner/src": "outer/inner", // the nearest project wins over the outer one
	} {
		if got := ProjectRoot(filepath.Join(tmp, start)); got != filepath.Join(tmp, want) {
			t.Errorf("ProjectRoot(%s) = %s, want %s", start, got, filepath.Join(tmp, want))
		}
	}
}

func TestProjectRoot_outsideAProjectReturnsTheStart(t *testing.T) {
	t.Run("PRRPR-B02: With no project above the start, the start comes back unchanged", func(t *testing.T) {})
	tmp := rootFixture(t)
	if _, err := os.Stat(filepath.Join(filepath.Dir(tmp), DefaultFile)); err == nil {
		t.Skip("the temporary directory sits inside a project")
	}
	start := filepath.Join(tmp, "loose/dir")
	if got := ProjectRoot(start); got != start {
		t.Errorf("ProjectRoot(%s) = %s, want the start itself", start, got)
	}
	t.Chdir(start)
	if got := ProjectRoot("."); got != "." {
		t.Errorf("ProjectRoot(.) = %s, want the relative start unchanged", got)
	}
}

func TestAbsRoot_defaultWalksUp(t *testing.T) {
	t.Run("PRRPR-B03: With no root given, the root is found by walking up from the working directory", func(t *testing.T) {})
	tmp := rootFixture(t)
	t.Chdir(filepath.Join(tmp, "outer/pkg/deep"))
	for _, root := range []string{".", ""} {
		got, err := AbsRoot(root)
		if err != nil || got != filepath.Join(tmp, "outer") {
			t.Errorf("AbsRoot(%q) = %s, %v; want %s", root, got, err, filepath.Join(tmp, "outer"))
		}
	}
}

func TestAbsRoot_explicitRootIsRespected(t *testing.T) {
	t.Run("PRRPR-X01: An explicit root is made absolute and never walked above", func(t *testing.T) {})
	tmp := rootFixture(t)
	t.Chdir(filepath.Join(tmp, "outer"))
	got, err := AbsRoot("pkg/deep")
	if err != nil || got != filepath.Join(tmp, "outer/pkg/deep") {
		t.Errorf("AbsRoot(pkg/deep) = %s, %v; want %s", got, err, filepath.Join(tmp, "outer/pkg/deep"))
	}
}

func TestAbsRoot_alwaysAbsolute(t *testing.T) {
	t.Run("PRRPR-I01: The resolved root is always absolute, even outside any project", func(t *testing.T) {})
	tmp := rootFixture(t)
	if _, err := os.Stat(filepath.Join(filepath.Dir(tmp), DefaultFile)); err == nil {
		t.Skip("the temporary directory sits inside a project")
	}
	t.Chdir(filepath.Join(tmp, "loose/dir"))
	for _, root := range []string{".", "", "sub"} {
		got, err := AbsRoot(root)
		if err != nil || !filepath.IsAbs(got) {
			t.Errorf("AbsRoot(%q) = %q, %v; want an absolute path", root, got, err)
		}
	}
}
