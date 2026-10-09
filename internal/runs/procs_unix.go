// @anchors
//   code: RPURN
//   ref: PRCRN

//go:build !windows

package runs

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// ListProcs lists the system's processes, without their working directories (Cwds reads
// those, for the few that matter).
func ListProcs() ([]Proc, error) {
	out, err := exec.Command("ps", "-axo", "pid=,ppid=,time=,command=").Output()
	if err != nil {
		return nil, err
	}
	return parsePS(string(out)), nil
}

// Cwds reads the working directory of each process given: from `/proc` on Linux, through
// `lsof` elsewhere. A process it cannot read is left out.
func Cwds(pids []int) map[int]string {
	out := map[int]string{}
	if len(pids) == 0 {
		return out
	}
	if runtime.GOOS == "linux" {
		for _, pid := range pids {
			if d, err := os.Readlink("/proc/" + strconv.Itoa(pid) + "/cwd"); err == nil {
				out[pid] = d
			}
		}
		return out
	}
	list := make([]string, len(pids))
	for i, p := range pids {
		list[i] = strconv.Itoa(p)
	}
	b, _ := exec.Command("lsof", "-a", "-d", "cwd", "-Fpn", "-p", strings.Join(list, ",")).Output()
	return parseLsofCwd(string(b))
}

// Alive says whether the process exists.
func Alive(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

// LoadAverage is the one-minute load average, and false where the system gives none: from
// `/proc` on Linux, from `sysctl` elsewhere.
func LoadAverage() (float64, bool) {
	if b, err := os.ReadFile("/proc/loadavg"); err == nil {
		return firstFloat(string(b))
	}
	b, err := exec.Command("sysctl", "-n", "vm.loadavg").Output()
	if err != nil {
		return 0, false
	}
	return firstFloat(strings.Trim(strings.TrimSpace(string(b)), "{} "))
}
