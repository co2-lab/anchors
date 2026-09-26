package common

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// AgentCard é um card que carrega o nome deste agente no último anchors-owner.
type AgentCard struct {
	Numero string
	Titulo string
	Estado string
}

// AgentCards devolve os cards cujo último dono é este agente.
func AgentCards(cfg *config.Config) []AgentCard {
	agente := strings.TrimSpace(os.Getenv("ANCHORS_AGENT"))
	if agente == "" || cfg == nil || cfg.Workflow == nil {
		return nil
	}
	if _, err := exec.LookPath("gh"); err != nil {
		return nil
	}
	out, err := exec.Command("gh", "issue", "list",
		"--repo", cfg.Workflow.Repo,
		"--state", "open", "--label", cfg.Workflow.Labels[0], "--limit", "200",
		"--json", "number,title,labels,comments",
		"--jq", `[.[] | select([.comments[].body | select(startswith("anchors-owner:"))] | last == "anchors-owner: `+agente+`")
		         | {n: .number, t: .title, e: ([.labels[].name | select(startswith("anchors:"))] | .[0] // "")}] | .[] | "\(.n)\t\(.t)\t\(.e)"`,
	).Output()
	if err != nil {
		return nil
	}
	var cards []AgentCard
	// Trim only the line breaks: the last line of a card with no state label ENDS in a
	// tab (its empty third field), and TrimSpace over the whole output ate it — that card
	// alone, in that position alone, was dropped.
	for _, l := range strings.Split(strings.Trim(string(out), "\r\n"), "\n") {
		p := strings.Split(strings.TrimRight(l, "\r"), "\t")
		if len(p) != 3 || p[0] == "" {
			continue
		}
		cards = append(cards, AgentCard{
			Numero: p[0],
			Titulo: p[1],
			Estado: strings.TrimPrefix(p[2], "anchors:"),
		})
	}
	return cards
}

// FirstLineOfReason faz o título da issue, que é uma linha.
func FirstLineOfReason(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	if len(s) > 70 {
		return s[:67] + "..."
	}
	return s
}

// NumeroDaIssue extrai o número da issue de uma URL retornada pelo gh.
func NumeroDaIssue(url string) string {
	i := strings.LastIndex(url, "/")
	if i < 0 || i+1 >= len(url) {
		return ""
	}
	n := strings.TrimSpace(url[i+1:])
	for _, r := range n {
		if r < '0' || r > '9' {
			return ""
		}
	}
	return n
}

var vinculoNoCorpoRE = regexp.MustCompile(`(?mi)^\s*(?:refs|closes|fixes|resolves)\s+#(\d+)`)

// CardDoPR lê o card que um PR declara no corpo, pelo Refs/Closes.
func CardDoPR(repo, pr string) string {
	pr = strings.TrimPrefix(strings.TrimSpace(pr), "#")
	if repo == "" || pr == "" {
		return ""
	}
	out, err := exec.Command("gh", "pr", "view", pr,
		"--repo", repo, "--json", "body", "--jq", ".body // \"\"").Output()
	if err != nil {
		return ""
	}
	m := vinculoNoCorpoRE.FindStringSubmatch(string(out))
	if m == nil {
		return ""
	}
	return m[1]
}
