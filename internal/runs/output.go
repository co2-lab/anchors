// @anchors
//   code: RNOTR
//   ref: MNTRS

package runs

import (
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

// Runner is a kind of long process the monitor recognizes: how its command line reads, how
// long it may stay quiet before it counts as stalled, and how its output says it passed,
// failed or advanced. Nothing here is tied to one tool: a project declares its runners in
// `anchors.yaml` (`monitor.runners`).
type Runner struct {
	Name     string
	Kind     string
	Match    *regexp.Regexp
	Stall    time.Duration
	Pass     *regexp.Regexp
	Fail     *regexp.Regexp
	Progress *regexp.Regexp
}

func re(s string) *regexp.Regexp {
	if s == "" {
		return nil
	}
	return regexp.MustCompile(s)
}

// OwnRunners are the runners the monitor knows with no configuration: Anchors' own long
// commands. Every other tool is the project's to declare (`monitor.runners`), the examples in
// config.ExampleRunners.
func OwnRunners() []Runner {
	return []Runner{
		{Name: "anchors", Kind: "anchors", Match: re(`(^|[\s/])anchors (test|mutation|map|check|verify|ingest|docs)(\s|$)`), Stall: 5 * time.Minute},
	}
}

// RunnerFor is the first runner whose pattern a command line matches.
func RunnerFor(runners []Runner, cmd string) (Runner, bool) {
	for _, r := range runners {
		if r.Match != nil && r.Match.MatchString(cmd) {
			return r, true
		}
	}
	return Runner{}, false
}

// tailBytes is how much of the end of an output the monitor reads: the summary a runner
// prints when it ends is there.
const tailBytes = 64 << 10

// OutputInfo is what the monitor reads of a run's output file: its size and its end.
type OutputInfo struct {
	Size int64
	Tail string
	OK   bool
}

// ReadOutput reads the size and the end of an output file.
func ReadOutput(path string) OutputInfo {
	if path == "" {
		return OutputInfo{}
	}
	f, err := os.Open(path)
	if err != nil {
		return OutputInfo{}
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return OutputInfo{}
	}
	start := st.Size() - tailBytes
	if start < 0 {
		start = 0
	}
	if _, err := f.Seek(start, io.SeekStart); err != nil {
		return OutputInfo{}
	}
	b, _ := io.ReadAll(f)
	return OutputInfo{Size: st.Size(), Tail: string(b), OK: true}
}

// Verdict reads from the end of an output whether the run failed or passed: a failure
// anywhere in its summary wins over a pass. It says false when the output says neither — a
// run cut off before its summary.
func Verdict(r Runner, tail string) (failed bool, ok bool) {
	if r.Fail != nil && r.Fail.MatchString(tail) {
		return true, true
	}
	if r.Pass != nil && r.Pass.MatchString(tail) {
		return false, true
	}
	return false, false
}

// LastLines are the last non-blank lines of an output, at most n.
func LastLines(tail string, n int) []string {
	var out []string
	lines := strings.Split(strings.ReplaceAll(tail, "\r\n", "\n"), "\n")
	for i := len(lines) - 1; i >= 0 && len(out) < n; i-- {
		if l := strings.TrimSpace(lines[i]); l != "" {
			out = append([]string{l}, out...)
		}
	}
	return out
}

// SummaryLine is the line of the output that carries its verdict, when one does.
func SummaryLine(r Runner, tail string) string {
	for _, x := range []*regexp.Regexp{r.Fail, r.Pass} {
		if x == nil {
			continue
		}
		if loc := x.FindStringIndex(tail); loc != nil {
			start := strings.LastIndex(tail[:loc[0]], "\n") + 1
			end := strings.Index(tail[loc[0]:], "\n")
			if end < 0 {
				return strings.TrimSpace(tail[start:])
			}
			return strings.TrimSpace(tail[start : loc[0]+end])
		}
	}
	return ""
}
