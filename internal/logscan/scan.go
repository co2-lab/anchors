// Package logscan scans the project's logs and identifies OCCURRENCES of declared failure.
//
// THE LOG FORMAT DOES NOT MATTER, and that is what makes scanning possible without
// dictating anything. If the failure code is on the line, it is the same sequence of
// characters in structured JSON, in plain text, in syslog, or in whatever anybody invents:
//
//	{"level":"error","rule":"CRED-E01","msg":"insufficient balance"}
//	2026-09-21 10:32:11 ERROR CRED-E01 insufficient balance
//	<134>1 2026-09-21T10:32:11Z app - - - CRED-E01 refused
//
// The three lines say the same thing, and Anchors reads them alike — because it looks for
// the CODE, not the format. The requirement is one, and it belongs to the project: that
// the handling path record the code of the failure it is handling. The `failure-logged`
// gate already charges that the handling logs; what this axis adds is that it logs WITH
// IDENTITY.
//
// The alternative — parsing each project's format — would be the heuristic that ages at
// every logger swap, and errs both ways: it matches what is not, and misses what is.
package logscan

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// ruleCodeRE validates the extracted code — the bridge between the log line and the `-E`
// rule.
//
// The validation exists because a loose pattern captures anything, and whatever it
// captures would enter the map as an occurrence of a failure nobody declared. The gate
// downstream would then report a failure that does not exist, and whoever chased it would
// find nothing.
var ruleCodeRE = regexp.MustCompile(`^[A-Z0-9]{3,6}-E[0-9]{2}$`)

// Result is what the scan found, ready to enter the map.
type Result struct {
	Occurrences []mapx.FailureSignal
	// Files is how many files were read, Lines how many lines. They go into the report
	// because "zero occurrences" has two opposite meanings — there was no failure, or there
	// was no log — and only the count tells them apart.
	Files int
	Lines int
	// Unknown holds the codes found in the log that NO spec declares.
	//
	// It is this layer's most valuable finding, and the one no observability tool gives:
	// the log shows an error the spec never foresaw. Without this field it would be
	// discarded in silence — which is exactly what nobody is watching.
	Unknown map[string]int
}

// Scan reads the declared logs and returns the occurrences by failure code.
//
// The graph comes in to separate what is DECLARED from what is not — and that separation
// produces this layer's most valuable finding: the log shows an error no spec foresaw.
// Without it the unknown code would be discarded in silence, which is exactly the case
// nobody is watching.
func Scan(root string, cfg *config.Config, g *mapx.Graph) (*Result, error) {
	if cfg == nil || cfg.Logs == nil || len(cfg.Logs.Paths) == 0 {
		// With no `logs:` declared Anchors scans nothing, and that is deliberate: hunting
		// for files that look like logs would read what it should not — a log usually
		// carries exactly what must not leak.
		return nil, nil
	}
	aliases, err := aliasMatchers(cfg.Logs)
	if err != nil {
		return nil, err
	}
	// The DECLARED SET is what defines the pattern — see `codesRE`.
	declared := map[string]bool{}
	if g != nil {
		declared = declaredCodes(root, g)
	}
	open, close := DefaultDelimiters[0], DefaultDelimiters[1]
	if len(cfg.Logs.Delimiters) == 2 {
		open, close = cfg.Logs.Delimiters[0], cfg.Logs.Delimiters[1]
	}
	known := codesRE(declared, open, close)
	unknown := unknownRE(open, close)
	var tsRE *regexp.Regexp
	if cfg.Logs.Timestamp != "" {
		if tsRE, err = regexp.Compile(cfg.Logs.Timestamp); err != nil {
			return nil, err
		}
	}

	res := &Result{Unknown: map[string]int{}}
	acc := map[string]*mapx.FailureSignal{}

	for _, pattern := range cfg.Logs.Paths {
		paths, _ := filepath.Glob(filepath.Join(root, pattern))
		for _, p := range paths {
			f, err := os.Open(p)
			if err != nil {
				continue
			}
			res.Files++
			sc := bufio.NewScanner(f)
			// A log line can be long (a whole stack trace in one JSON field), and the
			// Scanner's default buffer cuts at 64K — silently, losing exactly the line
			// that carries the most context.
			sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
			for sc.Scan() {
				line := sc.Text()
				res.Lines++
				code := ""
				switch {
				case known != nil && known.MatchString(line):
					code = known.FindStringSubmatch(line)[1]
				default:
					for _, a := range aliases {
						if a.re.MatchString(line) {
							code = a.code
							break
						}
					}
					// SECOND PASS, and the only place the measured risk lives: the code no
					// spec declares. Without a delimiter it demands that the line say it is
					// reporting a failure — otherwise a build id on an INFO line would
					// become an occurrence.
					// The level is only demanded when there is NO delimiter: `#[ORPHA-E88]`
					// is already unambiguous, and demanding `ERROR` there would lose the
					// failure a project records without a level word. With no delimiter,
					// the level is the only defence left against the measured noise.
					if code != "" {
						break // the alias answered: the line is legacy, and it has an owner
					}
					semDelim := open == "" && close == ""
					if !semDelim || failureLevelRE.MatchString(line) {
						if m := unknown.FindStringSubmatch(line); len(m) > 1 {
							res.Unknown[m[1]]++
						}
					}
					continue
				}
				if !ruleCodeRE.MatchString(code) {
					continue
				}
				o, ok := acc[code]
				if !ok {
					o = &mapx.FailureSignal{Rule: code}
					acc[code] = o
				}
				o.Count++
				if tsRE != nil {
					if m := tsRE.FindStringSubmatch(line); len(m) > 1 {
						o.First = earliest(o.First, m[1])
						o.Last = latest(o.Last, m[1])
					}
				}
			}
			f.Close()
		}
	}

	for _, o := range acc {
		res.Occurrences = append(res.Occurrences, *o)
	}
	sort.Slice(res.Occurrences, func(i, j int) bool {
		return res.Occurrences[i].Rule < res.Occurrences[j].Rule
	})

	return res, nil
}

// SpecFailureCodes returns the `-E` rules catalogued in a spec's content.
func SpecFailureCodes(content string) []string {
	var out []string
	for _, m := range specFailureRE.FindAllStringSubmatch(content, -1) {
		out = append(out, m[1])
	}
	return out
}

// declaredCodes reads from the map every `-E` rule the project's specs declare.
func declaredCodes(root string, g *mapx.Graph) map[string]bool {
	out := map[string]bool{}
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindSpec {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, n.ID))
		if err != nil {
			continue
		}
		for _, m := range specFailureRE.FindAllStringSubmatch(string(b), -1) {
			out[m[1]] = true
		}
	}
	return out
}

// specFailureRE finds a failure rule catalogued in a spec, in the three valid forms.
var specFailureRE = regexp.MustCompile(`(?m)(?:^#{1,6}\s+|^\s*\|\s*` + "`?" + `|^\s*-\s+\*\*)([A-Z0-9]{3,6}-E[0-9]{2})`)

// codesRE builds the pattern from the codes the project ACTUALLY declares — never from
// the generic shape of a code.
//
// Measured before this: a bare `\b[A-Z0-9]{3,6}-E[0-9]{2}\b` captured five codes in a
// five-line log, and one was a failure. The other four were a build id (`BUILD-E01`), a
// cache key (`USER-E42`), a SKU (`ITEM-E07`) and a trace id inside a URL (`ABCD-E99`).
// Eighty per cent false positive — worse than the language heuristic this project already
// discarded for erring nine times in ten.
//
// The fix is not a cleverer pattern: it is asking the right question. Anchors KNOWS which
// `-E` rules exist, so it does not need to recognise the SHAPE of a code — it needs to
// recognise THOSE. A pattern built from the declared set cannot invent a failure, because
// it only matches what some spec already catalogued.
//
// The consequence is deliberate: a log carrying a code nobody declared is NOT found here.
// That would lose the layer's most valuable finding, so the unknown is looked for in a
// second pass — see `unknownRE`, which only runs where the generic shape is safe.
func codesRE(declared map[string]bool, open, close string) *regexp.Regexp {
	if len(declared) == 0 {
		return nil
	}
	codes := make([]string, 0, len(declared))
	for c := range declared {
		codes = append(codes, regexp.QuoteMeta(c))
	}
	// STABLE order: without it the pattern changes between runs, and a regex that varies
	// makes two scans of the same log incomparable.
	sort.Strings(codes)
	body := `(` + strings.Join(codes, "|") + `)`
	if open == "" && close == "" {
		return regexp.MustCompile(`\b` + body + `\b`)
	}
	return regexp.MustCompile(regexp.QuoteMeta(open) + body + regexp.QuoteMeta(close))
}

// unknownRE finds a code in the generic shape — the second pass, and the one that carries
// the risk measured above.
//
// It only runs on lines the declared pattern did NOT match, and only where the context
// says the line is reporting a failure: the level word. A build id in an INFO line is not
// a failure, and it was exactly that kind of line that produced the four false positives.
//
// It is still a heuristic, and the verdict it produces says so — the unknown code is
// REPORTED, never bound to the map. Binding it would invent an owner; discarding it would
// lose the finding.
func unknownRE(open, close string) *regexp.Regexp {
	body := `([A-Z0-9]{3,6}-E[0-9]{2})`
	if open == "" && close == "" {
		return regexp.MustCompile(`\b` + body + `\b`)
	}
	return regexp.MustCompile(regexp.QuoteMeta(open) + body + regexp.QuoteMeta(close))
}

// DefaultDelimiters wrap the code in the log: `#[CRED-E01]`.
//
// The choice is measured, and the leading `#` is what closes the last gap. Against six
// noise lines that look like a failure:
//
//	bare code     5 false positives    build id, cache key, SKU, trace id
//	[CRED-E01]    2 false positives    markdown in prose, array index (`logs[...]`)
//	#[CRED-E01]   0 false positives
//
// All three catch the real failures alike — what changes is only what they capture
// BESIDES. And writing `#[` instead of `[` costs whoever logs exactly the same.
var DefaultDelimiters = [2]string{"#[", "]"}

// failureLevelRE recognises that the line is reporting a failure.
//
// Deliberately narrow: `error`, `fatal`, `severe`, `critical` and their common
// abbreviations, in any case. A `warn` does not enter — a warning is not the failure the
// spec declared, and widening this is the cheapest way to bring the false positives back.
var failureLevelRE = regexp.MustCompile(`(?i)\b(error|err|fatal|severe|critical|crit|panic|exception)\b`)

// aliasMatchers builds the LEGACY text patterns — the log that already exists and nobody
// will rewrite to carry the code.
//
// It is the transition bridge, and it has the problem of every heuristic: it matches what
// is not, and misses what is. That is why it comes AFTER the code: a line that already
// carries `CRED-E01` needs no guess, and letting the guess override the identity would
// trade certainty for approximation.
func aliasMatchers(lc *config.LogsConfig) ([]aliasPattern, error) {
	var out []aliasPattern
	for code, pattern := range lc.Aliases {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		out = append(out, aliasPattern{code: code, re: re})
	}
	// STABLE order: a Go map iterates randomly, and two scans of the same log would assign
	// the same line to different codes whenever two aliases matched.
	sort.Slice(out, func(i, j int) bool { return out[i].code < out[j].code })
	return out, nil
}

type aliasPattern struct {
	code string
	re   *regexp.Regexp
}

// earliest and latest keep the widest window. Dates are compared as text because ISO-8601
// sorts lexicographically — and demanding a parse would refuse a `2026-09` that some
// pipeline emits and that is perfectly usable.
func earliest(a, b string) string {
	if a == "" || (b != "" && b < a) {
		return b
	}
	return a
}

func latest(a, b string) string {
	if b > a {
		return b
	}
	return a
}
