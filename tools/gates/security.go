package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// runSecrets is `no-secret-leaked` — gitleaks (MIT).
//
// A leaked secret is the one finding with no way back: rotating a key is expensive and git
// history keeps the original. So it is the one born blocking. Before a commit it scans what
// is STAGED (0.7s measured); in CI it scans the history (29s measured) — hence the gate's
// `cost: fast`, and the mode chosen by the phase.
func runSecrets(_ []string) int {
	tool := findTool("gitleaks")
	if tool == "" {
		return skipMissing("gitleaks", "brew install gitleaks")
	}
	ci := os.Getenv("ANCHORS_PHASE") == "ci"
	args := []string{"protect", "--staged", "--no-banner", "--redact", "--exit-code", "1"}
	where := "staged"
	if ci {
		args = []string{"detect", "--no-banner", "--redact", "--exit-code", "1"}
		where = "history"
	}
	if passthrough(tool, args...) == 0 {
		fmt.Printf("no secret in the %s\n", where)
		return 0
	}
	fmt.Println("\nSecret detected. If it is a false positive, declare it in .gitleaks.toml (allowlist)")
	fmt.Println("with the REASON — the allowlist is the statement that someone looked and approved.")
	return 1
}

// runDepsVuln is `dependency-vulnerable` — osv-scanner (Apache-2.0, Google).
//
// It reads the module's lockfile directly: deterministic, and the only network is the OSV
// database. Informative by decision: most CVEs of a small project are TRANSITIVE, through
// the toolchain, and blocking a commit for a test dependency's CVE stops the work without
// reducing risk. The metric evolves; when it is low and stable, it becomes blocking.
func runDepsVuln(_ []string) int {
	tool := findTool("osv-scanner")
	if tool == "" {
		return skipMissing("osv-scanner", "brew install osv-scanner")
	}
	// osv-scanner exits 1 when it FINDS something; its JSON is the answer either way.
	out, _ := output(tool, "scan", "source", "--lockfile", "go.mod", "--format", "json")
	fmt.Print(summarizeOSV(out))
	return 0
}

// summarizeOSV is the report of an osv-scanner JSON: the vulnerabilities by severity, with
// the packages that carry them.
func summarizeOSV(data []byte) string {
	var d struct {
		Results []struct {
			Packages []struct {
				Package struct {
					Name string `json:"name"`
				} `json:"package"`
				Vulnerabilities []struct {
					DatabaseSpecific struct {
						Severity string `json:"severity"`
					} `json:"database_specific"`
				} `json:"vulnerabilities"`
			} `json:"packages"`
		} `json:"results"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return "osv-scanner produced no readable output — gate SKIPPED\n"
	}
	count := map[string]int{}
	pkgs := map[string]map[string]bool{}
	all := map[string]bool{}
	total := 0
	for _, r := range d.Results {
		for _, p := range r.Packages {
			for _, v := range p.Vulnerabilities {
				s := v.DatabaseSpecific.Severity
				if s == "" {
					s = "UNKNOWN"
				}
				count[s]++
				if pkgs[s] == nil {
					pkgs[s] = map[string]bool{}
				}
				pkgs[s][p.Package.Name] = true
				all[p.Package.Name] = true
				total++
			}
		}
	}
	if total == 0 {
		return "no known vulnerability in go.mod\n"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d vulnerabilities in %d packages:\n", total, len(all))
	for _, s := range []string{"CRITICAL", "HIGH", "MODERATE", "LOW", "UNKNOWN"} {
		if count[s] > 0 {
			fmt.Fprintf(&b, "  %-9s %4d  (%s)\n", s, count[s], firstNames(pkgs[s], 6))
		}
	}
	b.WriteString("\nInformative: check whether the CVE reaches an EXECUTION path or only the toolchain —\n")
	b.WriteString("`go mod why <package>` shows who pulls it in.\n")
	return b.String()
}

// runLicenses is `license-compatible` — go-licenses (Apache-2.0, Google).
//
// Anchors ships as a STATIC BINARY: everything it imports goes into the executable, so
// copyleft in a dependency reaches the whole binary. Blocking for strong copyleft
// (AGPL/SSPL/GPL) because the consequence is legal and the count today is ZERO: blocking a
// number that is already zero stops no one, and prevents the FIRST entry.
func runLicenses(_ []string) int {
	tool := findTool("go-licenses")
	if tool == "" {
		return skipMissing("go-licenses", "go install github.com/google/go-licenses@latest")
	}
	out, err := output(tool, "report", "./cmd/anchors")
	if err != nil && len(out) == 0 {
		fmt.Println("go-licenses failed. Gate SKIPPED (not approved).")
		return 0
	}
	mod, _ := goModule()
	text, bad := classifyLicenses(out, mod)
	fmt.Print(text)
	if bad {
		return 1
	}
	return 0
}

var (
	// strongCopyleft obliges whoever distributes the binary to open its code: incompatible
	// with this project's licence. An explicit exception (GPL with a linking exception) is
	// not that obligation.
	strongCopyleft = regexp.MustCompile(`(?i)\b(AGPL|SSPL|GPL-[23])`)
	// weakCopyleft obliges opening only the modified file: living with it is possible, but
	// whoever redistributes must know the obligation exists.
	weakCopyleft = regexp.MustCompile(`(?i)\b(LGPL|MPL|EPL|CDDL)`)
)

// classifyLicenses reads go-licenses' CSV (package, url, licence) and says whether a
// strong copyleft licence is among the third-party dependencies. The project's own module
// is no third party: go-licenses reports its source-available licence as "Unknown".
func classifyLicenses(csvData []byte, ownModule string) (string, bool) {
	rows, _ := csv.NewReader(bytes.NewReader(csvData)).ReadAll()
	count := map[string]int{}
	var strong, weak, unknown [][2]string
	for _, r := range rows {
		if len(r) < 3 {
			continue
		}
		pkg, lic := r[0], r[2]
		if ownModule != "" && strings.HasPrefix(pkg, ownModule) {
			continue
		}
		count[lic]++
		switch {
		case strongCopyleft.MatchString(lic) && !strings.Contains(strings.ToLower(lic), "exception"):
			strong = append(strong, [2]string{pkg, lic})
		case weakCopyleft.MatchString(lic):
			weak = append(weak, [2]string{pkg, lic})
		case lic == "" || strings.EqualFold(lic, "unknown"):
			unknown = append(unknown, [2]string{pkg, lic})
		}
	}
	total := 0
	lics := make([]string, 0, len(count))
	for l, n := range count {
		total += n
		lics = append(lics, l)
	}
	sort.Slice(lics, func(i, j int) bool {
		if count[lics[i]] != count[lics[j]] {
			return count[lics[i]] > count[lics[j]]
		}
		return lics[i] < lics[j]
	})
	var b strings.Builder
	fmt.Fprintf(&b, "%d third-party dependencies in the binary, %d distinct licences\n", total, len(count))
	for i, l := range lics {
		if i == 6 {
			break
		}
		fmt.Fprintf(&b, "  %-42s %d\n", truncate(l, 40), count[l])
	}
	list := func(title string, items [][2]string, all bool) {
		if len(items) == 0 {
			return
		}
		b.WriteString("\n" + title + "\n")
		for i, it := range items {
			if !all && i == 6 {
				break
			}
			l := it[1]
			if l == "" {
				l = "(empty)"
			}
			fmt.Fprintf(&b, "  %-54s %s\n", truncate(it[0], 52), l)
		}
	}
	list(fmt.Sprintf("weak copyleft (%d) — the obligation reaches the modified file:", len(weak)), weak, false)
	list(fmt.Sprintf("licence NOT declared (%d) — with no explicit licence the default is 'all rights\nreserved', which forbids redistribution:", len(unknown)), unknown, false)
	if len(strong) > 0 {
		list(fmt.Sprintf("STRONG copyleft (%d) — INCOMPATIBLE with this project's licence:", len(strong)), strong, true)
		b.WriteString("\nA GPL/AGPL dependency would require the whole binary to be GPL.\n")
		return b.String(), true
	}
	b.WriteString("\nno strong copyleft licence.\n")
	return b.String(), false
}

// runSBOM is `sbom-generated` — syft (Apache-2.0).
//
// The SBOM is an audit ARTEFACT, not a rule that can fail code — hence informative.
// Scanning the tree costs time that does not belong to a pre-commit, so it runs in CI.
func runSBOM(_ []string) int {
	tool := findTool("syft")
	if tool == "" {
		return skipMissing("syft", "brew install syft")
	}
	if err := os.MkdirAll("reports", 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	const out = "reports/sbom-cyclonedx.json"
	// The module only: the site's node_modules and agent worktrees are not what ships.
	if _, err := output(tool, "scan", "dir:.", "--exclude", "./site/**", "--exclude", "./**/node_modules/**",
		"--exclude", "./.claude/**", "--output", "cyclonedx-json="+out, "-q"); err != nil {
		fmt.Println("syft failed. Gate SKIPPED (not approved).")
		return 0
	}
	b, err := os.ReadFile(out)
	if err != nil {
		fmt.Println("syft wrote no SBOM. Gate SKIPPED (not approved).")
		return 0
	}
	fmt.Print(summarizeSBOM(b, out))
	return 0
}

// summarizeSBOM counts a CycloneDX SBOM's components by type.
func summarizeSBOM(data []byte, where string) string {
	var d struct {
		Components []struct {
			Type string `json:"type"`
		} `json:"components"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return "syft produced no readable SBOM — gate SKIPPED\n"
	}
	types := map[string]int{}
	for _, c := range d.Components {
		t := c.Type
		if t == "" {
			t = "?"
		}
		types[t]++
	}
	var b strings.Builder
	fmt.Fprintf(&b, "CycloneDX SBOM generated: %d components\n", len(d.Components))
	for _, t := range topByCount(types, 6) {
		fmt.Fprintf(&b, "  %-14s %d\n", t, types[t])
	}
	fmt.Fprintf(&b, "  → %s\n", where)
	return b.String()
}

// runSemgrep is `insecure-pattern` — semgrep (LGPL-2.1) with the `p/gosec` ruleset.
//
// Informative: a gosec finding is a question ("is this safe in your context?"), not a
// verdict. With no argument it scans the project (the `scope_full` mode); with arguments,
// the files the engine passes.
func runSemgrep(args []string) int {
	tool := findTool("semgrep")
	if tool == "" {
		return skipMissing("semgrep", "brew install semgrep")
	}
	targets := args
	if len(targets) == 0 {
		targets = []string{"."}
	}
	// semgrep exits non-zero when it FINDS something; its JSON is the answer either way.
	out, _ := output(tool, append([]string{"--config=p/gosec", "--json", "--quiet"}, targets...)...)
	fmt.Print(summarizeSemgrep(out))
	return 0
}

// summarizeSemgrep is the report of a semgrep JSON: the findings by rule, with an example.
func summarizeSemgrep(data []byte) string {
	var d struct {
		Results []struct {
			CheckID string `json:"check_id"`
			Path    string `json:"path"`
			Start   struct {
				Line int `json:"line"`
			} `json:"start"`
		} `json:"results"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		return "semgrep produced no readable output — gate SKIPPED\n"
	}
	if len(d.Results) == 0 {
		return "no known insecure pattern (p/gosec)\n"
	}
	byRule := map[string]int{}
	example := map[string]string{}
	files := map[string]bool{}
	for _, r := range d.Results {
		rule := r.CheckID[strings.LastIndex(r.CheckID, ".")+1:]
		byRule[rule]++
		if _, ok := example[rule]; !ok {
			example[rule] = fmt.Sprintf("%s:%d", r.Path, r.Start.Line)
		}
		files[r.Path] = true
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d gosec finding(s) in %d file(s):\n", len(d.Results), len(files))
	for _, rule := range topByCount(byRule, 10) {
		fmt.Fprintf(&b, "  %3dx %s  (e.g. %s)\n", byRule[rule], rule, example[rule])
	}
	b.WriteString("\nInformative: each finding is a QUESTION about the context, not a verdict.\n")
	b.WriteString("Mark what is deliberate with `// nosemgrep: <rule>` and the reason.\n")
	return b.String()
}

// topByCount are the keys by count, most first, ties by name, at most n.
func topByCount(m map[string]int, n int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]] != m[keys[j]] {
			return m[keys[i]] > m[keys[j]]
		}
		return keys[i] < keys[j]
	})
	if len(keys) > n {
		keys = keys[:n]
	}
	return keys
}

// firstNames are a set's names in order, at most n, with an ellipsis for the rest.
func firstNames(set map[string]bool, n int) string {
	names := make([]string, 0, len(set))
	for k := range set {
		names = append(names, k)
	}
	sort.Strings(names)
	more := ""
	if len(names) > n {
		names, more = names[:n], "…"
	}
	return strings.Join(names, ", ") + more
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
