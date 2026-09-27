package mapx

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"
)

// The map is ONE file, and every command that changes it rewrote it whole: it read the
// map, changed its part, and wrote everything back. Two processes doing that at once
// lost one of the two changes — reported from a project running three mutation ingests
// in parallel with an `anchors test` and the pre-commit check: ten ingestions of one
// agent disappeared, and the map went back to the old scores.
//
// Each command now says only what it changes, and applies it to the map as it is ON DISK
// at that moment, under an exclusive lock: re-read, apply, write. The heavy work — running
// a suite, parsing a report, running the gates — happens before, outside the lock, so the
// lock is held for the time of a read and a write, not of a run.

// LockTimeout is how long a writer waits for another one to finish.
var LockTimeout = 2 * time.Minute

// lockStaleAfter is how old a lock must be to be taken over when its owner cannot be
// checked (another machine, or Windows). A write under the lock takes seconds at most; a
// lock this old was left by a process that died holding it.
var lockStaleAfter = time.Minute

// lockPoll is how often a waiting writer tries again.
const lockPoll = 50 * time.Millisecond

// LockPath is the lock file of a map: beside it, with `.lock` appended.
func LockPath(mapPath string) string { return mapPath + ".lock" }

// Lock takes the map's write lock: a file created only if it does not exist, holding the
// owner's pid and host. A lock whose owner is a dead process of this machine, or one
// older than a minute, is taken over. It waits up to LockTimeout, and fails saying that
// another process is writing the map. The returned function releases it.
func Lock(mapPath string) (func(), error) {
	lp := LockPath(mapPath)
	host, _ := os.Hostname()
	me := fmt.Sprintf("%d %s", os.Getpid(), host)
	deadline := time.Now().Add(LockTimeout)
	for {
		f, err := os.OpenFile(lp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			_, werr := f.WriteString(me + "\n")
			cerr := f.Close()
			if werr != nil || cerr != nil {
				_ = os.Remove(lp)
				return nil, fmt.Errorf("write the map lock %s: %v", lp, errors.Join(werr, cerr))
			}
			return func() { _ = os.Remove(lp) }, nil
		}
		if !errors.Is(err, fs.ErrExist) {
			return nil, fmt.Errorf("take the map lock %s: %w", lp, err)
		}
		if lockAbandoned(lp, host) {
			// Only the lock read as abandoned is removed: if another waiter took it over
			// first, the file is now theirs, and Remove of a fresh lock would break it —
			// so the removal is by rename to a name only this attempt uses.
			tmp := fmt.Sprintf("%s.stale.%d.%d", lp, os.Getpid(), time.Now().UnixNano())
			if os.Rename(lp, tmp) == nil {
				if !lockAbandoned(tmp, host) {
					// It was taken over in between: give it back — by link, which fails
					// rather than overwrite a lock a third writer created meanwhile.
					_ = os.Link(tmp, lp)
				}
				_ = os.Remove(tmp)
			}
			continue
		}
		if time.Now().After(deadline) {
			owner, _ := os.ReadFile(lp)
			return nil, fmt.Errorf("another process is writing the map (%s holds %s) and did not finish in %s — "+
				"if no anchors process is running, remove the lock file",
				strings.TrimSpace(string(owner)), lp, LockTimeout)
		}
		time.Sleep(lockPoll)
	}
}

// lockAbandoned says whether the lock at lp was left by a process that no longer holds it:
// its owner is a dead process of this host, or the lock is older than lockStaleAfter.
func lockAbandoned(lp, host string) bool {
	fi, err := os.Stat(lp)
	if err != nil {
		return false // gone: the next attempt creates it
	}
	if time.Since(fi.ModTime()) > lockStaleAfter {
		return true
	}
	b, err := os.ReadFile(lp)
	if err != nil {
		return false
	}
	fields := strings.Fields(string(b))
	if len(fields) < 2 || fields[1] != host {
		return false // being written right now, or another machine's: only age decides
	}
	pid, err := strconv.Atoi(fields[0])
	if err != nil || pid <= 0 {
		return false
	}
	return !processAlive(pid)
}

// WithLock runs fn holding the map's write lock.
func WithLock(mapPath string, fn func() error) error {
	unlock, err := Lock(mapPath)
	if err != nil {
		return err
	}
	defer unlock()
	return fn()
}

// Update applies a change to the map as it is on disk now: under the lock it re-reads the
// map, applies the change and writes it. What another process wrote in between is in the
// map re-read, so it is kept; the change must say only what it changes.
func Update(mapPath string, apply func(g *Graph) error) error {
	return WithLock(mapPath, func() error {
		g, err := Load(mapPath)
		if err != nil {
			return err
		}
		if err := apply(g); err != nil {
			return err
		}
		return Save(g, mapPath)
	})
}
