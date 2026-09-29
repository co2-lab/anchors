package gate

import (
	"sort"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/testlist"
)

// projectTests reads the project's tests through the dialect's `tests` source — a pattern
// or a script the project declares — over the map's test files, once per scan.
//
// How a test is written belongs to the project's test library, not to the engine: the
// gates that read a test's title (feature-test-match, test-traceable, scenario-coverage)
// used to carry Jest's and Go's call syntax themselves, so a project on any other library
// had no title read at all. `declared` is false when the project says nothing and its
// family has no default; the callers then fall back to what they can do without titles.
// An error — a script that fails, output outside the contract — is the caller's to
// report: the tests could not be read, and answering as if there were none would be the
// silence the gates exist to break.
func projectTests(root string, g *mapx.Graph, cfg *config.Config) (tests []testlist.Test, declared bool, err error) {
	src := testsSource(cfg)
	if !src.Declared() {
		return nil, false, nil
	}
	ptMu.Lock()
	defer ptMu.Unlock()
	if ptOK && ptRoot == root && ptGraph == g && ptCfg == cfg {
		return ptTests, true, ptErr
	}
	var files []string
	if g != nil {
		for _, n := range g.Nodes {
			if n.Kind == mapx.KindTest && !n.Support {
				files = append(files, n.ID)
			}
		}
	}
	// Under `--index` the tests are listed from what the commit records, as the gates read
	// its files: the lines a test is found at must be lines of the content they index.
	dir, _, err := indexWorkdir(root)
	if err != nil {
		return nil, true, err
	}
	tests, err = testlist.List(dir, files, src)
	ptRoot, ptGraph, ptCfg, ptTests, ptErr, ptOK = root, g, cfg, tests, err, true
	return tests, true, err
}

// testsSource is the effective `tests` source of the project's dialect: its own, or its
// family's.
func testsSource(cfg *config.Config) testlist.Source {
	if cfg == nil {
		return testlist.Source{}
	}
	if t := cfg.DialectFor().Tests; t != nil {
		return testlist.Source{Pattern: t.Pattern, Script: t.Script}
	}
	return testlist.Source{}
}

// testsIn keeps the tests of the given files, in the order the files are given and, in
// each file, in the order of their lines.
func testsIn(tests []testlist.Test, files []string) []testlist.Test {
	rank := make(map[string]int, len(files))
	for i, f := range files {
		if _, seen := rank[f]; !seen {
			rank[f] = i
		}
	}
	var out []testlist.Test
	for _, t := range tests {
		if _, ok := rank[t.File]; ok {
			out = append(out, t)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if rank[out[i].File] != rank[out[j].File] {
			return rank[out[i].File] < rank[out[j].File]
		}
		return out[i].Line < out[j].Line
	})
	return out
}

// ONE READING PER SCAN, keyed like docs-fresh and duplication: the root, the map and the
// configuration of the scan.
var (
	ptMu    sync.Mutex
	ptRoot  string
	ptGraph *mapx.Graph
	ptCfg   *config.Config
	ptTests []testlist.Test
	ptErr   error
	ptOK    bool
)

// resetProjectTestsCache clears the memory between scans, for the tests.
func resetProjectTestsCache() {
	ptMu.Lock()
	defer ptMu.Unlock()
	ptOK, ptRoot, ptGraph, ptCfg, ptTests, ptErr = false, "", nil, nil, nil, nil
}

// realTests drops, from test paths, the support files of the map: files of a test layer
// that serve the tests (sub-flows, aggregators, helpers) without being tests. None of them
// proves a scenario, so none of them is confronted as a test.
func realTests(g *mapx.Graph, paths []string) []string {
	if g == nil {
		return paths
	}
	support := map[string]bool{}
	for _, n := range g.Nodes {
		if n.Support {
			support[n.ID] = true
		}
	}
	if len(support) == 0 {
		return paths
	}
	var out []string
	for _, p := range paths {
		if !support[p] {
			out = append(out, p)
		}
	}
	return out
}

// listedFiles is the set of files in which the source lists at least one test.
//
// A source describes the files it can read, not every test file of the map: the ts
// family's `it(` describes a Jest file and not a Maestro flow in YAML, which the same
// project keeps under a test layer. A file where the source lists nothing is a file the
// source does not describe, and the gates read it as they do without a source — the code
// anywhere in it. Taking "no test listed" for "no test cites the code" failed 654 flows in
// the reference app that name their scenario in `name:` and in `tags:`.
func listedFiles(tests []testlist.Test) map[string]bool {
	out := make(map[string]bool, len(tests))
	for _, t := range tests {
		out[t.File] = true
	}
	return out
}
