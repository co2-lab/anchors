package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// runMutation is the mutation run as `anchors mutation` runs it: gremlins over every
// package of the module, merged into the single report declared in anchors.yaml
// (`mutation:`).
//
// One gremlins run per package, because a run over a subpackage writes each `file_name`
// relative to THAT package: `queue.go`, not `internal/queue/queue.go`. Ingested as is, a
// bare name matches every file of the map with that basename, and `cmd/anchors/flow/queue.go`
// would receive the score of `internal/queue/queue.go`. Each report is prefixed with its
// package directory before the merge.
//
// gremlins also DESCENDS into subdirectories: `./cmd/anchors` mutated every package under
// it, which each has its own run here. `-E /` keeps each run to the files of its own
// package (the regexp is matched against paths relative to the target, and only a file of
// a subdirectory has a `/` in it).
//
// `--timeout-coefficient 10`: gremlins times each mutant against the coverage run, which
// is fast when the build cache is warm. A mutant that needs a rebuild after the cache was
// trimmed then exceeds it, and a TIMED OUT mutant counts as killed — measured on
// internal/testlist: 32 of 32 timed out (a false 100%), and 30 killed / 1 survived with
// the larger coefficient. ANCHORS_MUTATION_TIMEOUT_COEFFICIENT overrides the 10.
//
// `--workers`: gremlins runs one mutant per CPU at once, and a package whose tests build
// binaries and call git adds its own load to a machine other sessions share — measured: 11
// of 21 mutants of keep_evidence.go timed out with 10 workers, 0 of 21 with 2. The default
// is a quarter of the CPUs; ANCHORS_MUTATION_WORKERS overrides it.
//
// `GOFLAGS=-count=1`: gremlins sets each mutant's time limit from its coverage run, and
// with Go's test cache warm that run is answered from the cache in milliseconds — almost
// every mutant of a slow package "timed out". Measured on cmd/anchors/quality/map_sync.go:
// 15 of 18 timed out with the cache; 15 killed, 3 lived, 0 timed out without it.
//
// Usage: go run ./tools/gates mutation [package-dir | file.go ...]   (default: every package)
func runMutation(args []string) int {
	if findTool("gremlins") == "" {
		return skipMissing("gremlins", "go install github.com/go-gremlins/gremlins/cmd/gremlins@latest")
	}
	if err := os.MkdirAll(".anchors", 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	mod, err := goModule()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	entries := mutationEntries(args)
	if len(args) == 0 {
		if entries, err = modulePackages(mod); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	tmp, err := os.MkdirTemp("", "anchors-mutation")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer os.RemoveAll(tmp)

	coefficient := envOr("ANCHORS_MUTATION_TIMEOUT_COEFFICIENT", "10")
	workers := envOr("ANCHORS_MUTATION_WORKERS", strconv.Itoa(max(runtime.NumCPU()/4, 1)))
	var reports []map[string]any
	for i, e := range entries {
		fmt.Printf("[%d/%d] %s\n", i+1, len(entries), e.label())
		rep, ok := mutatePackage(e, mod, tmp, coefficient, workers)
		if ok {
			reports = append(reports, rep)
		}
	}
	if len(reports) == 0 {
		fmt.Fprintln(os.Stderr, "no package produced a report")
		return 1
	}
	merged := mergeMutationReports(reports)
	b, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	const out = ".anchors/mutation.json"
	if err := os.WriteFile(out, append(b, '\n'), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Printf("merged %d package report(s) into %s\n", len(reports), out)
	return 0
}

// mutationEntry is one gremlins run: a package directory, and the one file of it to mutate
// when the run is for a file.
type mutationEntry struct {
	pkg, only string
}

func (e mutationEntry) label() string {
	if e.only != "" {
		return e.pkg + "|" + e.only
	}
	return e.pkg
}

// mutationEntries reads the arguments: a `.go` file (relative, or absolute under the
// repository, as `run_changed: {{files}}` passes it) mutates that file alone; anything
// else is a package directory.
func mutationEntries(args []string) []mutationEntry {
	root, _ := os.Getwd()
	var out []mutationEntry
	for _, a := range args {
		if rel, err := filepath.Rel(root, a); err == nil && filepath.IsAbs(a) {
			a = rel
		}
		a = filepath.ToSlash(a)
		if strings.HasSuffix(a, ".go") {
			out = append(out, mutationEntry{pkg: path.Dir(a), only: path.Base(a)})
		} else {
			out = append(out, mutationEntry{pkg: a})
		}
	}
	return out
}

// modulePackages are the module's packages below its root, as directories. The root
// package and this repository's own tooling (tools/) are left out: the tooling is no unit
// of the product, and the map has no node to take its score.
func modulePackages(mod string) ([]mutationEntry, error) {
	b, err := output("go", "list", "./...")
	if err != nil {
		return nil, fmt.Errorf("go list: %w", err)
	}
	var out []mutationEntry
	for _, p := range strings.Fields(string(b)) {
		rel := strings.TrimPrefix(p, mod+"/")
		if p == mod || strings.HasPrefix(rel, "tools/") {
			continue
		}
		out = append(out, mutationEntry{pkg: rel})
	}
	return out, nil
}

// excludesFor are gremlins' exclusions for a run: the subdirectories always, and, for a
// file's run, every sibling file by its exact name.
func excludesFor(e mutationEntry, siblings []string) []string {
	out := []string{"-E", "/"}
	if e.only == "" {
		return out
	}
	for _, name := range siblings {
		if name != e.only {
			out = append(out, "-E", "^"+regexp.QuoteMeta(name)+"$")
		}
	}
	return out
}

// mutatePackage runs gremlins over one entry and returns its report with each file named
// from the repository root. A run that fails is reported and left out.
func mutatePackage(e mutationEntry, mod, tmp, coefficient, workers string) (map[string]any, bool) {
	siblings, _ := filepath.Glob(filepath.Join(filepath.FromSlash(e.pkg), "*.go"))
	for i := range siblings {
		siblings[i] = filepath.Base(siblings[i])
	}
	raw := filepath.Join(tmp, "raw.json")
	_ = os.Remove(raw)
	args := append([]string{"unleash", "./" + e.pkg}, excludesFor(e, siblings)...)
	args = append(args, "--timeout-coefficient", coefficient, "--workers", workers, "--output", raw)
	cmd := exec.Command(findTool("gremlins"), args...)
	cmd.Env = append(os.Environ(), "GOFLAGS="+strings.TrimSpace(os.Getenv("GOFLAGS")+" -count=1"))
	var log bytes.Buffer
	cmd.Stdout, cmd.Stderr = &log, &log
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "  gremlins failed on %s:\n%s", e.pkg, lastLines(log.String(), 5))
		return nil, false
	}
	sc := bufio.NewScanner(&log)
	for sc.Scan() {
		if strings.Contains(sc.Text(), "Test efficacy") {
			fmt.Println("  " + strings.TrimSpace(sc.Text()))
		}
	}

	var rep map[string]any
	if b, err := os.ReadFile(raw); err == nil && len(bytes.TrimSpace(b)) > 0 {
		dec := json.NewDecoder(bytes.NewReader(b))
		dec.UseNumber() // the report's numbers stay as they were written
		if err := dec.Decode(&rep); err != nil {
			fmt.Fprintf(os.Stderr, "  gremlins wrote an unreadable report for %s: %v\n", e.pkg, err)
			return nil, false
		}
	} else {
		// Nothing to mutate (a file of data tables, say): gremlins succeeds and writes no
		// report. That is an answer, not a failure — the file is listed with no mutation,
		// which the ingest reads as "nothing to mutate".
		fmt.Println("  nothing to mutate")
		rep = emptyMutationReport(mod, e, siblings)
	}
	return prefixMutationReport(rep, e.pkg), true
}

// emptyMutationReport lists the entry's files with no mutation: the file a file's run is
// for, or every non-test file of the package.
func emptyMutationReport(mod string, e mutationEntry, siblings []string) map[string]any {
	var files []any
	for _, name := range siblings {
		if (e.only != "" && name != e.only) || strings.HasSuffix(name, "_test.go") {
			continue
		}
		files = append(files, map[string]any{"file_name": name, "mutations": []any{}})
	}
	return map[string]any{"go_module": mod, "files": files}
}

// prefixMutationReport names each file of a package's report from the repository root.
func prefixMutationReport(rep map[string]any, pkg string) map[string]any {
	files, _ := rep["files"].([]any)
	for _, f := range files {
		if m, ok := f.(map[string]any); ok {
			if name, ok := m["file_name"].(string); ok {
				m["file_name"] = pkg + "/" + name
			}
		}
	}
	return rep
}

// mergeMutationReports is the one report: the first's module, and every package's files.
func mergeMutationReports(reports []map[string]any) map[string]any {
	var files []any
	for _, r := range reports {
		if fs, ok := r["files"].([]any); ok {
			files = append(files, fs...)
		}
	}
	if files == nil {
		files = []any{}
	}
	return map[string]any{"go_module": reports[0]["go_module"], "files": files}
}

func goModule() (string, error) {
	b, err := output("go", "list", "-m")
	if err != nil {
		return "", fmt.Errorf("go list -m: %w", err)
	}
	return strings.TrimSpace(string(b)), nil
}

func envOr(name, def string) string {
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		return v
	}
	return def
}

func lastLines(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n") + "\n"
}
