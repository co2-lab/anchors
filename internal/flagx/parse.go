package flagx

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Dir is where flag files live — outside the target tree, on the `product/` precedent.
const Dir = "flags"

// Suffix marks a flag file.
const Suffix = ".flag.md"

// Scenario is one row of the flag's scenario table: a code, a condition, and what holds
// when it applies.
//
// The Condition is PARSED, unlike the prose condition of a flow transition, and the
// difference is deliberate. A flow exit is judged by a person reading it; a flag scenario
// is confronted against code and tests, and confronting prose is not possible. It is the
// decision that produced the fixed grammar in `grammar.go`.
type Scenario struct {
	Code string
	Cond Condition
	Then string
	File string
	Line int
	// Err records a condition the grammar refused. The scenario is kept rather than
	// dropped: a row nobody can parse is a finding to report, and dropping it would make
	// it vanish silently — the failure mode the gates exist to end.
	Err error
}

// Flag is one feature flag and the scenarios its value opens.
type Flag struct {
	Code      string
	Name      string
	File      string
	Scenarios []Scenario
}

// Absent reports whether the flag declares what happens when it does not answer.
//
// It is the case that most often breaks and that almost nobody writes: flag service down,
// a new environment, a local test. Every tool in the benchmark can express it (Flagsmith
// names it `Is Not Set`), and none require it — which is why it gets its own question here.
func (f Flag) Absent() bool {
	for _, s := range f.Scenarios {
		if s.Cond.Op == OpAbsent {
			return true
		}
	}
	return false
}

// scenarioRowRE matches a scenario row of the table:
// `| ` + "`CHKUT-G01`" + ` | ` + "`= \"off\"`" + ` | the old checkout answers |`
//
// The TABLE form only. A scenario is a row because its three parts (code, condition,
// consequence) are columns of one statement — unlike a spec rule, which may span
// paragraphs and therefore accepts three forms.
var scenarioRowRE = regexp.MustCompile(
	"^\\s*\\|\\s*`?([A-Z0-9]{3,6}-G[0-9]{2})`?\\s*\\|([^|]*)\\|(.*)$")

// ParseFile reads one flag file.
func ParseFile(path, rel string) (Flag, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Flag{}, err
	}
	return parse(string(b), rel), nil
}

// ParseContent reads a flag file already in memory.
//
// The gates receive the CONTENT, never the path — `check` reads each file once and hands
// the bytes to every gate that applies, and re-reading from disk here would both cost a
// second read and risk confronting a different version than the one the run is judging.
func ParseContent(content, rel string) Flag { return parse(content, rel) }

func parse(content, rel string) Flag {
	f := Flag{File: rel}
	f.Name = strings.TrimSuffix(filepath.Base(rel), Suffix)

	for i, line := range strings.Split(content, "\n") {
		m := scenarioRowRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		code := m[1]
		cond, err := Parse(m[2])
		then := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(m[3]), "|"))
		f.Scenarios = append(f.Scenarios, Scenario{
			Code: code, Cond: cond, Then: then,
			File: rel, Line: i + 1, Err: err,
		})
		if f.Code == "" {
			if unit, _, ok := strings.Cut(code, "-"); ok {
				f.Code = unit
			}
		}
	}
	return f
}

// Load scans the project's `flags/` folder.
//
// A project with no `flags/` is not an error — it is a project that has not declared a
// flag yet, the same reading `flowx.Build` gives an absent `flows/`.
func Load(root string) ([]Flag, error) {
	entries, err := os.ReadDir(filepath.Join(root, Dir))
	if err != nil {
		return nil, nil
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), Suffix) {
			names = append(names, e.Name())
		}
	}
	// STABLE order across runs, as in `flowx.Build`: without it two scans of the same
	// repository report findings in different order, and the diff becomes noise.
	sort.Strings(names)

	var out []Flag
	for _, name := range names {
		rel := filepath.ToSlash(filepath.Join(Dir, name))
		f, err := ParseFile(filepath.Join(root, Dir, name), rel)
		if err != nil {
			continue
		}
		out = append(out, f)
	}
	return out, nil
}

// ByCode indexes scenarios by code, for resolving a `@gated-by` declaration.
func ByCode(flags []Flag) map[string]Scenario {
	out := map[string]Scenario{}
	for _, f := range flags {
		for _, s := range f.Scenarios {
			out[s.Code] = s
		}
	}
	return out
}
