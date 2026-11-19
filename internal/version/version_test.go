// @anchors
//   code: VRSNC
//   ref: VRSNE

package version

import (
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/changelog"
)

func mustScheme(t *testing.T, kind, format, pattern string) Scheme {
	t.Helper()
	s, err := NewScheme(kind, format, pattern)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func mustParse(t *testing.T, s Scheme, raw string) V {
	t.Helper()
	v, err := s.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// ascending asserts each version is below the next.
func ascending(t *testing.T, s Scheme, vs ...string) {
	t.Helper()
	for i := 0; i+1 < len(vs); i++ {
		if c := Compare(mustParse(t, s, vs[i]), mustParse(t, s, vs[i+1])); c != -1 {
			t.Errorf("%s vs %s: %d, want -1", vs[i], vs[i+1], c)
		}
	}
}

func refused(t *testing.T, s Scheme, vs ...string) {
	t.Helper()
	for _, v := range vs {
		if _, err := s.Parse(v); err == nil {
			t.Errorf("%q was read under %s", v, s.Kind)
		}
	}
}

func commit(subject, body string) changelog.Commit {
	return changelog.Commit{Hash: "h" + subject[:3], Subject: subject, Body: body}
}

var today = time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC)

func next(t *testing.T, in Input) Result {
	t.Helper()
	if in.Today.IsZero() {
		in.Today = today
	}
	r, err := Next(in)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestSchemes(t *testing.T) {
	t.Run("VRSNE-B01: Semver versions are ordered as SemVer 2.0 says", func(t *testing.T) {
		s := mustScheme(t, SemVer, "", "")
		ascending(t, s, "1.0.0-alpha", "1.0.0-alpha.1", "1.0.0-beta.2", "1.0.0-beta.11", "1.0.0-rc.1", "1.0.0", "1.10.0")
		if Compare(mustParse(t, s, "1.0.0+build"), mustParse(t, s, "1.0.0")) != 0 {
			t.Error("build metadata takes no part in the order")
		}
		refused(t, s, "1.2", "01.2.3", "v1.2.3")
	})
	t.Run("VRSNE-B02: A calver version is read in the widths its format gives", func(t *testing.T) {
		s := mustScheme(t, CalVer, "YYYY.0M.MICRO", "")
		ascending(t, s, "2026.09.0", "2026.11.2")
		refused(t, s, "2026.9.0", "2026.11")
	})
	t.Run("VRSNE-B03: A number version is a non-negative integer", func(t *testing.T) {
		s := mustScheme(t, Number, "", "")
		ascending(t, s, "9", "412")
		refused(t, s, "-1", "4a")
	})
	t.Run("VRSNE-B04: A pattern version is ordered by its groups", func(t *testing.T) {
		s := mustScheme(t, Pattern, "", `r(\d+)-(\d+)`)
		ascending(t, s, "r2-99", "r14-3", "r14-10")
	})
}

func TestIncrement(t *testing.T) {
	t.Run("VRSNE-B05: Semver raises a part and resets the ones after it, and a pre-release is the release it precedes", func(t *testing.T) {
		s := mustScheme(t, SemVer, "", "")
		for _, c := range []struct{ from, part, want string }{
			{"1.4.2", "major", "2.0.0"}, {"1.4.2", "minor", "1.5.0"}, {"1.4.2", "patch", "1.4.3"},
			{"1.1.0-rc.1", "major", "2.0.0"}, {"1.1.0-rc.1", "minor", "1.1.0"}, {"1.1.0-rc.1", "patch", "1.1.0"},
		} {
			if got, err := s.Increment(mustParse(t, s, c.from), c.part, today); err != nil || got != c.want {
				t.Errorf("%s %s: %q %v, want %s", c.from, c.part, got, err, c.want)
			}
		}
	})
	t.Run("VRSNE-B06: Calver takes the date from today and its counter goes up only on the same date", func(t *testing.T) {
		s := mustScheme(t, CalVer, "YYYY.0M.MICRO", "")
		v := mustParse(t, s, "2026.11.2")
		if got, _ := s.Increment(v, "", today); got != "2026.11.3" {
			t.Errorf("same month: %s", got)
		}
		if got, _ := s.Increment(v, "", time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)); got != "2026.12.0" {
			t.Errorf("new month: %s", got)
		}
	})
	t.Run("VRSNE-B07: A pattern raises the group it names and resets the later ones", func(t *testing.T) {
		s := mustScheme(t, Pattern, "", `r(?P<major>\d+)-(\d+)`)
		v := mustParse(t, s, "r14-3")
		if got, _ := s.Increment(v, "major", today); got != "r15-0" {
			t.Errorf("major: %s", got)
		}
		if got, _ := s.Increment(v, "2", today); got != "r14-4" {
			t.Errorf("group 2: %s", got)
		}
	})
}

func TestNext(t *testing.T) {
	semver, _ := NewScheme(SemVer, "", "")
	fix, feat := commit("fix: a crash", ""), commit("feat: a screen", "")
	t.Run("VRSNE-B08: Manual is the version the sources hold", func(t *testing.T) {
		r := next(t, Input{Scheme: semver, Strategy: Strategy{Name: Manual}, Current: "3.1.0", Released: "3.0.0"})
		if r.Version != "3.1.0" {
			t.Errorf("got %+v", r)
		}
	})
	t.Run("VRSNE-B09: Increment raises the latest release, and stays there with no commit since it", func(t *testing.T) {
		in := Input{Scheme: semver, Strategy: Strategy{Name: Increment, Part: "patch"}, Released: "1.4.2", Commits: []changelog.Commit{commit("chore: deps", "")}}
		if r := next(t, in); r.Version != "1.4.3" {
			t.Errorf("one commit: %+v", r)
		}
		in.Commits = nil
		if r := next(t, in); r.Version != "1.4.2" {
			t.Errorf("no commit: %+v", r)
		}
	})
	t.Run("VRSNE-B10: Semantic raises the release by the largest change since it", func(t *testing.T) {
		for _, c := range []struct {
			commits []changelog.Commit
			want    string
			reason  string
		}{
			{[]changelog.Commit{feat, fix, commit("feat(x): another", "")}, "1.5.0", "2 feat, 1 fix"},
			{[]changelog.Commit{fix}, "1.4.3", "1 fix"},
			{[]changelog.Commit{fix, commit("refactor!: the store", "")}, "2.0.0", "1 breaking"},
			{[]changelog.Commit{commit("chore: deps", "")}, "1.4.2", "nothing to release"},
		} {
			r := next(t, Input{Scheme: semver, Strategy: Strategy{Name: Semantic}, Released: "1.4.2", Commits: c.commits})
			if r.Version != c.want || !strings.Contains(r.Reason, c.reason) {
				t.Errorf("got %+v, want %s saying %q", r, c.want, c.reason)
			}
		}
	})
	t.Run("VRSNE-B11: Semantic with initial bumps the minor for a breaking change, and types give a type its level", func(t *testing.T) {
		st := Strategy{Name: Semantic, Initial: true, Types: map[string]string{"perf": "patch", "feat": "none"}}
		for _, c := range []struct {
			c    changelog.Commit
			want string
		}{
			{commit("feat!: a new store", ""), "0.4.0"}, {commit("perf: a faster load", ""), "0.3.2"}, {feat, "0.3.1"},
		} {
			if r := next(t, Input{Scheme: semver, Strategy: st, Released: "0.3.1", Commits: []changelog.Commit{c.c}}); r.Version != c.want {
				t.Errorf("%s: %+v, want %s", c.c.Subject, r, c.want)
			}
		}
	})
	t.Run("VRSNE-B12: Before any release, the sources hold the version, or the scheme gives its first one", func(t *testing.T) {
		number, _ := NewScheme(Number, "", "")
		for _, c := range []struct {
			in   Input
			want string
		}{
			{Input{Scheme: semver, Strategy: Strategy{Name: Semantic}, Current: "0.9.0"}, "0.9.0"},
			{Input{Scheme: semver, Strategy: Strategy{Name: Semantic}}, "1.0.0"},
			{Input{Scheme: semver, Strategy: Strategy{Name: Semantic, Initial: true}}, "0.1.0"},
			{Input{Scheme: number, Strategy: Strategy{Name: Increment}}, "1"},
		} {
			if r := next(t, c.in); r.Version != c.want {
				t.Errorf("%+v: %+v, want %s", c.in.Strategy, r, c.want)
			}
		}
	})
	t.Run("VRSNE-B13: A Version footer forces the next version, the newest one winning", func(t *testing.T) {
		cs := []changelog.Commit{commit("feat: one", "Version: 2.0.0"), commit("fix: two", "Version: 3.0.0\nBug: seen")}
		r := next(t, Input{Scheme: semver, Strategy: Strategy{Name: Semantic}, Released: "1.4.2", Commits: cs})
		if r.Version != "3.0.0" || !r.Forced || !strings.Contains(r.Reason, cs[1].Hash) {
			t.Errorf("got %+v", r)
		}
	})
	t.Run("VRSNE-B14: None gives no version, whatever the footers say", func(t *testing.T) {
		r := next(t, Input{Scheme: semver, Strategy: Strategy{Name: None}, Commits: []changelog.Commit{commit("feat: one", "Version: 2.0.0")}})
		if r.Version != "" || !strings.Contains(r.Reason, "never versions") {
			t.Errorf("got %+v", r)
		}
	})
}

func TestErrors(t *testing.T) {
	t.Run("VRSNE-E01: A scheme that reads nothing is refused", func(t *testing.T) {
		for _, c := range []struct{ kind, format, pattern, says string }{
			{"date", "", "", `"date"`}, {CalVer, "", "", "needs a format"}, {CalVer, "YYYY.MICRO.0M", "", "MICRO is the last"},
			{CalVer, "YYYY.WW", "", `"WW"`}, {Pattern, "", "(", "not a regular expression"}, {Pattern, "", `r\d+`, "no capture group"},
		} {
			if _, err := NewScheme(c.kind, c.format, c.pattern); err == nil || !strings.Contains(err.Error(), c.says) {
				t.Errorf("%s %q %q: %v, want it to say %s", c.kind, c.format, c.pattern, err, c.says)
			}
		}
	})
	semver, _ := NewScheme(SemVer, "", "")
	fails := func(t *testing.T, in Input, says ...string) {
		t.Helper()
		in.Today = today
		_, err := Next(in)
		for _, s := range says {
			if err == nil || !strings.Contains(err.Error(), s) {
				t.Errorf("%+v: %v, want it to say %q", in.Strategy, err, s)
			}
		}
	}
	t.Run("VRSNE-E02: A version the scheme cannot read is an error, never smaller nor equal", func(t *testing.T) {
		fails(t, Input{Scheme: semver, Strategy: Strategy{Name: Manual}, Current: "1.2"}, `"1.2"`, "semver")
		fails(t, Input{Scheme: semver, Strategy: Strategy{Name: Manual}, Current: "1.2.0", Released: "v1.0.0"}, `"v1.0.0"`, "semver")
	})
	t.Run("VRSNE-E03: A next version below the latest release is refused", func(t *testing.T) {
		fails(t, Input{Scheme: semver, Strategy: Strategy{Name: Manual}, Current: "1.9.0", Released: "2.0.0"}, "1.9.0", "2.0.0")
		fails(t, Input{Scheme: semver, Strategy: Strategy{Name: Semantic}, Released: "2.0.0", Commits: []changelog.Commit{commit("fix: x", "Version: 1.0.0")}}, "1.0.0", "2.0.0")
	})
	t.Run("VRSNE-E04: A strategy that would have to guess is refused", func(t *testing.T) {
		number, _ := NewScheme(Number, "", "")
		pattern, _ := NewScheme(Pattern, "", `r(\d+)`)
		calver, _ := NewScheme(CalVer, "YYYY.0M", "")
		one := []changelog.Commit{commit("fix: x", "")}
		fails(t, Input{Scheme: semver, Strategy: Strategy{Name: Increment, Part: "build"}, Released: "1.0.0", Commits: one}, `"build"`)
		fails(t, Input{Scheme: number, Strategy: Strategy{Name: Semantic}, Released: "4"}, "semantic", "semver")
		fails(t, Input{Scheme: semver, Strategy: Strategy{Name: Manual}}, "manual", "no")
		fails(t, Input{Scheme: pattern, Strategy: Strategy{Name: Increment, Part: "1"}}, "increment", "no release")
		fails(t, Input{Scheme: calver, Strategy: Strategy{Name: Increment}, Released: "2026.11", Commits: one}, "no counter")
		fails(t, Input{Scheme: semver, Strategy: Strategy{Name: "latest"}, Released: "1.0.0"}, `"latest"`)
	})
}

func TestIntentsStrategy(t *testing.T) {
	semver, _ := NewScheme(SemVer, "", "")
	t.Run("VRSNE-B15: intents raises the latest release by the largest bump of the intents waiting; none or no intent releases nothing", func(t *testing.T) {
		in := func(bumps ...string) Input {
			return Input{Scheme: semver, Strategy: Strategy{Name: Intents}, Released: "1.2.3", Bumps: bumps}
		}
		if r := next(t, in("patch", "minor", "none")); r.Version != "1.3.0" || !strings.Contains(r.Reason, "minor, the largest bump of 2 intent(s) since 1.2.3") {
			t.Errorf("minor: %+v", r)
		}
		if r := next(t, in("patch", "major")); r.Version != "2.0.0" {
			t.Errorf("major: %+v", r)
		}
		if r := next(t, in("none")); r.Version != "1.2.3" || !strings.Contains(r.Reason, "nothing to release") {
			t.Errorf("none: %+v", r)
		}
		num, _ := NewScheme(Number, "", "")
		if r := next(t, Input{Scheme: num, Strategy: Strategy{Name: Intents}, Released: "41", Bumps: []string{"number"}}); r.Version != "42" {
			t.Errorf("number: %+v", r)
		}
		if r := next(t, Input{Scheme: semver, Strategy: Strategy{Name: Intents}, Current: "0.1.0"}); r.Version != "0.1.0" || !strings.Contains(r.Reason, "no release yet") {
			t.Errorf("first: %+v", r)
		}
	})
}
