// @anchors
//   code: RPWRN
//   ref: PRCRN

//go:build windows

package runs

import (
	"os/exec"
	"strconv"
	"strings"
)

// winProcsScript prints each process as `pid<TAB>ppid<TAB>cpu seconds<TAB>command line`.
const winProcsScript = `Get-CimInstance Win32_Process | ForEach-Object { "{0}` + "`t" + `{1}` + "`t" + `{2}` + "`t" + `{3}" -f $_.ProcessId, $_.ParentProcessId, [int64](($_.KernelModeTime + $_.UserModeTime) / 10000000), $_.CommandLine }`

// ListProcs lists the system's processes. Windows gives no working directory to another
// process: a process belongs to the project by its command line, or by descending from one
// that does (InProject).
func ListProcs() ([]Proc, error) {
	out, err := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", winProcsScript).Output()
	if err != nil {
		return nil, err
	}
	return parseWinProcs(string(out)), nil
}

func parseWinProcs(out string) []Proc {
	var procs []Proc
	for _, line := range strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n") {
		f := strings.SplitN(line, "\t", 4)
		if len(f) < 4 {
			continue
		}
		pid, err1 := strconv.Atoi(strings.TrimSpace(f[0]))
		ppid, err2 := strconv.Atoi(strings.TrimSpace(f[1]))
		cpu, _ := strconv.ParseFloat(strings.TrimSpace(f[2]), 64)
		if err1 != nil || err2 != nil {
			continue
		}
		procs = append(procs, Proc{PID: pid, PPID: ppid, CPU: cpu, Cmd: strings.TrimSpace(f[3])})
	}
	return procs
}

// Cwds reads no working directory on Windows.
func Cwds(pids []int) map[int]string { return map[int]string{} }

// Alive says whether the process exists.
func Alive(pid int) bool {
	if pid <= 0 {
		return false
	}
	out, err := exec.Command("tasklist", "/FI", "PID eq "+strconv.Itoa(pid), "/NH").Output()
	return err == nil && strings.Contains(string(out), " "+strconv.Itoa(pid)+" ")
}

// LoadAverage: Windows has none.
func LoadAverage() (float64, bool) { return 0, false }
