// @anchors
//   code: RNPRR
//   ref: PRCRN

package runs

import (
	"path/filepath"
	"strconv"
	"strings"
)

// Proc is one process as the operating system lists it: its id, its parent, the CPU time it
// has used, its command line and — when asked for — its working directory.
type Proc struct {
	PID  int
	PPID int
	CPU  float64 // seconds
	Cmd  string
	Cwd  string
}

// parseCPUTime reads the CPU time `ps` prints — `[dd-][hh:]mm:ss[.ss]` — as seconds; a
// field it cannot read is zero.
func parseCPUTime(s string) float64 {
	s = strings.TrimSpace(s)
	days := 0.0
	if d, rest, ok := strings.Cut(s, "-"); ok {
		n, err := strconv.Atoi(d)
		if err != nil {
			return 0
		}
		days, s = float64(n), rest
	}
	total := 0.0
	for _, p := range strings.Split(s, ":") {
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return 0
		}
		total = total*60 + v
	}
	return days*86400 + total
}

// parsePS reads `ps -axo pid=,ppid=,time=,command=`: one process a line.
func parsePS(out string) []Proc {
	var procs []Proc
	for _, line := range strings.Split(out, "\n") {
		f := strings.Fields(line)
		if len(f) < 4 {
			continue
		}
		pid, err1 := strconv.Atoi(f[0])
		ppid, err2 := strconv.Atoi(f[1])
		if err1 != nil || err2 != nil {
			continue
		}
		// The command keeps its own spacing: it starts after the third field.
		rest := strings.TrimLeft(line, " \t")
		for i := 0; i < 3; i++ {
			rest = strings.TrimLeft(rest[strings.IndexAny(rest, " \t"):], " \t")
		}
		procs = append(procs, Proc{PID: pid, PPID: ppid, CPU: parseCPUTime(f[2]), Cmd: rest})
	}
	return procs
}

// parseLsofCwd reads `lsof -a -d cwd -Fpn`: a `p<pid>` line, then its `n<path>`.
func parseLsofCwd(out string) map[int]string {
	cwd := map[int]string{}
	pid := 0
	for _, line := range strings.Split(out, "\n") {
		switch {
		case strings.HasPrefix(line, "p"):
			pid, _ = strconv.Atoi(line[1:])
		case strings.HasPrefix(line, "n") && pid > 0:
			cwd[pid] = line[1:]
		}
	}
	return cwd
}

// InProject says whether a process belongs to the project at root: its working directory is
// the root or under it, or — where the system gives no working directory — its command line
// names the root.
func InProject(p Proc, root string) bool {
	root = filepath.Clean(root)
	if p.Cwd != "" {
		c := filepath.Clean(p.Cwd)
		return c == root || strings.HasPrefix(c, root+string(filepath.Separator))
	}
	return strings.Contains(strings.ToLower(p.Cmd), strings.ToLower(root))
}

// Tree indexes processes by id and by parent.
type Tree struct {
	ByPID    map[int]Proc
	Children map[int][]int
}

// NewTree indexes a listing.
func NewTree(procs []Proc) Tree {
	t := Tree{ByPID: map[int]Proc{}, Children: map[int][]int{}}
	for _, p := range procs {
		t.ByPID[p.PID] = p
		t.Children[p.PPID] = append(t.Children[p.PPID], p.PID)
	}
	return t
}

// Descendants are the processes under pid, at any depth.
func (t Tree) Descendants(pid int) []int {
	var out []int
	queue := []int{pid}
	seen := map[int]bool{pid: true}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, c := range t.Children[cur] {
			if !seen[c] {
				seen[c] = true
				out = append(out, c)
				queue = append(queue, c)
			}
		}
	}
	return out
}

// TreeCPU is the CPU time of a process and every process under it.
func (t Tree) TreeCPU(pid int) float64 {
	total := t.ByPID[pid].CPU
	for _, d := range t.Descendants(pid) {
		total += t.ByPID[d].CPU
	}
	return total
}

// Ancestors are the processes above pid, nearest first.
func (t Tree) Ancestors(pid int) []int {
	var out []int
	seen := map[int]bool{pid: true}
	for cur := pid; ; {
		p, ok := t.ByPID[cur]
		if !ok {
			return out
		}
		if _, listed := t.ByPID[p.PPID]; !listed || seen[p.PPID] {
			return out
		}
		seen[p.PPID] = true
		out = append(out, p.PPID)
		cur = p.PPID
	}
}

// firstFloat is the first figure of a reading, and false when it has none.
func firstFloat(s string) (float64, bool) {
	f := strings.Fields(s)
	if len(f) == 0 {
		return 0, false
	}
	v, err := strconv.ParseFloat(f[0], 64)
	return v, err == nil
}
