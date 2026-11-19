// @anchors
//   code: LDQUV
//   ref: LDQUS

// Package load admits the project's heavy runs by what they cost, against what the machine
// holds, so they wait for one another instead of fighting over one machine. A project
// reported two worktrees' suites running at once at a load of 30 to 40 on 10 cores, with
// timing tests breaking under it; a fixed queue would make a light check wait behind a
// mutation run of ten minutes.
//
// A run states what it costs in CPUs and memory and starts when that fits in what the
// machine has free. The leases of every run on the machine — every worktree, every project —
// live in one file, under a lock, so the processes see each other. A run that waited past
// its limit stops the admission of new runs until it fits: a stream of light runs never
// starves a heavy one. A lease whose process died is reclaimed.
package load

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/runs"
)

// Cost is what a run holds while it runs.
type Cost struct {
	CPU    float64 `json:"cpu"`
	Memory int64   `json:"memory"` // bytes
}

// Capacity is what the machine lends its runs: its CPUs and memory, minus a reserve left to
// whoever uses it. A Memory of 0 does not bound the queue (it could not be read).
type Capacity struct {
	CPU    float64 `json:"cpu"`
	Memory int64   `json:"memory"`
}

// Lease is a run admitted: what it costs, whose it is, and since when.
type Lease struct {
	ID      string    `json:"id"`
	PID     int       `json:"pid"`
	Project string    `json:"project"`
	Task    string    `json:"task"`
	Level   string    `json:"level"`
	Cost    Cost      `json:"cost"`
	Since   time.Time `json:"since"`
	Limit   float64   `json:"limit_seconds,omitempty"` // a waiter's wait limit
}

// State is the machine's leases and the runs waiting for one.
type State struct {
	Leases  []Lease `json:"leases"`
	Waiters []Lease `json:"waiters"`
}

// Request is a run asking to start.
type Request struct {
	Project, Task, Level string
	Cost                 Cost
	WaitLimit            time.Duration
	// Notify says, while the run waits, what it waits for; nil says nothing.
	Notify func(string)
}

// pollEvery is how often a waiting run asks again.
var pollEvery = 500 * time.Millisecond

var seq atomic.Int64

// Dir is where the machine's leases live: $ANCHORS_LOAD_DIR, else ~/.anchors/load.
func Dir() string {
	if d := os.Getenv("ANCHORS_LOAD_DIR"); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "anchors-load")
	}
	return filepath.Join(home, ".anchors", "load")
}

func stateFile() string { return filepath.Join(Dir(), "leases.json") }

// Enabled says whether the queue admits a project's runs: it declares `load:`, nothing
// turned the queue off ($ANCHORS_LOAD=off, which the commit hook sets so a commit never
// waits), and, in CI, the project asked for it there too.
func Enabled(w *config.Workload) bool {
	if w == nil {
		return false
	}
	if v := strings.ToLower(os.Getenv("ANCHORS_LOAD")); v == "off" || v == "0" || v == "false" {
		return false
	}
	if inCI() && !w.CI {
		return false
	}
	return true
}

func inCI() bool {
	for _, v := range []string{"CI", "GITHUB_ACTIONS"} {
		if s := strings.ToLower(os.Getenv(v)); s == "true" || s == "1" {
			return true
		}
	}
	return false
}

// MachineCapacity is what the machine lends: $ANCHORS_LOAD_CAPACITY (`cpu=8,memory=16GB`)
// when set, else its CPUs and memory minus a reserve — one CPU, and 2 GB or a quarter of the
// memory, whichever is less — so the machine stays usable while the runs go.
func MachineCapacity() Capacity {
	if v := os.Getenv("ANCHORS_LOAD_CAPACITY"); v != "" {
		if c, err := ParseCapacity(v); err == nil {
			return c
		}
	}
	c := Capacity{CPU: float64(runtime.NumCPU()) - 1}
	if c.CPU < 1 {
		c.CPU = 1
	}
	if total := totalMemory(); total > 0 {
		reserve := int64(2 << 30)
		if total/4 < reserve {
			reserve = total / 4
		}
		c.Memory = total - reserve
	}
	return c
}

// ParseCapacity reads `cpu=8,memory=16GB`; either part may be left out.
func ParseCapacity(s string) (Capacity, error) {
	var c Capacity
	for _, part := range strings.Split(s, ",") {
		k, v, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok {
			return c, fmt.Errorf("%q: write cpu=<n>,memory=<size>", s)
		}
		switch strings.TrimSpace(k) {
		case "cpu":
			n, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
			if err != nil || n <= 0 {
				return c, fmt.Errorf("%q: cpu is a positive number", s)
			}
			c.CPU = n
		case "memory":
			b, err := config.ParseBytes(v)
			if err != nil {
				return c, err
			}
			c.Memory = b
		default:
			return c, fmt.Errorf("%q: write cpu=<n>,memory=<size>", s)
		}
	}
	return c, nil
}

// Admit says whether a waiting run may start now: no other run that waited past its limit
// waited longer, and its cost fits in what the leases leave — or nothing runs at all, so a
// run costlier than the machine runs alone instead of never.
func Admit(st State, capacity Capacity, me Lease, now time.Time) bool {
	for _, w := range st.Waiters {
		if w.ID == me.ID {
			continue
		}
		starving := now.Sub(w.Since).Seconds() > w.Limit
		if starving && w.Since.Before(me.Since) {
			return false
		}
	}
	if len(st.Leases) == 0 {
		return true
	}
	used := Used(st)
	if used.CPU+me.Cost.CPU > capacity.CPU {
		return false
	}
	return capacity.Memory == 0 || used.Memory+me.Cost.Memory <= capacity.Memory
}

// Used is what the leases hold.
func Used(st State) Cost {
	var c Cost
	for _, l := range st.Leases {
		c.CPU += l.Cost.CPU
		c.Memory += l.Cost.Memory
	}
	return c
}

// Acquire waits until the run is admitted, and returns what releases its lease. A run never
// admitted for a reason the queue cannot answer — its file cannot be locked or read — runs
// anyway, saying so: the queue spaces runs out, it must not stop them.
func Acquire(r Request) (release func(), err error) {
	me := Lease{
		ID:      fmt.Sprintf("%d-%d", os.Getpid(), seq.Add(1)),
		PID:     os.Getpid(),
		Project: r.Project, Task: r.Task, Level: r.Level, Cost: r.Cost,
		Since: time.Now(), Limit: r.WaitLimit.Seconds(),
	}
	var told time.Time
	for {
		admitted, why, err := try(me)
		if err != nil {
			return func() {}, err
		}
		if admitted {
			return func() { _ = update(func(st *State) { st.Leases = without(st.Leases, me.ID) }) }, nil
		}
		if r.Notify != nil && time.Since(told) > 30*time.Second {
			r.Notify(why)
			told = time.Now()
		}
		time.Sleep(pollEvery)
	}
}

// try admits the run or records it as waiting, under the lock, and says what it waits for.
func try(me Lease) (bool, string, error) {
	admitted := false
	why := ""
	err := update(func(st *State) {
		now := time.Now()
		if !Admit(*st, MachineCapacity(), me, now) {
			if !has(st.Waiters, me.ID) {
				st.Waiters = append(st.Waiters, me)
			}
			why = waitingFor(*st, me, now)
			return
		}
		st.Waiters = without(st.Waiters, me.ID)
		lease := me
		lease.Since, lease.Limit = now, 0
		st.Leases = append(st.Leases, lease)
		admitted = true
	})
	return admitted, why, err
}

// waitingFor says what a waiting run waits for: the runs holding the machine.
func waitingFor(st State, me Lease, now time.Time) string {
	used, capacity := Used(st), MachineCapacity()
	var holders []string
	for _, l := range st.Leases {
		holders = append(holders, fmt.Sprintf("%s %s (%s, %s)", filepath.Base(l.Project), l.Task, l.Level, now.Sub(l.Since).Round(time.Second)))
	}
	return fmt.Sprintf("load: %s %s (%s: %.0f CPU, %s) waits — %.0f of %.0f CPU and %s of %s held by %s",
		filepath.Base(me.Project), me.Task, me.Level, me.Cost.CPU, Size(me.Cost.Memory),
		used.CPU, capacity.CPU, Size(used.Memory), Size(capacity.Memory), strings.Join(holders, "; "))
}

// update reads the state under the lock, drops the leases and waiters of dead processes,
// applies the change and writes it back.
func update(change func(*State)) error {
	if err := os.MkdirAll(Dir(), 0o755); err != nil {
		return err
	}
	unlock, err := mapx.LockNamed(stateFile(), "the load leases")
	if err != nil {
		return err
	}
	defer unlock()
	st, err := Read()
	if err != nil {
		return err
	}
	st.Leases = alive(st.Leases)
	st.Waiters = alive(st.Waiters)
	change(&st)
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	tmp := stateFile() + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, stateFile())
}

// Read is the machine's state as the file holds it; empty when there is none.
func Read() (State, error) {
	var st State
	b, err := os.ReadFile(stateFile())
	if errors.Is(err, fs.ErrNotExist) {
		return st, nil
	}
	if err != nil {
		return st, err
	}
	if len(b) == 0 {
		return st, nil
	}
	if err := json.Unmarshal(b, &st); err != nil {
		return State{}, fmt.Errorf("read %s: %w", stateFile(), err)
	}
	return st, nil
}

func alive(ls []Lease) []Lease {
	out := ls[:0]
	for _, l := range ls {
		if runs.Alive(l.PID) {
			out = append(out, l)
		}
	}
	return out
}

func without(ls []Lease, id string) []Lease {
	out := ls[:0]
	for _, l := range ls {
		if l.ID != id {
			out = append(out, l)
		}
	}
	return out
}

func has(ls []Lease, id string) bool {
	for _, l := range ls {
		if l.ID == id {
			return true
		}
	}
	return false
}

// Size writes bytes in the largest unit they reach, to a tenth: 6GB, 1.5GB, 512MB; "—" for 0.
func Size(b int64) string {
	if b <= 0 {
		return "—"
	}
	units := []struct {
		name string
		n    int64
	}{{"TB", 1 << 40}, {"GB", 1 << 30}, {"MB", 1 << 20}, {"KB", 1 << 10}}
	for _, u := range units {
		if b >= u.n {
			return strings.TrimSuffix(strconv.FormatFloat(float64(b)/float64(u.n), 'f', 1, 64), ".0") + u.name
		}
	}
	return strconv.FormatInt(b, 10) + "B"
}

// Sorted is the state with its leases and waiters oldest first.
func Sorted(st State) State {
	sort.SliceStable(st.Leases, func(i, j int) bool { return st.Leases[i].Since.Before(st.Leases[j].Since) })
	sort.SliceStable(st.Waiters, func(i, j int) bool { return st.Waiters[i].Since.Before(st.Waiters[j].Since) })
	return st
}
