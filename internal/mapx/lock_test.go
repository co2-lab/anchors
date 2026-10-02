// @anchors
//   ref: MPLCK

package mapx

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func lockedMap(t *testing.T, nodes ...string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), DefaultPath)
	g := &Graph{}
	for _, id := range nodes {
		g.Nodes = append(g.Nodes, Node{ID: id, Kind: KindCode})
	}
	if err := Save(g, p); err != nil {
		t.Fatal(err)
	}
	return p
}

func shortTimeout(t *testing.T, d time.Duration) {
	t.Helper()
	old := LockTimeout
	LockTimeout = d
	t.Cleanup(func() { LockTimeout = old })
}

func TestLock_FileBesideTheMap(t *testing.T) {
	t.Run("MPLCK-B01: The lock is a file beside the map with its owner", func(t *testing.T) {})
	p := lockedMap(t)
	unlock, err := Lock(p)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(p + ".lock")
	host, _ := os.Hostname()
	if err != nil || strings.TrimSpace(string(b)) != fmt.Sprintf("%d %s", os.Getpid(), host) || LockPath(p) != p+".lock" {
		t.Fatalf("the lock holds the pid and host, got %q %v", b, err)
	}
	unlock()
	if _, err := os.Stat(p + ".lock"); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the release removes the lock file, got %v", err)
	}
}

func TestLock_SecondWriterWaits(t *testing.T) {
	t.Run("MPLCK-B02: A second writer waits for the first", func(t *testing.T) {})
	p := lockedMap(t)
	unlock, err := Lock(p)
	if err != nil {
		t.Fatal(err)
	}
	var released atomic.Bool
	got := make(chan bool)
	go func() {
		u, err := Lock(p)
		if err == nil {
			u()
		}
		got <- err == nil && released.Load()
	}()
	time.Sleep(200 * time.Millisecond)
	released.Store(true)
	unlock()
	if !<-got {
		t.Fatal("the second writer must take the lock, and only after the release")
	}
}

func deadPid(t *testing.T) int {
	t.Helper()
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		t.Skip("no `true` to spawn a process that exits")
	}
	return cmd.Process.Pid
}

func TestLock_AbandonedLocks(t *testing.T) {
	t.Run("MPLCK-B03: An abandoned lock is taken over, a live one is not", func(t *testing.T) {})
	shortTimeout(t, 300*time.Millisecond)
	host, _ := os.Hostname()
	p := lockedMap(t)
	write := func(owner string, age time.Duration) {
		t.Helper()
		if err := os.WriteFile(p+".lock", []byte(owner+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		at := time.Now().Add(-age)
		if err := os.Chtimes(p+".lock", at, at); err != nil {
			t.Fatal(err)
		}
	}
	take := func() error {
		u, err := Lock(p)
		if err == nil {
			u()
		}
		return err
	}
	write(fmt.Sprintf("%d %s", deadPid(t), host), 0)
	if err := take(); err != nil {
		t.Errorf("a dead process's lock is taken over at once: %v", err)
	}
	write("1 some-other-host", 2*time.Minute)
	if err := take(); err != nil {
		t.Errorf("a lock older than a minute is taken over: %v", err)
	}
	write(fmt.Sprintf("%d some-other-host", deadPid(t)), 0)
	if err := take(); err == nil {
		t.Error("a fresh lock of another host is waited for, not taken — its pid means nothing here")
	}
	write(fmt.Sprintf("%d %s", os.Getpid(), host), 0)
	if err := take(); err == nil {
		t.Error("a live process's lock is waited for, not taken")
	}
}

// The child: records `MPLCK_N` run times of its own on its node of the shared map.
func TestLockChildWriter(t *testing.T) {
	p, node := os.Getenv("MPLCK_MAP"), os.Getenv("MPLCK_NODE")
	if p == "" {
		t.Skip("only runs as a child of TestLock_ParallelProcesses")
	}
	for j := 0; j < 20; j++ {
		key := fmt.Sprintf("suite-%d", j)
		if err := Update(p, func(g *Graph) error {
			g.RecordRunSeconds(node, key, 1)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLock_ParallelProcesses(t *testing.T) {
	t.Run("MPLCK-B04: Parallel processes each changing their part all reach the map", func(t *testing.T) {})
	nodes := []string{"n0", "n1", "n2", "n3"}
	p := lockedMap(t, nodes...)
	var cmds []*exec.Cmd
	for _, n := range nodes {
		c := exec.Command(os.Args[0], "-test.run=^TestLockChildWriter$", "-test.count=1")
		c.Env = append(os.Environ(), "MPLCK_MAP="+p, "MPLCK_NODE="+n)
		if err := c.Start(); err != nil {
			t.Fatal(err)
		}
		cmds = append(cmds, c)
	}
	for _, c := range cmds {
		if err := c.Wait(); err != nil {
			t.Fatalf("a child writer failed: %v", err)
		}
	}
	g, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, n := range g.Nodes {
		if n.Signal != nil {
			total += len(n.Signal.SecondsBySuite)
		}
	}
	if total != 80 {
		t.Fatalf("every write of every process reaches the map: want 80, got %d", total)
	}
}

func TestWithLock_ReleasesOnError(t *testing.T) {
	t.Run("MPLCK-B05: The lock is released when the function returns", func(t *testing.T) {})
	p := lockedMap(t)
	boom := errors.New("boom")
	if err := WithLock(p, func() error {
		if _, err := os.Stat(p + ".lock"); err != nil {
			t.Error("the function runs holding the lock")
		}
		return boom
	}); !errors.Is(err, boom) {
		t.Fatalf("the function's error comes back, got %v", err)
	}
	if _, err := os.Stat(p + ".lock"); !errors.Is(err, os.ErrNotExist) {
		t.Error("the lock is released after an error")
	}
}

func TestLock_Timeout(t *testing.T) {
	t.Run("MPLCK-E01: A lock held past the timeout", func(t *testing.T) {})
	shortTimeout(t, 200*time.Millisecond)
	p := lockedMap(t)
	unlock, err := Lock(p)
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	_, err = Lock(p)
	if err == nil || !strings.Contains(err.Error(), "another process is writing the map") || !strings.Contains(err.Error(), p+".lock") ||
		!strings.Contains(err.Error(), fmt.Sprint(os.Getpid())) {
		t.Fatalf("the error says another process writes the map, names the holder and the file, got %v", err)
	}
	if _, err := os.Stat(p + ".lock"); err != nil {
		t.Error("the holder's lock is left as it is")
	}
}

func TestUpdate_WritesNothingOnError(t *testing.T) {
	t.Run("MPLCK-E02: No map, or a refused change, writes nothing", func(t *testing.T) {})
	missing := filepath.Join(t.TempDir(), DefaultPath)
	if err := Update(missing, func(*Graph) error { return nil }); err == nil {
		t.Error("with no map, the load error comes back")
	}
	if _, err := os.Stat(missing); !errors.Is(err, os.ErrNotExist) {
		t.Error("with no map, nothing is written")
	}
	p := lockedMap(t, "a")
	before, _ := os.ReadFile(p)
	refused := errors.New("refused")
	if err := Update(p, func(g *Graph) error { g.Nodes = nil; return refused }); !errors.Is(err, refused) {
		t.Fatalf("the change's refusal comes back, got %v", err)
	}
	if after, _ := os.ReadFile(p); string(after) != string(before) {
		t.Error("a refused change writes nothing")
	}
	for _, f := range []string{missing, p} {
		if _, err := os.Stat(f + ".lock"); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("the lock of %s is free after the error", f)
		}
	}
}

func TestLock_CannotCreate(t *testing.T) {
	t.Run("MPLCK-E03: A lock that cannot be created names the file", func(t *testing.T) {})
	p := filepath.Join(t.TempDir(), "no-such-dir", DefaultPath)
	start := time.Now()
	_, err := Lock(p)
	if err == nil || !strings.Contains(err.Error(), p+".lock") {
		t.Fatalf("the error names the lock file, got %v", err)
	}
	if time.Since(start) > 5*time.Second {
		t.Error("a lock that cannot be created fails at once, not at the timeout")
	}
}

// The liveness a takeover rests on, asked of the real system: this process is alive, and
// a child that has exited is not. On Windows it used to answer "alive" for every pid.
func TestProcessAlive(t *testing.T) {
	t.Run("MPLCK-B03: An abandoned lock is taken over, a live one is not", func(t *testing.T) {})
	if !processAlive(os.Getpid()) {
		t.Error("this process is alive")
	}
	c := exec.Command(os.Args[0], "-test.run=^$")
	if err := c.Run(); err != nil {
		t.Fatal(err)
	}
	if processAlive(c.Process.Pid) {
		t.Errorf("a child that exited (pid %d) is not alive", c.Process.Pid)
	}
}
