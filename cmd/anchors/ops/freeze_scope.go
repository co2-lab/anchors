// @anchors
//   code: FSCFR
//   ref: FRZEX

package ops

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// --- the freeze with a scope: a part of the project kept as it is (DESIGN-frozen-units.md) ---
//
// The brake (`freeze` with no scope) stops the whole project. With a scope — a unit, a rule,
// a file — the same verbs keep one part as it is, by decision, until someone thaws it. The
// freeze is recorded once, with its reason, as an entry of a doctrine document coded with
// the freeze letter, and each place it covers carries a mark with the code: `frozen:` in a
// header, `@frozen:` on a rule's line. With no doctrine, the marks carry the reason, and
// nothing ties one to another.

// freezeScope are the places a freeze covers, as the flags name them.
type freezeScope struct {
	units, rules, files []string
	in                  string
}

func (s freezeScope) empty() bool { return len(s.units)+len(s.rules)+len(s.files) == 0 }

// freezeSection is the heading the freeze entries live under in a doctrine document.
const freezeSection = "## Freezes"

// freezeGraph is the project's map as the files are now: a freeze is written and thawed
// against the text, never against a map an edit since has made old.
func freezeGraph(root string) (*mapx.Graph, error) {
	cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", config.DefaultFile, err)
	}
	files, err := scan.Walk(root, cfg)
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return mapx.Build(files, cfg, nil), nil
}

// doctrineDocs are the project's doctrine documents with a code, by path.
func doctrineDocs(g *mapx.Graph) []mapx.Node {
	var out []mapx.Node
	for _, n := range g.Nodes {
		if n.Kind == mapx.KindProduct && n.CodeDeclarado && n.Code != "" {
			out = append(out, n)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// pickDoctrine is the document a new freeze is recorded in: the one `--in` names, or the
// project's only one. None means the freeze is not traceable.
func pickDoctrine(g *mapx.Graph, in string) (mapx.Node, bool, error) {
	docs := doctrineDocs(g)
	if in != "" {
		in = filepath.ToSlash(strings.TrimPrefix(in, "./"))
		for _, d := range docs {
			if d.ID == in || d.Code == in {
				return d, true, nil
			}
		}
		return mapx.Node{}, false, fmt.Errorf("--in %s is not a doctrine document of the project", in)
	}
	switch len(docs) {
	case 0:
		return mapx.Node{}, false, nil
	case 1:
		return docs[0], true, nil
	}
	var names []string
	for _, d := range docs {
		names = append(names, d.ID)
	}
	return mapx.Node{}, false, fmt.Errorf("the project has %d doctrine documents — say which records the freeze with --in: %s", len(docs), strings.Join(names, ", "))
}

// nextFreezeCode is the doctrine document's next free freeze code.
func nextFreezeCode(doc mapx.Node) string {
	max := 0
	for c := range doc.Freezes {
		if i := strings.LastIndex(c, "-"+config.FreezeLetter); i >= 0 {
			if n, err := strconv.Atoi(c[i+2:]); err == nil && n > max {
				max = n
			}
		}
	}
	return fmt.Sprintf("%s-%s%02d", doc.Code, config.FreezeLetter, max+1)
}

func readLines(root, rel string) ([]string, error) {
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(b), "\n"), nil
}

func writeLines(root, rel string, lines []string) error {
	return os.WriteFile(filepath.Join(root, filepath.FromSlash(rel)), []byte(strings.Join(lines, "\n")), 0o644)
}

// addFreezeEntry records the freeze in the doctrine document, under its Freezes section — at
// the document's end when it has none yet.
func addFreezeEntry(root, rel, code, reason string) error {
	lines, err := readLines(root, rel)
	if err != nil {
		return err
	}
	entry := fmt.Sprintf("- **%s** — %s", code, strings.Join(strings.Fields(reason), " "))
	start := -1
	for i, l := range lines {
		if strings.EqualFold(strings.TrimSpace(l), freezeSection) {
			start = i
			break
		}
	}
	if start < 0 {
		for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
			lines = lines[:len(lines)-1]
		}
		lines = append(lines, "", freezeSection, "", entry, "")
		return writeLines(root, rel, lines)
	}
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") || strings.HasPrefix(lines[i], "# ") {
			end = i
			break
		}
	}
	at := end
	for at > start+1 && strings.TrimSpace(lines[at-1]) == "" {
		at--
	}
	if at == start+1 {
		lines = append(lines[:at], append([]string{"", entry}, lines[at:]...)...)
	} else {
		lines = append(lines[:at], append([]string{entry}, lines[at:]...)...)
	}
	return writeLines(root, rel, lines)
}

// removeFreezeEntry deletes the freeze's entry from the doctrine document.
func removeFreezeEntry(root, rel, code string) (bool, error) {
	lines, err := readLines(root, rel)
	if err != nil {
		return false, err
	}
	re := scan.RuleDefinitionRE()
	var out []string
	removed := false
	for _, l := range lines {
		if m := re.FindStringSubmatch(l); m != nil && m[1] == code {
			removed = true
			continue
		}
		out = append(out, l)
	}
	if !removed {
		return false, nil
	}
	return true, writeLines(root, rel, dropEmptyFreezeSection(out))
}

// dropEmptyFreezeSection removes the Freezes section once its last entry is gone — its
// heading and the blank lines around it —, so a freeze and its thaw give the document back
// as it was.
func dropEmptyFreezeSection(lines []string) []string {
	start := -1
	for i, l := range lines {
		if strings.EqualFold(strings.TrimSpace(l), freezeSection) {
			start = i
			break
		}
	}
	if start < 0 {
		return lines
	}
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "#") {
			end = i
			break
		}
		if strings.TrimSpace(lines[i]) != "" {
			return lines
		}
	}
	from := start
	for from > 0 && strings.TrimSpace(lines[from-1]) == "" {
		from--
	}
	if end < len(lines) {
		// A section follows: keep one blank line between what came before and it.
		return append(append(lines[:from:from], ""), lines[end:]...)
	}
	return append(lines[:from:from], "")
}

// headerFieldRE is a field line of the `@anchors` header, capturing its comment prefix.
var headerFieldRE = regexp.MustCompile(`^(\s*(?://|#|--|<!--|\*)?\s*)(code|ref):`)

// headerFrozenLineRE is the header's `frozen:` line, capturing its value.
var headerFrozenLineRE = regexp.MustCompile(`^\s*(?://|#|--|<!--|\*)?\s*frozen:\s*(.*?)\s*(?:-->)?\s*$`)

// headerSpan is the index range of the file's `@anchors` header lines.
func headerSpan(lines []string) (int, int) {
	for i, l := range lines {
		if i > 40 {
			break
		}
		if strings.Contains(l, "@anchors") {
			end := i + 1
			for end < len(lines) && end < i+40 && strings.TrimSpace(lines[end]) != "" &&
				!strings.HasPrefix(strings.TrimSpace(lines[end]), "-->") && headerish(lines[end]) {
				end++
			}
			return i, end
		}
	}
	return -1, -1
}

// headerish says whether a line can belong to a header block: a comment line, or an
// indented field of a markdown header.
func headerish(l string) bool {
	t := strings.TrimSpace(l)
	return strings.HasPrefix(t, "//") || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "--") ||
		strings.HasPrefix(t, "*") || strings.HasPrefix(l, " ") || strings.HasPrefix(l, "\t")
}

// markHeader writes `frozen: <value>` into the file's header, under its code or ref line.
func markHeader(root, rel, value string) (string, error) {
	lines, err := readLines(root, rel)
	if err != nil {
		return "", err
	}
	start, end := headerSpan(lines)
	if start < 0 {
		return "", fmt.Errorf("%s has no @anchors header to carry the freeze", rel)
	}
	field := -1
	for i := start; i < end; i++ {
		if m := headerFrozenLineRE.FindStringSubmatch(lines[i]); m != nil {
			if m[1] == value {
				return "already marked", nil
			}
			return "", fmt.Errorf("%s is already frozen by %q — thaw it first", rel, m[1])
		}
		if field < 0 && headerFieldRE.MatchString(lines[i]) {
			field = i
		}
	}
	if field < 0 {
		return "", fmt.Errorf("%s has no code: or ref: line in its header", rel)
	}
	prefix := headerFieldRE.FindStringSubmatch(lines[field])[1]
	prefix = strings.Replace(prefix, "<!--", "    ", 1)
	lines = append(lines[:field+1], append([]string{prefix + "frozen: " + value}, lines[field+1:]...)...)
	return "marked", writeLines(root, rel, lines)
}

// markRule writes `@frozen: <value>` on the line that defines the rule.
func markRule(root, rel, rule, value string) (string, error) {
	lines, err := readLines(root, rel)
	if err != nil {
		return "", err
	}
	re := scan.RuleDefinitionRE()
	for i, l := range lines {
		m := re.FindStringSubmatch(l)
		if m == nil || m[1] != rule {
			continue
		}
		if strings.Contains(l, "@frozen:") {
			if strings.Contains(l, "@frozen: "+value) {
				return "already marked", nil
			}
			return "", fmt.Errorf("%s is already frozen — thaw it first", rule)
		}
		mark := "@frozen: " + value
		if t := strings.TrimRight(l, " \t"); strings.HasSuffix(t, "|") {
			lines[i] = strings.TrimRight(t[:len(t)-1], " \t") + " " + mark + " |"
		} else {
			lines[i] = t + " " + mark
		}
		return "marked", writeLines(root, rel, lines)
	}
	return "", fmt.Errorf("rule %s is not defined in %s", rule, rel)
}

// ruleMarkRE is a `@frozen:` mark on a line, up to the form's punctuation that closes it.
var ruleMarkRE = regexp.MustCompile("[ \\t]*@frozen:[ \\t]*([^|\\n]*?)[ \\t]*(\\||-->|\\*/|$)")

// unmark removes the freeze marks of a file — every one, or the ones of a value —: its
// header's when header is set, its lines' when body is set (only the rule's line when rule
// is given). A mark cited between backticks is prose, and stays. It says how many.
func unmark(root, rel, value, rule string, header, body bool) (int, error) {
	lines, err := readLines(root, rel)
	if err != nil {
		return 0, err
	}
	start, end := headerSpan(lines)
	re := scan.RuleDefinitionRE()
	removed := 0
	var out []string
	for i, l := range lines {
		if header && i >= start && i < end {
			if m := headerFrozenLineRE.FindStringSubmatch(l); m != nil && (value == "" || m[1] == value) {
				removed++
				continue
			}
		}
		if !body || strings.Contains(l, "`@frozen:") {
			out = append(out, l)
			continue
		}
		if rule != "" {
			if m := re.FindStringSubmatch(l); m == nil || m[1] != rule {
				out = append(out, l)
				continue
			}
		}
		l = ruleMarkRE.ReplaceAllStringFunc(l, func(mark string) string {
			sub := ruleMarkRE.FindStringSubmatch(mark)
			if value != "" && strings.TrimSpace(sub[1]) != value {
				return mark
			}
			removed++
			if sub[2] == "" {
				return ""
			}
			return " " + sub[2]
		})
		out = append(out, l)
	}
	if removed == 0 {
		return 0, nil
	}
	return removed, writeLines(root, rel, out)
}

// unitOwner is the file that owns a unit's code: its spec, or its doctrine document.
func unitOwner(g *mapx.Graph, code string) (string, bool) {
	best := ""
	for _, n := range g.Nodes {
		if !n.CodeDeclarado || n.Code != code {
			continue
		}
		switch {
		case n.Kind == mapx.KindSpec || n.Kind == mapx.KindProduct:
			return n.ID, true
		case best == "":
			best = n.ID
		}
	}
	return best, best != ""
}

// placesOf resolves the scope's flags to files and rules of the map.
type place struct{ path, rule string }

func placesOf(root string, g *mapx.Graph, s freezeScope) ([]place, error) {
	var out []place
	for _, u := range s.units {
		p, ok := unitOwner(g, u)
		if !ok {
			return nil, fmt.Errorf("unit %s has no owner in the project — a spec or a doctrine document that declares it", u)
		}
		out = append(out, place{path: p})
	}
	for _, r := range s.rules {
		p, ok := unitOwner(g, mapx.RuleRoot(r))
		if !ok {
			return nil, fmt.Errorf("rule %s: unit %s has no owner in the project", r, mapx.RuleRoot(r))
		}
		out = append(out, place{path: p, rule: r})
	}
	for _, f := range s.files {
		rel := filepath.ToSlash(f)
		if abs, err := filepath.Abs(f); err == nil {
			if r, err := filepath.Rel(root, abs); err == nil && !strings.HasPrefix(r, "..") {
				rel = filepath.ToSlash(r)
			}
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			return nil, fmt.Errorf("file %s does not exist", f)
		}
		out = append(out, place{path: rel})
	}
	return out, nil
}

// scopedFreeze creates a freeze, or adds a scope to an existing one (code), and marks every
// place of the scope.
func scopedFreeze(w io.Writer, root, code, reason string, s freezeScope) error {
	if s.empty() {
		return errors.New("name what the freeze covers: --unit, --rule or --file")
	}
	g, err := freezeGraph(root)
	if err != nil {
		return err
	}
	places, err := placesOf(root, g, s)
	if err != nil {
		return err
	}
	value := code
	switch {
	case code != "":
		if _, ok := g.FreezeRecords()[code]; !ok {
			return fmt.Errorf("%s is not a freeze any doctrine document records — `anchors freeze show` lists them", code)
		}
	case strings.TrimSpace(reason) == "":
		return errors.New("a freeze needs its reason: --reason \"…\"")
	default:
		doc, traceable, err := pickDoctrine(g, s.in)
		if err != nil {
			return err
		}
		if traceable {
			value = nextFreezeCode(doc)
			if err := addFreezeEntry(root, doc.ID, value, reason); err != nil {
				return err
			}
			fmt.Fprintf(w, "✓ %s recorded in %s: %s\n", value, doc.ID, reason)
		} else {
			value = strings.Join(strings.Fields(reason), " ")
			fmt.Fprintln(w, "· the project has no doctrine document: the marks carry the reason, and the freeze is not traceable")
		}
	}
	for _, p := range places {
		var what string
		if p.rule != "" {
			what, err = markRule(root, p.path, p.rule, value)
		} else {
			what, err = markHeader(root, p.path, value)
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "✓ %s %s%s\n", what, p.path, ruleSuffix(p.rule))
	}
	fmt.Fprintln(w, "the map takes the freeze at the next `anchors map build` (the commit hook runs it)")
	return nil
}

func ruleSuffix(rule string) string {
	if rule == "" {
		return ""
	}
	return " " + rule
}

// thawScoped removes a freeze: by code, every mark of it and its entry; with a scope, the
// marks at those places, whatever they carry.
func thawScoped(w io.Writer, root, code string, s freezeScope) error {
	g, err := freezeGraph(root)
	if err != nil {
		return err
	}
	if code == "" {
		places, err := placesOf(root, g, s)
		if err != nil {
			return err
		}
		for _, p := range places {
			n, err := unmark(root, p.path, "", p.rule, p.rule == "", p.rule != "")
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "✓ %d mark(s) removed from %s%s\n", n, p.path, ruleSuffix(p.rule))
		}
		return nil
	}
	rec, recorded := g.FreezeRecords()[code]
	total := 0
	for _, n := range g.Nodes {
		if n.Frozen != code && !hasValue(n.FrozenRules, code) {
			continue
		}
		k, err := unmark(root, n.ID, code, "", true, true)
		if err != nil {
			return err
		}
		if k > 0 {
			total += k
			fmt.Fprintf(w, "✓ %d mark(s) removed from %s\n", k, n.ID)
		}
	}
	if recorded {
		if _, err := removeFreezeEntry(root, rec.Doc, code); err != nil {
			return err
		}
		fmt.Fprintf(w, "✓ %s removed from %s\n", code, rec.Doc)
	}
	if total == 0 && !recorded {
		return fmt.Errorf("%s is neither recorded nor marked anywhere", code)
	}
	fmt.Fprintf(w, "%s thawed: %d mark(s) removed\n", code, total)
	return nil
}

func hasValue(m map[string]string, v string) bool {
	for _, x := range m {
		if x == v {
			return true
		}
	}
	return false
}

// showFreezes lists the freezes — each with its reason and the places it covers — or, with
// a code, every place that freeze covers, inherited ones included.
func showFreezes(w io.Writer, root, code string) error {
	g, err := freezeGraph(root)
	if err != nil {
		return err
	}
	parts := g.FrozenParts(code)
	records := g.FreezeRecords()
	if code != "" {
		if _, ok := records[code]; !ok && len(parts) == 0 {
			return fmt.Errorf("%s is neither recorded nor marked anywhere", code)
		}
		fmt.Fprintln(w, g.DescribeFreeze(code))
		for _, p := range parts {
			fmt.Fprintf(w, "  %s%s\n", p.Path, ruleSuffix(p.Rule))
		}
		return nil
	}
	if len(parts) == 0 && len(records) == 0 {
		fmt.Fprintln(w, "nothing is frozen")
		return nil
	}
	byValue := map[string][]mapx.FrozenPart{}
	for _, p := range parts {
		byValue[p.Value] = append(byValue[p.Value], p)
	}
	for c := range records {
		if _, ok := byValue[c]; !ok {
			byValue[c] = nil
		}
	}
	values := make([]string, 0, len(byValue))
	for v := range byValue {
		values = append(values, v)
	}
	sort.Strings(values)
	for _, v := range values {
		fmt.Fprintf(w, "❄ %s\n", g.DescribeFreeze(v))
		if len(byValue[v]) == 0 {
			fmt.Fprintln(w, "  (marks nothing)")
		}
		for _, p := range byValue[v] {
			fmt.Fprintf(w, "  %s%s\n", p.Path, ruleSuffix(p.Rule))
		}
	}
	return nil
}

// --- the agent hook: before the agent edits a frozen part, it is told ---

// freezeHookCommand is what the agent's hook runs before an edit.
const freezeHookCommand = "anchors freeze hook"

// freezeHookInput is what the agent's hook receives on stdin, the fields this hook reads.
type freezeHookInput struct {
	Event string `json:"hook_event_name"`
	Cwd   string `json:"cwd"`
	Input struct {
		FilePath     string `json:"file_path"`
		NotebookPath string `json:"notebook_path"`
	} `json:"tool_input"`
}

// freezeHook answers one hook call: for a file a freeze covers — whole, or some of its
// rules —, the context that tells the agent so, with the freeze and its reason; for any
// other, nothing. It never blocks the edit.
func freezeHook(raw []byte) string {
	var in freezeHookInput
	if json.Unmarshal(raw, &in) != nil {
		return ""
	}
	path := in.Input.FilePath
	if path == "" {
		path = in.Input.NotebookPath
	}
	if path == "" {
		return ""
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(in.Cwd, path)
	}
	root := config.ProjectRoot(filepath.Dir(path))
	if _, err := os.Stat(filepath.Join(root, config.DefaultFile)); err != nil {
		return ""
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return ""
	}
	rel = filepath.ToSlash(rel)
	msg := ""
	g, _ := mapx.Load(filepath.Join(root, mapx.DefaultPath))
	describe := func(v string) string { return v }
	var node *mapx.Node
	if g != nil {
		describe = g.DescribeFreeze
		for i := range g.Nodes {
			if g.Nodes[i].ID == rel {
				node = &g.Nodes[i]
				break
			}
		}
	}
	whole := ""
	if b, err := os.ReadFile(path); err == nil {
		whole = scan.FrozenOf(b)
	}
	if whole == "" && node != nil {
		whole = node.Frozen
	}
	switch {
	case whole != "":
		msg = fmt.Sprintf("FROZEN: %s is frozen (%s) — it is not changed without an order from a person. This edit is NOT blocked and goes through: if no person ordered it, undo it and ask. `anchors freeze show` lists the freezes; `anchors thaw` releases one.", rel, describe(whole))
	case node != nil && len(node.FrozenRules) > 0:
		var rules []string
		for r, v := range node.FrozenRules {
			rules = append(rules, r+" ("+describe(v)+")")
		}
		sort.Strings(rules)
		msg = fmt.Sprintf("FROZEN: %s holds frozen rules — they are not changed without an order from a person: %s. This edit is NOT blocked and goes through: if it changes one of them and no person ordered it, undo that part and ask.", rel, strings.Join(rules, "; "))
	default:
		return ""
	}
	event := in.Event
	if event == "" {
		event = "PreToolUse"
	}
	out, _ := json.Marshal(map[string]any{"hookSpecificOutput": map[string]any{
		"hookEventName":     event,
		"additionalContext": msg,
	}})
	return string(out)
}
