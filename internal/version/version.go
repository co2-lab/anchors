// @anchors
//   code: VRSNA
//   ref: VRSNE

// Package version reads, orders and moves a product's version, whatever its scheme, and
// decides the next one by the strategy a branch resolves to.
//
// Nothing here knows a language, a package manager or a store: a version is a string a
// scheme reads, and the next one comes from the releases and the commits git gives. A
// version the scheme cannot read is an error, never "smaller" nor "equal": a guess there
// would let anything through.
package version

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/changelog"
)

// The schemes a product orders its versions by.
const (
	SemVer  = "semver"
	CalVer  = "calver"
	Number  = "number"
	Pattern = "pattern"
)

// Scheme reads and orders the versions of one product.
type Scheme struct {
	Kind   string
	Format string // calver: the fields, `YYYY.0M.MICRO`
	fields []string
	seps   []string
	re     *regexp.Regexp // pattern
}

// V is a version a scheme read: its numbers, left to right, and, for semver, the
// pre-release identifiers.
type V struct {
	Raw  string
	Nums []int
	Pre  []string
}

var semverRE = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?$`)

var calverFields = map[string]bool{"YYYY": true, "YY": true, "0M": true, "MM": true, "0D": true, "DD": true, "MICRO": true}

// NewScheme builds a scheme: `calver` needs its format, `pattern` a regular expression
// whose capture groups are the version's numbers.
func NewScheme(kind, format, pattern string) (Scheme, error) {
	s := Scheme{Kind: kind, Format: format}
	switch kind {
	case SemVer, Number:
	case CalVer:
		if format == "" {
			return s, fmt.Errorf("scheme calver needs a format, such as YYYY.0M.MICRO")
		}
		fs, seps := splitFormat(format)
		for _, f := range fs {
			if !calverFields[f] {
				return s, fmt.Errorf("calver format %q: %q is not a field (YYYY, YY, 0M, MM, 0D, DD, MICRO)", format, f)
			}
		}
		for i, f := range fs {
			if f == "MICRO" && i != len(fs)-1 {
				return s, fmt.Errorf("calver format %q: MICRO is the last field", format)
			}
		}
		s.fields, s.seps = fs, seps
	case Pattern:
		re, err := regexp.Compile("^(?:" + pattern + ")$")
		if err != nil {
			return s, fmt.Errorf("pattern %q is not a regular expression: %v", pattern, err)
		}
		if re.NumSubexp() == 0 {
			return s, fmt.Errorf("pattern %q has no capture group: the groups are the version's numbers", pattern)
		}
		s.re = re
	default:
		return s, fmt.Errorf("scheme %q: write semver, calver, number or pattern", kind)
	}
	return s, nil
}

// splitFormat cuts a calver format into its fields and the separators between them.
func splitFormat(format string) (fields, seps []string) {
	cur := ""
	for _, r := range format {
		if r == '.' || r == '-' || r == '_' {
			fields, seps, cur = append(fields, cur), append(seps, string(r)), ""
			continue
		}
		cur += string(r)
	}
	return append(fields, cur), seps
}

// Parse reads a version, or says why the scheme cannot.
func (s Scheme) Parse(raw string) (V, error) {
	v := V{Raw: raw}
	bad := func() (V, error) { return v, fmt.Errorf("%q is not a %s version", raw, s.describe()) }
	switch s.Kind {
	case SemVer:
		m := semverRE.FindStringSubmatch(raw)
		if m == nil {
			return bad()
		}
		v.Nums = atois(m[1:4])
		if m[4] != "" {
			v.Pre = strings.Split(m[4], ".")
		}
	case Number:
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 || strings.HasPrefix(raw, "+") {
			return bad()
		}
		v.Nums = []int{n}
	case CalVer:
		rest := raw
		for i := range s.fields {
			part := rest
			if i < len(s.seps) {
				var ok bool
				part, rest, ok = strings.Cut(rest, s.seps[i])
				if !ok {
					return bad()
				}
			}
			n, err := strconv.Atoi(part)
			if err != nil || n < 0 || !calverWidthOK(s.fields[i], part) {
				return bad()
			}
			v.Nums = append(v.Nums, n)
		}
	case Pattern:
		m := s.re.FindStringSubmatch(raw)
		if m == nil {
			return bad()
		}
		for _, g := range m[1:] {
			n, err := strconv.Atoi(g)
			if err != nil {
				return v, fmt.Errorf("%q: the group %q of the pattern is not a number", raw, g)
			}
			v.Nums = append(v.Nums, n)
		}
	}
	return v, nil
}

func (s Scheme) describe() string {
	switch s.Kind {
	case CalVer:
		return "calver " + s.Format
	case Pattern:
		return "pattern " + strings.TrimSuffix(strings.TrimPrefix(s.re.String(), "^(?:"), ")$")
	}
	return s.Kind
}

// calverWidthOK refuses a field written in a width its format does not allow: `0M` is
// two digits, `YYYY` four.
func calverWidthOK(field, part string) bool {
	switch field {
	case "YYYY":
		return len(part) == 4
	case "YY", "0M", "0D":
		return len(part) == 2
	case "MM", "DD":
		return len(part) <= 2 && (len(part) == 1 || part[0] != '0')
	}
	return true
}

func atois(ss []string) []int {
	out := make([]int, len(ss))
	for i, s := range ss {
		out[i], _ = strconv.Atoi(s)
	}
	return out
}

// Compare orders two versions of the scheme: -1, 0 or 1. A semver pre-release ranks below
// its release, and its identifiers are ordered as SemVer 2.0 says.
func Compare(a, b V) int {
	for i := 0; i < len(a.Nums) && i < len(b.Nums); i++ {
		if a.Nums[i] != b.Nums[i] {
			return sign(a.Nums[i] - b.Nums[i])
		}
	}
	if len(a.Nums) != len(b.Nums) {
		return sign(len(a.Nums) - len(b.Nums))
	}
	switch {
	case len(a.Pre) == 0 && len(b.Pre) == 0:
		return 0
	case len(a.Pre) == 0:
		return 1
	case len(b.Pre) == 0:
		return -1
	}
	for i := 0; i < len(a.Pre) && i < len(b.Pre); i++ {
		if c := comparePre(a.Pre[i], b.Pre[i]); c != 0 {
			return c
		}
	}
	return sign(len(a.Pre) - len(b.Pre))
}

func comparePre(a, b string) int {
	na, ea := strconv.Atoi(a)
	nb, eb := strconv.Atoi(b)
	switch {
	case ea == nil && eb == nil:
		return sign(na - nb)
	case ea == nil:
		return -1
	case eb == nil:
		return 1
	}
	return strings.Compare(a, b)
}

func sign(n int) int {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	}
	return 0
}

// Parts are the parts `increment` may raise in this scheme. An empty part names the only
// one, where the scheme has a single choice.
func (s Scheme) Parts() []string {
	switch s.Kind {
	case SemVer:
		return []string{"major", "minor", "patch"}
	case Pattern:
		var out []string
		for i, n := range s.re.SubexpNames()[1:] {
			out = append(out, strconv.Itoa(i+1))
			if n != "" {
				out = append(out, n)
			}
		}
		return out
	}
	return []string{""}
}

// ValidPart says whether `increment` may raise the part in this scheme.
func (s Scheme) ValidPart(part string) bool {
	if s.Kind == Number || s.Kind == CalVer {
		return part == "" || part == s.Kind || (s.Kind == CalVer && part == "micro")
	}
	return slices.Contains(s.Parts(), part)
}

// Increment raises the part and resets the parts after it. In calver the date fields come
// from today, and the counter goes up only while the date is the same.
func (s Scheme) Increment(v V, part string, today time.Time) (string, error) {
	if !s.ValidPart(part) {
		return "", fmt.Errorf("%s has no part %q to increment (%s)", s.describe(), part, strings.Join(s.Parts(), ", "))
	}
	switch s.Kind {
	case SemVer:
		n := append([]int(nil), v.Nums...)
		i := map[string]int{"major": 0, "minor": 1, "patch": 2}[part]
		// A pre-release already stands for the release it precedes: 1.1.0-rc.1 raised
		// in its minor is 1.1.0, and in its patch 1.1.0 too.
		precedes := len(v.Pre) > 0
		for j := i + 1; j < 3; j++ {
			precedes = precedes && n[j] == 0
		}
		if !precedes {
			n[i]++
		}
		for j := i + 1; j < 3; j++ {
			n[j] = 0
		}
		return fmt.Sprintf("%d.%d.%d", n[0], n[1], n[2]), nil
	case Number:
		return strconv.Itoa(v.Nums[0] + 1), nil
	case CalVer:
		var b strings.Builder
		same := true
		for i, f := range s.fields {
			if i > 0 {
				b.WriteString(s.seps[i-1])
			}
			if f == "MICRO" {
				n := 0
				if same {
					n = v.Nums[i] + 1
				}
				b.WriteString(strconv.Itoa(n))
				continue
			}
			d := dateField(f, today)
			if n, _ := strconv.Atoi(d); n != v.Nums[i] {
				same = false
			}
			b.WriteString(d)
		}
		out := b.String()
		if same && s.fields[len(s.fields)-1] != "MICRO" {
			return "", fmt.Errorf("calver %s has no counter, and today is the date of %s", s.Format, v.Raw)
		}
		return out, nil
	}
	// pattern: the group named by its index or its name goes up, the later ones reset
	idx := 0
	names := s.re.SubexpNames()
	for i := 1; i < len(names); i++ {
		if part == strconv.Itoa(i) || part == names[i] {
			idx = i
			break
		}
	}
	loc := s.re.FindStringSubmatchIndex(v.Raw)
	var b strings.Builder
	last := 0
	for g := 1; g < len(names); g++ {
		st, en := loc[2*g], loc[2*g+1]
		if st < 0 {
			continue
		}
		b.WriteString(v.Raw[last:st])
		n := v.Nums[g-1]
		switch {
		case g == idx:
			n++
		case g > idx:
			n = 0
		}
		b.WriteString(strconv.Itoa(n))
		last = en
	}
	b.WriteString(v.Raw[last:])
	return b.String(), nil
}

func dateField(f string, t time.Time) string {
	switch f {
	case "YYYY":
		return t.Format("2006")
	case "YY":
		return t.Format("06")
	case "0M":
		return t.Format("01")
	case "MM":
		return strconv.Itoa(int(t.Month()))
	case "0D":
		return t.Format("02")
	case "DD":
		return strconv.Itoa(t.Day())
	}
	return ""
}

// The strategies a branch decides its next version by.
const (
	Manual    = "manual"
	Increment = "increment"
	Semantic  = "semantic"
	None      = "none"
)

// Strategy is how a branch decides its next version.
type Strategy struct {
	Name    string
	Part    string            // increment
	Initial bool              // semantic: below 1.0.0, a breaking change bumps the minor
	Types   map[string]string // semantic: commit type → major, minor, patch or none
}

// Input is what the next version is decided from.
type Input struct {
	Scheme   Scheme
	Strategy Strategy
	Current  string             // the version the sources hold; empty with no source
	Released string             // the latest released version; empty before the first release
	Commits  []changelog.Commit // since the latest release, oldest first
	Today    time.Time
}

// Result is the next version and why.
type Result struct {
	Version string // empty: the branch never versions
	Reason  string
	Forced  bool // a `Version:` footer decided it
}

// Levels a commit type may bump.
var levels = map[string]int{"none": -1, "patch": 2, "minor": 1, "major": 0}

// ValidLevel says whether a `types:` value names a level.
func ValidLevel(l string) bool { _, ok := levels[l]; return ok }

// Next decides the next version: a `Version:` footer forces it, else the strategy computes
// it. The result is read by the scheme and, once a release exists, must be above it.
func Next(in Input) (Result, error) {
	if in.Strategy.Name == None {
		return Result{Reason: "the branch never versions (next: none)"}, nil
	}
	var released *V
	if in.Released != "" {
		r, err := in.Scheme.Parse(in.Released)
		if err != nil {
			return Result{}, fmt.Errorf("the latest release: %w", err)
		}
		released = &r
	}
	res, err := compute(in, released)
	if err != nil {
		return res, err
	}
	if forced, hash := footerVersion(in.Commits); forced != "" {
		res = Result{Version: forced, Forced: true, Reason: "forced by the Version: footer of " + hash}
	}
	v, err := in.Scheme.Parse(res.Version)
	if err != nil {
		return res, fmt.Errorf("%s: %w", res.Reason, err)
	}
	if released != nil && Compare(v, *released) < 0 {
		return res, fmt.Errorf("%s: %s is below the latest release %s", res.Reason, res.Version, in.Released)
	}
	return res, nil
}

func compute(in Input, released *V) (Result, error) {
	switch in.Strategy.Name {
	case Manual, "":
		if in.Current == "" {
			return Result{}, fmt.Errorf("strategy manual reads the version from the sources, and the product declares none")
		}
		return Result{Version: in.Current, Reason: "manual: the version the sources hold"}, nil
	case Increment:
		if released == nil {
			return firstRelease(in, "increment")
		}
		if len(in.Commits) == 0 {
			return Result{Version: in.Released, Reason: "no commit since " + in.Released}, nil
		}
		next, err := in.Scheme.Increment(*released, in.Strategy.Part, in.Today)
		if err != nil {
			return Result{}, err
		}
		part := in.Strategy.Part
		if part == "" {
			part = in.Scheme.Kind
		}
		return Result{Version: next, Reason: fmt.Sprintf("increment %s of %s", part, in.Released)}, nil
	case Semantic:
		if in.Scheme.Kind != SemVer {
			return Result{}, fmt.Errorf("strategy semantic needs scheme semver, not %s", in.Scheme.Kind)
		}
		if released == nil {
			return firstRelease(in, "semantic")
		}
		level, why := semanticLevel(in.Commits, in.Strategy)
		if level < 0 {
			return Result{Version: in.Released, Reason: "semantic: nothing to release since " + in.Released}, nil
		}
		if in.Strategy.Initial && level == 0 && released.Nums[0] == 0 {
			level = 1
		}
		next, err := in.Scheme.Increment(*released, []string{"major", "minor", "patch"}[level], in.Today)
		if err != nil {
			return Result{}, err
		}
		return Result{Version: next, Reason: fmt.Sprintf("semantic: %s since %s", why, in.Released)}, nil
	}
	return Result{}, fmt.Errorf("strategy %q: write manual, increment, semantic or none", in.Strategy.Name)
}

// firstRelease is the version before any release: what the sources hold, or, with no
// source, the scheme's first version.
func firstRelease(in Input, name string) (Result, error) {
	if in.Current != "" {
		return Result{Version: in.Current, Reason: name + ": no release yet, the version the sources hold"}, nil
	}
	first := map[string]string{SemVer: "1.0.0", Number: "1"}[in.Scheme.Kind]
	if in.Scheme.Kind == SemVer && in.Strategy.Initial {
		first = "0.1.0"
	}
	if first == "" {
		return Result{}, fmt.Errorf("%s: no release and no source, so there is no version to start from", name)
	}
	return Result{Version: first, Reason: name + ": the first release"}, nil
}

// semanticLevel is the largest bump the commits ask for — 0 major, 1 minor, 2 patch, -1
// none — and a count of what asked for it.
func semanticLevel(commits []changelog.Commit, st Strategy) (int, string) {
	r := changelog.Classify("", "", commits)
	level := -1
	raise := func(l int) {
		if l >= 0 && (level < 0 || l < level) {
			level = l
		}
	}
	counts := map[string]int{}
	if len(r.Breaking) > 0 {
		raise(0)
		counts["breaking"] = len(r.Breaking)
	}
	typed := map[string]int{"feat": 1, "fix": 2}
	for t, l := range st.Types {
		typed[t] = levels[l]
	}
	for _, c := range commits {
		m := headerType(c.Subject)
		if m == "" {
			continue
		}
		if l, ok := typed[m]; ok && l >= 0 {
			raise(l)
			counts[m]++
		}
	}
	var parts []string
	for _, k := range []string{"breaking", "feat", "fix"} {
		if counts[k] > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", counts[k], k))
			delete(counts, k)
		}
	}
	var rest []string
	for k := range counts {
		rest = append(rest, k)
	}
	sort.Strings(rest)
	for _, k := range rest {
		parts = append(parts, fmt.Sprintf("%d %s", counts[k], k))
	}
	return level, strings.Join(parts, ", ")
}

var typeRE = regexp.MustCompile(`^([a-z]+)(\([^)]*\))?(!)?: +.+$`)

func headerType(subject string) string {
	t, _ := CommitType(subject)
	return t
}

// CommitType is a Conventional Commits subject's type, and whether a `!` marks it breaking;
// empty for a subject outside the convention.
func CommitType(subject string) (string, bool) {
	m := typeRE.FindStringSubmatch(strings.TrimSpace(subject))
	if m == nil {
		return "", false
	}
	return m[1], m[3] == "!"
}

// footerVersion is the version the newest `Version:` footer forces, and its commit.
func footerVersion(commits []changelog.Commit) (string, string) {
	for i := len(commits) - 1; i >= 0; i-- {
		for _, l := range changelog.FooterLines(commits[i].Body) {
			k, v, _ := strings.Cut(l, ":")
			if strings.TrimSpace(k) == "Version" && strings.TrimSpace(v) != "" {
				return strings.TrimSpace(v), commits[i].Hash
			}
		}
	}
	return "", ""
}
