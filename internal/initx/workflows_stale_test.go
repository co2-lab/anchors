package initx

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// staleFakeGH answers the stale script: the issue list runs the script's own `--jq`
// through the real `jq` over $FAKE_ISSUES; the comments of card N come from
// $FAKE_COMMENTS_N (a JSON array), and $FAKE_API_FAIL names a card whose read fails.
const staleFakeGH = `#!/usr/bin/env bash
echo "gh $*" >> "$FAKE_LOG"
call="$1 $2"; path="$2"
jqx=""
while [ $# -gt 0 ]; do
  case "$1" in --jq) jqx="$2"; shift ;; esac
  shift
done
case "$call" in
  "issue list") printf '%s' "$FAKE_ISSUES" | jq -r "$jqx" ;;
  "api "*)
    n=$(echo "$path" | sed -E 's#.*/issues/([0-9]+)/comments#\1#')
    if [ "${FAKE_API_FAIL:-}" = "$n" ]; then echo "gh: HTTP 502" >&2; exit 1; fi
    var="FAKE_COMMENTS_$n"
    printf '%s' "${!var:-[]}" | jq -r "$jqx" ;;
esac
exit 0
`

func runStale(t *testing.T, env ...string) (out, calls string) {
	t.Helper()
	for _, bin := range []string{"bash", "jq"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("no %s on PATH", bin)
		}
	}
	b, err := workflowsFS.ReadFile("workflows/anchors-stale.yml")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Jobs map[string]struct {
			Steps []struct {
				Run string `yaml:"run"`
			} `yaml:"steps"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	script := ""
	for _, j := range doc.Jobs {
		for _, s := range j.Steps {
			if strings.Contains(s.Run, "anchors-owner: (liberado)") {
				script = s.Run
			}
		}
	}
	if script == "" {
		t.Fatal("the stale step vanished")
	}
	// Fixed cut-offs: macOS `date` has no `-d`, and the test needs a known clock.
	script = strings.Replace(script, `limite_espera=$(date -u -d "-${HORAS_ESPERA} hours" +%Y-%m-%dT%H:%M:%SZ)`, `limite_espera=2026-09-24T08:00:00Z`, 1)
	script = strings.Replace(script, `limite_andamento=$(date -u -d "-${HORAS_ANDAMENTO} hours" +%Y-%m-%dT%H:%M:%SZ)`, `limite_andamento=2026-09-23T10:00:00Z`, 1)
	script = strings.Replace(script, `limite_retrabalho=$(date -u -d "-${HORAS_RETRABALHO} hours" +%Y-%m-%dT%H:%M:%SZ)`, `limite_retrabalho=2026-09-24T09:00:00Z`, 1)
	if strings.Contains(script, "date -u -d") {
		t.Fatal("the cut-off lines changed shape; update the test's replacement")
	}
	dir := t.TempDir()
	bin := filepath.Join(dir, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bin, "gh"), []byte(staleFakeGH), 0o755); err != nil {
		t.Fatal(err)
	}
	logFile := filepath.Join(dir, "gh.log")
	cmd := exec.Command("bash", "-c", script)
	cmd.Env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"FAKE_LOG="+logFile, "GH_REPO=o/r", "LABEL=anchors", "HORAS_ESPERA=2", "HORAS_ANDAMENTO=24", "HORAS_RETRABALHO=1")
	cmd.Env = append(cmd.Env, env...)
	o, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("the stale script failed: %v\n%s", err, o)
	}
	c, _ := os.ReadFile(logFile)
	return string(o), string(c)
}

// THE STALE JOB, RUN. Measured in blue-eyes: the list asked for the comments of 200 cards
// and GitHub's GraphQL answered 502/504 for 60 runs in a row. The list now carries no
// comments; the owner comes per card from REST; a card already released is not
// released again (the marker is the PREFIX); and an API failure skips the card.
func TestStaleReleasesByRESTOwnerAndOnlyOnce(t *testing.T) {
	old := "2026-09-20T10:00:00Z"
	issues := `[
	 {"number":5,"updatedAt":"` + old + `","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}]},
	 {"number":6,"updatedAt":"` + old + `","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}]},
	 {"number":7,"updatedAt":"` + old + `","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}]},
	 {"number":8,"updatedAt":"2026-09-24T09:00:00Z","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}]}
	]`
	out, calls := runStale(t, "FAKE_ISSUES="+issues,
		`FAKE_COMMENTS_5=[{"body":"anchors-owner: a/one"},{"body":"work log"}]`,
		`FAKE_COMMENTS_6=[{"body":"anchors-owner: a/one"},{"body":"anchors-owner: (liberado) — sem progresso há mais de 24h"}]`,
		`FAKE_COMMENTS_7=[{"body":"anchors-owner: a/one"}]`, "FAKE_API_FAIL=7",
		`FAKE_COMMENTS_8=[{"body":"anchors-owner: a/two"}]`)

	for _, line := range strings.Split(calls, "\n") {
		if strings.HasPrefix(line, "gh issue list") && strings.Contains(line, "comments") {
			t.Errorf("the issue list asks for comments again — the query GitHub timed out on: %s", line)
		}
	}
	if !strings.Contains(calls, "gh issue comment 5 --body anchors-owner: (liberado)") ||
		!strings.Contains(calls, "dono anterior: a/one") {
		t.Errorf("the stale owned card #5 was not released:\n%s\n%s", calls, out)
	}
	if strings.Contains(calls, "gh issue comment 6 ") {
		t.Error("card #6 was already released, and it was released again")
	}
	if strings.Contains(calls, "gh issue comment 7 ") || !strings.Contains(out, "#7") {
		t.Errorf("an API failure must skip card #7 and say so:\n%s", out)
	}
	if strings.Contains(calls, "gh issue comment 8 ") {
		t.Error("card #8 had progress inside the window and was released")
	}
}

// A REJECTED card goes back to `to-do` owned by its author, and the author's turn is short:
// 1h, not the 2h of an ordinary waiting card. The clock here: now 10:00, waiting cut-off
// 08:00, rework cut-off 09:00. A card last touched at 08:30 is released only if its owner
// came from a rejection (the "↩ PR #N rejected" line after the owner line).
func TestStaleReleasesARejectedCardAfterOneHour(t *testing.T) {
	at := "2026-09-24T08:30:00Z"
	issues := `[
	 {"number":5,"updatedAt":"` + at + `","labels":[{"name":"anchors"},{"name":"anchors:to-do"}]},
	 {"number":6,"updatedAt":"` + at + `","labels":[{"name":"anchors"},{"name":"anchors:to-do"}]},
	 {"number":7,"updatedAt":"2026-09-24T07:00:00Z","labels":[{"name":"anchors"},{"name":"anchors:to-do"}]}
	]`
	out, calls := runStale(t, "FAKE_ISSUES="+issues,
		`FAKE_COMMENTS_5=[{"body":"anchors-owner: a/author"},{"body":"↩ PR #9 rejected by a/rev — the card is back in to-do"}]`,
		`FAKE_COMMENTS_6=[{"body":"anchors-owner: a/other"}]`,
		`FAKE_COMMENTS_7=[{"body":"anchors-owner: a/other"}]`)
	if !strings.Contains(calls, "gh issue comment 5 --body anchors-owner: (liberado)") {
		t.Errorf("a rejected card idle for 1h30 must go to any agent:\n%s\n%s", calls, out)
	}
	if strings.Contains(calls, "gh issue comment 6 ") {
		t.Errorf("an ordinary waiting card idle for 1h30 keeps its 2h:\n%s", calls)
	}
	if !strings.Contains(calls, "gh issue comment 7 --body anchors-owner: (liberado)") {
		t.Errorf("an ordinary waiting card idle for 3h must still be released:\n%s", calls)
	}
}
