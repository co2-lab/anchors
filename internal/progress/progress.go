// @anchors
//   code: PFCPR
//   ref: DMLDD

// Package progress keeps the project's progress file, anchors.progress.md: its demands, their
// optional stages and their weighted tasks, in a delimited format that only Anchors writes
// (DESIGN-progress-tracking.md).
package progress

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// FileName is the progress file, at the project's root — the name the scan keeps out of the
// map (scan.ProgressFile).
const FileName = scan.ProgressFile

// DefaultUnit names the weight when the project names none.
const DefaultUnit = "points"

// headerLine opens the file and says who writes it.
const headerLine = "<!-- anchors:progress v1 — written by `anchors progress`; never edit by hand -->"

// State is where a task stands, written as the mark inside its box.
type State string

const (
	Todo      State = " "
	Doing     State = "~"
	Done      State = "x"
	Cancelled State = "-"
)

// Task is one piece of work, with its weight — an estimate of its effort.
type Task struct {
	ID       string  `json:"id"`
	num      int     // its number in the demand
	Title    string  `json:"title"`
	Weight   float64 `json:"weight"`
	Estimate bool    `json:"estimate,omitempty"` // a weight nobody estimated yet (a migrated task)
	// Unplanned is work the plan did not list, done for it: kept where it arose, marked so it
	// is seen (DESIGN-unplanned-work.md).
	Unplanned bool   `json:"unplanned,omitempty"`
	State     State  `json:"state"`
	At        string `json:"at,omitempty"`     // when it was started or done
	By        string `json:"by,omitempty"`     // who started or did it
	Reason    string `json:"reason,omitempty"` // why it was cancelled
}

// Stage is an optional step of a demand, holding some of its tasks.
type Stage struct {
	ID    string  `json:"id"`
	num   int     // its number in the demand
	Title string  `json:"title"`
	Tasks []*Task `json:"tasks"`
}

// Demand is a piece of work the project tracks — a plan in development, or any other.
type Demand struct {
	ID     string   `json:"id"`
	num    int      // its number in the file
	Title  string   `json:"title"`
	Plan   string   `json:"plan,omitempty"`
	Tasks  []*Task  `json:"tasks"` // the tasks of no stage
	Stages []*Stage `json:"stages"`
}

// File is the whole progress file.
type File struct {
	Unit    string    `json:"unit"`
	Demands []*Demand `json:"demands"`
}

// ErrEdited is the error of a progress file whose content does not match its sum: it was
// changed by hand, and Anchors does not write over what it cannot read as its own.
var ErrEdited = errors.New(FileName + " was edited by hand (its sum does not match its content) — " +
	"restore it from git, and change it only through `anchors progress`")

var (
	sumRE    = regexp.MustCompile(`^<!-- anchors:progress sum=([0-9a-f]{12}) -->$`)
	totalRE  = regexp.MustCompile(`^# Progress — (?:\d+%|—) \([0-9.]+/[0-9.]+ (.+)\)$`)
	demandRE = regexp.MustCompile(`^## D(\d+) — (.+?)(?: \(plan: ([^)]+)\))? — (?:\d+%|—) \([0-9.]+/[0-9.]+\)$`)
	stageRE  = regexp.MustCompile(`^### D(\d+)\.S(\d+) — (.+) — (?:\d+%|—) \([0-9.]+/[0-9.]+\)$`)
	taskRE   = regexp.MustCompile(`^- \[([ ~x-])\] D(\d+)\.T(\d+) — (.+)$`)
	doingRE  = regexp.MustCompile(`^doing since (\d{4}-\d{2}-\d{2}(?: \d{2}:\d{2})?) by (\S+)$`)
	doneRE   = regexp.MustCompile(`^done (\S+) by (\S+)$`)
)

// sumOf is the sum the file carries of its content.
func sumOf(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])[:12]
}

// Parse reads a progress file. Its content must match its sum (ErrEdited otherwise), and
// every line must be one Anchors writes.
func Parse(text string) (*File, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	body := strings.TrimRight(text, "\n")
	cut := strings.LastIndex(body, "\n")
	if cut < 0 {
		return nil, fmt.Errorf("%s: no sum line", FileName)
	}
	m := sumRE.FindStringSubmatch(body[cut+1:])
	if m == nil {
		return nil, ErrEdited
	}
	content := body[:cut+1]
	if sumOf(content) != m[1] {
		return nil, ErrEdited
	}
	f := &File{Unit: DefaultUnit}
	var d *Demand
	var st *Stage
	for i, line := range strings.Split(strings.TrimRight(content, "\n"), "\n") {
		n := i + 1
		switch {
		case line == "" || (n == 1 && line == headerLine):
		case totalRE.MatchString(line):
			f.Unit = totalRE.FindStringSubmatch(line)[1]
		case demandRE.MatchString(line):
			m := demandRE.FindStringSubmatch(line)
			num, _ := strconv.Atoi(m[1])
			d = &Demand{ID: demandID(num), num: num, Title: m[2], Plan: m[3]}
			st = nil
			f.Demands = append(f.Demands, d)
		case stageRE.MatchString(line):
			m := stageRE.FindStringSubmatch(line)
			dn, _ := strconv.Atoi(m[1])
			sn, _ := strconv.Atoi(m[2])
			if d == nil || d.num != dn {
				return nil, fmt.Errorf("%s:%d: stage %s outside its demand", FileName, n, m[0])
			}
			st = &Stage{ID: stageID(dn, sn), num: sn, Title: m[3]}
			d.Stages = append(d.Stages, st)
		case taskRE.MatchString(line):
			t, dn, err := parseTask(line)
			if err != nil {
				return nil, fmt.Errorf("%s:%d: %w", FileName, n, err)
			}
			if d == nil || d.num != dn {
				return nil, fmt.Errorf("%s:%d: task %s outside its demand", FileName, n, t.ID)
			}
			if st != nil {
				st.Tasks = append(st.Tasks, t)
			} else {
				d.Tasks = append(d.Tasks, t)
			}
		default:
			return nil, fmt.Errorf("%s:%d: a line Anchors does not write: %q", FileName, n, line)
		}
	}
	return f, nil
}

func parseTask(line string) (*Task, int, error) {
	m := taskRE.FindStringSubmatch(line)
	dn, _ := strconv.Atoi(m[2])
	tn, _ := strconv.Atoi(m[3])
	parts := strings.Split(m[4], " · ")
	if len(parts) < 2 {
		return nil, 0, fmt.Errorf("task D%02d.T%d has no weight", dn, tn)
	}
	t := &Task{ID: taskID(dn, tn), num: tn, Title: parts[0], State: State(m[1])}
	w := parts[1]
	// The weight's marks, in parentheses: `2 (estimate)`, `2 (unplanned)`, `2 (estimate, unplanned)`.
	if i := strings.Index(w, " ("); i >= 0 && strings.HasSuffix(w, ")") {
		for _, mark := range strings.Split(w[i+2:len(w)-1], ", ") {
			switch mark {
			case "estimate":
				t.Estimate = true
			case "unplanned":
				t.Unplanned = true
			default:
				return nil, 0, fmt.Errorf("task %s: the weight's mark %q is neither estimate nor unplanned", t.ID, mark)
			}
		}
		w = w[:i]
	}
	weight, err := strconv.ParseFloat(w, 64)
	if err != nil || weight <= 0 {
		return nil, 0, fmt.Errorf("task %s: weight %q is not a positive number", t.ID, w)
	}
	t.Weight = weight
	if len(parts) > 2 {
		rest := strings.Join(parts[2:], " · ")
		switch {
		case t.State == Doing && doingRE.MatchString(rest):
			mm := doingRE.FindStringSubmatch(rest)
			t.At, t.By = mm[1], mm[2]
		case t.State == Done && doneRE.MatchString(rest):
			mm := doneRE.FindStringSubmatch(rest)
			t.At, t.By = mm[1], mm[2]
		case t.State == Cancelled && strings.HasPrefix(rest, "cancelled: "):
			t.Reason = strings.TrimPrefix(rest, "cancelled: ")
		default:
			return nil, 0, fmt.Errorf("task %s: %q does not match its state", t.ID, rest)
		}
	}
	return t, dn, nil
}

func demandID(d int) string       { return fmt.Sprintf("D%02d", d) }
func stageID(d, s int) string     { return fmt.Sprintf("D%02d.S%d", d, s) }
func taskID(d, t int) string      { return fmt.Sprintf("D%02d.T%d", d, t) }
func weightText(w float64) string { return strconv.FormatFloat(w, 'f', -1, 64) }

// Totals are the done weight and the weight that counts — every task not cancelled.
func Totals(tasks []*Task) (done, total float64) {
	for _, t := range tasks {
		if t.State == Cancelled {
			continue
		}
		total += t.Weight
		if t.State == Done {
			done += t.Weight
		}
	}
	return done, total
}

// AllTasks are a demand's tasks, its stages' included.
func (d *Demand) AllTasks() []*Task {
	out := append([]*Task{}, d.Tasks...)
	for _, s := range d.Stages {
		out = append(out, s.Tasks...)
	}
	return out
}

// AllTasks are every task of the file.
func (f *File) AllTasks() []*Task {
	var out []*Task
	for _, d := range f.Demands {
		out = append(out, d.AllTasks()...)
	}
	return out
}

// Percent is the done share of the total, rounded down; "—" when nothing counts.
func Percent(done, total float64) string {
	if total == 0 {
		return "—"
	}
	return fmt.Sprintf("%d%%", int(math.Floor(done*100/total+1e-9)))
}

func ratio(done, total float64) string {
	return Percent(done, total) + " (" + weightText(done) + "/" + weightText(total)
}

// Render writes the file: the percentages computed from the weights, and the sum of the
// content on the last line.
func Render(f *File) string {
	unit := f.Unit
	if unit == "" {
		unit = DefaultUnit
	}
	var b strings.Builder
	b.WriteString(headerLine + "\n")
	done, total := Totals(f.AllTasks())
	b.WriteString("# Progress — " + ratio(done, total) + " " + unit + ")\n")
	for _, d := range f.Demands {
		b.WriteString("\n")
		dd, dt := Totals(d.AllTasks())
		plan := ""
		if d.Plan != "" {
			plan = " (plan: " + d.Plan + ")"
		}
		b.WriteString("## " + d.ID + " — " + d.Title + plan + " — " + ratio(dd, dt) + ")\n")
		for _, t := range d.Tasks {
			b.WriteString(taskLine(t) + "\n")
		}
		for _, s := range d.Stages {
			sd, stt := Totals(s.Tasks)
			b.WriteString("### " + s.ID + " — " + s.Title + " — " + ratio(sd, stt) + ")\n")
			for _, t := range s.Tasks {
				b.WriteString(taskLine(t) + "\n")
			}
		}
	}
	b.WriteString("\n")
	content := b.String()
	return content + "<!-- anchors:progress sum=" + sumOf(content) + " -->\n"
}

func taskLine(t *Task) string {
	line := "- [" + string(t.State) + "] " + t.ID + " — " + t.Title + " · " + weightText(t.Weight)
	var marks []string
	if t.Estimate {
		marks = append(marks, "estimate")
	}
	if t.Unplanned {
		marks = append(marks, "unplanned")
	}
	if len(marks) > 0 {
		line += " (" + strings.Join(marks, ", ") + ")"
	}
	switch t.State {
	case Doing:
		line += " · doing since " + t.At + " by " + t.By
	case Done:
		line += " · done " + t.At + " by " + t.By
	case Cancelled:
		line += " · cancelled: " + t.Reason
	}
	return line
}

// --- reading and writing it, one writer at a time ---

// Path is the progress file of a project.
func Path(root string) string { return filepath.Join(root, FileName) }

// Load reads the project's progress file; an absent one is an empty file.
func Load(root string) (*File, error) {
	b, err := os.ReadFile(Path(root))
	if errors.Is(err, os.ErrNotExist) {
		return &File{Unit: DefaultUnit}, nil
	}
	if err != nil {
		return nil, err
	}
	return Parse(string(b))
}

// Update changes the progress file under its lock: it reads the file as it is on disk now,
// applies the change, and writes it whole and atomically — so two agents changing it at
// once lose nothing of each other's. The unit is the project's.
func Update(root, unit string, change func(*File) error) error {
	if err := os.MkdirAll(filepath.Join(root, ".anchors"), 0o755); err != nil {
		return err
	}
	unlock, err := mapx.LockNamed(filepath.Join(root, ".anchors", "progress"), "the progress file")
	if err != nil {
		return err
	}
	defer unlock()
	f, err := Load(root)
	if err != nil {
		return err
	}
	if err := change(f); err != nil {
		return err
	}
	if unit != "" {
		f.Unit = unit
	}
	tmp := fmt.Sprintf("%s.%d.tmp", Path(root), os.Getpid())
	if err := os.WriteFile(tmp, []byte(Render(f)), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, Path(root))
}

// --- the changes ---

// checkText refuses a title or a reason the format could not read back.
func checkText(what, s string) error {
	switch {
	case strings.TrimSpace(s) == "":
		return fmt.Errorf("the %s is empty", what)
	case strings.ContainsAny(s, "\n\r"):
		return fmt.Errorf("the %s spans lines", what)
	case strings.Contains(s, " · "):
		return fmt.Errorf("the %s holds ` · `, the format's separator", what)
	}
	return nil
}

// AddDemand adds a demand, optionally tied to a plan, and gives its id.
func (f *File) AddDemand(title, plan string) (string, error) {
	if err := checkText("title", title); err != nil {
		return "", err
	}
	if strings.ContainsAny(plan, "()") {
		return "", fmt.Errorf("the plan path holds a parenthesis")
	}
	n := 0
	for _, d := range f.Demands {
		n = max(n, d.num)
	}
	d := &Demand{ID: demandID(n + 1), num: n + 1, Title: strings.TrimSpace(title), Plan: plan}
	f.Demands = append(f.Demands, d)
	return d.ID, nil
}

// Demand is the demand of that id, or nil.
func (f *File) Demand(id string) *Demand {
	for _, d := range f.Demands {
		if d.ID == id {
			return d
		}
	}
	return nil
}

// AddStage adds a stage to a demand, and gives its id.
func (f *File) AddStage(demand, title string) (string, error) {
	d := f.Demand(demand)
	if d == nil {
		return "", fmt.Errorf("no demand %s", demand)
	}
	if err := checkText("title", title); err != nil {
		return "", err
	}
	n := 0
	for _, s := range d.Stages {
		n = max(n, s.num)
	}
	s := &Stage{ID: stageID(d.num, n+1), num: n + 1, Title: strings.TrimSpace(title)}
	d.Stages = append(d.Stages, s)
	return s.ID, nil
}

// AddTask adds a task to a demand or to one of its stages, with its weight, and gives its id.
func (f *File) AddTask(parent, title string, weight float64, estimate bool) (string, error) {
	if err := checkText("title", title); err != nil {
		return "", err
	}
	if weight <= 0 {
		return "", fmt.Errorf("a task needs a positive weight, even a rough estimate")
	}
	demand, stage, _ := strings.Cut(parent, ".")
	d := f.Demand(demand)
	if d == nil {
		return "", fmt.Errorf("no demand %s", demand)
	}
	n := 0
	for _, t := range d.AllTasks() {
		n = max(n, t.num)
	}
	t := &Task{ID: taskID(d.num, n+1), num: n + 1, Title: strings.TrimSpace(title), Weight: weight, Estimate: estimate, State: Todo}
	if stage == "" {
		d.Tasks = append(d.Tasks, t)
		return t.ID, nil
	}
	for _, s := range d.Stages {
		if s.ID == parent {
			s.Tasks = append(s.Tasks, t)
			return t.ID, nil
		}
	}
	return "", fmt.Errorf("no stage %s", parent)
}

// Task is the task of that id, or nil.
func (f *File) Task(id string) *Task {
	for _, t := range f.AllTasks() {
		if t.ID == id {
			return t
		}
	}
	return nil
}

func (f *File) task(id string) (*Task, error) {
	if t := f.Task(id); t != nil {
		return t, nil
	}
	return nil, fmt.Errorf("no task %s", id)
}

// StartedAt is the form of the moment a task is started: the date and the minute, so the
// monitor can tell how long a held task has gone without activity (DESIGN-agent-activity.md).
const StartedAt = "2006-01-02 15:04"

// Since is when a task in progress was started, in local time: a file written before the
// minute was recorded holds only the date, which counts from the start of that day.
func (t *Task) Since() (time.Time, bool) {
	if at, err := time.ParseInLocation(StartedAt, t.At, time.Local); err == nil {
		return at, true
	}
	at, err := time.ParseInLocation("2006-01-02", t.At, time.Local)
	return at, err == nil
}

// Start marks a task in progress, held by who. A task another agent holds is refused
// unless take.
func (f *File) Start(id, who, date string, take bool) error {
	t, err := f.task(id)
	if err != nil {
		return err
	}
	switch {
	case t.State == Done || t.State == Cancelled:
		return fmt.Errorf("%s is %s — reopen it first", id, stateName(t.State))
	case t.State == Doing && t.By != who && !take:
		return fmt.Errorf("%s is held by %s since %s — take it over with --take", id, t.By, t.At)
	}
	t.State, t.By, t.At, t.Reason = Doing, who, date, ""
	return nil
}

// Done marks a task done by who.
func (f *File) Done(id, who, date string) error {
	t, err := f.task(id)
	if err != nil {
		return err
	}
	if t.State == Cancelled {
		return fmt.Errorf("%s is cancelled — reopen it first", id)
	}
	t.State, t.By, t.At, t.Reason = Done, who, date, ""
	return nil
}

// Reopen puts a task back to do.
func (f *File) Reopen(id string) error {
	t, err := f.task(id)
	if err != nil {
		return err
	}
	t.State, t.By, t.At, t.Reason = Todo, "", "", ""
	return nil
}

// Cancel takes a task out of the count, with the reason. Done work is never deleted.
func (f *File) Cancel(id, reason string) error {
	t, err := f.task(id)
	if err != nil {
		return err
	}
	if err := checkText("reason", reason); err != nil {
		return err
	}
	t.State, t.By, t.At, t.Reason = Cancelled, "", "", strings.TrimSpace(reason)
	return nil
}

// Weigh sets a task's weight; it is an estimate no longer.
func (f *File) Weigh(id string, weight float64) error {
	t, err := f.task(id)
	if err != nil {
		return err
	}
	if weight <= 0 {
		return fmt.Errorf("the weight must be positive")
	}
	t.Weight, t.Estimate = weight, false
	return nil
}

// UnplannedAdvice is what the flow tells an agent about work the plan does not list, in the
// plan guide, the work prompt and the generated CONTRIBUTING.md and AGENTS.md: recorded
// where it arises, marked, before it is done — unrecorded, it never enters the progress
// (DESIGN-unplanned-work.md).
const UnplannedAdvice = "Work the plan does not list, found on the way and done for it — a fix, a missing helper, " +
	"a test the plan forgot — is recorded before it is done: " +
	"`anchors progress add task <demand|stage> \"<what>\" --weight <n> --unplanned`. " +
	"The plan is not rewritten for it; the progress counts it, marked so it is seen."

// MarkUnplanned marks a task as work the plan did not list, or, with on false, takes the
// mark off (DESIGN-unplanned-work.md).
func (f *File) MarkUnplanned(id string, on bool) error {
	t, err := f.task(id)
	if err != nil {
		return err
	}
	t.Unplanned = on
	return nil
}

// Rename gives a demand, a stage or a task a new title; its id, its tasks and their history
// stay. Without it, the only way to change a name was to cancel the tasks and add them again
// (reported from MIF).
func (f *File) Rename(id, title string) error {
	if err := checkText("title", title); err != nil {
		return err
	}
	title = strings.TrimSpace(title)
	if d := f.Demand(id); d != nil {
		d.Title = title
		return nil
	}
	for _, d := range f.Demands {
		for _, s := range d.Stages {
			if s.ID == id {
				s.Title = title
				return nil
			}
		}
	}
	if t := f.Task(id); t != nil {
		t.Title = title
		return nil
	}
	return fmt.Errorf("no demand, stage or task %s", id)
}

func stateName(s State) string {
	switch s {
	case Doing:
		return "in progress"
	case Done:
		return "done"
	case Cancelled:
		return "cancelled"
	}
	return "to do"
}

// --- the summary ---

// Row is one line of the summary: a demand, or one of its stages.
type Row struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Plan     string  `json:"plan,omitempty"`
	Stage    bool    `json:"stage,omitempty"`
	Done     float64 `json:"done"`
	Total    float64 `json:"total"`
	Percent  string  `json:"percent"`
	Open     int     `json:"open"`               // tasks to do or in progress
	Doing    int     `json:"doing"`              // tasks in progress
	Estimate int     `json:"estimate,omitempty"` // tasks whose weight is still an estimate
	// Unplanned is the weight of the work the plan did not list, cancelled tasks left out.
	Unplanned float64 `json:"unplanned,omitempty"`
}

// Summary is the short view of what is in progress: one row per demand, and one per stage
// when the demand has stages — never the tasks.
func Summary(f *File) []Row {
	var out []Row
	row := func(id, title, plan string, stage bool, tasks []*Task) Row {
		done, total := Totals(tasks)
		r := Row{ID: id, Title: title, Plan: plan, Stage: stage, Done: done, Total: total, Percent: Percent(done, total)}
		for _, t := range tasks {
			switch t.State {
			case Todo:
				r.Open++
			case Doing:
				r.Open++
				r.Doing++
			}
			if t.Estimate && t.State != Cancelled {
				r.Estimate++
			}
			if t.Unplanned && t.State != Cancelled {
				r.Unplanned += t.Weight
			}
		}
		return r
	}
	for _, d := range f.Demands {
		out = append(out, row(d.ID, d.Title, d.Plan, false, d.AllTasks()))
		for _, s := range d.Stages {
			out = append(out, row(s.ID, s.Title, "", true, s.Tasks))
		}
	}
	return out
}

// --- the merge of two branches ---

// stateRank orders the states for a merge where both sides moved a task: work done is never
// undone by a merge, and work cancelled stays out of the count.
func stateRank(s State) int {
	switch s {
	case Done:
		return 3
	case Cancelled:
		return 2
	case Doing:
		return 1
	}
	return 0
}

// Merge unites two branches' progress over their common base. What one side added is kept;
// a task both sides changed takes the further state (done beats cancelled beats in
// progress beats to do), and a weight changed on one side only takes that side's. An id
// both sides created for different work — two agents adding the next task on two branches
// — is renumbered on their side, so both tasks stay.
func Merge(base, ours, theirs *File) *File {
	if base == nil {
		base = &File{}
	}
	out := ours
	for _, td := range theirs.Demands {
		od := out.Demand(td.ID)
		bd := base.Demand(td.ID)
		if od != nil && bd == nil && od.Title != td.Title {
			// Both sides created this id for different demands: theirs takes the next one.
			n := 0
			for _, d := range out.Demands {
				n = max(n, d.num)
			}
			renumberDemand(td, n+1)
			od = nil
		}
		if od == nil {
			out.Demands = append(out.Demands, td)
			continue
		}
		mergeDemand(bd, od, td)
	}
	if theirs.Unit != "" && ours.Unit == "" {
		out.Unit = theirs.Unit
	}
	return out
}

func renumberDemand(d *Demand, n int) {
	d.num, d.ID = n, demandID(n)
	for _, s := range d.Stages {
		s.ID = stageID(n, s.num)
	}
	for _, t := range d.AllTasks() {
		t.ID = taskID(n, t.num)
	}
}

// mergeDemand brings their side of one demand into ours.
func mergeDemand(bd, od, td *Demand) {
	baseTask := func(id string) *Task {
		if bd == nil {
			return nil
		}
		for _, t := range bd.AllTasks() {
			if t.ID == id {
				return t
			}
		}
		return nil
	}
	ourTask := func(id string) *Task {
		for _, t := range od.AllTasks() {
			if t.ID == id {
				return t
			}
		}
		return nil
	}
	nextTask := func() int {
		n := 0
		for _, t := range od.AllTasks() {
			n = max(n, t.num)
		}
		return n + 1
	}
	stageOf := func(ts *Stage) *Stage {
		if ts == nil {
			return nil
		}
		for _, s := range od.Stages {
			if s.ID == ts.ID {
				if s.Title == ts.Title {
					return s
				}
				break
			}
		}
		// A stage only their side has — or the same id for another stage — joins ours.
		n := 0
		for _, s := range od.Stages {
			n = max(n, s.num)
			if s.Title == ts.Title {
				return s
			}
		}
		ns := &Stage{ID: stageID(od.num, n+1), num: n + 1, Title: ts.Title}
		od.Stages = append(od.Stages, ns)
		return ns
	}
	add := func(ts *Stage, t *Task) {
		if s := stageOf(ts); s != nil {
			s.Tasks = append(s.Tasks, t)
			return
		}
		od.Tasks = append(od.Tasks, t)
	}
	visit := func(ts *Stage, tt *Task) {
		ot := ourTask(tt.ID)
		bt := baseTask(tt.ID)
		switch {
		case ot == nil:
			add(ts, tt)
		case bt == nil && ot.Title != tt.Title:
			// Both sides created this id for different tasks: theirs takes the next one.
			tt.num = nextTask()
			tt.ID = taskID(od.num, tt.num)
			add(ts, tt)
		default:
			mergeTask(bt, ot, tt)
		}
	}
	for _, t := range td.Tasks {
		visit(nil, t)
	}
	for _, s := range td.Stages {
		for _, t := range s.Tasks {
			visit(s, t)
		}
		stageOf(s) // a stage their side added with no task yet
	}
}

// mergeTask brings their change of one task into ours.
func mergeTask(bt, ot, tt *Task) {
	theirsMoved := bt == nil || tt.State != bt.State || tt.By != bt.By || tt.Reason != bt.Reason
	oursMoved := bt == nil || ot.State != bt.State || ot.By != bt.By || ot.Reason != bt.Reason
	if theirsMoved && (!oursMoved || stateRank(tt.State) > stateRank(ot.State)) {
		ot.State, ot.By, ot.At, ot.Reason = tt.State, tt.By, tt.At, tt.Reason
	}
	if bt != nil && (tt.Weight != bt.Weight || tt.Estimate != bt.Estimate) && ot.Weight == bt.Weight && ot.Estimate == bt.Estimate {
		ot.Weight, ot.Estimate = tt.Weight, tt.Estimate
	}
	if bt != nil && tt.Unplanned != bt.Unplanned && ot.Unplanned == bt.Unplanned {
		ot.Unplanned = tt.Unplanned
	}
}

// IsFile says whether a text is a progress file Anchors wrote: it opens with its header.
func IsFile(text string) bool {
	return strings.HasPrefix(strings.ReplaceAll(text, "\r\n", "\n"), headerLine+"\n")
}
