package initx

import (
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// The review outcome as the `anchors/review` commit status (blue-eyes #835), and the merge
// that happens without it (blue-eyes #782).
//
// These tests RUN the scripts of `anchors-pr-checks.yml` against a fake `gh` that serves
// JSON fixtures through the real `jq` — the same filters the pipeline uses in CI. Reading
// the YAML for strings would prove the words are there; running it proves which status a
// given card history produces.

const prChecksFile = "workflows/anchors-pr-checks.yml"

type prChecksDoc struct {
	On   map[string]any `yaml:"on"`
	Perm map[string]any `yaml:"permissions"`
	Jobs map[string]struct {
		If      string         `yaml:"if"`
		Needs   any            `yaml:"needs"`
		Outputs map[string]any `yaml:"outputs"`
		Steps   []struct {
			Name string            `yaml:"name"`
			Env  map[string]string `yaml:"env"`
			Run  string            `yaml:"run"`
		} `yaml:"steps"`
	} `yaml:"jobs"`
}

func loadPRChecks(t *testing.T) prChecksDoc {
	t.Helper()
	b, err := fs.ReadFile(workflowsFS, prChecksFile)
	if err != nil {
		t.Fatal(err)
	}
	var doc prChecksDoc
	if err := yaml.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	return doc
}

// reviewFakeGH is a `gh` that answers reads from JSON fixtures (filtered by the real `jq`, as
// `gh --jq` does) and records every write in calls.log. The fixture of a read is named
// after its positional arguments: `gh pr view 7` → pr_view_7.json, `gh api
// repos/o/r/issues/12/events` → api_repos_o_r_issues_12_events.json.
const reviewFakeGH = `#!/usr/bin/env bash
dir="$(cd "$(dirname "$0")" && pwd)"
jq_expr=""; pos=(); write=""; body=""; fields=()
while [ $# -gt 0 ]; do
  case "$1" in
    --jq) jq_expr="$2"; shift 2 ;;
    --add-label) fields+=("add-label=$2"); shift 2 ;;
    --json|--label|--limit|--state|--repo|--remove-label) shift 2 ;;
    --method) [ "$2" = "POST" ] && write=1; shift 2 ;;
    --body) body="$2"; shift 2 ;;
    -f) fields+=("$2"); shift 2 ;;
    --paginate) shift ;;
    *) pos+=("$1"); shift ;;
  esac
done
case "${pos[0]} ${pos[1]}" in
  "issue comment"|"issue edit") write=1 ;;
esac
if [ -n "$write" ]; then
  { printf 'CALL %s\n' "${pos[*]}"; for f in "${fields[@]}"; do printf 'FIELD %s\n' "$f"; done
    [ -n "$body" ] && printf 'BODY %s\n' "$body"; } >> "$dir/calls.log"
  exit 0
fi
key=$(IFS=_; echo "${pos[*]}" | tr '/' '_')
f="$dir/fixtures/$key.json"
[ -f "$f" ] || { echo "fake gh: no fixture $key" >&2; exit 1; }
if [ -n "$jq_expr" ]; then jq -r "$jq_expr" "$f"; else cat "$f"; fi
`

type ghWorld struct {
	t   *testing.T
	dir string
}

func newGHWorld(t *testing.T) *ghWorld {
	t.Helper()
	for _, bin := range []string{"bash", "jq", "base64"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s not on the PATH", bin)
		}
	}
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "fixtures"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(reviewFakeGH), 0o755); err != nil {
		t.Fatal(err)
	}
	return &ghWorld{t: t, dir: dir}
}

func (w *ghWorld) fixture(key string, v any) {
	w.t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		w.t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.dir, "fixtures", key+".json"), b, 0o644); err != nil {
		w.t.Fatal(err)
	}
}

// run executes a step script with the fake `gh` first on the PATH. It returns the script
// output, the recorded writes and the GITHUB_OUTPUT content.
func (w *ghWorld) run(script string, env map[string]string) (out, calls, outputs string) {
	w.t.Helper()
	ghOut := filepath.Join(w.dir, "github_output")
	_ = os.WriteFile(ghOut, nil, 0o644)
	_ = os.Remove(filepath.Join(w.dir, "calls.log"))
	cmd := exec.Command("bash", "-c", script)
	cmd.Env = append(os.Environ(),
		"PATH="+w.dir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"GITHUB_OUTPUT="+ghOut,
		"GITHUB_STEP_SUMMARY="+filepath.Join(w.dir, "summary"),
		"GH_REPO=o/r", "LABEL=anchors",
		"RUN_URL=https://example.test/run/1",
	)
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		w.t.Fatalf("script failed: %v\n%s", err, b)
	}
	c, _ := os.ReadFile(filepath.Join(w.dir, "calls.log"))
	o, _ := os.ReadFile(ghOut)
	return string(b), string(c), string(o)
}

// Card #12, linked by the PR body, assigned to `agent-b` for review at 10:00.
type reviewWorld struct {
	labels     []string
	owners     []map[string]string // card comments
	assignedAt string              // "" = never moved to in-review
	prComments []map[string]any
	prBody     string
	perms      map[string]string // login → repository permission, as the API reports it
}

func defaultReviewWorld() reviewWorld {
	return reviewWorld{
		labels: []string{"anchors", "anchors:in-review"},
		owners: []map[string]string{
			{"createdAt": "2026-09-20T08:00:00Z", "body": "anchors-owner: agent-a"},
			{"createdAt": "2026-09-20T09:00:00Z", "body": "anchors-owner: (liberado) — implementação concluída"},
			{"createdAt": "2026-09-20T10:00:00Z", "body": "anchors-owner: agent-b"},
		},
		assignedAt: "2026-09-20T10:00:00Z",
		prBody:     "Implements the thing.\n\nRefs #12\n",
	}
}

func prComment(at, who, body string) map[string]any {
	return prCommentBy(at, who, "someone", body)
}

// prCommentBy names the author: the permission fallback looks the LOGIN up.
func prCommentBy(at, who, login, body string) map[string]any {
	return map[string]any{"created_at": at, "author_association": who, "user": map[string]string{"login": login}, "body": body}
}

func (w *ghWorld) seedReview(rw reviewWorld) {
	w.fixture("pr_view_7", map[string]any{
		"body": rw.prBody, "title": "some work", "headRefOid": "abc123",
	})
	var labels []map[string]string
	for _, l := range rw.labels {
		labels = append(labels, map[string]string{"name": l})
	}
	w.fixture("issue_view_12", map[string]any{"labels": labels, "comments": rw.owners})
	w.fixture("issue_list", []any{})
	events := []map[string]any{
		{"event": "labeled", "label": map[string]string{"name": "anchors:in-progress"}, "created_at": "2026-09-20T08:00:00Z"},
	}
	if rw.assignedAt != "" {
		events = append(events, map[string]any{
			"event": "labeled", "label": map[string]string{"name": "anchors:in-review"}, "created_at": rw.assignedAt,
		})
	}
	w.fixture("api_repos_o_r_issues_12_events", events)
	if rw.prComments == nil {
		rw.prComments = []map[string]any{}
	}
	w.fixture("api_repos_o_r_issues_7_comments", rw.prComments)
	for login, perm := range rw.perms {
		w.fixture("api_repos_o_r_collaborators_"+login+"_permission", map[string]string{"permission": perm})
	}
}

func reviewScript(t *testing.T) string {
	t.Helper()
	job, ok := loadPRChecks(t).Jobs["review"]
	if !ok || len(job.Steps) == 0 {
		t.Fatal("`anchors-pr-checks.yml` has no `review` job — nothing publishes `anchors/review`")
	}
	return job.Steps[0].Run
}

var statusState = regexp.MustCompile(`FIELD state=(\w+)`)

func publishedStatus(calls string) (state, desc string) {
	if !strings.Contains(calls, "CALL api repos/o/r/statuses/abc123") {
		return "", ""
	}
	if m := statusState.FindStringSubmatch(calls); m != nil {
		state = m[1]
	}
	for _, l := range strings.Split(calls, "\n") {
		if strings.HasPrefix(l, "FIELD description=") {
			desc = strings.TrimPrefix(l, "FIELD description=")
		}
	}
	return state, desc
}

func TestPRChecksPublishesTheReviewStatus(t *testing.T) {
	script := reviewScript(t)
	env := map[string]string{"PR": "7", "HEAD_SHA": "abc123"}

	cases := []struct {
		name      string
		edit      func(*reviewWorld)
		wantState string // "" = no status published
		wantDesc  string
	}{
		{
			name: "success only on the assigned reviewer's approved line",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "## Revisão\n\nAll rules hold.\n\nanchors-review: approved by agent-b\n"),
				}
			},
			wantState: "success", wantDesc: "approved by agent-b",
		},
		{
			name: "failure on the assigned reviewer's rejected line",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: rejected by agent-b\r\n"),
				}
			},
			wantState: "failure", wantDesc: "rejected by agent-b",
		},
		{
			name:      "pending while the assigned reviewer has not posted a verdict",
			edit:      func(rw *reviewWorld) {},
			wantState: "pending", wantDesc: "awaiting the verdict of agent-b",
		},
		{
			name: "another agent's approval does not count",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-a"),
				}
			},
			wantState: "pending", wantDesc: "agent-b",
		},
		{
			name: "an approval from before the assignment belongs to an earlier round",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T09:30:00Z", "OWNER", "anchors-review: approved by agent-b"),
				}
			},
			wantState: "pending", wantDesc: "agent-b",
		},
		{
			name: "a verdict line inside a code block is an example",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "Write this when done:\n\n```\nanchors-review: approved by agent-b\n```\n"),
				}
			},
			wantState: "pending",
		},
		{
			name: "a line quoted mid-sentence is not a verdict",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "I will post anchors-review: approved by agent-b later"),
				}
			},
			wantState: "pending",
		},
		{
			name: "a commenter without write access cannot sign the review",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "NONE", "anchors-review: approved by agent-b"),
				}
			},
			wantState: "pending",
		},
		{
			// blue-eyes #982: an org with no PUBLIC members. The job's token sees the member's
			// comment as CONTRIBUTOR, and every approval was dropped.
			name: "a private org member's line counts by repository permission",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prCommentBy("2026-09-20T11:00:00Z", "CONTRIBUTOR", "dev", "anchors-review: approved by agent-b"),
				}
				rw.perms = map[string]string{"dev": "write"}
			},
			wantState: "success",
		},
		{
			name: "a contributor with only read access still cannot sign the review",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prCommentBy("2026-09-20T11:00:00Z", "CONTRIBUTOR", "outsider", "anchors-review: approved by agent-b"),
				}
				rw.perms = map[string]string{"outsider": "read"}
			},
			wantState: "pending",
		},
		{
			name: "the last verdict wins",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b"),
					prComment("2026-09-20T12:00:00Z", "MEMBER", "anchors-review: rejected by agent-b"),
				}
			},
			wantState: "failure",
		},
		{
			name: "the reviewer is the owner AT the assignment, not a later one",
			edit: func(rw *reviewWorld) {
				rw.owners = append(rw.owners, map[string]string{
					"createdAt": "2026-09-20T13:00:00Z", "body": "anchors-owner: agent-c",
				})
				rw.prComments = []map[string]any{
					prComment("2026-09-20T14:00:00Z", "OWNER", "anchors-review: approved by agent-c"),
				}
			},
			wantState: "pending", wantDesc: "agent-b",
		},
		{
			// Timestamps have one-second resolution: a release stamped in the same second
			// as the assignment sorts before it, and `(liberado)` is not a reviewer's name.
			name: "a released owner is never taken for the reviewer",
			edit: func(rw *reviewWorld) {
				rw.owners = append(rw.owners, map[string]string{
					"createdAt": "2026-09-20T10:00:00Z", "body": "anchors-owner: (liberado) — race",
				})
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b"),
				}
			},
			wantState: "success", wantDesc: "approved by agent-b",
		},
		{
			name: "a card back in ready-to-review is pending, whatever the earlier round said",
			edit: func(rw *reviewWorld) {
				rw.labels = []string{"anchors", "anchors:ready-to-review"}
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b"),
				}
			},
			wantState: "pending", wantDesc: "awaiting a reviewer",
		},
		{
			name: "a card never assigned for review gets no status",
			edit: func(rw *reviewWorld) {
				rw.labels = []string{"anchors", "anchors:in-progress"}
				rw.assignedAt = ""
			},
			wantState: "",
		},
		{
			name: "a PR that declares no card gets no status",
			edit: func(rw *reviewWorld) {
				rw.prBody = "Pipeline fix, no card.\n"
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b"),
				}
			},
			wantState: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			w := newGHWorld(t)
			rw := defaultReviewWorld()
			c.edit(&rw)
			w.seedReview(rw)
			out, calls, outputs := w.run(script, env)
			state, desc := publishedStatus(calls)
			if state != c.wantState {
				t.Fatalf("anchors/review = %q, want %q\noutput:\n%s\ncalls:\n%s", state, c.wantState, out, calls)
			}
			if c.wantState != "" {
				if !strings.Contains(calls, "FIELD context=anchors/review") {
					t.Errorf("the status must be published under the context `anchors/review`:\n%s", calls)
				}
				if !strings.Contains(outputs, "state="+c.wantState) {
					t.Errorf("the job output must carry the state for `mover`, got:\n%s", outputs)
				}
			}
			if c.wantDesc != "" && !strings.Contains(desc, c.wantDesc) {
				t.Errorf("description %q should name %q", desc, c.wantDesc)
			}
		})
	}
}

// The job must actually be wired: a script nobody triggers publishes nothing.
func TestPRChecksReviewJobIsWired(t *testing.T) {
	doc := loadPRChecks(t)
	if _, ok := doc.On["issue_comment"]; !ok {
		t.Error("`anchors-pr-checks.yml` does not listen to `issue_comment` — the verdict line " +
			"is a PR comment, and without the trigger `anchors/review` never leaves pending")
	}
	if doc.Perm["statuses"] != "write" {
		t.Error("`anchors-pr-checks.yml` needs `statuses: write` to publish `anchors/review`")
	}
	review := doc.Jobs["review"]
	if !strings.Contains(review.If, "issue_comment") || !strings.Contains(review.If, "anchors-review:") {
		t.Errorf("the `review` job should run on PR comments carrying a verdict line, if: %q", review.If)
	}
	mover := doc.Jobs["mover"]
	if needs, _ := mover.Needs.(string); needs != "review" {
		t.Errorf("`mover` must run after `review` (needs: review) to know the outcome at the merge, got %v", mover.Needs)
	}
	// `always()`: a review job that failed or was skipped must not hold the card.
	if !strings.Contains(mover.If, "always()") {
		t.Errorf("`mover` must run even when `review` failed or was skipped, if: %q", mover.If)
	}
	// A comment is a verdict, and moves no card.
	if !strings.Contains(mover.If, "!= 'issue_comment'") {
		t.Errorf("`mover` must not run on a PR comment, if: %q", mover.If)
	}
	env := mover.Steps[0].Env
	if !strings.Contains(env["REVIEW_STATE"], "needs.review.outputs.state") ||
		!strings.Contains(env["REVIEWER"], "needs.review.outputs.reviewer") {
		t.Errorf("`mover` must receive the review outcome from the `review` job, env: %v", env)
	}
}

// The verdict line the pipeline parses is the one the review guide teaches. Two copies of
// one format drift the first time one of them changes, and here drifting means every
// reviewed PR stays pending.
func TestPRChecksVerdictLineMatchesTheReviewGuide(t *testing.T) {
	guide, err := os.ReadFile(filepath.Join("..", "..", "cmd", "anchors", "governance", "guide_review.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range []string{"anchors-review: approved by <you>", "anchors-review: rejected by <you>"} {
		if !strings.Contains(string(guide), line) {
			t.Errorf("the review guide no longer teaches %q — the line `anchors-pr-checks.yml` parses", line)
		}
	}
	if !strings.Contains(reviewScript(t), `^anchors-review: (approved|rejected) by `) {
		t.Error("the `review` job no longer parses `anchors-review: approved|rejected by <reviewer>`")
	}
}

// moverScript returns the `mover` script with the `${{ … }}` expressions it inlines replaced
// by the values of one event.
func moverScript(t *testing.T, expr map[string]string) string {
	t.Helper()
	job := loadPRChecks(t).Jobs["mover"]
	script := job.Steps[0].Run
	return regexp.MustCompile(`\$\{\{\s*([^}]*?)\s*\}\}`).ReplaceAllStringFunc(script, func(m string) string {
		k := strings.TrimSpace(m[3 : len(m)-2])
		return expr[k]
	})
}

func (w *ghWorld) seedMerge(labels ...string) {
	w.fixture("pr_view_7", map[string]any{
		"body": "Refs #12\n", "title": "some work", "headRefOid": "abc123",
	})
	var ls []map[string]string
	for _, l := range append([]string{"anchors"}, labels...) {
		ls = append(ls, map[string]string{"name": l})
	}
	w.fixture("issue_view_12", map[string]any{"labels": ls, "comments": []any{}})
	w.fixture("pr_checks_7", []map[string]string{{"bucket": "pass"}})
}

var mergedEvent = map[string]string{
	"inputs.pr": "7", "github.event.action": "closed", "github.event.pull_request.merged": "true",
}

func TestPRChecksMergeWithoutReviewOutcomeIsSaidOnTheCard(t *testing.T) {
	script := moverScript(t, mergedEvent)

	t.Run("merged with anchors/review success: the usual move, no warning", func(t *testing.T) {
		w := newGHWorld(t)
		w.seedMerge("anchors:in-review")
		_, calls, _ := w.run(script, map[string]string{"REVIEW_STATE": "success", "REVIEWER": "agent-b"})
		if !strings.Contains(calls, "FIELD add-label=anchors:ready-to-test") {
			t.Fatalf("the merge should move the card to ready-to-test:\n%s", calls)
		}
		if strings.Contains(calls, "WITHOUT the review outcome") {
			t.Errorf("a reviewed merge must not be flagged:\n%s", calls)
		}
	})

	t.Run("merged while pending: the card moves AND says who was reviewing", func(t *testing.T) {
		w := newGHWorld(t)
		w.seedMerge("anchors:in-review")
		_, calls, _ := w.run(script, map[string]string{"REVIEW_STATE": "pending", "REVIEWER": "agent-b"})
		if !strings.Contains(calls, "FIELD add-label=anchors:ready-to-test") {
			t.Fatalf("the card must still move to ready-to-test — the work landed:\n%s", calls)
		}
		if !strings.Contains(calls, "PR #7 merged WITHOUT the review outcome") {
			t.Fatalf("the card must say the PR merged without the review outcome:\n%s", calls)
		}
		if !strings.Contains(calls, "reviewer: **agent-b**") || !strings.Contains(calls, "`pending`") {
			t.Errorf("the comment must name the reviewer and the status it had:\n%s", calls)
		}
	})

	t.Run("merged while rejected: flagged too", func(t *testing.T) {
		w := newGHWorld(t)
		w.seedMerge("anchors:in-review")
		_, calls, _ := w.run(script, map[string]string{"REVIEW_STATE": "failure", "REVIEWER": "agent-b"})
		if !strings.Contains(calls, "WITHOUT the review outcome") || !strings.Contains(calls, "`failure`") {
			t.Errorf("a merge over a rejection must be flagged on the card:\n%s", calls)
		}
	})

	t.Run("merged before any review was assigned: flagged, saying no reviewer", func(t *testing.T) {
		w := newGHWorld(t)
		w.seedMerge("anchors:ready-to-review")
		_, calls, _ := w.run(script, map[string]string{"REVIEW_STATE": "", "REVIEWER": ""})
		if !strings.Contains(calls, "WITHOUT the review outcome") {
			t.Fatalf("a merge with no review status must be flagged:\n%s", calls)
		}
		if !strings.Contains(calls, "no reviewer had been assigned") || !strings.Contains(calls, "not reported") {
			t.Errorf("the comment must say no reviewer was assigned and no status reported:\n%s", calls)
		}
	})
}

// The move to `ready-to-review` is when the review becomes OWED: `mover` publishes
// `anchors/review` pending right there, because the `review` job ran before the move.
func TestPRChecksGreenPRPublishesReviewPending(t *testing.T) {
	script := moverScript(t, map[string]string{"inputs.pr": "7", "github.event.action": "synchronize"})
	w := newGHWorld(t)
	w.seedMerge("anchors:in-progress")
	_, calls, _ := w.run(script, map[string]string{"REVIEW_STATE": "none", "REVIEWER": ""})
	if !strings.Contains(calls, "FIELD add-label=anchors:ready-to-review") {
		t.Fatalf("a green PR should move the card to ready-to-review:\n%s", calls)
	}
	state, desc := publishedStatus(calls)
	if state != "pending" || !strings.Contains(calls, "FIELD context=anchors/review") {
		t.Errorf("the move to ready-to-review must publish anchors/review pending, got %q:\n%s", state, calls)
	}
	if !strings.Contains(desc, "awaiting a reviewer") {
		t.Errorf("the pending status should say it awaits a reviewer, got %q", desc)
	}
}

// The claim hands the reviewer the instructions, and they must teach the SAME line the
// `review` job parses. The claim used to say "aceito → move the card to ready-to-test",
// which bypasses the status and, done before the merge, hides the #782 comment.
func TestClaimTeachesTheReviewVerdictLine(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{"anchors-review: approved by <você>", "anchors-review: rejected by <você>"} {
		if !strings.Contains(s, want) {
			t.Errorf("the reviewer instructions do not teach %q", want)
		}
	}
	if strings.Contains(s, "**aceito** → mova o card") {
		t.Error("the reviewer is still told to move the card on approval")
	}
}

// THE VERDICT FREES THE REVIEWER: after an approved/rejected line, the card no longer names
// the reviewer as owner, so `anchors next` stops resuming it. Pending changes nothing, and
// a card already released is not released twice.
func TestPRChecksVerdictReleasesTheReviewer(t *testing.T) {
	script := reviewScript(t)
	env := map[string]string{"PR": "7", "HEAD_SHA": "abc123"}
	release := "BODY anchors-owner: (liberado)"

	for _, verdict := range []string{"approved", "rejected"} {
		w := newGHWorld(t)
		rw := defaultReviewWorld()
		rw.prComments = []map[string]any{prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: "+verdict+" by agent-b")}
		w.seedReview(rw)
		_, calls, _ := w.run(script, env)
		if !strings.Contains(calls, "CALL issue comment 12") || !strings.Contains(calls, release) {
			t.Errorf("%s verdict did not release the reviewer:\n%s", verdict, calls)
		}
	}

	w := newGHWorld(t)
	w.seedReview(defaultReviewWorld()) // no verdict yet
	if _, calls, _ := w.run(script, env); strings.Contains(calls, release) {
		t.Errorf("a pending review released the reviewer:\n%s", calls)
	}

	w = newGHWorld(t)
	rw := defaultReviewWorld()
	rw.owners = append(rw.owners, map[string]string{"createdAt": "2026-09-20T11:30:00Z", "body": "anchors-owner: (liberado) — revisão concluída: approved por agent-b"})
	rw.prComments = []map[string]any{prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b")}
	w.seedReview(rw)
	if _, calls, _ := w.run(script, env); strings.Contains(calls, release) {
		t.Errorf("an already released card was released again:\n%s", calls)
	}
}

// blue-eyes #988/#989: the body opens with `Refs #13` (a related card, itself in the
// queue) and ends with `Closes #12` (the card under review). The closing line wins in the
// review job AND in the mover, whatever the order in the body.
func TestPRChecksClosingLineWinsOverRefs(t *testing.T) {
	const body = "Refs #13 · the pattern is there\n\nImplements the thing.\n\nCloses #12\n"
	related := map[string]any{
		"labels":   []map[string]string{{"name": "anchors"}, {"name": "anchors:ready-to-review"}},
		"comments": []any{},
	}

	t.Run("review: the verdict counts for the closed card", func(t *testing.T) {
		w := newGHWorld(t)
		rw := defaultReviewWorld()
		rw.prBody = body
		rw.prComments = []map[string]any{
			prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b"),
		}
		w.seedReview(rw)
		w.fixture("issue_view_13", related)
		out, calls, _ := w.run(reviewScript(t), map[string]string{"PR": "7", "HEAD_SHA": "abc123"})
		if state, desc := publishedStatus(calls); state != "success" || !strings.Contains(desc, "card #12") {
			t.Fatalf("anchors/review = %q (%q), want success on card #12\noutput:\n%s", state, desc, out)
		}
	})

	t.Run("mover: the merge moves the closed card, not the referenced one", func(t *testing.T) {
		w := newGHWorld(t)
		w.seedMerge("anchors:in-review")
		w.fixture("pr_view_7", map[string]any{"body": body, "title": "some work", "headRefOid": "abc123"})
		w.fixture("issue_view_13", related)
		_, calls, _ := w.run(moverScript(t, mergedEvent), map[string]string{"REVIEW_STATE": "success", "REVIEWER": "agent-b"})
		if !strings.Contains(calls, "CALL issue edit 12") || strings.Contains(calls, "CALL issue edit 13") {
			t.Fatalf("the merge must move #12 and leave #13 alone:\n%s", calls)
		}
	})
}

// REJECTED sends the work back to its author: the card leaves `in-review` for
// `in-progress`, owned by the author again (the owner before the reviewer), so their
// `anchors next` resumes it. Before, it stayed in `in-review` with no owner and was
// offered to no one (blue-eyes, 2026-09-24). APPROVED leaves the card for the merge.
func TestPRChecksRejectedGoesBackToTheAuthor(t *testing.T) {
	run := func(verdict string) string {
		w := newGHWorld(t)
		rw := defaultReviewWorld()
		rw.prComments = []map[string]any{
			prComment("2026-09-20T11:00:00Z", "OWNER", "## Revisão\n\nFails B02.\n\nanchors-review: "+verdict+" by agent-b\n"),
		}
		w.seedReview(rw)
		_, calls, _ := w.run(reviewScript(t), map[string]string{"PR": "7", "HEAD_SHA": "abc123"})
		return calls
	}

	calls := run("rejected")
	for _, want := range []string{
		"BODY anchors-owner: (liberado) — revisão concluída: rejected por agent-b",
		"FIELD add-label=anchors:in-progress",
		"BODY anchors-owner: agent-a",
	} {
		if !strings.Contains(calls, want) {
			t.Errorf("a rejection should do %q:\n%s", want, calls)
		}
	}

	calls = run("approved")
	if strings.Contains(calls, "anchors:in-progress") || strings.Contains(calls, "BODY anchors-owner: agent-a") {
		t.Errorf("an approval must leave the card for the merge:\n%s", calls)
	}
}
