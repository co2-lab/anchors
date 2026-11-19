// @anchors
//   code: LDQUT
//   ref: LDQUS

package load

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/config"
)

const gb = int64(1) << 30

func lease(id string, cpu float64, mem int64, since time.Time, limit time.Duration) Lease {
	return Lease{ID: id, PID: os.Getpid(), Cost: Cost{CPU: cpu, Memory: mem}, Since: since, Limit: limit.Seconds()}
}

func TestAdmit(t *testing.T) {
	t.Run("LDQUS-B01: A run starts when its cost fits in CPU and memory beside the leases; nothing running admits any cost", func(t *testing.T) {
		now := time.Now()
		capacity := Capacity{CPU: 8, Memory: 16 * gb}
		running := State{Leases: []Lease{lease("a", 4, 6*gb, now, 0)}}
		for _, c := range []struct {
			me   Lease
			want bool
		}{
			{lease("b", 4, 6*gb, now, time.Minute), true},
			{lease("c", 5, 1*gb, now, time.Minute), false},
			{lease("d", 1, 11*gb, now, time.Minute), false},
		} {
			if got := Admit(running, capacity, c.me, now); got != c.want {
				t.Errorf("%s: %v, want %v", c.me.ID, got, c.want)
			}
		}
		if !Admit(State{}, capacity, lease("huge", 32, 64*gb, now, time.Minute), now) {
			t.Errorf("a run costlier than the machine runs alone")
		}
		if !Admit(running, Capacity{CPU: 8}, lease("e", 4, 100*gb, now, time.Minute), now) {
			t.Errorf("an unread memory does not bound the queue")
		}
	})
	t.Run("LDQUS-B02: A run that waited past its limit stops the admission of every run that came after it, until it fits", func(t *testing.T) {
		now := time.Now()
		capacity := Capacity{CPU: 8, Memory: 16 * gb}
		heavy := lease("heavy", 6, 8*gb, now.Add(-10*time.Minute), 5*time.Minute)
		st := State{Leases: []Lease{lease("a", 4, 2*gb, now, 0)}, Waiters: []Lease{heavy}}
		light := lease("light", 1, 1*gb, now.Add(-time.Second), 5*time.Minute)
		if Admit(st, capacity, light, now) {
			t.Errorf("a light run that came after a heavy one past its limit waits")
		}
		heavy.Limit = (time.Hour).Seconds()
		st.Waiters = []Lease{heavy}
		if !Admit(st, capacity, light, now) {
			t.Errorf("within its limit, the heavy run does not hold the others")
		}
		st.Leases = nil
		heavy.Limit = 0
		st.Waiters = []Lease{heavy}
		if !Admit(st, capacity, heavy, now) {
			t.Errorf("the run past its limit starts once it fits")
		}
	})
}

func TestEnabled(t *testing.T) {
	t.Run("LDQUS-B03: The queue admits a project that declares load:, unless $ANCHORS_LOAD turns it off or the run is in CI and load.ci is not set", func(t *testing.T) {
		t.Setenv("CI", "")
		t.Setenv("GITHUB_ACTIONS", "")
		t.Setenv("ANCHORS_LOAD", "")
		w := &config.Workload{}
		if Enabled(nil) || !Enabled(w) {
			t.Errorf("declared or not")
		}
		t.Setenv("ANCHORS_LOAD", "off")
		if Enabled(w) {
			t.Errorf("turned off")
		}
		t.Setenv("ANCHORS_LOAD", "")
		t.Setenv("GITHUB_ACTIONS", "true")
		if Enabled(w) || !Enabled(&config.Workload{CI: true}) {
			t.Errorf("in CI only when asked")
		}
	})
}

func TestCapacity(t *testing.T) {
	t.Run("LDQUS-B04: The capacity is the machine's CPUs and memory minus a reserve, or what $ANCHORS_LOAD_CAPACITY says", func(t *testing.T) {
		t.Setenv("ANCHORS_LOAD_CAPACITY", "")
		c := MachineCapacity()
		if c.CPU < 1 || (totalMemory() > 0 && (c.Memory <= 0 || c.Memory >= totalMemory())) {
			t.Errorf("detected: %+v of %d", c, totalMemory())
		}
		t.Setenv("ANCHORS_LOAD_CAPACITY", "cpu=3,memory=4GB")
		if c := MachineCapacity(); c.CPU != 3 || c.Memory != 4*gb {
			t.Errorf("set: %+v", c)
		}
		if _, err := ParseCapacity("cores=3"); err == nil {
			t.Errorf("an unknown part is refused")
		}
		if Size(6*gb) != "6GB" || Size(3*gb/2) != "1.5GB" || Size(0) != "—" {
			t.Errorf("sizes: %s %s %s", Size(6*gb), Size(3*gb/2), Size(0))
		}
	})
}

func TestAcquire(t *testing.T) {
	t.Run("LDQUS-B05: Runs share the machine's leases: a run that does not fit waits, saying for what, until a lease is released", func(t *testing.T) {
		t.Setenv("ANCHORS_LOAD_DIR", t.TempDir())
		t.Setenv("ANCHORS_LOAD_CAPACITY", "cpu=4,memory=8GB")
		prev := pollEvery
		pollEvery = 20 * time.Millisecond
		t.Cleanup(func() { pollEvery = prev })
		first, err := Acquire(Request{Project: "/p/app", Task: "test unit", Level: "heavy", Cost: Cost{CPU: 3, Memory: 6 * gb}, WaitLimit: time.Minute})
		if err != nil {
			t.Fatal(err)
		}
		st, _ := Read()
		if len(st.Leases) != 1 || st.Leases[0].Task != "test unit" {
			t.Fatalf("the lease is recorded: %+v", st)
		}
		said := make(chan string, 1)
		done := make(chan struct{})
		go func() {
			second, _ := Acquire(Request{Project: "/p/web", Task: "mutation unit", Level: "medium", Cost: Cost{CPU: 2, Memory: 2 * gb}, WaitLimit: time.Minute,
				Notify: func(s string) {
					select {
					case said <- s:
					default:
					}
				}})
			second()
			close(done)
		}()
		select {
		case s := <-said:
			if !strings.Contains(s, "web mutation unit (medium") || !strings.Contains(s, "held by app test unit (heavy") {
				t.Errorf("the wait says for what: %s", s)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("the second run should wait")
		}
		if st, _ := Read(); len(st.Waiters) != 1 {
			t.Errorf("the waiting run is recorded: %+v", st)
		}
		first()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("the released lease admits the waiting run")
		}
		if st, _ := Read(); len(st.Leases)+len(st.Waiters) != 0 {
			t.Errorf("everything released: %+v", st)
		}
	})
	t.Run("LDQUS-B06: A lease whose process died is reclaimed", func(t *testing.T) {
		t.Setenv("ANCHORS_LOAD_DIR", t.TempDir())
		t.Setenv("ANCHORS_LOAD_CAPACITY", "cpu=4,memory=8GB")
		dead := Lease{ID: "x", PID: 999999, Cost: Cost{CPU: 4, Memory: 8 * gb}, Since: time.Now()}
		if err := update(func(st *State) { st.Leases = append(st.Leases, dead) }); err != nil {
			t.Fatal(err)
		}
		release, err := Acquire(Request{Task: "test", Cost: Cost{CPU: 4, Memory: 8 * gb}, WaitLimit: time.Minute})
		if err != nil {
			t.Fatal(err)
		}
		release()
	})
}
